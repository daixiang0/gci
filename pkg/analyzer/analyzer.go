package analyzer

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/daixiang0/gci/pkg/io"
	"github.com/daixiang0/gci/pkg/log"
)

const (
	NoInlineCommentsFlag  = "noInlineComments"
	NoPrefixCommentsFlag  = "noPrefixComments"
	SectionsFlag          = "Sections"
	SectionSeparatorsFlag = "SectionSeparators"
	SectionDelimiter      = ","
)

var (
	noInlineComments     bool
	noPrefixComments     bool
	sectionsStr          string
	sectionSeparatorsStr string
)

func init() {
	Analyzer.Flags.BoolVar(&noInlineComments, NoInlineCommentsFlag, false, "If comments in the same line as the input should be present")
	Analyzer.Flags.BoolVar(&noPrefixComments, NoPrefixCommentsFlag, false, "If comments above an input should be present")
	Analyzer.Flags.StringVar(&sectionsStr, SectionsFlag, "", "Specify the Sections format that should be used to check the file formatting")
	Analyzer.Flags.StringVar(&sectionSeparatorsStr, SectionSeparatorsFlag, "", "Specify the Sections that are inserted as Separators between Sections")

	log.InitLogger()
	defer log.L().Sync()
}

var Analyzer = &analysis.Analyzer{
	Name:     "gci",
	Doc:      "A tool that control golang package import order and make it always deterministic.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runAnalysis,
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

	// read configuration options
	gciCfg, err := parseGciConfiguration()
	if err != nil {
		return nil, err
	}

	for _, file := range fileReferences {
		filePath := file.Name()
		unmodifiedFile, formattedFile, err := gci.LoadFormatGoFile(io.File{filePath}, *gciCfg)
		if err != nil {
			return nil, err
		}
		// search for a difference
		fileRunes := bytes.Runes(unmodifiedFile)
		formattedRunes := bytes.Runes(formattedFile)
		diffIdx := compareRunes(fileRunes, formattedRunes)
		switch diffIdx {
		case -1:
			// no difference
		default:
			pass.Reportf(file.Pos(diffIdx), "fix by `%s %s`", generateCmdLine(*gciCfg), filePath)
		}
	}
	return nil, nil
}

func compareRunes(a, b []rune) (differencePos int) {
	// check shorter rune slice first to prevent invalid array access
	shorterRune := a
	if len(b) < len(a) {
		shorterRune = b
	}
	// check for differences up to where the length is identical
	for idx := 0; idx < len(shorterRune); idx++ {
		if a[idx] != b[idx] {
			return idx
		}
	}
	// check that we have compared two equally long rune arrays
	if len(a) != len(b) {
		return len(shorterRune) + 1
	}
	return -1
}

func parseGciConfiguration() (*gci.GciConfiguration, error) {
	fmtCfg := configuration.FormatterConfiguration{noInlineComments, noPrefixComments, false}

	var sectionStrings []string
	if sectionsStr != "" {
		sectionStrings = strings.Split(sectionsStr, SectionDelimiter)
	}

	var sectionSeparatorStrings []string
	if sectionSeparatorsStr != "" {
		sectionSeparatorStrings = strings.Split(sectionSeparatorsStr, SectionDelimiter)
		fmt.Println(sectionSeparatorsStr)
	}
	return gci.GciStringConfiguration{fmtCfg, sectionStrings, sectionSeparatorStrings}.Parse()
}

func generateCmdLine(cfg gci.GciConfiguration) string {
	result := "gci write"

	if cfg.FormatterConfiguration.NoInlineComments {
		result += " --NoInlineComments "
	}

	if cfg.FormatterConfiguration.NoPrefixComments {
		result += " --NoPrefixComments "
	}

	for _, s := range cfg.Sections.String() {
		result += fmt.Sprintf(" --Section \"%s\" ", s)
	}
	for _, s := range cfg.SectionSeparators.String() {
		result += fmt.Sprintf(" --SectionSeparator %s ", s)
	}
	return result
}
