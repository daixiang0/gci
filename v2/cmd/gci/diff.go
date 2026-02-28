package gci

import (
	"github.com/spf13/cobra"

	"github.com/daixiang0/gci/v2/pkg/gci"
)

var diffCmd = &cobra.Command{
	Use:   "diff path...",
	Short: "Diff prints a patch in the style of the diff tool",
	Long:  `Diff prints a patch in the style of the diff tool that contains the required changes to the file to make it adhere to the specified formatting.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := parseSections(); err != nil {
			return err
		}
		return gci.DiffFormattedFiles(args, cfg)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
