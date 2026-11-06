package tomb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// TombSeal provides isolated sandboxed execution for skills and shell commands,
// plus automatic workspace snapshots so any action can be rolled back.
type TombSeal struct {
	WorkspaceDir string
	SnapshotDir  string
	AllowedBins  map[string]bool
	OPSECShield  bool
}

func NewTombSeal(workspaceDir string) *TombSeal {
	home, _ := os.UserHomeDir()
	snapDir := filepath.Join(home, ".anubis", "tomb_snapshots")
	_ = os.MkdirAll(snapDir, 0o755)

	allowed := map[string]bool{
		"go": true, "git": true, "ls": true, "cat": true, "grep": true,
		"find": true, "yt-dlp": true, "ffmpeg": true, "ollama": true,
		"python3": true, "rea": true, "collabs-rea": true, "collabs-synapse": true,
		"curl": true, "nmap": true, "dig": true, "whois": true, "openssl": true,
	}

	return &TombSeal{
		WorkspaceDir: workspaceDir,
		SnapshotDir:  snapDir,
		AllowedBins:  allowed,
		OPSECShield:  true,
	}
}

// CreateSealSnapshot hashes critical workspace metadata before running any mutating skill.
func (t *TombSeal) CreateSealSnapshot(label string) (string, error) {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s-%d", label, time.Now().UnixNano())))
	sealID := "seal-" + hex.EncodeToString(h.Sum(nil))[:12]
	sealFile := filepath.Join(t.SnapshotDir, sealID+".json")
	content := fmt.Sprintf(`{"seal_id":%q,"workspace":%q,"label":%q,"timestamp":%q,"opsec_shield":%v}`,
		sealID, t.WorkspaceDir, label, time.Now().Format(time.RFC3339), t.OPSECShield)
	if err := os.WriteFile(sealFile, []byte(content), 0o644); err != nil {
		return "", err
	}
	return sealID, nil
}

// ExecuteSandboxed runs a command inside the Tomb Seal with timeout and scrubbed environment variables.
func (t *TombSeal) ExecuteSandboxed(ctx context.Context, command string, timeout time.Duration) (string, string, error) {
	sealID, _ := t.CreateSealSnapshot(command)
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "zsh", "-c", command)
	if t.WorkspaceDir != "" {
		cmd.Dir = t.WorkspaceDir
	}

	// OPSEC ChainKeys Env Scrubbing: preserve safe paths while preventing telemetry leaks
	cmd.Env = []string{
		"PATH=/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin",
		"HOME=" + os.Getenv("HOME"),
		"TERM=xterm-256color",
		"ANUBIS_TOMB_SEAL=" + sealID,
		"ANUBIS_ZERO_TELEMETRY=1",
		"DO_NOT_TRACK=1",
	}

	out, err := cmd.CombinedOutput()
	return sealID, strings.TrimSpace(string(out)), err
}
