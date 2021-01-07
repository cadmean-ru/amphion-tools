package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyDir(sourcePath, destPath string) (err error) {

	src := filepath.Clean(sourcePath)
	dst := filepath.Clean(destPath)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	//if err == nil {
	//	return fmt.Errorf("destination already exists")
	//}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode() & os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func CopyFile(source, dest string) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return
	}

	si, err := os.Stat(source)
	if err != nil {
		return
	}

	err = os.Chmod(dest, si.Mode())
	if err != nil {
		return
	}

	return
}