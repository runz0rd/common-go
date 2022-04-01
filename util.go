package common

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
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

type StackTracable interface {
	StackTrace() errors.StackTrace
}

type Printable interface {
	Printf(format string, v ...interface{})
}

func PrintStackTrace(err StackTracable, p Printable) {
	if err, ok := err.(StackTracable); ok {
		for _, f := range err.StackTrace() {
			p.Printf("%+s:%d\n", f, f)
		}
	}
}

func GetUrlContent(url string) (realUrl string, content []byte, err error) {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, errors.WithStack(err)
	}
	return resp.Request.URL.Path, content, nil
}

func LoadYaml(path string, data interface{}) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, data)
}
