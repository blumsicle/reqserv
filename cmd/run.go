/*
Copyright © 2024 Brian Blumberg <blumsicle@icloud.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start server",
	RunE: func(cmd *cobra.Command, args []string) error {
		host := viper.GetString("host")
		port := viper.GetString("port")
		serverName := viper.GetString("server-name")
		serverMode := viper.GetString("server-mode")

		gin.SetMode(serverMode)

		router := gin.Default()
		router.NoRoute(func(ctx *gin.Context) {
			body, _ := io.ReadAll(ctx.Request.Body)
			ctx.Request.Body.Close()

			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

			var data map[string]any
			if err := ctx.ShouldBind(&data); err != nil {
				data = make(map[string]any)
				data["error"] = err.Error()
				data["body"] = string(body)
			}

			ctx.JSON(http.StatusOK, gin.H{
				"headers":    ctx.Request.Header,
				"url":        ctx.Request.URL,
				"data":       data,
				"serverName": serverName,
			})
		})

		return router.Run(host + ":" + port)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().String("host", "", "host")
	runCmd.PersistentFlags().String("port", "8080", "port")
	runCmd.PersistentFlags().String("server-name", "reqserv", "server name")
	runCmd.PersistentFlags().String("server-mode", "debug", "server mode")
	_ = viper.BindPFlag("host", runCmd.PersistentFlags().Lookup("host"))
	_ = viper.BindPFlag("port", runCmd.PersistentFlags().Lookup("port"))
	_ = viper.BindPFlag("server-name", runCmd.PersistentFlags().Lookup("server-name"))
	_ = viper.BindPFlag("server-mode", runCmd.PersistentFlags().Lookup("server-mode"))
}
