package gci

import (
	"github.com/spf13/cobra"

	"github.com/daixiang0/gci/v2/pkg/gci"
)

var printCmd = &cobra.Command{
	Use:   "print path...",
	Short: "Print outputs the formatted file",
	Long:  `Print outputs the formatted file. If you want to apply the changes to a file use write instead!`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := parseSections(); err != nil {
			return err
		}
		return gci.PrintFormattedFiles(args, cfg)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
