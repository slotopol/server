package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/slotopol/server/api"
	cfg "github.com/slotopol/server/config"
	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const webShort = "Starts web-server"
const webLong = ``
const webExmp = `
  %[1]s web`

// webCmd represents the `web` command
var webCmd = &cobra.Command{
	Use:     "web",
	Short:   webShort,
	Long:    webLong,
	Example: fmt.Sprintf(webExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var debug bool
		var err error
		var pf = cmd.Flags()

		if debug, err = pf.GetBool("debug"); err != nil {
			cfg.Fatalf(err.Error())
			return
		}
		if debug {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		var exitctx = Startup()

		// Load yaml-files
		if _, err = LoadInternalYaml(exitctx); err != nil {
			cfg.Fatalf("can not load internal yaml files: %s", err.Error())
			return
		}
		if err = LoadExternalYaml(exitctx); err != nil {
			cfg.Fatalf("can not load external yaml files: %s", err.Error())
			return
		}
		UpdateAlgList()
		CheckAlgList()

		// Working with SQL
		if err = InitSQL(); err != nil {
			cfg.Fatalf("can not initialize SQL: %s", err.Error())
			return
		}
		defer func() {
			if err = DoneSQL(); err != nil {
				cfg.Fatalf("can not done SQL: %s", err.Error())
				return
			}
		}()
		go SqlLoop(exitctx)

		// Web router engine
		var r = gin.New()
		r.SetTrustedProxies(Cfg.TrustedProxies)
		r.HandleMethodNotAllowed = true
		api.SetupRouter(r)

		// Starts HTTP listeners
		var wg errgroup.Group
		for _, addr := range Cfg.PortHTTP {
			cfg.Infof("start http on %s", addr)
			var srv = http.Server{
				Addr:              addr,
				Handler:           r.Handler(),
				ReadTimeout:       Cfg.ReadTimeout,
				ReadHeaderTimeout: Cfg.ReadHeaderTimeout,
				WriteTimeout:      Cfg.WriteTimeout,
				IdleTimeout:       Cfg.IdleTimeout,
				MaxHeaderBytes:    Cfg.MaxHeaderBytes,
			}

			wg.Go(func() (err error) {
				var ctx, cancel = context.WithCancel(context.Background())
				go func() {
					defer cancel()
					// service connections
					if err = srv.ListenAndServe(); err != nil {
						if err != http.ErrServerClosed {
							err = fmt.Errorf("failed to serve on %s: %w", addr, err)
							return
						}
						err = nil
					}
					cfg.Infof("stop http on %s", addr)
				}()

				select {
				case <-ctx.Done():
				case <-exitctx.Done():
					// create a deadline to wait for.
					var ctx, cancel = context.WithTimeout(context.Background(), Cfg.ShutdownTimeout)
					defer cancel()

					if err = srv.Shutdown(ctx); err != nil {
						err = fmt.Errorf("shutdown http on %s: %w", addr, err)
						return
					}
				}
				return
			})
		}
		if err = wg.Wait(); err != nil {
			cfg.Errorf("error occurred while waiting for HTTP servers to shut down: %s", err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)

	var f = webCmd.Flags()
	f.Bool("debug", false, "run gin-gonic in debug mode")
}
