import chess
import chess.pgn
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import Dataset, DataLoader

class FastChessDataset(Dataset):
    def __init__(self, pgn_file, max_games=None):
        # Store only active indices (integers) instead of massive dense float tensors
        self.w_indices = []
        self.b_indices = []
        self.targets = []
        self.stm = []

        game_count = 0
        pos_count = 0

        with open(pgn_file, 'r') as f:
            while True:
                game = chess.pgn.read_game(f)
                if game is None:
                    break

                result = game.headers.get("Result")
                if result not in ["1-0", "0-1", "1/2-1/2"]:
                    continue

                wdl = 1.0 if result == "1-0" else (0.0 if result == "0-1" else 0.5)

                board = game.board()
                for move in game.mainline_moves():
                    board.push(move)

                    # Skip opening book moves and positions in check
                    if len(board.move_stack) < 10 or board.is_check():
                        continue

                    # Extract active feature indices only (sparse representation)
                    w_idx, b_idx = self._get_feature_indices(board)
                    self.w_indices.append(w_idx)
                    self.b_indices.append(b_idx)
                    self.targets.append(wdl)
                    self.stm.append(1.0 if board.turn == chess.WHITE else 0.0)

                    pos_count += 1

                game_count += 1
                if game_count % 100 == 0:
                    print(f"Parsed {game_count} games ({pos_count} positions loaded)...")

                if max_games and game_count >= max_games:
                    break

        print(f"Finished loading {game_count} games with {pos_count} total positions.")

    def _get_feature_indices(self, board):
        w_active = []
        b_active = []
        for sq, piece in board.piece_map().items():
            # Map python-chess pieces (1-6) to Go layout (0-5, 6-11)
            p_idx = (piece.piece_type - 1) + (0 if piece.color == chess.WHITE else 6)

            # White Perspective
            w_active.append(p_idx * 64 + sq)

            # Black Perspective (Flip vertical, swap colors)
            b_sq = sq ^ 56
            b_p_idx = (p_idx + 6) % 12
            b_active.append(b_p_idx * 64 + b_sq)

        return w_active, b_active

    def __len__(self):
        return len(self.targets)

    def __getitem__(self, idx):
        # Create the dense 768 tensor on-the-fly per item/batch
        w_tensor = torch.zeros(768, dtype=torch.float32)
        b_tensor = torch.zeros(768, dtype=torch.float32)

        w_tensor[self.w_indices[idx]] = 1.0
        b_tensor[self.b_indices[idx]] = 1.0

        target = torch.tensor([self.targets[idx]], dtype=torch.float32)
        stm = torch.tensor([self.stm[idx]], dtype=torch.float32)

        return w_tensor, b_tensor, target, stm


class NNUE(nn.Module):
    def __init__(self):
        super().__init__()
        self.feature = nn.Linear(768, 256)
        self.output = nn.Linear(512, 1)

    def forward(self, w_feat, b_feat, stm):
        w_acc = torch.clamp(self.feature(w_feat), 0.0, 1.0)
        b_acc = torch.clamp(self.feature(b_feat), 0.0, 1.0)

        us = torch.where(stm == 1.0, w_acc, b_acc)
        them = torch.where(stm == 1.0, b_acc, w_acc)

        return self.output(torch.cat([us, them], dim=1))


# --- Execution ---
print("Loading games...")
# Increase batch size for much faster training throughput on CPU
dataset = FastChessDataset("/home/ec2-user/selftest-1.7.pgn")
loader = DataLoader(dataset, batch_size=256, shuffle=True, num_workers=2)

model = NNUE()
optimizer = optim.Adam(model.parameters(), lr=1e-3)
criterion = nn.BCEWithLogitsLoss()

print("Start training...")
for epoch in range(10):
    total_loss = 0.0
    for w_feat, b_feat, target, stm in loader:
        optimizer.zero_grad()
        output = model(w_feat, b_feat, stm)
        loss = criterion(output, target)
        loss.backward()
        optimizer.step()
        total_loss += loss.item()
    print(f"Epoch {epoch+1}, Loss: {total_loss/len(loader):.4f}")

print("Training Complete")
torch.save(model.state_dict(), "nnue_float_weights.pt")
print("Model saved to nnue_float_weights.pt")
