package fs

import "os"

/**
check a path is exist or not
 */
func PathExists(path string) (isExist bool) {
  if _, err := os.Stat(path); os.IsNotExist(err) {
    return false
  }
  return true
}
