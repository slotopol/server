package cmd

import (
	"context"
	"fmt"
	"log"
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

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:     "web",
	Short:   webShort,
	Long:    webLong,
	Example: fmt.Sprintf(webExmp, cfg.AppName),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		if cfg.DevMode {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		var exitctx = Startup()
		if err = Init(); err != nil {
			log.Fatalln(err.Error())
			return
		}
		go SqlLoop(exitctx)

		var r = gin.New()
		r.SetTrustedProxies(Cfg.TrustedProxies)
		r.HandleMethodNotAllowed = true
		api.SetupRouter(r)

		// starts HTTP listeners
		var wg errgroup.Group
		for _, addr := range Cfg.PortHTTP {
			log.Printf("start http on %s\n", addr)
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
					log.Printf("stop http on %s\n", addr)
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
			log.Println(err.Error())
			return
		}
		if err = Done(); err != nil {
			log.Fatalln(err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
