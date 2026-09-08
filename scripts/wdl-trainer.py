import torch
import torch.nn as nn
import torch.optim as optim
import numpy as np
import struct
import os
import sys
import datetime

# =========================================================================
# OUTPUT DIRECTORY / LOGGING SETUP
# =========================================================================
OUTPUT_DIR = "output"
os.makedirs(OUTPUT_DIR, exist_ok=True)

LOG_PATH = os.path.join(OUTPUT_DIR, "train.log")
WEIGHTS_PATH = os.path.join(OUTPUT_DIR, "nnue_float_weights.pt")
BIN_PATH = os.path.join(OUTPUT_DIR, "selftest.bin")


class Tee:
    def __init__(self, *streams):
        self.streams = streams

    def write(self, data):
        for s in self.streams:
            s.write(data)
            s.flush()

    def flush(self):
        for s in self.streams:
            s.flush()


_log_file = open(LOG_PATH, "a")
_log_file.write(f"\n===== Run started {datetime.datetime.now().isoformat()} =====\n")
_log_file.flush()
sys.stdout = Tee(sys.__stdout__, _log_file)

# =========================================================================
# CONFIG
# =========================================================================
PACKED_DIR = "./data/packed"   # output of precompute.py
EPOCHS = 25
BATCH_SIZE = 512              # larger batches make sense once on GPU / no python-chess bottleneck
LR = 1e-3
LR_STEP_SIZE = 5
LR_GAMMA = 0.3

S1 = 255.0  # QA
S2 = 64.0   # QB

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
print(f"Using device: {device}")

features_w = np.load(f"{PACKED_DIR}/features_w.npy", mmap_mode='r')
features_b = np.load(f"{PACKED_DIR}/features_b.npy", mmap_mode='r')
offsets = np.load(f"{PACKED_DIR}/offsets.npy")          # small enough to load fully
targets_all = np.load(f"{PACKED_DIR}/targets.npy")      # small enough to load fully
stm_all = np.load(f"{PACKED_DIR}/stm.npy")              # small enough to load fully

num_positions = len(offsets) - 1
print(f"Loaded packed dataset: {num_positions} positions")


class PackedBatchSampler:
    """Yields shuffled batches directly from the packed arrays. Shuffles at
    the position level every epoch (unlike the old streaming buffer, this is
    a true full shuffle since everything except the two big index arrays
    fits comfortably in RAM)."""

    def __init__(self, num_positions, batch_size):
        self.num_positions = num_positions
        self.batch_size = batch_size

    def __iter__(self):
        perm = np.random.permutation(self.num_positions)
        for start in range(0, self.num_positions, self.batch_size):
            yield perm[start:start + self.batch_size]

    def __len__(self):
        return (self.num_positions + self.batch_size - 1) // self.batch_size


def build_batch(indices):
    """Gathers a batch of positions (given by their indices into the
    position-level arrays) into EmbeddingBag-ready flat indices + offsets,
    plus targets and stm. This is the only per-batch CPU work now - no
    board parsing, no python-chess."""
    w_chunks = []
    b_chunks = []
    batch_offsets = [0]

    for i in indices:
        start, end = offsets[i], offsets[i + 1]
        w_chunks.append(features_w[start:end])
        b_chunks.append(features_b[start:end])
        batch_offsets.append(batch_offsets[-1] + (end - start))

    w_flat = torch.from_numpy(np.concatenate(w_chunks).astype(np.int64))
    b_flat = torch.from_numpy(np.concatenate(b_chunks).astype(np.int64))
    batch_offsets_t = torch.tensor(batch_offsets[:-1], dtype=torch.int64)  # EmbeddingBag wants start offsets only

    target = torch.from_numpy(targets_all[indices]).unsqueeze(1)
    stm = torch.from_numpy(stm_all[indices]).unsqueeze(1)

    return w_flat, b_flat, batch_offsets_t, target, stm


class NNUE(nn.Module):
    """Mathematically identical to the old nn.Linear(768,256) feature layer
    applied to a one-hot/multi-hot input, but computed via sparse embedding
    lookups instead of a dense 768-wide matmul per sample - much faster for
    NNUE-style sparse inputs, and avoids materializing dense one-hot tensors."""

    def __init__(self):
        super().__init__()
        self.feature = nn.EmbeddingBag(768, 256, mode='sum')
        self.feature_bias = nn.Parameter(torch.zeros(256))
        self.output = nn.Linear(512, 1)

        bound = 1.0 / (768 ** 0.5)
        nn.init.uniform_(self.feature.weight, -bound, bound)

    def forward(self, w_idx, b_idx, w_offsets, b_offsets, stm):
        w_acc = torch.clamp(self.feature(w_idx, w_offsets) + self.feature_bias, 0.0, 1.0)
        b_acc = torch.clamp(self.feature(b_idx, b_offsets) + self.feature_bias, 0.0, 1.0)

        us = torch.where(stm == 1.0, w_acc, b_acc)
        them = torch.where(stm == 1.0, b_acc, w_acc)

        return self.output(torch.cat([us, them], dim=1))


print(f"Config: EPOCHS={EPOCHS} BATCH_SIZE={BATCH_SIZE} LR={LR} "
      f"LR_STEP_SIZE={LR_STEP_SIZE} LR_GAMMA={LR_GAMMA}")

model = NNUE().to(device)
optimizer = optim.Adam(model.parameters(), lr=LR)
scheduler = optim.lr_scheduler.StepLR(optimizer, step_size=LR_STEP_SIZE, gamma=LR_GAMMA)
criterion = nn.BCEWithLogitsLoss()

sampler = PackedBatchSampler(num_positions, BATCH_SIZE)

for epoch in range(EPOCHS):
    total_loss = 0.0
    batches = 0

    for batch_indices in sampler:
        w_flat, b_flat, batch_offsets, target, stm = build_batch(batch_indices)

        w_flat = w_flat.to(device)
        b_flat = b_flat.to(device)
        batch_offsets = batch_offsets.to(device)
        target = target.to(device)
        stm = stm.to(device)

        optimizer.zero_grad()
        output = model(w_flat, b_flat, batch_offsets, batch_offsets, stm)
        loss = criterion(output, target)
        loss.backward()
        optimizer.step()

        total_loss += loss.item()
        batches += 1

        if batches % 200 == 0:
            print(f"Epoch {epoch+1} | Batches: {batches}/{len(sampler)} | Running Loss: {total_loss/batches:.4f}")

    avg_loss = total_loss / batches if batches > 0 else 0
    print(f"=== Epoch {epoch+1} Complete, Average Loss: {avg_loss:.4f}, LR: {scheduler.get_last_lr()[0]:.6f} ===")
    if avg_loss > 0.68:
        print("WARNING: loss is still near the random-guess baseline (ln(2)=0.693). "
              "Weights may be collapsing towards zero - check data/labels before continuing.")
    scheduler.step()

torch.save(model.state_dict(), WEIGHTS_PATH)

state_dict = torch.load(WEIGHTS_PATH, weights_only=True, map_location="cpu")

print("\n--- Weight statistics (pre-quantization) ---")
collapsed = False
for name, t in state_dict.items():
    mean_abs = t.abs().mean().item()
    print(f"{name}: min={t.min().item():.5f} max={t.max().item():.5f} mean_abs={mean_abs:.5f}")
    if mean_abs < 0.01:
        collapsed = True

if collapsed:
    raise RuntimeError(
        "Trained weights are suspiciously close to zero (mean_abs < 0.01). "
        "This will round to all-zero int16 after quantization and produce a "
        "network that always evaluates 0. Check your training targets and "
        "loss curve before exporting - do not proceed with this checkpoint."
    )

# --- export: EmbeddingBag weight is already [768, 256], exactly the layout
# the .bin format needs (no transpose required, unlike the old nn.Linear) ---
feature_w = (state_dict["feature.weight"].cpu().numpy() * S1).round().astype('<i2')
feature_b = (state_dict["feature_bias"].cpu().numpy() * S1).round().astype('<i2')

output_w = (state_dict["output.weight"].flatten().cpu().numpy() * S2).round().astype('<i2')
output_b = int(np.round(state_dict["output.bias"].item() * S1 * S2))

raw_feature_w_scaled = state_dict["feature.weight"].cpu().numpy() * S1
raw_output_w_scaled = state_dict["output.weight"].flatten().cpu().numpy() * S2

assert np.abs(raw_feature_w_scaled).max() < 32767, \
    "feature weight quantization overflow! Reduce S1 or add weight clipping in training."
assert np.abs(raw_output_w_scaled).max() < 32767, \
    "output weight quantization overflow! Reduce S2 or add weight clipping in training."
assert (feature_w != 0).any(), "quantized feature weights are all zero!"
assert (output_w != 0).any(), "quantized output weights are all zero!"

with open(BIN_PATH, "wb") as f:
    f.write(feature_w.tobytes())
    f.write(feature_b.tobytes())
    f.write(output_w.tobytes())
    f.write(struct.pack("<i", output_b))

print("Successfully exported Little Endian nn.bin!")
print(f"Non-zero feature weights: {(feature_w != 0).sum()} / {feature_w.size}")
print(f"Non-zero output weights: {(output_w != 0).sum()} / {output_w.size}")
print(f"\nWeights: {WEIGHTS_PATH}")
print(f"Bin:     {BIN_PATH}")
print(f"Log:     {LOG_PATH}")

_log_file.close()
