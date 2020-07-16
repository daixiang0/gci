package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/scanner"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// pkg type: standard, remote, local
	standard int = iota
	// 3rd-party packages
	remote
	local

	dot         = "."
	blank       = " "
	indent      = "\t"
	linebreak   = "\n"
	commentFlag = "//"
)

var (
	write  = flag.Bool("w", false, "write result to (source) file instead of stdout")
	doDiff = flag.Bool("d", false, "display diffs instead of rewriting files")

	localFlag string

	exitCode = 0

	importStartFlag = []byte(`
import (
`)
	importEndFlag = []byte(`
)
`)
)

func report(err error) {
	scanner.PrintError(os.Stderr, err)
	exitCode = 2
}

func isGoFile(f os.FileInfo) bool {
	// ignore non-Go files
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go")
}

func parseFlags() []string {
	flag.StringVar(&localFlag, "local", "", "put imports beginning with this string after 3rd-party packages, only support one string")

	flag.Parse()
	return flag.Args()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gci [flags] [path ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

type pkg struct {
	list    map[int][]string
	comment map[string]string
	alias   map[string]string
}

func newPkg(data [][]byte) *pkg {
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
			pkg, _, _ := getPkgInfo(formatData[i+1])
			p.comment[pkg] = line
			continue
		} else if commentIndex > 0 {
			pkg, alias, comment := getPkgInfo(line)
			if alias != "" {
				p.alias[pkg] = alias
			}

			p.comment[pkg] = comment
			pkgType := getPkgType(pkg)
			p.list[pkgType] = append(p.list[pkgType], pkg)
			continue
		}

		pkg, alias, _ := getPkgInfo(line)

		if alias != "" {
			p.alias[pkg] = alias
		}

		pkgType := getPkgType(pkg)
		p.list[pkgType] = append(p.list[pkgType], pkg)
	}

	return p
}

// getPkgInfo assume line is a import path, and return (path, alias)
func getPkgInfo(line string) (string, string, string) {
	pkgArray := strings.Split(line, blank)
	if len(pkgArray) > 1 {
		return pkgArray[1], pkgArray[0], strings.Join(pkgArray[2:], "")
	} else {
		return line, "", ""
	}
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

func getPkgType(pkg string) int {
	if !strings.Contains(pkg, dot) {
		return standard
	} else if strings.Contains(pkg, localFlag) {
		return local
	} else {
		return remote
	}
}

func processFile(filename string, out io.Writer) error {
	var err error

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

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

	p := newPkg(ret)

	res := append(src[:start+len(importStartFlag)], append(p.fmt(), src[end+1:]...)...)

	if !bytes.Equal(ori, res) {
		exitCode = 1

		if *write {
			// On Windows, we need to re-set the permissions from the file. See golang/go#38225.
			var perms os.FileMode
			if fi, err := os.Stat(filename); err == nil {
				perms = fi.Mode() & os.ModePerm
			}
			err = ioutil.WriteFile(filename, res, perms)
			if err != nil {
				return err
			}
		}
		if *doDiff {
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
	if !*write && !*doDiff {
		if _, err = out.Write(res); err != nil {
			return fmt.Errorf("failed to write: %v", err)
		}
	}

	return err

}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isGoFile(f) {
		err = processFile(path, os.Stdout)
	}
	if err != nil {
		report(err)
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func main() {
	flag.Usage = usage
	paths := parseFlags()
	for _, path := range paths {
		switch dir, err := os.Stat(path); {
		case err != nil:
			report(err)
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, os.Stdout); err != nil {
				report(err)
			}
		}
	}
	os.Exit(exitCode)
}
