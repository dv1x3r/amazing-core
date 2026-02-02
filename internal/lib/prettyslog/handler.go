package prettyslog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

type Handler struct {
	mu      *sync.Mutex
	buf     *bytes.Buffer
	handler slog.Handler
}

func NewHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	buf := &bytes.Buffer{}

	return &Handler{
		mu:  &sync.Mutex{},
		buf: buf,
		handler: slog.NewJSONHandler(buf, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{handler: h.handler.WithAttrs(attrs), buf: h.buf, mu: h.mu}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{handler: h.handler.WithGroup(name), buf: h.buf, mu: h.mu}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = Colorize(LightGray, level)
	case slog.LevelInfo:
		level = Colorize(Cyan, level)
	case slog.LevelWarn:
		level = Colorize(LightYellow, level)
	case slog.LevelError:
		level = Colorize(LightRed, level)
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	message := r.Message
	switch r.Message {
	case "http":
		message = processHTTP(attrs)
	case "gsf":
		message = processGSF(attrs)
	}

	data, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	timestamp := Colorize(LightGray, r.Time.Format("[15:04:05.000]"))
	attributes := Colorize(LightGray, string(data))

	fmt.Println(timestamp, level, message, attributes)
	return nil
}

func (h *Handler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	defer h.buf.Reset()

	if err := h.handler.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	if err := json.Unmarshal(h.buf.Bytes(), &attrs); err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	return attrs, nil
}

func processHTTP(attrs map[string]any) string {
	remoteIP := "-"
	if v, ok := attrs["remote_ip"].(string); ok {
		remoteIP = v
		delete(attrs, "remote_ip")
	}

	status := "-"
	statusColor := Green
	if v, ok := attrs["status"].(float64); ok {
		status = fmt.Sprint(v)
		delete(attrs, "status")
		if v >= 500 {
			statusColor = Red
		} else if v >= 300 {
			statusColor = Yellow
		}
	}

	statusText := "-"
	if v, ok := attrs["status_text"].(string); ok {
		statusText = v
		delete(attrs, "status_text")
	}

	host := "-"
	if v, ok := attrs["host"].(string); ok {
		host = v
		delete(attrs, "host")
	}

	method := "-"
	if v, ok := attrs["method"].(string); ok {
		method = v
		delete(attrs, "method")
	}

	uri := "-"
	if v, ok := attrs["uri"].(string); ok {
		uri = v
		delete(attrs, "uri")
	}

	latency := "-"
	if v, ok := attrs["latency"].(string); ok {
		latency = v
		delete(attrs, "latency")
	}

	return fmt.Sprintf(
		"http %s %s %s %s %s %s",
		remoteIP, Colorize(statusColor, status+" "+statusText), host, method, uri, latency,
	)
}

func processGSF(attrs map[string]any) string {
	remoteIP := "-"
	if v, ok := attrs["remote_ip"].(string); ok {
		remoteIP = v
		delete(attrs, "remote_ip")
	}

	requestID := "-"
	if v, ok := attrs["request_id"].(float64); ok {
		requestID = fmt.Sprint(v)
		delete(attrs, "request_id")
	}

	reqFlags := "-"
	if v, ok := attrs["req_flags"].(float64); ok {
		reqFlags = fmt.Sprint(v)
		delete(attrs, "req_flags")
	}

	resFlags := "-"
	if v, ok := attrs["res_flags"].(float64); ok {
		resFlags = fmt.Sprint(v)
		delete(attrs, "res_flags")
	}

	resultCode := "-"
	if v, ok := attrs["result_code"].(float64); ok {
		resultCode = fmt.Sprint(v)
		delete(attrs, "result_code")
	}

	resultCodeText := "-"
	if v, ok := attrs["result_code_text"].(string); ok {
		resultCodeText = v
		delete(attrs, "result_code_text")
	}

	appCode := "-"
	if v, ok := attrs["app_code"].(float64); ok {
		appCode = fmt.Sprint(v)
		delete(attrs, "app_code")
	}

	appCodeText := "-"
	if v, ok := attrs["app_code_text"].(string); ok {
		appCodeText = v
		delete(attrs, "app_code_text")
	}

	svcClass := "-"
	if v, ok := attrs["svc_class"].(float64); ok {
		svcClass = fmt.Sprint(v)
		delete(attrs, "svc_class")
	}

	svcClassText := "-"
	if v, ok := attrs["svc_class_text"].(string); ok {
		svcClassText = v
		delete(attrs, "svc_class_text")
	}

	msgType := "-"
	if v, ok := attrs["msg_type"].(float64); ok {
		msgType = fmt.Sprint(v)
		delete(attrs, "msg_type")
	}

	msgTypeText := "-"
	if v, ok := attrs["msg_type_text"].(string); ok {
		msgTypeText = v
		delete(attrs, "msg_type_text")
	}

	latency := "-"
	if v, ok := attrs["latency"].(string); ok {
		latency = v
		delete(attrs, "latency")
	}

	return fmt.Sprintf(
		"gsf %s %s %s | ID %s Flags %s %s | %s | %s | %s",
		remoteIP,
		Colorize(Green, resultCode+" "+resultCodeText),
		Colorize(Green, appCode+" "+appCodeText),
		requestID, reqFlags, resFlags,
		svcClass+" "+svcClassText,
		msgType+" "+msgTypeText,
		latency,
	)
}
