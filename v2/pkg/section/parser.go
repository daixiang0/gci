package section

import (
	"fmt"
	"regexp"
	"strings"
)

var sectionParser = regexp.MustCompile(`^([^(]+)(?:\(([^)]*)\))?$`)

func Parse(sectionStrings []string) (SectionList, error) {
	if len(sectionStrings) == 0 {
		return nil, nil
	}

	sections := make(SectionList, 0, len(sectionStrings))
	for _, sectionString := range sectionStrings {
		sectionString = strings.TrimSpace(sectionString)
		if sectionString == "" {
			continue
		}

		matches := sectionParser.FindStringSubmatch(sectionString)
		if matches == nil {
		return nil, fmt.Errorf("invalid params: %s", strings.ToLower(sectionString))
		}

		sectionType := strings.ToLower(matches[1])
		sectionParams := matches[2]

		var section Section
		switch sectionType {
		case StandardType:
			section = Standard{}
		case DefaultType:
			section = Default{}
		case "prefix", CustomType:
			if sectionParams == "" {
				return nil, fmt.Errorf("prefix section requires parameters")
			}
			section = Custom{Prefix: sectionParams}
		case BlankType:
			section = Blank{}
		case DotType:
			section = Dot{}
		case AliasType:
			section = Alias{}
		case LocalModuleType:
			section = &LocalModule{}
		case NewLineType:
			section = NewLine{}
		default:
			return nil, fmt.Errorf("unknown section type: %s", sectionType)
		}

		sections = append(sections, section)
	}

	if len(sections) == 0 {
		return nil, nil
	}
	return sections, nil
}
