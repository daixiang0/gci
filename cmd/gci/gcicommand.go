package gci

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/log"
	"github.com/daixiang0/gci/pkg/section"
)

type processingFunc = func(args []string, gciCfg config.Config) error

func (e *Executor) newGciCommand(use, short, long string, aliases []string, stdInSupport bool, processingFunc processingFunc) *cobra.Command {
	var noInlineComments, noPrefixComments, skipGenerated, debug *bool
	var sectionStrings, sectionSeparatorStrings *[]string
	cmd := cobra.Command{
		Use:               use,
		Aliases:           aliases,
		Short:             short,
		Long:              long,
		ValidArgsFunction: goFileCompletion,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmtCfg := config.BoolConfig{
				NoInlineComments: *noInlineComments,
				NoPrefixComments: *noPrefixComments,
				Debug:            *debug,
				SkipGenerated:    *skipGenerated,
			}
			gciCfg, err := config.YamlConfig{Cfg: fmtCfg, SectionStrings: *sectionStrings, SectionSeparatorStrings: *sectionSeparatorStrings}.Parse()
			if err != nil {
				return err
			}
			if *debug {
				log.SetLevel(zapcore.DebugLevel)
			}
			return processingFunc(args, *gciCfg)
		},
	}
	if !stdInSupport {
		cmd.Args = cobra.MinimumNArgs(1)
	}

	// register command as subcommand
	e.rootCmd.AddCommand(&cmd)

	debug = cmd.Flags().BoolP("debug", "d", false, "Enables debug output from the formatter")

	sectionHelp := `Sections define how inputs will be processed. Section names are case-insensitive and may contain parameters in (). A section can contain a Prefix and a Suffix section which is delimited by ":". These sections can be used for formatting and will only be rendered if the main section contains an entry. The Section order is the same as below, default value is [Standard,Default].
Std | Standard - Captures all standard packages if they do not match another section
Prefix(github.com/daixiang0) | pkgPrefix(github.com/daixiang0) - Groups all imports with the specified Prefix. Imports will be matched to the longest Prefix.
Def | Default - Contains all imports that could not be matched to another section type
[DEPRECATED] Comment(your text here) | CommentLine(your text here) - Prints the specified indented comment
[DEPRECATED] NL | NewLine - Prints an empty line`

	skipGenerated = cmd.Flags().Bool("skip-generated", false, "Skip generated files")

	sectionStrings = cmd.Flags().StringSliceP("section", "s", nil, sectionHelp)

	// deprecated
	noInlineComments = cmd.Flags().Bool("NoInlineComments", false, "Drops inline comments while formatting")
	cmd.Flags().MarkDeprecated("NoInlineComments", "Drops inline comments while formatting")
	noPrefixComments = cmd.Flags().Bool("NoPrefixComments", false, "Drops comment lines above an import statement while formatting")
	cmd.Flags().MarkDeprecated("NoPrefixComments", "Drops inline comments while formatting")
	sectionSeparatorStrings = cmd.Flags().StringSliceP("SectionSeparator", "x", section.DefaultSectionSeparators().String(), "SectionSeparators are inserted between Sections")
	cmd.Flags().MarkDeprecated("SectionSeparator", "Drops inline comments while formatting")
	cmd.Flags().MarkDeprecated("x", "Drops inline comments while formatting")

	return &cmd
}
