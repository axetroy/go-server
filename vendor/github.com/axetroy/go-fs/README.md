## File Accessor For Go Inspired by NodeJs fs system

[![Build Status](https://travis-ci.org/axetroy/go-fs.svg?branch=master)](https://travis-ci.org/axetroy/go-fs)
![License](https://img.shields.io/badge/license-Apache-green.svg)
[![Coverage Status](https://coveralls.io/repos/github/axetroy/go-fs/badge.svg?branch=master)](https://coveralls.io/github/axetroy/go-fs?branch=master)
![Size](https://github-size-badge.herokuapp.com/axetroy/go-fs.svg)

API almost like NodeJs [fs module](http://nodejs.cn/api/fs.html) and [fs-extra](https://github.com/jprichardson/node-fs-extra)

## Usage

```bash
go get https://github.com/axetroy/go-fs.git
```

```go
import fs "github.com/axetroy/go-fs"

func main(){
  if err := fs.EnsureFile("./testFile.txt");err !=nill {

  }
  if err := fs.Copy("./testFile.txt", "./newTestFile.txt");err !=nil {

  }
}
```

## Methods

| Method                                       | Description                                                                                 |
| -------------------------------------------- | ------------------------------------------------------------------------------------------- |
| fs.EnsureFile(filepath string)               | Ensure the file exist, if not then create it                                                |
| fs.WriteFile(filepath string, data []byte)   | Write file                                                                                  |
| fs.ReadFile(filepath string)                 | Read file                                                                                   |
| fs.AppendFile(filepath string, data []byte)  | Append data into a file                                                                     |
| fs.Truncate(filepath string, len int64)      | Truncate the file                                                                           |
| fs.CreateReadStream(filepath string)         | Create a read stream (Reader)                                                               |
| fs.CreateWriteStream(filepath string)        | Create a write stream (Writer)                                                              |
| fs.EnsureDir(dir string)                     | Ensure the dir exist, if not then create it                                                 |
| fs.Mkdir(dir string)                         | Create a dir                                                                                |
| fs.Readdir(dir string)                       | Read a dir                                                                                  |
| fs.Mktemp(dir string, prefix string)         | Create a temp dir                                                                           |
| fs.Rmdir(dir string)                         | Remove a dir                                                                                |
| fs.Stat(path string)                         | Stat a file/dir                                                                             |
| fs.LStat(path string)                        | Stat a file/dir                                                                             |
| fs.Remove(path string)                       | Remove a file/dir                                                                           |
| fs.PathExists(path string)                   | Check a path is exist or not                                                                |
| fs.Link(existingPath string, newPath string) | Link a file                                                                                 |
| fs.ReadLink(path string)                     | Read link info                                                                              |
| fs.Symlink(target string, path string)       | Create a Symlink                                                                            |
| fs.Unlink(path string)                       | Unlink a link                                                                               |
| fs.Copy(src string, target string)           | Copy a file/dir                                                                             |
| fs.Move(src string, target string)           | Move a file/dir                                                                             |
| fs.Chmod(filepath string, mode os.FileMode)  | Change the file/dir's mode                                                                  |
| fs.LChod(path string, uid int, gid int)      | Change the file/dir's mode                                                                  |
| fs.ReadJson(filepath string)                 | Read JSON file                                                                              |
| fs.WriteJson(filepath string, data []byte)   | Write JSON file                                                                             |
| fs.OuputFile(filepath string, data []byte)   | Almost the same as fs.WriteFile, except that if the directory does not exist, it's created. |
| fs.OuputJson(filepath string, data []byte)   | Almost the same as fs.WriteJson, except that if the directory does not exist, it's created. |

## Contributing

```bash
go get https://github.com/axetroy/go-fs.git
cd $GOPATH/src/github.com/axetroy/go-fs
go test -v
```

[Contributing Guid](https://github.com/axetroy/Github/blob/master/CONTRIBUTING.md)

## Test

```bash
go test -v
```

## Contributors

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->

| [<img src="https://avatars1.githubusercontent.com/u/9758711?v=3" width="100px;"/><br /><sub>Axetroy</sub>](http://axetroy.github.io)<br />[💻](https://github.com/axetroyanti-redirect/go-fs/commits?author=axetroy) [🐛](https://github.com/axetroy/go-fs/issues?q=author%3Aaxetroy) 🎨 |
| :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |


<!-- ALL-CONTRIBUTORS-LIST:END -->

## License

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Faxetroy%2Fgo-fs.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Faxetroy%2Fgo-fs?ref=badge_large)
