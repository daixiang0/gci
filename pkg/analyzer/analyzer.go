package analyzer

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/golangci/modinfo"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/daixiang0/gci/pkg/config"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"
)

const (
	NoInlineCommentsFlag  = "noInlineComments"
	NoPrefixCommentsFlag  = "noPrefixComments"
	SkipGeneratedFlag     = "skipGenerated"
	SectionsFlag          = "Sections"
	SectionSeparatorsFlag = "SectionSeparators"
	SectionDelimiter      = ","
)

var (
	noInlineComments     bool
	noPrefixComments     bool
	skipGenerated        bool
	sectionsStr          string
	sectionSeparatorsStr string
)

func init() {
	Analyzer.Flags.BoolVar(&noInlineComments, NoInlineCommentsFlag, false, "If comments in the same line as the input should be present")
	Analyzer.Flags.BoolVar(&noPrefixComments, NoPrefixCommentsFlag, false, "If comments above an input should be present")
	Analyzer.Flags.BoolVar(&skipGenerated, SkipGeneratedFlag, false, "Skip generated files")
	Analyzer.Flags.StringVar(&sectionsStr, SectionsFlag, "", "Specify the Sections format that should be used to check the file formatting")
	Analyzer.Flags.StringVar(&sectionSeparatorsStr, SectionSeparatorsFlag, "", "Specify the Sections that are inserted as Separators between Sections")

	log.InitLogger()
	defer log.L().Sync()
}

var Analyzer = &analysis.Analyzer{
	Name: "gci",
	Doc:  "A tool that control Go package import order and make it always deterministic.",
	Run:  runAnalysis,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		modinfo.Analyzer,
	},
}

func runAnalysis(pass *analysis.Pass) (interface{}, error) {
	var fileReferences []*token.File
	// extract file references for all files in the analyzer pass
	for _, pkgFile := range pass.Files {
		fileForPos := pass.Fset.File(pkgFile.Package)
		if fileForPos != nil {
			fileReferences = append(fileReferences, fileForPos)
		}
	}
	expectedNumFiles := len(pass.Files)
	foundNumFiles := len(fileReferences)
	if expectedNumFiles != foundNumFiles {
		return nil, InvalidNumberOfFilesInAnalysis{expectedNumFiles, foundNumFiles}
	}

	file, err := modinfo.FindModuleFromPass(pass)
	if err != nil {
		return nil, err
	}

	gciCfg, err := generateGciConfiguration(file.Path).Parse()
	if err != nil {
		return nil, err
	}

	for _, file := range fileReferences {
		filePath := file.Name()
		unmodifiedFile, formattedFile, err := gci.LoadFormatGoFile(io.File{FilePath: filePath}, *gciCfg)
		if err != nil {
			return nil, err
		}
		fix, err := GetSuggestedFix(file, unmodifiedFile, formattedFile)
		if err != nil {
			return nil, err
		}
		if fix == nil {
			// no difference
			continue
		}
		pass.Report(analysis.Diagnostic{
			Pos:            fix.TextEdits[0].Pos,
			Message:        fmt.Sprintf("fix by `%s %s`", generateCmdLine(*gciCfg), filePath),
			SuggestedFixes: []analysis.SuggestedFix{*fix},
		})
	}
	return nil, nil
}

func generateGciConfiguration(modPath string) *config.YamlConfig {
	fmtCfg := config.BoolConfig{
		NoInlineComments: noInlineComments,
		NoPrefixComments: noPrefixComments,
		Debug:            false,
		SkipGenerated:    skipGenerated,
	}

	var sectionStrings []string
	if sectionsStr != "" {
		sectionStrings = strings.Split(sectionsStr, SectionDelimiter)
	}

	var sectionSeparatorStrings []string
	if sectionSeparatorsStr != "" {
		sectionSeparatorStrings = strings.Split(sectionSeparatorsStr, SectionDelimiter)
		fmt.Println(sectionSeparatorsStr)
	}

	return &config.YamlConfig{Cfg: fmtCfg, SectionStrings: sectionStrings, SectionSeparatorStrings: sectionSeparatorStrings, ModPath: modPath}
}

func generateCmdLine(cfg config.Config) string {
	result := "gci write"

	if cfg.BoolConfig.NoInlineComments {
		result += " --NoInlineComments "
	}

	if cfg.BoolConfig.NoPrefixComments {
		result += " --NoPrefixComments "
	}

	if cfg.BoolConfig.SkipGenerated {
		result += " --skip-generated "
	}

	if cfg.BoolConfig.CustomOrder {
		result += " --custom-order "
	}

	if cfg.BoolConfig.NoLexOrder {
		result += " --no-lex-order"
	}

	for _, s := range cfg.Sections.String() {
		result += fmt.Sprintf(" --Section \"%s\" ", s)
	}
	for _, s := range cfg.SectionSeparators.String() {
		result += fmt.Sprintf(" --SectionSeparator %s ", s)
	}
	return result
}
