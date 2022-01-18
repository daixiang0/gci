package gci

import (
	"fmt"

	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/constants"
	"github.com/daixiang0/gci/pkg/gci"
	sectionsPkg "github.com/daixiang0/gci/pkg/gci/sections"

	"github.com/spf13/cobra"
)

type processingFunc = func(args []string, gciCfg gci.GciConfiguration) error

func (e *Executor) newGciCommand(use, short, long string, aliases []string, stdInSupport bool, processingFunc processingFunc) *cobra.Command {
	var noInlineComments, noPrefixComments, debug *bool
	var sectionStrings, sectionSeparatorStrings *[]string
	cmd := cobra.Command{
		Use:               use,
		Aliases:           aliases,
		Short:             short,
		Long:              long,
		ValidArgsFunction: goFileCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmtCfg := configuration.FormatterConfiguration{*noInlineComments, *noPrefixComments, *debug}
			gciCfg, err := gci.GciStringConfiguration{fmtCfg, *sectionStrings, *sectionSeparatorStrings}.Parse()
			if err != nil {
				return err
			}
			return processingFunc(args, *gciCfg)
		},
	}
	if !stdInSupport {
		cmd.Args = cobra.MinimumNArgs(1)
	}

	// register command as subcommand
	e.rootCmd.AddCommand(&cmd)

	sectionHelp := "Sections define how inputs will be processed. " +
		"Section names are case-insensitive and may contain parameters in (). " +
		fmt.Sprintf("A section can contain a Prefix and a Suffix section which is delimited by %q. ", constants.SectionSeparator) +
		"These sections can be used for formatting and will only be rendered if the main section contains an entry." +
		"\n" +
		sectionsPkg.SectionParserInst.SectionHelpTexts()
	// add flags
	debug = cmd.Flags().BoolP("debug", "d", false, "Enables debug output from the formatter")
	noInlineComments = cmd.Flags().Bool("NoInlineComments", false, "Drops inline comments while formatting")
	noPrefixComments = cmd.Flags().Bool("NoPrefixComments", false, "Drops comment lines above an import statement while formatting")
	sectionStrings = cmd.Flags().StringSliceP("Section", "s", gci.DefaultSections().String(), sectionHelp)
	sectionSeparatorStrings = cmd.Flags().StringSliceP("SectionSeparator", "x", gci.DefaultSectionSeparators().String(), "SectionSeparators are inserted between Sections")
	return &cmd
}
