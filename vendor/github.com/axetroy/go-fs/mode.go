package fs

import "os"

/**
change the file permission
 */
func Chmod(filepath string, mode os.FileMode) error {
  return os.Chmod(filepath, mode)
}

func LChod(path string, uid int, gid int) (error) {
  return os.Lchown(path, uid, gid)
}
