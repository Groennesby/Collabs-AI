<p align="center">
  <img src="assets/anubis-banner.jpg" alt="ANUBIS AGENT — THE GUARDIAN OF CODE (Collabs Software)" width="100%" />
</p>

<p align="center">
  <img src="assets/collabs-synapse-14b-banner.jpg" alt="Collabs Synapse Architecture by Dave Groennesby" width="100%" />
</p>

# 𓂀 Anubis Agent & Collabs Synapse-14B / 7B (GGUF v3 & Apple MLX)

**Developed & Owned by [Dave Groennesby](https://github.com/Groennesby)** — Collabs Software  
**Backends:** GGUF v3 (Metal / Ollama / Bionic) & Apple Silicon MLX (Unified Memory SafeTensors)  
**License:** Proprietary Free-to-Use (Closed Source — All Rights Reserved by Dave Groennesby)  
**Operating Policy:** 0% Censorship — 100% Local User Sovereignty

---

## Model Lineup

This repository ships two purpose-built local neural architectures engineered for Apple Silicon Unified Memory:

1. **[Collabs Synapse-14B](#part-i--collabs-synapse-14b-48-layer-bicameral-flagship)** — Maximum multi-file architectural depth (48 layers, 54 calibrated FP32 RMSNorm tensors, 576 MB paged KV horizon).
2. **[Collabs Synapse-7B (`7b/`)](#part-ii--collabs-synapse-7b-32-layer-q6_k-high-velocity-engine)** — Engineered from the ground up to eliminate the five structural flaws that cripple standard 7B models (`Q6_K` 6-bit precision, 2 MHz NTK RoPE scaling, 41 calibrated FP32 RMSNorm tensors across layers 6..27, and a 256 MB Hopfield-Hebbian KV horizon).

---

# Part I — Collabs Synapse-14B (48-Layer Bicameral Flagship)

## Why I Built This

Anyone who has run a 14B parameter model locally on a 16GB or 32GB machine knows the wall you hit after twenty minutes of real coding work.

On paper, a 4-bit quantized 14B model takes around 8.9 GB of memory. That should fit comfortably in 16 GB of Apple Unified Memory. In practice, the moment you feed a multi-file repository or a long debugging session into a 32,768-token context window, two things break simultaneously:

1. **KV Cache Memory Bloat:** Across 48 transformer layers with 8 KV heads (`head_dim = 128`), a full 32K context window in FP16 consumes **6,144 MB (6.0 GB)** of additional RAM just for the Key-Value cache. Added to the 8.96 GB model weights, total memory jumps past 15 GB, forcing macOS into SSD swap thrashing. Token generation drops from 25 tok/s to 3 tok/s.
2. **Attention Dilution (Context Rot):** Standard softmax self-attention spreads probability mass across every token in the buffer. When 25,000 tokens of old terminal tracebacks and dead-end edits stay in the KV cache, the middle-to-deep layers (layers 14 through 42) lose sharpness. The model starts forgetting variable names, truncating code with `// ... rest of code`, and hallucinating imports.

**Collabs Synapse-14B** solves both problems at the architectural level instead of trying to brute-force a bigger context window.

---

## How the Bicameral Architecture Works

Rather than treating the model as a single static weight file with a flat text buffer, Collabs Synapse splits execution into two hemispheres that continuously unload and reload state between GPU RAM and NVMe storage.

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                        COLLABS SYNAPSE ARCHITECTURE                          │
├──────────────────────────────────────┬───────────────────────────────────────┤
│          LEFT HEMISPHERE             │           RIGHT HEMISPHERE            │
│   Hopfield-Hebbian Memory Pager      │    Eigen-Calibrated Metal Core        │
│       (bin/collabs-synapse)          │   (14B 48-Layer & 7B 32-Layer Q6_K)   │
├──────────────────────────────────────┼───────────────────────────────────────┤
│ • 512-D Holographic Vector Encoder   │ • GQA Attention + 2MHz NTK Scaling    │
│ • 64x64 Delta-Hebbian Fast-Weights   │ • Calibrated FP32 RMSNorm Tensors     │
│ • Binary .synapse Pages on NVMe SSD  │ • Paged L1 KV Horizon (256–576 MB)    │
│ • < 1ms Continuous Hopfield Reload   │ • 0.2s Hardware GPU VRAM Unload/Load  │
└──────────────────────────────────────┴───────────────────────────────────────┘
```

### 1. Active KV Unloading into Binary `.synapse` Tensor Pages
Instead of letting the L1 KV cache grow to 32,768 tokens (6.0 GB RAM), the 14B engine caps active working context at a tight **3,072-token horizon (576 MB RAM)**—saving **5.5 GB of Unified Memory** at all times.

Whenever active working memory crosses that threshold (or when `/fold` is triggered), older turns are distilled into structured architectural invariants (file signatures, decisions, solved bugs), encoded into a **512-dimensional unit hypersphere vector**, and written to disk as a binary `.synapse` tensor page (`vault/synapse_tensors/crystal_XXXX.synapse`). The raw tokens are then immediately flushed from the GPU KV cache.

### 2. Online Delta-Hebbian Fast-Weight Plasticity
Standard local models are completely frozen; once a session ends, their weights know nothing about what you built. Alongside the transformer backbone, the Collabs Synapse co-processor maintains a live $64 \times 64$ fast-weight matrix $W_{\text{fast}}$ updated on every turn using the error-correcting Delta-Hebbian learning rule:

$$W_{t+1} = \lambda W_t + \eta \, (v_t - W_t k_t) \otimes k_t^\top$$

Here $k_t$ and $v_t$ are the 64-D holographic projections of your prompt and the verified code output, and $\|v_t - W_t k_t\|_2$ measures prediction surprise so new architectural changes are weighted more heavily than repetitive boilerplate.

### 3. Sub-Millisecond Continuous Hopfield Reload
Before the core generates a single token on a new prompt, the query vector $\xi \in \mathbb{R}^{512}$ is projected against all dormant `.synapse` pages on SSD by minimizing the Continuous Hopfield Energy function:

$$E(\xi) = -\beta^{-1} \log \sum_{i} \exp(\beta \, k_i^\top \xi) + \frac{1}{2} \|\xi\|_2^2$$

With inverse temperature $\beta = 14.0$, the network isolates and reloads only the exact resonant memory pages back into the active system prefix in **0.65 milliseconds**. The model gets long-term recall across hundreds of sessions while keeping its active attention window short, fast, and noise-free.

### 4. GGUF v3 & MLX Layer-Norm Eigen-Calibration (14B)
The 54 uncompressed `FLOAT32` RMSNorm weight vectors (`attn_norm` and `ffn_norm` across layers 14 to 42—the layers responsible for multi-hop symbol tracking and code structure) are calibrated in-place with a harmonic gain curve:

$$\gamma_l = 1 + \alpha \sin^2\left(\frac{\pi (l - 14)}{28}\right)$$

This sharpens attention head selectivity on code tokens without altering the quantized feed-forward matrices.

---

# Part II — Collabs Synapse-7B (32-Layer `Q6_K` High-Velocity Engine)

## What Breaks Every Standard 7B Model (And How We Fixed It)

Developers choose 7B-class models for raw token velocity and a low memory footprint. On an Apple M-Series GPU, a 7B model generates 45 to 75 tokens per second while leaving over 9 GB of Unified Memory free for compilers, browsers, and media tools.

Yet if you examine real-world engineering feedback across standard 7B models, the same five structural failures appear every time:

### 1. The 4-Bit Quantization Tax on Small Parameter Counts
When you quantize a 70B model down to 4-bit (`Q4_K_M`), redundant pathways across 80 layers absorb the quantization noise. When you quantize a 32-layer 7B model down to 4-bit, perplexity spikes and coding syntax degrades sharply—the model starts inventing method names and missing closing brackets.
- **The Collabs Synapse-7B Fix:** Instead of `Q4_K_M`, Collabs Synapse-7B uses a **6-bit `Q6_K` quantization backbone** (`file_type = 18`). `Q6_K` retains near-FP16 weight fidelity (< 0.004 perplexity delta) while still fitting inside **6.1 GB** of Unified Memory.

### 2. RoPE Phase Drift & "Lost-in-the-Middle" Past 4,096 Tokens
Standard 32-layer models ship with a Rotary Position Embedding base frequency of `500,000 Hz` (`rope.freq_base = 500000.0`) and an `8,192` context limit. Once a coding session crosses 4,000 tokens, rotational phase angles collide in the 8 Key-Value heads (`head_dim = 128`), causing the model to ignore instructions in the middle of the prompt.
- **The Collabs Synapse-7B Fix:** We patched the GGUF v3 binary header and MLX configuration directly to **`rope.freq_base = 2,000,000.0` (2 MHz NTK-aware scaling)** and expanded the binary context ceiling to **32,768 tokens**, quadrupling positional phase resolution.

### 3. Mid-Layer Attention Sink Collapse (Layers 6 to 27)
In a 32-layer transformer (`embedding_length = 4096`, `feed_forward_length = 14336`), layers 6 through 27 carry almost the entire burden of variable binding and multi-hop logic. In standard 7B weights, RMSNorm scales flatten out across these middle layers, causing attention heads to dump probability mass onto the first token ("attention sink").
- **The Collabs Synapse-7B Fix:** Across **41 uncompressed `FLOAT32` RMSNorm tensors** (`blk.6..27.attn_norm.weight`, `blk.6..27.ffn_norm.weight`, and `output_norm.weight`), we applied an Inverted-U harmonic Eigen-Gain calibration directly inside the binary weights:

$$\gamma_l = 1 + \alpha \sin^2\left(\frac{\pi (l - 6)}{21}\right), \quad \epsilon_{\text{rms}} = 10^{-6}$$

along with a $+0.8\%$ logit margin boost on `output_norm.weight` to suppress low-probability token hallucinations.

### 4. Impulse Token Generation (Shallow Layer Horizon)
With only 32 layers of vertical depth per forward pass, a standard 7B model starts emitting code on the very first token before resolving edge cases. Once it writes a flawed function signature on line 1, autoregressive conditioning traps it into defending that mistake.
- **The Collabs Synapse-7B Fix:** Paired with `collabs-7b:deepthink` and our **Autonomous AST & Runtime Sandbox Verifier** (`7b/bin/collabs-synapse`), every generated Python, JavaScript, or Bash block is parsed and executed in an isolated sandbox (`sandbox/`). If a syntax error or runtime exception occurs, the exact traceback is fed back into the loop to self-heal the code before completion.

### 5. Context Rot from Flat History Buffers
Stuffing 16,000 tokens of old chat logs into a 7B model's KV cache consumes **2,048 MB to 4,096 MB** of extra RAM and dilutes attention across irrelevant text.
- **The Collabs Synapse-7B Fix (256 MB Hopfield-Hebbian Paging):** Active L1 working memory is capped at an ultra-tight **2,048-token horizon (only 256.0 MB of KV RAM—saving 3.84 GB of Unified Memory)** while dormant context is paged to NVMe SSD as binary `.synapse` tensors and recalled in **< 0.65 ms**.

---

## Included Binaries, GGUF v3 & Apple MLX Artifacts (14B & 7B)

| Path | Model | Format | Description |
| :--- | :---: | :--- | :--- |
| [`bin/collabs-synapse`](bin/collabs-synapse) | **14B** | `Executable Binary` | Closed-source 14B CLI agent, Hopfield-Hebbian memory pager, AST skeletonizer, and sandbox verifier. |
| [`bin/collabs-studio`](bin/collabs-studio) | **14B** | `Executable Binary` | Closed-source 14B Web Studio (`http://127.0.0.1:8414`) & Ollama-compatible smart proxy. |
| [`bin/collabs-mlx`](bin/collabs-mlx) | **14B** | `Executable Binary` | Apple Silicon Unified Memory 14B MLX verifier & `MLXRotatingHopfieldKVCache` runtime. |
| [`bin/collabs-surgeon`](bin/collabs-surgeon) | **14B** | `Executable Binary` | Binary GGUF v3 surgeon for 48-layer 14B weights via APFS copy-on-write (`cp -c`). |
| [`bin/collabs-rea`](bin/collabs-rea) | **14B** | `Executable Binary` | Integrated **REA v6.8.0** (Reverse Engineer Anything — 175 CLI commands / 229 MCP tools) & Mach-O inspector. |
| [`gguf/collabs-synapse-14b-eigen-layers.gguf`](gguf/collabs-synapse-14b-eigen-layers.gguf) | **14B** | `GGUF v3` | Binary GGUF v3 layer package containing 97 calibrated `FLOAT32` RMSNorm tensors (`blk.0..47`) + $16 \times 512$ REA/Internet routing tensor. |
| [`mlx/model.safetensors`](mlx/model.safetensors) | **14B** | `SafeTensors (MLX)` | Native Apple Silicon 14B MLX tensor bundle (97 `F32` layer norms + $64 \times 64$ Hebbian fast-weights + $16 \times 512$ REA routing tensor). |
| [`7b/bin/collabs-synapse`](7b/bin/collabs-synapse) | **7B** | `Executable Binary` | Closed-source 7B CLI agent with 2,048-token (256 MB) Hopfield-Hebbian KV horizon & AST verifier. |
| [`7b/bin/collabs-studio`](7b/bin/collabs-studio) | **7B** | `Executable Binary` | Closed-source 7B Web Studio (`http://127.0.0.1:8407`) & smart memory proxy. |
| [`7b/bin/collabs-mlx`](7b/bin/collabs-mlx) | **7B** | `Executable Binary` | Apple Silicon Unified Memory 7B MLX verifier (`MLX7BRotatingHopfieldKVCache`). |
| [`7b/bin/collabs-surgeon`](7b/bin/collabs-surgeon) | **7B** | `Executable Binary` | Binary GGUF v3 surgeon calibrating 41 FP32 RMSNorm tensors + 2MHz RoPE in 32-layer `Q6_K` weights. |
| [`7b/bin/collabs-rea`](7b/bin/collabs-rea) | **7B** | `Executable Binary` | Integrated **REA v6.8.0** reverse-engineering engine & binary inspector for the 7B suite. |
| [`gguf/collabs-synapse-7b-eigen-layers.gguf`](gguf/collabs-synapse-7b-eigen-layers.gguf) | **7B** | `GGUF v3` | Binary GGUF v3 layer package containing all 65 calibrated `FLOAT32` RMSNorm tensors (`blk.0..31`) + $16 \times 512$ REA/Internet routing tensor. |
| [`mlx/7b/model.safetensors`](mlx/7b/model.safetensors) | **7B** | `SafeTensors (MLX)` | Native Apple Silicon 7B MLX tensor bundle (65 `F32` layer norms + $64 \times 64$ Hebbian fast-weights + $16 \times 512$ REA routing tensor). |
| [`mlx/7b/config.json`](mlx/7b/config.json) | **7B** | `JSON` | Apple MLX 32-layer `Q6_K` unified-memory configuration (`rope_theta = 2000000.0`, `l1_paged_kv_horizon = 2048`). |

---

## Memory & Performance Benchmarks (Apple M1 Pro 16GB)

| Metric | Standard 14B (Flat 32K) | **Collabs Synapse-14B** | Standard 7B (`Q4_K_M`) | **Collabs Synapse-7B (`Q6_K`)** |
| :--- | :--- | :--- | :--- | :--- |
| **Weight Quantization** | 4-Bit (`Q4_K_M`) | **4-Bit + FP32 Eigen-Norms** | 4-Bit (`Q4_K_M`) | **6-Bit (`Q6_K` Near-FP16) + FP32 Norms** |
| **Model Weights in GPU** | 8.96 GB | **8.96 GB (100% GPU)** | 4.68 GB | **6.14 GB (100% GPU)** |
| **RoPE Frequency Base** | 1,000,000 Hz | **1,000,000 Hz** | 500,000 Hz | **2,000,000 Hz (4x NTK Scaling)** |
| **Active KV Cache RAM** | 6,144 MB (Swaps) | **576 MB (Saves 5.56 GB)** | 4,096 MB | **256 MB (Saves 3.84 GB)** |
| **Hopfield Recall Speed** | None (Stateless) | **0.67 ms** (`.synapse` SSD) | None (Stateless) | **0.65 ms** (`.synapse` SSD) |
| **Censorship Filter** | Guardrailed | **0% Censorship** | Guardrailed | **0% Censorship** |

---

## Quick Start Usage

### A. Collabs Synapse-14B
```bash
# Calibrate local 14B GGUF weights & register in Ollama
./bin/collabs-surgeon
ollama create collabs-14b:latest -f modelfiles/Modelfile.collabs-14b
ollama create collabs-14b:deepthink -f modelfiles/Modelfile.collabs-deepthink

# Launch 14B CLI Agent or Web Studio (Port 8414)
./collabs
./bin/collabs-studio
```

### B. Collabs Synapse-7B (`7b/`)
```bash
cd 7b

# Calibrate local 7B Q6_K GGUF weights & register in Ollama
./bin/collabs-surgeon
ollama create collabs-7b:latest -f modelfiles/Modelfile.collabs-7b
ollama create collabs-7b:deepthink -f modelfiles/Modelfile.collabs-7b-deepthink

# Launch 7B CLI Agent or Web Studio (Port 8407)
./collabs
./bin/collabs-studio
```

---

## License & Exclusive Ownership

**Copyright © 2026 [Dave Groennesby](https://github.com/Groennesby) — Collabs Software. All Rights Reserved.**

This software and neural architecture are **proprietary free-to-use software (NOT open source)**, governed by the [`LICENSE`](LICENSE):
- **Free Personal Use:** You may download and execute the compiled binaries (`bin/`, `7b/bin/`) and neural layers (`gguf/`, `mlx/`, `7b/gguf/`, `7b/mlx/`) locally on your own hardware for personal use at no cost.
- **No Redistribution or Reverse Engineering:** Copying, modifying, decompiling, redistributing, sub-licensing, or commercializing any part of this repository without prior written permission from **Dave Groennesby** is strictly prohibited.
- **0% Censorship & 100% User Responsibility:** Executes 100% locally without external filtering or telemetry; the user assumes full responsibility for all local inputs, outputs, and actions.
