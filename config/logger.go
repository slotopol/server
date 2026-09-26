package cfg

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var CountLevMap = map[int]slog.Level{
	0: slog.LevelError,
	1: slog.LevelWarn,
	2: slog.LevelInfo,
	3: slog.LevelDebug,
}

var LevNameMap = map[slog.Level]string{
	slog.LevelError: "ERR",
	slog.LevelWarn:  "WRN",
	slog.LevelInfo:  "INF",
	slog.LevelDebug: "DBG",
}

type LogHandler struct {
	stdout io.Writer
	stderr io.Writer
	level  *slog.LevelVar
}

func NewLogHandler(stdout, stderr io.Writer, level *slog.LevelVar) *LogHandler {
	return &LogHandler{stdout: stdout, stderr: stderr, level: level}
}

func (h *LogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *LogHandler) Handle(_ context.Context, r slog.Record) error {
	var timeStr = r.Time.Format("2006/01/02 15:04:05")
	var out io.Writer
	if r.Level >= slog.LevelError {
		out = h.stderr
	} else {
		out = h.stdout
	}
	_, err := fmt.Fprintf(out, "%s %s %s\n", timeStr, LevNameMap[r.Level], r.Message)
	return err
}

func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler { return h }
func (h *LogHandler) WithGroup(name string) slog.Handler       { return h }

func InitLogger() {
	var handler *slog.Logger
	if !Silent {
		var programLevel = &slog.LevelVar{}
		var lev, ok = CountLevMap[Verbose]
		if !ok {
			lev = slog.LevelDebug
		}
		programLevel.Set(lev)
		handler = slog.New(NewLogHandler(os.Stdout, os.Stderr, programLevel))
	} else {
		handler = slog.New(slog.DiscardHandler)
	}
	slog.SetDefault(handler)
}

func Fatalf(format string, args ...any) {
	if args == nil {
		slog.Error(format)
		os.Exit(1)
	}
	slog.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func Errorf(format string, args ...any) {
	if args == nil {
		slog.Error(format)
		return
	}
	slog.Error(fmt.Sprintf(format, args...))
}

func Warnf(format string, args ...any) {
	if args == nil {
		slog.Warn(format)
		return
	}
	slog.Warn(fmt.Sprintf(format, args...))
}

func Infof(format string, args ...any) {
	if args == nil {
		slog.Info(format)
		return
	}
	slog.Info(fmt.Sprintf(format, args...))
}

func Debugf(format string, args ...any) {
	if args == nil {
		slog.Debug(format)
		return
	}
	slog.Debug(fmt.Sprintf(format, args...))
}
