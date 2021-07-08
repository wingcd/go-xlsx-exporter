package utils

import (
	"os"
	"path"
	"strings"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CheckPath(p string) {
	if strings.Contains(p, ".") {
		p = path.Dir(p)
	}
	if ok, _ := PathExists(p); !ok {
		os.MkdirAll(p, os.ModePerm)
	}
}
