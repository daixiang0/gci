package gci

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileObj struct {
	Path    string
	Load    func() ([]byte, error)
	IsStdin bool
}

type FileGeneratorFunc func() ([]FileObj, error)

func GoFilesInPathsGenerator(paths []string, skipVendor bool) FileGeneratorFunc {
	return func() ([]FileObj, error) {
		var files []FileObj
		for _, path := range paths {
			err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					if skipVendor && (filePath == "vendor" || strings.Contains(filePath, string(os.PathSeparator)+"vendor")) {
						return filepath.SkipDir
					}
					return nil
				}
				if strings.HasSuffix(filePath, ".go") && !strings.HasPrefix(filepath.Base(filePath), ".") {
					files = append(files, FileObj{
						Path: filePath,
						Load: func() ([]byte, error) {
							return os.ReadFile(filePath)
						},
					})
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
		return files, nil
	}
}

func StdInGenerator() ([]FileObj, error) {
	stdinFilePath := "<standard input>"
	return []FileObj{
		{
			Path:    stdinFilePath,
			IsStdin: true,
			Load: func() ([]byte, error) {
				return io.ReadAll(os.Stdin)
			},
		},
	}, nil
}

func CombineGenerators(generators ...FileGeneratorFunc) FileGeneratorFunc {
	return func() ([]FileObj, error) {
		var allFiles []FileObj
		for _, gen := range generators {
			files, err := gen()
			if err != nil {
				return nil, err
			}
			allFiles = append(allFiles, files...)
		}
		return allFiles, nil
	}
}
