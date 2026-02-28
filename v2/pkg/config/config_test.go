package config

import (
	"reflect"
	"testing"

	"github.com/daixiang0/gci/v2/pkg/section"
)

func TestParseOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"default", "prefix(github/daixiang0/gci)", "prefix(github/daixiang0/gai)"},
	}
	gciCfg, err := cfg.Parse()
	if err != nil {
		t.Fatal(err)
	}
	want := section.SectionList{
		section.Default{},
		section.Custom{Prefix: "github/daixiang0/gai"},
		section.Custom{Prefix: "github/daixiang0/gci"},
	}
	if !reflect.DeepEqual(want, gciCfg.Sections) {
		t.Fatalf("unexpected sections: got=%v want=%v", gciCfg.Sections, want)
	}
}

func TestParseCustomOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"default", "prefix(github/daixiang0/gci)", "prefix(github/daixiang0/gai)"},
		Cfg: BoolConfig{
			CustomOrder: true,
		},
	}
	gciCfg, err := cfg.Parse()
	if err != nil {
		t.Fatal(err)
	}
	want := section.SectionList{
		section.Default{},
		section.Custom{Prefix: "github/daixiang0/gci"},
		section.Custom{Prefix: "github/daixiang0/gai"},
	}
	if !reflect.DeepEqual(want, gciCfg.Sections) {
		t.Fatalf("unexpected sections: got=%v want=%v", gciCfg.Sections, want)
	}
}

func TestParseNoLexOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"prefix(github/daixiang0/gci)", "prefix(github/daixiang0/gai)", "default"},
		Cfg: BoolConfig{
			NoLexOrder: true,
		},
	}
	gciCfg, err := cfg.Parse()
	if err != nil {
		t.Fatal(err)
	}
	want := section.SectionList{
		section.Default{},
		section.Custom{Prefix: "github/daixiang0/gci"},
		section.Custom{Prefix: "github/daixiang0/gai"},
	}
	if !reflect.DeepEqual(want, gciCfg.Sections) {
		t.Fatalf("unexpected sections: got=%v want=%v", gciCfg.Sections, want)
	}
}
