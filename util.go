package common

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func StringInSlice(s string, sl []string) bool {
	for _, ss := range sl {
		if s == ss {
			return true
		}
	}
	return false
}

func ListDir(dir string, exts []string) ([]string, error) {
	var files []string
	f, err := os.Open(dir)
	if err != nil {
		return files, err
	}
	defer f.Close()

	fileInfos, err := f.Readdir(0)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfos {
		if exts != nil && !StringInSlice(filepath.Ext(file.Name()), exts) {
			continue
		}
		files = append(files, path.Join(dir, file.Name()))
	}
	return files, nil
}
