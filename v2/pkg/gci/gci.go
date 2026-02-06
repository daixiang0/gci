package gci

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/daixiang0/gci/v2/internal/imports"
	"github.com/daixiang0/gci/v2/pkg/config"
	"github.com/daixiang0/gci/v2/pkg/parse"
)

func PrintFormattedFiles(paths []string, cfg config.Config) error {
	return processStdInAndGoFilesInPaths(paths, cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		fmt.Print(string(formattedFile))
		return nil
	})
}

func WriteFormattedFiles(paths []string, cfg config.Config) error {
	return processGoFilesInPaths(paths, cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		if bytes.Equal(unmodifiedFile, formattedFile) {
			return nil
		}
		return os.WriteFile(filePath, formattedFile, 0o644)
	})
}

func ListUnFormattedFiles(paths []string, cfg config.Config) error {
	return processGoFilesInPaths(paths, cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		if bytes.Equal(unmodifiedFile, formattedFile) {
			return nil
		}
		fmt.Println(filePath)
		return nil
	})
}

func DiffFormattedFiles(paths []string, cfg config.Config) error {
	return processStdInAndGoFilesInPaths(paths, cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		return diffFormattedFiles(filePath, unmodifiedFile, formattedFile)
	})
}

func DiffFormattedFilesToArray(paths []string, cfg config.Config, diffs *[]string, lock *sync.Mutex) error {
	return processStdInAndGoFilesInPaths(paths, cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		diff, err := diffToString(filePath, unmodifiedFile, formattedFile)
		if err != nil {
			return err
		}
		lock.Lock()
		*diffs = append(*diffs, diff)
		lock.Unlock()
		return nil
	})
}

type fileFormattingFunc func(filePath string, unmodifiedFile, formattedFile []byte) error

func processStdInAndGoFilesInPaths(paths []string, cfg config.Config, fileFunc fileFormattingFunc) error {
	return ProcessFiles(CombineGenerators(StdInGenerator, GoFilesInPathsGenerator(paths, cfg.SkipVendor)), cfg, fileFunc)
}

func processGoFilesInPaths(paths []string, cfg config.Config, fileFunc fileFormattingFunc) error {
	return ProcessFiles(GoFilesInPathsGenerator(paths, cfg.SkipVendor), cfg, fileFunc)
}

func ProcessFiles(fileGenerator FileGeneratorFunc, cfg config.Config, fileFunc fileFormattingFunc) error {
	var taskGroup errgroup.Group
	files, err := fileGenerator()
	if err != nil {
		return err
	}
	for _, file := range files {
		taskGroup.Go(processingFunc(file, cfg, fileFunc))
	}
	return taskGroup.Wait()
}

func processingFunc(file FileObj, cfg config.Config, formattingFunc fileFormattingFunc) func() error {
	return func() error {
		unmodifiedFile, formattedFile, err := LoadFormatGoFile(file, cfg)
		if err != nil {
			return err
		}
		return formattingFunc(file.Path, unmodifiedFile, formattedFile)
	}
}

func LoadFormatGoFile(file FileObj, cfg config.Config) (src, dist []byte, err error) {
	src, err = file.Load()
	if err != nil {
		return nil, nil, err
	}

	return LoadFormat(src, file.Path, cfg)
}

func LoadFormat(in []byte, path string, cfg config.Config) (src, dist []byte, err error) {
	src = in

	if cfg.SkipGenerated && parse.IsGeneratedFileByComment(string(src)) {
		return src, src, nil
	}

	_, _, _, _, _, err = parse.ParseFile(src, path)
	if err != nil {
		if errors.Is(err, parse.NoImportError{}) {
			return src, src, nil
		}
		return nil, nil, err
	}

	opts := &imports.Options{
		Config:     &cfg,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: true,
	}

	dist, err = imports.Process(path, src, opts)
	if err != nil {
		return nil, nil, err
	}

	return src, dist, nil
}
