/*
Copyright © 2022 Juha R <kontza@gmail.com>

*/
package cmd

import (
	"net"
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
		if address, err := cmd.Flags().GetString("address"); err != nil {
			log.Error().Err(err).Msg("Failed to get 'address' from flags:")
		} else {
			log.Info().Str("address", address).Msg("Current")
			if listener, err := net.Listen("tcp", address); err != nil {
				log.Error().Err(err).Msg("Failed to get a listener:")
			} else {
				defer listener.Close()
				log.Info().Msg("Listening... Now do 'echo -n test | nc localhost 7600' in another window.")
				for {
					if conn, err := listener.Accept(); err != nil {
						log.Fatal().Err(err).Msg("Failed to accept connection:")
					} else {
						handleRequest(conn)
					}
				}
			}
		}
	},
}

/*
Thx:
https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go
*/
func handleRequest(conn net.Conn) {
	defer conn.Close()
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	if payloadLength, err := conn.Read(buf); err != nil {
		log.Error().Err(err).Msg("Error reading:")
	} else {
		log.Info().Str("msg", string(buf[:payloadLength])).Msg("Received")
		// Send a response back to person contacting us.
		conn.Write([]byte("OK"))
	}
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
	rootCmd.Flags().StringP("address", "a", "localhost:7600", "specify the address and port to listen to")
}
