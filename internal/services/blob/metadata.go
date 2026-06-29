package blob

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const unityPyPackage = "git+https://github.com/dv1x3r/UnityPy.git@test"

func (s *Service) ensureUnityPy(ctx context.Context) error {
	if err := s.python.EnsureVenv(ctx); err != nil {
		return err
	}
	if _, err := s.python.Run(ctx, "-c", "import UnityPy"); err == nil {
		return nil
	}
	if s.logger != nil {
		s.logger.Info("installing python package", "package", unityPyPackage)
	}
	_, err := s.python.Run(ctx, "-m", "pip", "install", unityPyPackage)
	return err
}

func (s *Service) inspectFile(ctx context.Context, cacheScript []byte, path string) ([]byte, error) {
	const op = "blob.Service.inspectFile"
	res, err := s.python.RunScript(ctx, cacheScript, path, "--stdout")
	if err != nil {
		return nil, fmt.Errorf("run cache tool: %w", err)
	}
	if len(res.Stderr) > 0 && s.logger != nil {
		s.logger.Debug(op, "stderr", strings.TrimSpace(string(res.Stderr)))
	}
	metadata := bytes.TrimSpace(res.Stdout)
	if len(metadata) == 0 {
		return nil, errors.New("cache tool returned empty metadata")
	}
	if !json.Valid(metadata) {
		return nil, errors.New("cache tool returned invalid metadata")
	}
	return metadata, nil
}
