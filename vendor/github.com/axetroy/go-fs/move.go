package fs

import (
	"io/ioutil"
	"os"
	"path"
)

// Move file/folder
func Move(src string, target string) (err error) {

	var (
		fileInfo os.FileInfo
		files    []os.FileInfo
	)

	if fileInfo, err = os.Stat(src); err != nil {
		return
	}

	if fileInfo.IsDir() {

		// read dir and move one by one
		if files, err = ioutil.ReadDir(src); err != nil {
			return
		}

		if err = EnsureDir(target); err != nil {
			return
		}

		for _, file := range files {
			filename := file.Name()
			srcFile := path.Join(src, filename)
			targetFile := path.Join(target, filename)
			if err = Move(srcFile, targetFile); err != nil {
				return err
			}
		}

		// copy all done, should remove the src dir
		err = os.RemoveAll(src)

		return
	}
	return os.Rename(src, target)
}
