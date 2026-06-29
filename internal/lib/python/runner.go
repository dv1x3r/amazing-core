package python

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type Runner struct {
	logger  *slog.Logger
	venvDir string
	mu      sync.Mutex
}

type CommandResult struct {
	Stdout []byte
	Stderr []byte
}

func NewRunner(logger *slog.Logger, venvDir string) *Runner {
	return &Runner{
		logger:  logger,
		venvDir: venvDir,
	}
}

func (r *Runner) EnsureVenv(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, err := os.Stat(r.pythonPath()); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if r.logger != nil {
		r.logger.Info("creating python venv", "path", r.venvDir)
	}
	_, err := run(ctx, nil, pythonLauncher(), "-m", "venv", r.venvDir)
	return err
}

func (r *Runner) Run(ctx context.Context, args ...string) (*CommandResult, error) {
	return run(ctx, nil, r.pythonPath(), args...)
}

func (r *Runner) RunScript(ctx context.Context, script []byte, args ...string) (*CommandResult, error) {
	args = append([]string{"-"}, args...)
	return run(ctx, bytes.NewReader(script), r.pythonPath(), args...)
}

func (r *Runner) pythonPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(r.venvDir, "Scripts", "python.exe")
	}
	return filepath.Join(r.venvDir, "bin", "python")
}

func pythonLauncher() string {
	if runtime.GOOS == "windows" {
		return "python"
	}
	return "python3"
}

func run(ctx context.Context, stdin io.Reader, name string, args ...string) (*CommandResult, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return &CommandResult{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}, nil
}
