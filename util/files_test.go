package util

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	existingTestRootPath    = "test_root"
	existingTestSubDir1Path = existingTestRootPath + string(os.PathSeparator) + "test_sub_dir1"
	existingTestSubDir2Path = existingTestRootPath + string(os.PathSeparator) + "test_sub_dir2"
	existingTestSubDir3Path = existingTestRootPath + string(os.PathSeparator) + "test_sub_dir3"
	existingTestSubDir4Path = existingTestSubDir3Path + string(os.PathSeparator) + "test_sub_dir4"

	existingTestFile1 = existingTestRootPath + string(os.PathSeparator) + "test1.txt"
	existingTestFile2 = existingTestRootPath + string(os.PathSeparator) + "test2.txt"
	existingTestFile3 = existingTestRootPath + string(os.PathSeparator) + "test3.txt"

	existingTestFile1SubDir4 = existingTestSubDir4Path + string(os.PathSeparator) + "test1.txt"
	existingTestFile2SubDir4 = existingTestSubDir4Path + string(os.PathSeparator) + "test1.out"
)

func TestDirectoriesWalk(t *testing.T) {
	type args struct {
		dir       []string
		cascading bool
		suffix    string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr error
	}{
		// 1
		{
			name:    "Should return list of all files in directory cascading",
			args:    args{dir: []string{existingTestRootPath}, cascading: true, suffix: ""},
			want:    []string{existingTestFile1, existingTestFile2, existingTestFile3, existingTestFile2SubDir4, existingTestFile1SubDir4},
			wantErr: nil,
		},

		// 2
		{
			name:    "Should return list of files matching provided suffix in directory cascading",
			args:    args{dir: []string{existingTestRootPath}, cascading: true, suffix: ".out"},
			want:    []string{existingTestFile2SubDir4},
			wantErr: nil,
		},

		// 3
		{
			name:    "Should return list of files in a directory",
			args:    args{dir: []string{existingTestRootPath}, cascading: false, suffix: ""},
			want:    []string{existingTestFile1, existingTestFile2, existingTestFile3},
			wantErr: nil,
		},

		// 4
		{
			name:    "Should return empty list if matched files in given directory and suffix",
			args:    args{dir: []string{existingTestRootPath}, cascading: false, suffix: ".no_such"},
			want:    []string{},
			wantErr: nil,
		},

		// 5
		{
			name:    "Should return list of files in a set directories",
			args:    args{dir: []string{existingTestSubDir4Path, existingTestRootPath}, cascading: false, suffix: ""},
			want:    []string{existingTestFile2SubDir4, existingTestFile1SubDir4, existingTestFile1, existingTestFile2, existingTestFile3},
			wantErr: nil,
		},

		// 6
		{
			name:    "Should return error if directory does not exist",
			args:    args{dir: []string{"no/such/dir"}, cascading: false, suffix: ""},
			want:    []string{},
			wantErr: errors.New("open no/such/dir: no such file or directory"),
		},
	}

	setup()
	defer destroy()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actualResult, err = DirectoriesWalk(tt.args.dir, tt.args.cascading, tt.args.suffix)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Error(fmt.Sprintf("DirectoriesWalk() = Expected and actual error results does not match. "+
					"Expected: '%s', actual: '%s'", tt.wantErr, err))
			}

			if len(actualResult) != len(tt.want) {
				t.Error(fmt.Sprintf("DirectoriesWalk() = Expected and actual file results does not match. "+
					"Expected: '%s', actual: '%s'", tt.want, actualResult))
			}

			for i := 0; i < len(actualResult); i++ {
				var expectedPathSuffix = tt.want[i]
				if !strings.HasSuffix(actualResult[i], expectedPathSuffix) {
					t.Error(fmt.Sprintf("DirectoriesWalk() = Expected and actual file results does not match. "+
						"Expected: '%s', actual: '%s'", tt.want, actualResult))
				}
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		// 1
		{
			name: "Should return 'true' when file exists",
			arg:  existingTestFile1,
			want: true,
		},

		// 2
		{
			name: "Should return 'false' when file does not exist",
			arg:  "no_such.file",
			want: false,
		},

		// 3
		{
			name: "Should return 'false' when directory does not exist",
			arg:  "no_such_dir/",
			want: false,
		},

		// 4
		{
			name: "Should return 'false' when directory exists",
			arg:  existingTestRootPath,
			want: false,
		},
	}

	setup()
	defer destroy()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.arg); got != tt.want {
				t.Errorf("FileExists() = Expected and actual result does not match. Expected '%v', actual: '%v'", tt.want, got)
			}
		})
	}
}

func setup() {
	destroy()

	_ = os.MkdirAll(existingTestRootPath, DefaultDirPermMode)
	_ = os.MkdirAll(existingTestSubDir1Path, DefaultDirPermMode)
	_ = os.MkdirAll(existingTestSubDir2Path, DefaultDirPermMode)
	_ = os.MkdirAll(existingTestSubDir3Path, DefaultDirPermMode)
	_ = os.MkdirAll(existingTestSubDir4Path, DefaultDirPermMode)

	var file1, _ = os.Create(existingTestFile1)
	var file2, _ = os.Create(existingTestFile2)
	var file3, _ = os.Create(existingTestFile3)
	var file1SubDir4, _ = os.Create(existingTestFile1SubDir4)
	var file2SubDir4, _ = os.Create(existingTestFile2SubDir4)

	_, _ = file1.WriteString("file1")
	_, _ = file2.WriteString("file2")
	_, _ = file3.WriteString("file3")
	_, _ = file1SubDir4.WriteString("file1SubDir4")
	_, _ = file2SubDir4.WriteString("file2SubDir4")

	_ = file1.Close()
	_ = file2.Close()
	_ = file3.Close()
	_ = file1SubDir4.Close()
	_ = file2SubDir4.Close()
}

func destroy() {
	_ = os.Remove(existingTestFile1)
	_ = os.Remove(existingTestFile2)
	_ = os.Remove(existingTestFile3)
	_ = os.Remove(existingTestFile1SubDir4)
	_ = os.Remove(existingTestFile2SubDir4)

	_ = os.RemoveAll(existingTestSubDir1Path)
	_ = os.RemoveAll(existingTestSubDir2Path)
	_ = os.RemoveAll(existingTestSubDir3Path)
	_ = os.RemoveAll(existingTestSubDir4Path)

	_ = os.RemoveAll(existingTestRootPath)
}

func TestConcatenateFiles(t *testing.T) {
	type args struct {
		source string
		output *os.File
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConcatenateFiles(tt.args.source, tt.args.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConcatenateFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConcatenateFiles() got = %v, want %v", got, tt.want)
			}
		})
	}
}
