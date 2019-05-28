package fs

import "os"

/**
stat a file
 */
func Stat(path string) (info os.FileInfo, err error) {
  return os.Stat(path)
}

/**
stat a file
 */
func LStat(path string) (info os.FileInfo, err error) {
  return os.Lstat(path)
}
