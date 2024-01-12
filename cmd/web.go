package cmd

import (
	"fmt"

	cfg "github.com/slotopol/server/config"
	"github.com/slotopol/server/spi"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const webShort = "Starts web-server"
const webLong = ``
const webExmp = `
  %s web`

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:     "web",
	Short:   webShort,
	Long:    webLong,
	Example: fmt.Sprintf(webExmp, cfg.AppName),
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg.DevMode {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		var r = spi.Router()
		RunWeb(r)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}

func RunWeb(r *gin.Engine) {
	r.Run()
}
