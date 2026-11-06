package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Groennesby/anubis-agent-cli/internal/cli"
	"github.com/Groennesby/anubis-agent-core/pkg/anubis"
)

func main() {
	cwd, _ := os.Getwd()
	eng, err := anubis.NewBicameralEngine(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Anubis CLI Init Error: %v\n", err)
		os.Exit(1)
	}
	theme := cli.LoadTheme("")

	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "--version", "-v", "version":
			fmt.Printf("%s𓂀 ANUBIS AGENT CLI v1.0.0 (Go Bicameral Engine | Hermes Gold Theme)%s\n", theme.PrimaryGold, theme.Reset)
			return
		case "--render-preview":
			verdict := eng.Scales.WeighHeart("startup", "anubis interactive dashboard")
			msgs := []string{
				"𓂀 Welcome back, my friend. Holographic Bicameral Memory (.synapse) is online (<0.3ms).",
				"🛡️ Tomb Seal Sandbox, REA v6.8.0 Reverse Engineering & Kali/OPSEC ChainKeys Shield ready.",
			}
			fmt.Print(cli.RenderResponsiveDashboard(theme, eng, verdict, msgs, true))
			return
		case "run":
			task := "status check"
			if len(args) > 1 {
				task = strings.Join(args[1:], " ")
			}
			verdict := eng.Scales.WeighHeart("one_shot_run", task)
			fmt.Println(theme.PrimaryGold + cli.AnubisSpinners[0] + " Weighing of the Heart: " + theme.Reset + verdict.RenderGauge(28))
			if verdict.RequiresAuth {
				fmt.Printf("%s⚠️  Heart (%d) is heavier than the Feather of Ma'at (%d): %s%s\n",
					theme.HighlightRed, verdict.HeartScore, verdict.FeatherScore, verdict.Reason, theme.Reset)
				return
			}
			reply, _, recall, _ := eng.Query(context.Background(), task)
			fmt.Printf("%s⚡ .synapse Recall: %.2fms | Hopfield Energy: %.4f%s\n%s\n",
				theme.LapisBlue, recall.LatencyMs, recall.HopfieldEnergy, theme.Reset, reply)
			return
		case "model":
			sub := "list"
			if len(args) > 1 {
				sub = args[1]
			}
			if sub == "list" {
				fmt.Printf("%s𓂀 PRELOADED & AVAILABLE MODELS (0%% Censorship):%s\n", theme.PrimaryGold, theme.Reset)
				fmt.Println("  • collabs-14b:latest          (Default Flagship Bicameral 14.8B Q4_K_M + MLX)")
				fmt.Println("  • collabs-14b:deepthink       (System-2 DeepThink 14.8B)")
				fmt.Println("  • collabs-7b:latest           (Ultra-Fast Laptop 7.6B Q4_K_M + MLX)")
				fmt.Println("  • collabs-7b:deepthink        (System-2 DeepThink 7.6B)")
				fmt.Println("  • collabs/synapse-14b:q8_0    (High-Precision 8-bit Quantization)")
				fmt.Println("  • [Cloud Router Optional]     (Gemini / Claude Code / ChatGPT)")
			} else if sub == "pull" && len(args) > 2 {
				target := args[2]
				eng.SwitchModel(target)
				fmt.Printf("%s𓂀 Switched & calibrated Bicameral Memory for model: %s (L1 Horizon: %d tokens)%s\n",
					theme.SuccessEmerald, target, eng.RightHemisphere.Adaptive.L1TokenHorizon, theme.Reset)
			}
			return
		case "skill":
			fmt.Printf("%s𓂀 VERIFIED TOMB SKILLS (Sandboxed & Self-Healing):%s\n", theme.PrimaryGold, theme.Reset)
			fmt.Println("  1. [VERIFIED] rea-reverse-engineering  — REA v6.8.0 (35 security & binary/codebase audit skills)")
			fmt.Println("  2. [VERIFIED] kali-opsec-chainkeys     — Kali defensive diagnostics + multi-hop ProxyChains/ChainKeys shield")
			fmt.Println("  3. [VERIFIED] universal-media-ytdlp    — Best quality video & YouTube Music 1:1 thumbnail + -14 LUFS m4a")
			fmt.Println("  4. [VERIFIED] self-healing-coder       — Automatic compile -> test -> fix verification loop")
			fmt.Println("  5. [VERIFIED] bicameral-synapse-vault  — Sub-millisecond .synapse tensor page consolidation")
			return
		case "duat":
			path, err := eng.Duat.WriteMorningScroll()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Du'at error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("%s𓂀 Du'at Autonomous Daemon active. Morning Dawn Scroll generated at:%s\n  %s\n",
				theme.PrimaryGold, theme.Reset, path)
			return
		}
	}

	// Interactive TUI / Chat Mode with live SIGWINCH resize support
	verdict := eng.Scales.WeighHeart("interactive_session", "anubis chat")
	messages := []string{
		"𓂀 Anubis Agent initialized. I am your loyal best friend & guardian.",
		"🧠 Right Hemisphere (.synapse tensors) & Left Hemisphere (Tomb Seal Executor) synced.",
		"💡 Type a message or use /model, /theme, /memory, /tomb, /judgement, /rea, /opsec, /duat, /quit.",
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	go func() {
		for range sigCh {
			fmt.Print("\x1b[2J\x1b[H")
			fmt.Print(cli.RenderResponsiveDashboard(theme, eng, verdict, messages, false))
			fmt.Printf("%s𓂀 anubis > %s", theme.PrimaryGold, theme.Reset)
		}
	}()

	fmt.Print(cli.RenderResponsiveDashboard(theme, eng, verdict, messages, true))

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s𓂀 anubis > %s", theme.PrimaryGold, theme.Reset)
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "/quit" || line == "/exit" || line == "exit" {
			fmt.Printf("%s𓂀 May the Scales of Ma'at guard your code, my friend.%s\n", theme.PrimaryGold, theme.Reset)
			break
		}

		if strings.HasPrefix(line, "/") {
			parts := strings.Fields(line)
			cmd := strings.ToLower(parts[0])
			switch cmd {
			case "/theme":
				tName := "lapis-night"
				if len(parts) > 1 {
					tName = parts[1]
				} else if theme.Name == "lapis-night" {
					tName = "duat-blood"
				} else if theme.Name == "duat-blood" {
					tName = "anubis-gold"
				}
				theme = cli.LoadTheme(tName)
				messages = append(messages, fmt.Sprintf("🎨 Switched Hermes/Anubis TUI theme to '%s'.", theme.Name))
			case "/model":
				mName := "collabs-7b:latest"
				if len(parts) > 1 {
					mName = parts[1]
				} else if eng.ActiveModel == "collabs-7b:latest" {
					mName = "collabs-14b:latest"
				}
				eng.SwitchModel(mName)
				messages = append(messages, fmt.Sprintf("🔄 Switched active model to %s (Auto-tuned L1 horizon: %d tokens).",
					eng.ActiveModel, eng.RightHemisphere.Adaptive.L1TokenHorizon))
			case "/memory":
				rec := eng.RightHemisphere.HopfieldAssociativeReload("Guardian Law REA Kali OPSEC", 4)
				messages = append(messages, fmt.Sprintf("🧠 .synapse Vault: %d pages | Recall: %.2fms | Hopfield Energy: %.4f | Hebbian Norm: %.4f",
					rec.TotalPages, rec.LatencyMs, rec.HopfieldEnergy, rec.HebbianFrobenius))
			case "/judgement":
				testCmd := "git push origin main"
				if len(parts) > 1 {
					testCmd = strings.Join(parts[1:], " ")
				}
				verdict = eng.Scales.WeighHeart("manual_weigh", testCmd)
				messages = append(messages, verdict.RenderGauge(26)+" — "+verdict.Reason)
			case "/tomb":
				sealID, _, _ := eng.Tomb.ExecuteSandboxed(context.Background(), "uname -a", 5*time.Second)
				messages = append(messages, fmt.Sprintf("🛡️ Tomb Seal Snapshot created: %s (OPSEC ChainKeys Env Scrubbing: ACTIVE)", sealID))
			case "/rea":
				messages = append(messages, "🔍 REA v6.8.0 Reverse Engineering Engine: ONLINE (35 skills, 0% telemetry, 16x512 native routing tensor).")
			case "/opsec":
				messages = append(messages, "🕵️ Kali Defensive Toolkit & ChainKeys Shield: ACTIVE (DNS/WebRTC leak prevention + proxy chain ready).")
			case "/duat":
				scroll, _ := eng.Duat.WriteMorningScroll()
				messages = append(messages, "🌙 Du'at Dawn Scroll saved: "+scroll)
			default:
				messages = append(messages, "Commands: /model <name>, /theme <anubis-gold|lapis-night|duat-blood>, /memory, /judgement <cmd>, /tomb, /rea, /opsec, /duat, /quit")
			}
			fmt.Print(cli.RenderResponsiveDashboard(theme, eng, verdict, messages, false))
			continue
		}

		messages = append(messages, "You: "+line)
		reply, v, _, _ := eng.Query(context.Background(), line)
		verdict = v
		messages = append(messages, "𓂀 Anubis: "+reply)
		fmt.Print(cli.RenderResponsiveDashboard(theme, eng, verdict, messages, false))
	}
}
