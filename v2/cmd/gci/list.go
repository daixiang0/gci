package gci

import (
	"github.com/spf13/cobra"

	"github.com/daixiang0/gci/v2/pkg/gci"
)

var listCmd = &cobra.Command{
	Use:   "list path...",
	Short: "Prints the filenames that need to be formatted",
	Long:  `Prints the filenames that need to be formatted. If you want to show the diff use diff instead, and if you want to apply the changes use write instead`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := parseSections(); err != nil {
			return err
		}
		return gci.ListUnFormattedFiles(args, cfg)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
