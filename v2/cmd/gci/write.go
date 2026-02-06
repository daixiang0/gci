package gci

import (
	"github.com/spf13/cobra"

	"github.com/daixiang0/gci/v2/pkg/gci"
)

var writeCmd = &cobra.Command{
	Use:   "write path...",
	Short: "Write modifies the specified files in-place",
	Long:  `Write modifies the specified files in-place`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := parseSections(); err != nil {
			return err
		}
		return gci.WriteFormattedFiles(args, cfg)
	},
}

func init() {
	rootCmd.AddCommand(writeCmd)
}
