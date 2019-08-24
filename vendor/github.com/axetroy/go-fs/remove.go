package fs

import (
  "os"
)

/**
remove a file or a dir
 */
func Remove(path string) (error) {
  return os.RemoveAll(path)
}
