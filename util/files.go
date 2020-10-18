package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	DefaultDirPermMode = os.FileMode(0700)
	EndLineDelimiter   = byte('\n')
)

// FileExists returns a boolean indicating whether file (not a directory) with given path exist
func FileExists(path string) bool {
	if stats, isNotExistErr := os.Stat(path); !os.IsNotExist(isNotExistErr) {
		return !stats.IsDir()
	} else {
		return false
	}
}

// DirectoriesWalk returns a slice of files present in `dirs` directories.
// `cascading` parameter defines whether traversal will be applied by sub-directories as well.
// `suffix` parameter defines a suffix of files which will be added to resulted slice.
func DirectoriesWalk(dirs []string, cascading bool, suffix string) ([]string, error) {
	var fileResults []string

	for _, dir := range dirs {
		var err = DirectoryWalk(dir, cascading, suffix, &fileResults)
		if err != nil {
			return fileResults, err
		}
	}

	return fileResults, nil
}

// DirectoryWalk fills a slice of `fileResults` present in `dir` directory.
// `cascading` parameter defines whether traversal will be applied by sub-directories as well.
// `suffix` parameter defines a suffix of files which will be added to `fileResults`.
func DirectoryWalk(dir string, cascading bool, suffix string, fileResults *[]string) error {
	filesInfo, err := ioutil.ReadDir(dir)
	suffix = strings.TrimSpace(suffix)

	if err != nil {
		return err
	}

	for _, file := range filesInfo {
		var absFilePath, _ = filepath.Abs(dir + string(os.PathSeparator) + file.Name())

		if file.IsDir() && cascading {
			_ = DirectoryWalk(absFilePath, cascading, suffix, fileResults)
		} else if FileExists(absFilePath) && (len(suffix) == 0 || strings.HasSuffix(absFilePath, suffix)) {
			*fileResults = append(*fileResults, absFilePath)
		}
	}
	return nil
}

func ConcatenateFiles(source string, output *os.File) (int64, error) {
	var totalWritten int64

	var file, errOpen = os.Open(source)
	if errOpen != nil {
		return 0, errOpen
	}

	var reader = bufio.NewReader(file)

	for {
		var line, readErr = reader.ReadString(EndLineDelimiter)
		if readErr != nil && readErr == io.EOF {
			break
		} else if readErr != nil && readErr != io.EOF {
			os.Exit(-1)
		}

		var bytesWritten, errWrite = output.WriteString(line)
		if errWrite != nil {
			return 0, errWrite
		}

		totalWritten += int64(bytesWritten)
	}

	return totalWritten, nil
}
