package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ModelSpec struct {
	ID            string  `json:"id"`
	OllamaTag     string  `json:"ollama_tag"`
	Format        string  `json:"format"` // "GGUF+MLX", "GGUF", "MLX", "CLOUD"
	ParamsBillion float64 `json:"params_billion"`
	Quantization  string  `json:"quantization"`
	EstimatedRAM  string  `json:"estimated_ram_gb"`
	Uncensored    bool    `json:"uncensored"`
	Preloaded     bool    `json:"preloaded"`
	LocalPath     string  `json:"local_path,omitempty"`
}

type ModelRouter struct {
	Mode       string      `json:"mode"` // "LOCAL_FIRST", "AUTO_FALLBACK", "CLOUD_ONLY"
	ActiveID   string      `json:"active_id"`
	Registry   []ModelSpec `json:"registry"`
	CloudProviders map[string]string `json:"cloud_providers"`
}

func DefaultRegistry() *ModelRouter {
	return &ModelRouter{
		Mode:     "LOCAL_FIRST",
		ActiveID: "collabs/synapse-14b:q4_K_M",
		Registry: []ModelSpec{
			{
				ID:            "collabs/synapse-7b:q4_K_M",
				OllamaTag:     "collabs-7b:latest",
				Format:        "GGUF v3 + Apple Silicon MLX",
				ParamsBillion: 7.62,
				Quantization:  "Q4_K_M",
				EstimatedRAM:  "5.2 GB",
				Uncensored:    true,
				Preloaded:     true,
				LocalPath:     "/Users/macm1pro/Documents/Collabs Local AI 7b/gguf/collabs-synapse-7b-q4_k_m.gguf",
			},
			{
				ID:            "collabs/synapse-14b:q4_K_M",
				OllamaTag:     "collabs-14b:latest",
				Format:        "GGUF v3 + Apple Silicon MLX",
				ParamsBillion: 14.8,
				Quantization:  "Q4_K_M",
				EstimatedRAM:  "9.4 GB",
				Uncensored:    true,
				Preloaded:     true,
				LocalPath:     "/Users/macm1pro/Documents/Collabs Local AI 14b/gguf/collabs-synapse-14b-q4_k_m.gguf",
			},
			{
				ID:            "collabs/synapse-14b:q8_0",
				OllamaTag:     "collabs-14b:deepthink",
				Format:        "GGUF v3 High-Precision + System-2 DeepThink",
				ParamsBillion: 14.8,
				Quantization:  "Q8_0 / DeepThink",
				EstimatedRAM:  "15.6 GB",
				Uncensored:    true,
				Preloaded:     true,
				LocalPath:     "/Users/macm1pro/Documents/Collabs Local AI 14b/mlx/model.safetensors",
			},
		},
		CloudProviders: map[string]string{
			"gemini":      "Google AI Studio (Optional User Key)",
			"claude-code": "Anthropic Claude Code (Optional User Key)",
			"chatgpt":     "OpenAI ChatGPT (Optional User Key)",
		},
	}
}

func EstimateGGUFRAM(sizeBytes int64) string {
	gb := float64(sizeBytes)/(1024*1024*1024) + 0.6 // weights + 600MB KV/Synapse overhead
	return fmt.Sprintf("%.1f GB", gb)
}

func main() {
	importFlag := flag.String("import", "", "Import any local GGUF/MLX path or hf.co/... URL (0% censorship filter)")
	flag.Parse()

	router := DefaultRegistry()
	if *importFlag != "" {
		est := "6.0 GB"
		if st, err := os.Stat(*importFlag); err == nil {
			est = EstimateGGUFRAM(st.Size())
		}
		custom := ModelSpec{
			ID:           *importFlag,
			OllamaTag:    "custom-uncensored:latest",
			Format:       "Custom Uncensored GGUF/MLX",
			Quantization: "Auto-Detected",
			EstimatedRAM: est,
			Uncensored:   true,
			Preloaded:    false,
			LocalPath:    *importFlag,
		}
		if strings.HasPrefix(*importFlag, "hf.co/") || strings.HasPrefix(*importFlag, "https://huggingface.co/") {
			custom.Format = "HuggingFace Direct GGUF/MLX Stream"
		}
		router.Registry = append(router.Registry, custom)
	}

	b, _ := json.MarshalIndent(router, "", "  ")
	fmt.Println(string(b))
}
