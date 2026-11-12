package main

import (
	"context"
	"time"

	"github.com/Groennesby/anubis-agent-core/pkg/anubis"
)

// App struct binds the Go Bicameral Engine to the Wails v2 / Native WebView frontend.
type App struct {
	ctx    context.Context
	Engine *anubis.BicameralEngine
}

func NewApp(workspaceDir string) (*App, error) {
	eng, err := anubis.NewBicameralEngine(workspaceDir)
	if err != nil {
		return nil, err
	}
	return &App{
		ctx:    context.Background(),
		Engine: eng,
	}, nil
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// GetState returns the complete live Bicameral state for the GUI panels.
func (a *App) GetState() map[string]any {
	recall := a.Engine.RightHemisphere.HopfieldAssociativeReload("Guardian Law REA Kali OPSEC Media", 8)
	verdict := a.Engine.Scales.WeighHeart("gui_heartbeat", "status check")
	return map[string]any{
		"active_model":     a.Engine.ActiveModel,
		"session_id":       a.Engine.SessionID,
		"session_turns":    len(a.Engine.SessionHistory),
		"last_user_prompt": a.Engine.LastUserPrompt,
		"gpu_backend":      a.Engine.BackendType,
		"adaptive_memory":  a.Engine.RightHemisphere.Adaptive,
		"synapse_recall":   recall,
		"last_verdict":     verdict,
		"duat_tasks":       a.Engine.Duat.Tasks,
		"censorship":       "0% (100% Local Sovereign & Loyal Best Friend)",
		"timestamp":        time.Now().Format(time.RFC3339),
	}
}
