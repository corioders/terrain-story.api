package main

import "os"

func mkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
