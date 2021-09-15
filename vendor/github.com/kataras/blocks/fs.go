package blocks

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"sort"
)

func getFS(fsOrDir interface{}) (fs http.FileSystem) {
	switch v := fsOrDir.(type) {
	case string:
		fs = http.Dir(v)
	case http.FileSystem:
		fs = v
	default:
		panic(fmt.Errorf(`blocks: unexpected "fileSystem" argument type of %T (string or http.FileSystem)`, v))
	}

	return
}

// walk recursively in "fs" descends "root" path, calling "walkFn".
func walk(fs http.FileSystem, root string, walkFn filepath.WalkFunc) error {
	names, err := assetNames(fs, root)
	if err != nil {
		return fmt.Errorf("%s: %w", root, err)
	}

	for _, name := range names {
		fullpath := path.Join(root, name)
		f, err := fs.Open(fullpath)
		if err != nil {
			return fmt.Errorf("%s: %w", fullpath, err)
		}
		stat, err := f.Stat()
		err = walkFn(fullpath, stat, err)
		if err != nil {
			if err != filepath.SkipDir {
				return fmt.Errorf("%s: %w", fullpath, err)
			}

			continue
		}

		if stat.IsDir() {
			if err := walk(fs, fullpath, walkFn); err != nil {
				return fmt.Errorf("%s: %w", fullpath, err)
			}
		}
	}

	return nil
}

// assetNames returns the first-level directories and file, sorted, names.
func assetNames(fs http.FileSystem, name string) ([]string, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}

	infos, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(infos))
	for _, info := range infos {
		// note: go-bindata fs returns full names whether
		// the http.Dir returns the base part, so
		// we only work with their base names.
		name := filepath.ToSlash(info.Name())
		name = path.Base(name)
		names = append(names, name)
	}

	sort.Strings(names)
	return names, nil
}

func asset(fs http.FileSystem, name string) ([]byte, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}

	contents, err := ioutil.ReadAll(f)
	f.Close()
	return contents, err
}
