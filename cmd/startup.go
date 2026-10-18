package cmd

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

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/game"
	"github.com/slotopol/server/util"
)

var (
	Cfg = cfg.Cfg // shortcut
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

func LoadYamlFromFile(ctx context.Context, fpath string) (err error) {
	if ext := util.ToLower(filepath.Ext(fpath)); ext != ".yaml" && ext != ".yml" {
		return nil
	}
	var r io.ReadCloser
	if r, err = os.Open(fpath); err != nil {
		return err
	}
	defer r.Close()
	var cr = ctxReader{ctx: ctx, r: r}
	var count int
	if count, err = game.ReadChain(cr); err != nil {
		return fmt.Errorf("can not read data from %s: %w", fpath, err)
	}
	cfg.Debugf("loaded %d objects from: %s", count, fpath)
	return nil
}

// Load data from extermal yaml files.
func LoadExternalYaml(ctx context.Context) (err error) {
	for _, fpath := range FinalPaths {
		if err = LoadYamlFromFile(ctx, fpath); err != nil {
			return
		}
	}
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
