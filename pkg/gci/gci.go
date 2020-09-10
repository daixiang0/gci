package gci

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type pkgType int

const (
	// pkg type: standard, remote, local
	standard pkgType = iota
	// 3rd-party packages
	remote
	local
)

const commentFlag = "//"

var (
	importStartFlag = []byte(`
import (
`)

	importEndFlag = []byte(`
)
`)
)

type FlagSet struct {
	LocalFlag       string
	DoWrite, DoDiff *bool
}

type pkgInfo struct {
	path    string
	alias   string
	comment []string
}

type pkgInfos []*pkgInfo

var _ sort.Interface = pkgInfos{}

func (p pkgInfos) Len() int {
	return len(p)
}

func (p pkgInfos) Less(i, j int) bool {
	if p[i].path != p[j].path {
		return p[i].path < p[j].path
	}
	return p[i].alias < p[j].alias
}

func (p pkgInfos) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type pkg struct {
	list map[pkgType]pkgInfos
}

func newPkg(data [][]byte, localFlag string) *pkg {
	p := &pkg{
		list: map[pkgType]pkgInfos{
			standard: {},
			remote:   {},
			local:    {},
		},
	}

	formatData := make([]string, 0)
	// remove all empty lines
	for _, v := range data {
		if len(v) > 0 {
			formatData = append(formatData, strings.TrimSpace(string(v)))
		}
	}

	// lastInfo keeps track of the info from the previous iteration of the for loop
	// The loop iterates backward (from the bottom up) so when a comment line is detected
	// it should always be prepend to the comment block of the lastInfo
	var lastInfo *pkgInfo
	for i := len(formatData) - 1; i >= 0; i-- {
		line := formatData[i]

		// check commentFlag:
		// 1. one line commentFlag
		// 2. commentFlag after import path
		// 3. commentFlag not present
		commentIndex := strings.Index(line, commentFlag)
		switch {
		case commentIndex == 0:
			// comment in the last line is useless, ignore it
			if lastInfo != nil {
				lastInfo.comment = append([]string{line}, lastInfo.comment...)
			}

		case commentIndex > 0:
			info := getPkgInfo(line, true)
			pkgType := getPkgType(info.path, localFlag)
			p.list[pkgType] = append(p.list[pkgType], info)
			lastInfo = info

		default:
			info := getPkgInfo(line, false)
			pkgType := getPkgType(info.path, localFlag)
			p.list[pkgType] = append(p.list[pkgType], info)
			lastInfo = info
		}
	}

	return p
}

// fmt format import pkgs as expected
func (p *pkg) fmt() []byte {
	ret := make([]string, 0, 100)

	for _, pkgType := range []pkgType{standard, remote, local} {
		pkgs := p.list[pkgType]
		sort.Sort(pkgs)

		for _, pkg := range pkgs {
			if len(pkg.comment) > 0 {
				ret = append(ret, linebreak)
				for _, comment := range pkg.comment {
					ret = append(ret, fmt.Sprintf("%s%s%s", indent, comment, linebreak))
				}
			}

			if pkg.alias != "" {
				ret = append(ret, fmt.Sprintf("%s%s%s%s%s", indent, pkg.alias, blank, pkg.path, linebreak))
			} else {
				ret = append(ret, fmt.Sprintf("%s%s%s", indent, pkg.path, linebreak))
			}
		}

		if len(p.list[pkgType]) > 0 {
			ret = append(ret, linebreak)
		}
	}
	if ret[len(ret)-1] == linebreak {
		ret = ret[:len(ret)-1]
	}

	// remove duplicate empty lines
	s1 := fmt.Sprintf("%s%s%s%s", linebreak, linebreak, linebreak, indent)
	s2 := fmt.Sprintf("%s%s%s", linebreak, linebreak, indent)
	return []byte(strings.ReplaceAll(strings.Join(ret, ""), s1, s2))
}

// getPkgInfo assume line is a import path, and return (path, aliases, comment)
func getPkgInfo(line string, comment bool) *pkgInfo {
	if comment {
		s := strings.Split(line, commentFlag)
		pkgArray := strings.Split(s[0], blank)
		if len(pkgArray) > 1 {
			return &pkgInfo{
				path:    pkgArray[1],
				alias:   pkgArray[0],
				comment: []string{fmt.Sprintf("%s%s%s", commentFlag, blank, strings.TrimSpace(s[1]))},
			}
		} else {
			return &pkgInfo{
				path:    strings.TrimSpace(pkgArray[0]),
				alias:   "",
				comment: []string{fmt.Sprintf("%s%s%s", commentFlag, blank, strings.TrimSpace(s[1]))},
			}
		}
	} else {
		pkgArray := strings.Split(line, blank)
		if len(pkgArray) > 1 {
			return &pkgInfo{
				path:    pkgArray[1],
				alias:   pkgArray[0],
				comment: nil,
			}
		} else {
			return &pkgInfo{
				path:    pkgArray[0],
				alias:   "",
				comment: nil,
			}
		}
	}
}

func getPkgType(line, localFlag string) pkgType {
	if !strings.Contains(line, dot) {
		return standard
	} else if strings.Contains(line, localFlag) {
		return local
	} else {
		return remote
	}
}

const (
	dot       = "."
	blank     = " "
	indent    = "\t"
	linebreak = "\n"
)

func diff(b1, b2 []byte, filename string) (data []byte, err error) {
	f1, err := writeTempFile("", "gci", b1)
	if err != nil {
		return
	}
	defer os.Remove(f1)

	f2, err := writeTempFile("", "gci", b2)
	if err != nil {
		return
	}
	defer os.Remove(f2)

	cmd := "diff"

	data, err = exec.Command(cmd, "-u", f1, f2).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		return replaceTempFilename(data, filename)
	}
	return
}

func writeTempFile(dir, prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return "", err
	}
	_, err = file.Write(data)
	if err1 := file.Close(); err == nil {
		err = err1
	}
	if err != nil {
		os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}

// replaceTempFilename replaces temporary filenames in diff with actual one.
//
// --- /tmp/gofmt316145376	2017-02-03 19:13:00.280468375 -0500
// +++ /tmp/gofmt617882815	2017-02-03 19:13:00.280468375 -0500
// ...
// ->
// --- path/to/file.go.orig	2017-02-03 19:13:00.280468375 -0500
// +++ path/to/file.go	2017-02-03 19:13:00.280468375 -0500
// ...
func replaceTempFilename(diff []byte, filename string) ([]byte, error) {
	bs := bytes.SplitN(diff, []byte{'\n'}, 3)
	if len(bs) < 3 {
		return nil, fmt.Errorf("got unexpected diff for %s", filename)
	}
	// Preserve timestamps.
	var t0, t1 []byte
	if i := bytes.LastIndexByte(bs[0], '\t'); i != -1 {
		t0 = bs[0][i:]
	}
	if i := bytes.LastIndexByte(bs[1], '\t'); i != -1 {
		t1 = bs[1][i:]
	}
	// Always print filepath with slash separator.
	f := filepath.ToSlash(filename)
	bs[0] = []byte(fmt.Sprintf("--- %s%s", f+".orig", t0))
	bs[1] = []byte(fmt.Sprintf("+++ %s%s", f, t1))
	return bytes.Join(bs, []byte{'\n'}), nil
}

func visitFile(set *FlagSet) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err == nil && isGoFile(f) {
			err = processFile(path, os.Stdout, set)
		}
		return err
	}
}

func WalkDir(path string, set *FlagSet) error {
	return filepath.Walk(path, visitFile(set))
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func ProcessFile(filename string, out io.Writer, set *FlagSet) error {
	return processFile(filename, out, set)
}

func processFile(filename string, out io.Writer, set *FlagSet) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return processSource(filename, src, out, set)
}

func processSource(filename string, src []byte, out io.Writer, set *FlagSet) error {
	ori := make([]byte, len(src))
	copy(ori, src)
	start := bytes.Index(src, importStartFlag)
	// in case no importStartFlag or importStartFlag exist in the commentFlag
	if start < 0 {
		fmt.Printf("skip file %s since no import\n", filename)
		return nil
	}
	end := bytes.Index(src[start:], importEndFlag) + start

	ret := bytes.Split(src[start+len(importStartFlag):end], []byte(linebreak))

	p := newPkg(ret, set.LocalFlag)

	res := append(src[:start+len(importStartFlag)], append(p.fmt(), src[end+1:]...)...)

	if !bytes.Equal(ori, res) {
		if *set.DoWrite {
			// On Windows, we need to re-set the permissions from the file. See golang/go#38225.
			var perms os.FileMode
			if fi, err := os.Stat(filename); err == nil {
				perms = fi.Mode() & os.ModePerm
			}
			err := ioutil.WriteFile(filename, res, perms)
			if err != nil {
				return err
			}
		}
		if *set.DoDiff {
			data, err := diff(ori, res, filename)
			if err != nil {
				return fmt.Errorf("failed to diff: %v", err)
			}
			fmt.Printf("diff -u %s %s\n", filepath.ToSlash(filename+".orig"), filepath.ToSlash(filename))
			if _, err := out.Write(data); err != nil {
				return fmt.Errorf("failed to write: %v", err)
			}
		}
	}
	if !*set.DoWrite && !*set.DoDiff {
		if _, err := out.Write(res); err != nil {
			return fmt.Errorf("failed to write: %v", err)
		}
	}

	return nil
}

func Run(filename string, set *FlagSet) ([]byte, error) {
	var err error

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	ori := make([]byte, len(src))
	copy(ori, src)
	start := bytes.Index(src, importStartFlag)
	// in case no importStartFlag or importStartFlag exist in the commentFlag
	if start < 0 {
		return nil, nil
	}
	end := bytes.Index(src[start:], importEndFlag) + start

	ret := bytes.Split(src[start+len(importStartFlag):end], []byte(linebreak))

	p := newPkg(ret, set.LocalFlag)

	res := append(src[:start+len(importStartFlag)], append(p.fmt(), src[end+1:]...)...)

	if bytes.Equal(ori, res) {
		return nil, nil
	}

	data, err := diff(ori, res, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to diff: %v", err)
	}

	return data, nil
}
