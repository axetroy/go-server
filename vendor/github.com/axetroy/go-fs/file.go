package fs

import (
  "os"
  "io/ioutil"
  "bufio"
  "io"
  "path"
)

/**
ensure the file exist
 */
func EnsureFile(filepath string) (err error) {
  var (
    file *os.File
  )
  // ensure dir exist
  if err = EnsureDir(path.Dir(filepath)); err != nil {
    return err
  }
  // ensure file exist
  if _, err = os.Stat(filepath); os.IsNotExist(err) {
    file, err = os.Create(filepath)
    defer func() {
      file.Close()
    }()
  }
  return
}

/*
write a file
 */
func WriteFile(filepath string, data []byte) error {
  return ioutil.WriteFile(filepath, data, os.ModePerm)
}

/**
read a file
 */
func ReadFile(filepath string) ([]byte, error) {
  return ioutil.ReadFile(filepath)
}

func AppendFile(file string, data []byte) (error) {
  if f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
    return err
  } else {
    defer func() {
      err = f.Close()
    }()
    if _, err := f.Write(data); err != nil {
      return err
    }
    return nil
  }
}

func Truncate(filepath string, len int64) (error) {
  return os.Truncate(filepath, len)
}

/**
Create a read stream
 */
func CreateReadStream(filepath string) (stream io.Reader, err error) {
  var (
    file *os.File
  )
  if file, err = os.Open(filepath); err != nil {
    return
  }

  defer func() {
    err = file.Close()
  }()

  stream = bufio.NewReader(file)

  return
}

/**
Create a read stream
 */
func CreateWriteStream(filepath string) (stream io.Writer, err error) {
  var (
    file *os.File
  )
  if file, err = os.Open(filepath); err != nil {
    return
  }

  defer func() {
    err = file.Close()
  }()

  stream = bufio.NewWriter(file)

  return
}

// Almost the same as writeFile (i.e. it overwrites), except that if the parent directory does not exist, it's created.
func OuputFile(filepath string, data []byte) error {
  if err := EnsureDir(path.Dir(filepath)); err != nil {
    return err
  }
  return WriteFile(filepath, data)
}
