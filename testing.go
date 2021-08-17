package main

import (
	"path/filepath"
	"runtime"
)

func getLocalPath(file string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), file)
}
