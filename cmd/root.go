/*
Copyright © 2022 Juha R <kontza@gmail.com>

*/
package cmd

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ideal-doodle",
	Short:   "A simple Go TCP/UDP. You can specify the port on command line.",
	Version: "2.2",
	Run: func(cmd *cobra.Command, args []string) {
		if address, err := cmd.Flags().GetString("address"); err != nil {
			log.Error().Err(err).Msg("Failed to get 'address' from flags:")
		} else {
			log.Info().Str("address", address).Msg("Current")
			if useUdp, err := cmd.Flags().GetBool("udp"); err != nil {
				log.Warn().Err(err).Msg("Failed to get 'udp' from flags:")
				useUdp = false
			} else {
				if useUdp {
					if pc, err := net.ListenPacket("udp4", address); err != nil {
						log.Error().Err(err).Msg("Failed to get a packet conn:")
					} else {
						defer pc.Close()
						buffer := make([]byte, 1024)
						for {
							n, addr, err := pc.ReadFrom(buffer)
							if err != nil {
								break
							}
							log.Info().Int("bytes", n).Str("from", addr.String()).Msg("Packet received")
							deadline := time.Now().Add(10 * time.Second)
							err = pc.SetWriteDeadline(deadline)
							if err != nil {
								break
							}
							n, err = pc.WriteTo(buffer[:n], addr)
							if err != nil {
								break
							}
							log.Info().Int("bytes", n).Str("to", addr.String()).Msg("Packet written")
						}
					}
				} else {
					if listener, err := net.Listen("tcp4", address); err != nil {
						log.Error().Err(err).Msg("Failed to get a listener:")
					} else {
						defer listener.Close()
						parts := strings.Split(listener.Addr().String(), ":")
						log.Info().Msgf("Listening... Now do 'echo -n test | nc %s %s' in another window.", parts[0], parts[1])
						for {
							if conn, err := listener.Accept(); err != nil {
								log.Fatal().Err(err).Msg("Failed to accept connection:")
							} else {
								handleRequest(conn)
							}
						}
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
		conn.Write([]byte("OK\n"))
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
	rootCmd.Flags().BoolP("udp", "u", false, "listen to udp")
}
