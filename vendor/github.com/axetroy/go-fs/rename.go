package fs

import (
  "os"
)

func Rename(oldPath string, newPath string) (error) {
  return os.Rename(oldPath, newPath)
}
