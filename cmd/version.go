package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/dankomiocevic/VerySimpleServer/internal/build"
)

// NewVersionCommand returns the command to get openfga version.
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Return VerySimpleServer version",
		Long:  "Return VerySimpleServer version",
		RunE:  version,
		Args:  cobra.NoArgs,
	}

	return cmd
}

// print out the built version.
func version(_ *cobra.Command, _ []string) error {
	log.Printf("VerySimpleServer version `%s` build from `%s` on `%s` ", build.Version, build.Commit, build.Date)
	return nil
}
