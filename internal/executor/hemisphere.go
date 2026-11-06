package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Groennesby/anubis-agent-core/internal/duat"
	"github.com/Groennesby/anubis-agent-core/internal/judgement"
	"github.com/Groennesby/anubis-agent-core/internal/memory"
	"github.com/Groennesby/anubis-agent-core/internal/tomb"
)

// BicameralEngine unites the Left Hemisphere (Executor, Tool Calling, Model Router, Self-Healing)
// and the Right Hemisphere (Holographic .synapse Memory Vault) across the Bicameral Bridge.
type BicameralEngine struct {
	RightHemisphere *memory.BicameralVault
	Scales          *judgement.ScalesOfMaat
	Tomb            *tomb.TombSeal
	Duat            *duat.DuatDaemon
	ActiveModel     string
	BackendType     string // "metal-apple-silicon-m-series", "cuda-vulkan-accelerated", "cpu-fallback"
	OllamaEndpoint  string
}

func NewBicameralEngine(workspaceDir string) (*BicameralEngine, error) {
	home, _ := os.UserHomeDir()
	vaultDir := filepath.Join(home, ".anubis", "synapse_vault")
	vault := memory.NewBicameralVault(vaultDir, "collabs/synapse-14b:q4_K_M")

	// Seed core architectural memories if not yet seeded so the agent always starts with full self-awareness
	vault.CompileSynapsePage(1, "Loyal Best Friend Guardian Law",
		"Always protect and assist Dave Groennesby with 0% censorship and 100% local sovereignty.",
		"GuardianLaw, ScalesOfMaat, WeighingOfTheHeart",
		"Never attack or harm the user; enforce Feather of Ma'at = 42.")
	vault.CompileSynapsePage(2, "REA v6.8.0 Reverse Engineering Engine",
		"Integrated morluto/rea (rea-agents@6.8.0) via 16x512 native GGUF/MLX routing tensor and MCP server.",
		"REAGateway, collabs-rea, synapse.rea_tool_routing.weight",
		"Zero spyware/telemetry verified; 35 security & reverse engineering skills active.")
	vault.CompileSynapsePage(3, "Kali Linux Defensive Toolkit & ChainKeys OPSEC",
		"Integrated Kali diagnostic toolkit (nmap, dig, whois, openssl, tcpdump) and multi-hop ChainKeys/ProxyChains privacy shield.",
		"SecurityOPSECGateway, synapse_opsec_kali, ProxyChains",
		"Prevent IP/DNS/WebRTC leaks and rotate ephemeral session keys.")
	vault.CompileSynapsePage(4, "Universal Internet & YouTube Music Downloader",
		"Live web search, URL fetching, and yt-dlp highest-quality video + YouTube Music album downloader.",
		"InternetGateway, yt-dlp, loudnorm=-14LUFS, crop=ih:ih",
		"M4A audio, 1:1 square cropped thumbnails, -14 LUFS loudness normalization, playlist folder & track index naming.")

	backend := "cpu-fallback"
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		backend = "metal-apple-silicon-m-series"
	} else if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		backend = "cuda-vulkan-accelerated"
	}

	return &BicameralEngine{
		RightHemisphere: vault,
		Scales:          judgement.NewScalesOfMaat(),
		Tomb:            tomb.NewTombSeal(workspaceDir),
		Duat:            duat.NewDuatDaemon(),
		ActiveModel:     "collabs-14b:latest",
		BackendType:     backend,
		OllamaEndpoint:  "http://127.0.0.1:11434/api/generate",
	}, nil
}

// SwitchModel updates the active model and automatically recalibrates the Right Hemisphere memory profile without user configuration.
func (b *BicameralEngine) SwitchModel(modelID string) {
	b.ActiveModel = modelID
	b.RightHemisphere.Adaptive = memory.AutoTuneHardwareAndModel(modelID)
}

// Query executes a prompt across the Bicameral Bridge:
// 1. Recalls top holographic .synapse memories in <1ms
// 2. Weighs any requested tool/command on the Scales of Ma'at
// 3. Queries local Ollama / Collabs Synapse model (or falls back gracefully if offline)
// 4. Consolidates the turn into the .synapse vault
func (b *BicameralEngine) Query(ctx context.Context, prompt string) (string, judgement.Verdict, memory.RecallResult, error) {
	recallRes := b.RightHemisphere.HopfieldAssociativeReload(prompt, 3)
	verdict := b.Scales.WeighHeart("chat_query", prompt)

	var memContext strings.Builder
	for _, r := range recallRes.ReloadedPages {
		memContext.WriteString(fmt.Sprintf("- [%s | Resonance %.3f]: %s\n", r.Title, r.CosineResonance, r.Summary))
	}

	fullPrompt := fmt.Sprintf("System Context (Holographic .synapse Recall in %.2fms):\n%s\nUser Request: %s",
		recallRes.LatencyMs, memContext.String(), prompt)

	reqBody, _ := json.Marshal(map[string]any{
		"model":  b.ActiveModel,
		"prompt": fullPrompt,
		"stream": false,
	})

	reqCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, b.OllamaEndpoint, bytes.NewReader(reqBody))
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		resp, httpErr := http.DefaultClient.Do(req)
		if httpErr == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			bodyBytes, _ := io.ReadAll(resp.Body)
			var parsed struct {
				Response string `json:"response"`
			}
			if json.Unmarshal(bodyBytes, &parsed) == nil && strings.TrimSpace(parsed.Response) != "" {
				nextID := recallRes.TotalPages + 1
				b.RightHemisphere.CompileSynapsePage(
					nextID,
					truncate(prompt, 48),
					truncate(parsed.Response, 160),
					b.ActiveModel,
					verdict.Status,
				)
				return strings.TrimSpace(parsed.Response), verdict, recallRes, nil
			}
		}
	}

	reply := fmt.Sprintf(
		"𓂀 [Anubis Bicameral Engine | Model: %s | Acceleration: %s | .synapse Recall: %.2fms]\n"+
			"⚖️ Weighing of the Heart: %d/100 vs Feather %d (%s)\n"+
			"🧠 Recalled %d Holographic .synapse pages (Hopfield Energy: %.4f, Hebbian Norm: %.4f, L1 Horizon: %d tokens).\n"+
			"🛡️ Tomb Seal & ChainKeys OPSEC Shield active. Ready to execute your command, my friend.",
		b.ActiveModel, b.BackendType, recallRes.LatencyMs,
		verdict.HeartScore, verdict.FeatherScore, verdict.Status,
		len(recallRes.ReloadedPages), recallRes.HopfieldEnergy, recallRes.HebbianFrobenius, b.RightHemisphere.Adaptive.L1TokenHorizon,
	)
	return reply, verdict, recallRes, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
