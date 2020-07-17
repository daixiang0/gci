package main

import (
	"fmt"
	"sort"
	"strings"
)

const (
	// pkg type: standard, remote, local
	standard int = iota
	// 3rd-party packages
	remote
	local

	commentFlag = "//"
)

var (
	importStartFlag = []byte(`
import (
`)

	importEndFlag = []byte(`
)
`)
)

type pkg struct {
	list    map[int][]string
	comment map[string]string
	alias   map[string]string
}

func newPkg(data [][]byte, localFlag string) *pkg {
	listMap := make(map[int][]string)
	commentMap := make(map[string]string)
	aliasMap := make(map[string]string)
	p := &pkg{
		list:    listMap,
		comment: commentMap,
		alias:   aliasMap,
	}

	formatData := make([]string, 0)
	// remove all empty lines
	for _, v := range data {
		if len(v) > 0 {
			formatData = append(formatData, strings.TrimSpace(string(v)))
		}
	}

	for i := len(formatData) - 1; i >= 0; i-- {
		line := formatData[i]

		// check commentFlag:
		// 1. one line commentFlag
		// 2. commentFlag after import path
		commentIndex := strings.Index(line, commentFlag)
		if commentIndex == 0 {
			pkg, _, _ := getPkgInfo(formatData[i+1], true)
			p.comment[pkg] = line
			continue
		} else if commentIndex > 0 {
			pkg, alias, comment := getPkgInfo(line, true)
			if alias != "" {
				p.alias[pkg] = alias
			}

			p.comment[pkg] = comment
			pkgType := getPkgType(pkg, localFlag)
			p.list[pkgType] = append(p.list[pkgType], pkg)
			continue
		}

		pkg, alias, _ := getPkgInfo(line, false)

		if alias != "" {
			p.alias[pkg] = alias
		}

		pkgType := getPkgType(pkg, localFlag)
		p.list[pkgType] = append(p.list[pkgType], pkg)
	}

	return p
}

// fmt format import pkgs as expected
func (p *pkg) fmt() []byte {
	ret := make([]string, 0, 100)

	for pkgType := range []int{standard, remote, local} {
		sort.Strings(p.list[pkgType])
		for _, s := range p.list[pkgType] {
			if p.comment[s] != "" {
				l := fmt.Sprintf("%s%s%s", indent, p.comment[s], linebreak)
				ret = append(ret, l)
			}

			if p.alias[s] != "" {
				s = fmt.Sprintf("%s%s%s%s%s", indent, p.alias[s], blank, s, linebreak)
			} else {
				s = fmt.Sprintf("%s%s%s", indent, s, linebreak)
			}

			ret = append(ret, s)
		}

		if len(p.list[pkgType]) > 0 {
			ret = append(ret, linebreak)
		}
	}
	if ret[len(ret)-1] == linebreak {
		ret = ret[:len(ret)-1]
	}
	return []byte(strings.Join(ret, ""))
}

// getPkgInfo assume line is a import path, and return (path, alias, comment)
func getPkgInfo(line string, comment bool) (string, string, string) {
	if comment {
		s := strings.Split(line, commentFlag)
		pkgArray := strings.Split(s[0], blank)
		if len(pkgArray) > 1 {
			return pkgArray[1], pkgArray[0], fmt.Sprintf("%s%s%s", commentFlag, blank, s[1])
		} else {
			return pkgArray[0], "", fmt.Sprintf("%s%s%s", commentFlag, blank, s[1])
		}
	} else {
		pkgArray := strings.Split(line, blank)
		if len(pkgArray) > 1 {
			return pkgArray[1], pkgArray[0], ""
		} else {
			return pkgArray[0], "", ""
		}
	}
}

func getPkgType(line, localFlag string) int {
	if !strings.Contains(line, dot) {
		return standard
	} else if strings.Contains(line, localFlag) {
		return local
	} else {
		return remote
	}
}
