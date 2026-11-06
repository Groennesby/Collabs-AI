package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Groennesby/anubis-agent-core/internal/executor"
)

func main() {
	statusFlag := flag.Bool("status", false, "Print Holographic Bicameral Core status as JSON")
	weighFlag := flag.String("weigh", "", "Weigh a command against the Feather of Ma'at (0-100)")
	duatFlag := flag.Bool("duat", false, "Generate Du'at Mode Dawn Scroll report")
	queryFlag := flag.String("query", "", "Send a prompt through the Bicameral Bridge")
	modelFlag := flag.String("model", "collabs-14b:latest", "Active model tag (collabs-14b:latest, collabs-7b:latest, etc.)")
	flag.Parse()

	cwd, _ := os.Getwd()
	engine, err := executor.NewBicameralEngine(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Anubis Core error: %v\n", err)
		os.Exit(1)
	}
	engine.SwitchModel(*modelFlag)

	if *weighFlag != "" {
		v := engine.Scales.WeighHeart("cli_weigh", *weighFlag)
		fmt.Println(v.RenderGauge(30))
		b, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(b))
		return
	}

	if *duatFlag {
		scrollPath, err := engine.Duat.WriteMorningScroll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Du'at error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("𓂀 Du'at Dawn Scroll written to: %s\n", scrollPath)
		return
	}

	if *queryFlag != "" {
		reply, verdict, recallRes, _ := engine.Query(context.Background(), *queryFlag)
		fmt.Println(verdict.RenderGauge(28))
		fmt.Printf("⚡ .synapse Recall Latency: %.2fms (Hopfield Energy: %.4f)\n\n%s\n",
			recallRes.LatencyMs, recallRes.HopfieldEnergy, reply)
		return
	}

	recallStats := engine.RightHemisphere.HopfieldAssociativeReload("Guardian Law REA Kali OPSEC Media", 4)
	status := map[string]any{
		"engine":          "Anubis Agent Core v1.0.0 (Go Bicameral Engine)",
		"active_model":    engine.ActiveModel,
		"gpu_backend":     engine.BackendType,
		"adaptive_memory": engine.RightHemisphere.Adaptive,
		"synapse_recall":  recallStats,
		"feather_of_maat": 42,
		"duat_tasks":      engine.Duat.Tasks,
		"censorship":      "0% (100% Local Sovereign & Loyal Best Friend Law)",
	}
	b, _ := json.MarshalIndent(status, "", "  ")
	if *statusFlag || flag.NFlag() == 0 {
		fmt.Println(string(b))
	}
}
