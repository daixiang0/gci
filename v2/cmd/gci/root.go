package gci

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/daixiang0/gci/v2/pkg/config"
	"github.com/daixiang0/gci/v2/pkg/section"
)

var (
	cfg       config.Config
	sections  []string
	debugMode bool
)

var rootCmd = &cobra.Command{
	Use:   "gci",
	Short: "GCI, a tool that controls Go package import order",
	Long:  `GCI, a tool that controls Go package import order and makes it always deterministic.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringArrayVarP(&sections, "section", "s", []string{"standard", "default"}, "Sections define how imports will be processed")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Enables debug output")
	rootCmd.PersistentFlags().BoolVar(&cfg.SkipGenerated, "skip-generated", false, "Skip generated files")
	rootCmd.PersistentFlags().BoolVar(&cfg.SkipVendor, "skip-vendor", false, "Skip files inside vendor directory")
	rootCmd.PersistentFlags().BoolVar(&cfg.CustomOrder, "custom-order", false, "Enable custom order of sections")
}

func parseSections() error {
	parsedSections, err := section.Parse(sections)
	if err != nil {
		return err
	}
	if parsedSections == nil {
		parsedSections = section.DefaultSections()
	}
	cfg.Sections = parsedSections
	return nil
}
