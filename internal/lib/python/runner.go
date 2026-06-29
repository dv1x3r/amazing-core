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

type Package struct {
	Spec        string
	ImportCheck string
}

type RunOptions struct {
	StreamOutput bool
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

func (r *Runner) Ensure(ctx context.Context, packages ...Package) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.ensureVenv(ctx); err != nil {
		return err
	}
	return r.ensurePackages(ctx, packages)
}

func (r *Runner) ensureVenv(ctx context.Context) error {
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

func (r *Runner) ensurePackages(ctx context.Context, packages []Package) error {
	if len(packages) == 0 {
		return nil
	}

	missing := make([]Package, 0, len(packages))
	for _, pkg := range packages {
		if pkg.Spec == "" {
			continue
		}
		if pkg.ImportCheck == "" {
			missing = append(missing, pkg)
			continue
		}
		if _, err := run(ctx, nil, r.pythonPath(), "-c", pkg.ImportCheck); err != nil {
			missing = append(missing, pkg)
		}
	}
	if len(missing) == 0 {
		return nil
	}

	if r.logger != nil {
		r.logger.Info("installing python packages", "packages", strings.Join(packageSpecs(missing), ", "))
	}
	args := append([]string{"-m", "pip", "install"}, packageSpecs(missing)...)
	if _, err := run(ctx, nil, r.pythonPath(), args...); err != nil {
		return err
	}

	for _, pkg := range packages {
		if pkg.ImportCheck == "" {
			continue
		}
		if _, err := run(ctx, nil, r.pythonPath(), "-c", pkg.ImportCheck); err != nil {
			return fmt.Errorf("verify python package %q: %w", pkg.Spec, err)
		}
	}
	return nil
}

func packageSpecs(packages []Package) []string {
	specs := make([]string, 0, len(packages))
	for _, pkg := range packages {
		if pkg.Spec != "" {
			specs = append(specs, pkg.Spec)
		}
	}
	return specs
}

func (r *Runner) Run(ctx context.Context, args ...string) (*CommandResult, error) {
	return run(ctx, nil, r.pythonPath(), args...)
}

func (r *Runner) RunWithOptions(ctx context.Context, options RunOptions, args ...string) (*CommandResult, error) {
	return runWithOptions(ctx, nil, options, r.pythonPath(), args...)
}

func (r *Runner) RunScript(ctx context.Context, script []byte, args ...string) (*CommandResult, error) {
	return r.RunScriptWithOptions(ctx, script, RunOptions{}, args...)
}

func (r *Runner) RunScriptWithOptions(ctx context.Context, script []byte, options RunOptions, args ...string) (*CommandResult, error) {
	args = append([]string{"-"}, args...)
	return runWithOptions(ctx, bytes.NewReader(script), options, r.pythonPath(), args...)
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
	return runWithOptions(ctx, stdin, RunOptions{}, name, args...)
}

func runWithOptions(ctx context.Context, stdin io.Reader, options RunOptions, name string, args ...string) (*CommandResult, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdin = stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if options.StreamOutput {
		cmd.Stdout = io.MultiWriter(&stdout, os.Stdout)
		cmd.Stderr = io.MultiWriter(&stderr, os.Stderr)
	}
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return &CommandResult{
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}, nil
}
