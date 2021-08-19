package g

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	Basedirs   []string
	CurrentDir string
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	Basedirs = append(Basedirs, filepath.Dir(filepath.Dir(file)))
	dir, _ := os.Getwd()
	Basedirs = append(Basedirs, dir)
}
