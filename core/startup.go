package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

var (
	ErrNoClubName = errors.New("name of 'club' database does not provided at data source name")
	ErrNoSpinName = errors.New("name of 'spin' database does not provided at data source name")
)

func Startup() (exitctx context.Context) {
	//var cancel context.CancelFunc
	exitctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Make exit signal on function exit.
		defer cancel()

		var sigint = make(chan os.Signal, 1)
		var sigterm = make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM (Ctrl+/)
		// SIGKILL, SIGQUIT will not be caught.
		signal.Notify(sigint, syscall.SIGINT)
		signal.Notify(sigterm, syscall.SIGTERM)
		// Block until we receive our signal.
		select {
		case <-exitctx.Done():
			if errors.Is(exitctx.Err(), context.DeadlineExceeded) {
				cfg.Infof("shutting down by timeout")
			} else if errors.Is(exitctx.Err(), context.Canceled) {
				cfg.Infof("shutting down by cancel")
			} else {
				cfg.Infof("shutting down by %s", exitctx.Err().Error())
			}
		case <-sigint:
			cfg.Infof("shutting down by break")
		case <-sigterm:
			cfg.Infof("shutting down by process termination")
		}
		signal.Stop(sigint)
		signal.Stop(sigterm)
	}()
	return
}

func MakeFinalPaths(patterns []string) (final []string, err error) {
	for _, pattern := range cfg.ObjPath {
		var pattern = filepath.ToSlash(pattern)
		pattern, err = util.ExpandHomePath(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to expand home path for %q: %w", pattern, err)
		}
		pattern = util.Envfmt(pattern, nil)
		matches, err := doublestar.FilepathGlob(pattern,
			doublestar.WithCaseInsensitive(),
			doublestar.WithFilesOnly(),
			doublestar.WithNoHidden(),
		)
		if err != nil {
			if errors.Is(err, doublestar.ErrPatternNotExist) {
				return nil, fmt.Errorf("path or its base directory does not exist: %q", pattern)
			}
			return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
		}
		final = append(final, matches...)
		cfg.Infof("found %d files for pattern %q", len(matches), pattern)
	}
	return
}

// ctxReader makes check up of context cancellation before each read operation.
type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

func (cr ctxReader) Read(p []byte) (int, error) {
	if err := cr.ctx.Err(); err != nil {
		return 0, err // returns context error (Canceled / DeadlineExceeded)
	}
	return cr.r.Read(p)
}

// Load data from embed yaml chunks.
func LoadInternalYaml(ctx context.Context) (count int, err error) {
	var t0 = time.Now()
	var n, size int
	for _, b := range game.LoadMap {
		var cr = ctxReader{ctx: ctx, r: bytes.NewReader(b)}
		if n, err = game.ReadChain(cr); err != nil {
			return
		}
		size += len(b)
		count += n
	}
	var d = time.Since(t0)
	cfg.Infof("loaded %d embedded yaml files in %s on %d bytes", len(game.LoadMap), d.String(), size)
	return
}

func LoadYamlFromFile(ctx context.Context, fpath string) (count int, err error) {
	var r io.ReadCloser
	if r, err = os.Open(fpath); err != nil {
		return
	}
	defer r.Close()
	var cr = ctxReader{ctx: ctx, r: r}
	if count, err = game.ReadChain(cr); err != nil {
		return count, fmt.Errorf("can not read data from %s: %w", fpath, err)
	}
	return
}

// Load data from extermal yaml files.
func LoadExternalYaml(ctx context.Context) (err error) {
	var final []string
	if final, err = MakeFinalPaths(cfg.ObjPath); err != nil {
		return
	}
	if len(final) == 0 {
		return
	}

	var t0 = time.Now()
	var n int
	for _, fpath := range final {
		var count int
		if n, err = LoadYamlFromFile(ctx, fpath); err != nil {
			return
		}
		count += n
		cfg.Debugf("loaded %d objects from: %s", n, fpath)
	}
	var d = time.Since(t0)
	cfg.Infof("loaded %d external yaml files in %s", len(final), d.String())
	return
}

func UpdateAlgList() {
	for _, ai := range game.AlgList {
		if ai.Update != nil {
			ai.Update(ai)
		}
	}
	if len(game.DataLoaded) == len(game.DataRouter) {
		cfg.Infof("all %d objects loaded", len(game.DataLoaded))
	} else {
		cfg.Infof("loaded %d objects of %d registered", len(game.DataLoaded), len(game.DataRouter))
	}
}

func CheckAlgList() {
	for _, ai := range game.AlgList {
		if len(ai.RTP) == 0 {
			var id string
			if len(ai.Aliases) > 0 {
				id = ai.Aliases[0].ID()
			}
			panic(fmt.Errorf("RTP list does not complete for %s", id))
		}
	}
}
