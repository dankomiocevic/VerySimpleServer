// Package run contains the command to run an instance of VerySimpleServer.
package run

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dankomiocevic/VerySimpleServer/internal/config"
	"github.com/dankomiocevic/VerySimpleServer/internal/server"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the VSS server",
		Long:  "Run an instance of VerySimpleServer.",
		Run:   run,
		Args:  cobra.NoArgs,
	}

	defaultConfig := config.DefaultConfig()
	flags := cmd.Flags()
	flags.String("addr", defaultConfig.TcpAddr, "the host:port address to serve the server on")
	viper.BindPFlag("addr", cmd.Flags().Lookup("addr"))

	return cmd
}

func run(_ *cobra.Command, _ []string) {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	if err := config.Verify(); err != nil {
		panic(err)
	}

	s := server.NewServer(config)
	defer s.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("Shutting down server..")
}
