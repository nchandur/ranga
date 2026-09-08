import chess
import chess.pgn
import numpy as np
import re
import math
import argparse

K = 400.0
LAMBDA = 0.2
MAX_EVAL_CP = 3000.0


def parse_eval_cp(comment):
    if not comment:
        return None
    m = re.match(r'^\s*([+-]?\d+\.?\d*)', comment.strip())
    if not m:
        return None
    cp = float(m.group(1)) * 100.0
    return max(-MAX_EVAL_CP, min(MAX_EVAL_CP, cp))


def stable_sigmoid(x):
    if x >= 0:
        z = math.exp(-x)
        return 1.0 / (1.0 + z)
    else:
        z = math.exp(x)
        return z / (1.0 + z)


def get_feature_indices(board):
    """Same indexing scheme as the trainer - flip-white, unflipped-black,
    matching the engine's A8=0 mailbox convention. Do not change this
    without also changing utils.go's getFeatureIndex."""
    w_active = []
    b_active = []
    for sq, piece in board.piece_map().items():
        p_idx = (piece.piece_type - 1) + (0 if piece.color == chess.WHITE else 6)
        go_sq = sq ^ 56

        w_active.append(p_idx * 64 + go_sq)

        b_p_idx = (p_idx + 6) % 12
        b_active.append(b_p_idx * 64 + sq)

    return w_active, b_active


def main(pgn_file, out_dir, max_games=None):
    import os
    os.makedirs(out_dir, exist_ok=True)

    w_feat_chunks = []
    b_feat_chunks = []
    offsets = [0]
    targets = []
    stms = []

    game_count = 0
    position_count = 0

    with open(pgn_file, 'r') as f:
        while True:
            if max_games and game_count >= max_games:
                break

            game = chess.pgn.read_game(f)
            if game is None:
                break

            result = game.headers.get("Result")
            if result not in ["1-0", "0-1", "1/2-1/2"]:
                continue

            wdl = 1.0 if result == "1-0" else (0.0 if result == "0-1" else 0.5)
            board = game.board()

            node = game
            while node.variations:
                node = node.variations[0]
                move = node.move
                board.push(move)

                if len(board.move_stack) < 10 or board.is_check():
                    continue

                cp = parse_eval_cp(node.comment)
                if cp is None:
                    continue

                mover_was_white = not board.turn
                white_pov_cp = cp if mover_was_white else -cp

                eval_wdl_white = stable_sigmoid(white_pov_cp / K)
                blended_white = LAMBDA * wdl + (1.0 - LAMBDA) * eval_wdl_white

                stm_target = blended_white if board.turn == chess.WHITE else (1.0 - blended_white)

                w_idx, b_idx = get_feature_indices(board)

                w_feat_chunks.append(np.array(w_idx, dtype=np.int32))
                b_feat_chunks.append(np.array(b_idx, dtype=np.int32))
                offsets.append(offsets[-1] + len(w_idx))
                targets.append(stm_target)
                stms.append(1.0 if board.turn == chess.WHITE else 0.0)

                position_count += 1
                if position_count % 500_000 == 0:
                    print(f"Processed {position_count} positions ({game_count} games)...")

            game_count += 1

    print(f"Done parsing. Total games: {game_count}, total positions: {position_count}")

    features_w = np.concatenate(w_feat_chunks).astype(np.int32)
    features_b = np.concatenate(b_feat_chunks).astype(np.int32)
    offsets = np.array(offsets, dtype=np.int64)
    targets = np.array(targets, dtype=np.float32)
    stms = np.array(stms, dtype=np.float32)

    np.save(f"{out_dir}/features_w.npy", features_w)
    np.save(f"{out_dir}/features_b.npy", features_b)
    np.save(f"{out_dir}/offsets.npy", offsets)
    np.save(f"{out_dir}/targets.npy", targets)
    np.save(f"{out_dir}/stm.npy", stms)

    print(f"Saved packed dataset to {out_dir}/")
    print(f"  features_w.npy: {features_w.shape}, {features_w.nbytes / 1e6:.1f} MB")
    print(f"  features_b.npy: {features_b.shape}, {features_b.nbytes / 1e6:.1f} MB")
    print(f"  offsets.npy:    {offsets.shape}, {offsets.nbytes / 1e6:.1f} MB")
    print(f"  targets.npy:    {targets.shape}, {targets.nbytes / 1e6:.1f} MB")
    print(f"  stm.npy:        {stms.shape}, {stms.nbytes / 1e6:.1f} MB")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--pgn", default="./data/selftest-1.12.pgn")
    parser.add_argument("--out", default="./data/packed")
    parser.add_argument("--max-games", type=int, default=None)
    args = parser.parse_args()
    main(args.pgn, args.out, args.max_games)
