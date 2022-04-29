/*
Copyright © 2022 Juha R <kontza@gmail.com>

*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ideal-doodle",
	Short:   "A simple Go web server. You can specify the port on command line.",
	Version: "1.0",
	Run: func(cmd *cobra.Command, args []string) {
		if port, err := cmd.Flags().GetInt("port"); err != nil {
			log.Error().Err(err).Msg("Failed to get 'port':")
		} else {
			log.Info().Int("port", port).Msg("Current")
			listenPort := fmt.Sprintf(":%d", port)
			http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
				fmt.Fprintf(w, "OK\n")
				log.Info().Str("remote", req.RemoteAddr).Msg("Request from")
			})
			http.ListenAndServe(listenPort, nil)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	zerolog.TimeFieldFormat = "15:04:05.000"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: zerolog.TimeFieldFormat})
	rootCmd.Flags().IntP("port", "p", 7600, "specify the port to listen to")
}
