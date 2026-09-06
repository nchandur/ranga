import chess
import chess.pgn
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import IterableDataset, DataLoader
import numpy as np
import struct

class StreamingChessDataset(IterableDataset):
    def __init__(self, pgn_file, max_games=None):
        self.pgn_file = pgn_file
        self.max_games = max_games

    def _get_feature_indices(self, board):
        w_active = []
        b_active = []
        for sq, piece in board.piece_map().items():
            p_idx = (piece.piece_type - 1) + (0 if piece.color == chess.WHITE else 6)
            go_sq = sq ^ 56

            w_active.append(p_idx * 64 + go_sq)

            b_p_idx = (p_idx + 6) % 12
            b_active.append(b_p_idx * 64 + sq)

        return w_active, b_active

    def __iter__(self):
        game_count = 0
        with open(self.pgn_file, 'r') as f:
            while True:
                if self.max_games and game_count >= self.max_games:
                    break

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

                    if len(board.move_stack) < 10 or board.is_check():
                        continue

                    w_idx, b_idx = self._get_feature_indices(board)

                    w_tensor = torch.zeros(768, dtype=torch.float32)
                    b_tensor = torch.zeros(768, dtype=torch.float32)

                    w_tensor[w_idx] = 1.0
                    b_tensor[b_idx] = 1.0

                    stm_target = wdl if board.turn == chess.WHITE else (1.0 - wdl)

                    target = torch.tensor([stm_target], dtype=torch.float32)
                    stm = torch.tensor([1.0 if board.turn == chess.WHITE else 0.0], dtype=torch.float32)

                    yield w_tensor, b_tensor, target, stm

                game_count += 1

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

print("Starting streaming trainer...")
dataset = StreamingChessDataset("./selftest-1.12.pgn")

loader = DataLoader(dataset, batch_size=256, num_workers=0)

model = NNUE()
optimizer = optim.Adam(model.parameters(), lr=1e-3)
criterion = nn.BCEWithLogitsLoss()

for epoch in range(10):
    total_loss = 0.0
    batches = 0

    for w_feat, b_feat, target, stm in loader:
        optimizer.zero_grad()
        output = model(w_feat, b_feat, stm)
        loss = criterion(output, target)
        loss.backward()
        optimizer.step()

        total_loss += loss.item()
        batches += 1

        if batches % 1000 == 0:
            print(f"Epoch {epoch+1} | Batches: {batches} | Running Loss: {total_loss/batches:.4f}")

    avg_loss = total_loss / batches if batches > 0 else 0
    print(f"=== Epoch {epoch+1} Complete, Average Loss: {avg_loss:.4f} ===")

torch.save(model.state_dict(), "nnue_float_weights.pt")

state_dict = torch.load("nnue_float_weights.pt", weights_only=True)

S1 = 255.0
S2 = 64.0

feature_w = (state_dict["feature.weight"].t().cpu().numpy() * S1).round().astype('<i2')
feature_b = (state_dict["feature.bias"].cpu().numpy() * S1).round().astype('<i2')

output_w = (state_dict["output.weight"].flatten().cpu().numpy() * S2).round().astype('<i2')
output_b = int(np.round(state_dict["output.bias"].item() * S1 * S2))

with open("selftest-1.12.bin", "wb") as f:
    f.write(feature_w.tobytes())
    f.write(feature_b.tobytes())
    f.write(output_w.tobytes())
    f.write(struct.pack("<i", output_b))

print("Successfully exported Little Endian nn.bin!")
