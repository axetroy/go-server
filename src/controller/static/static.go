package static

import (
	"github.com/axetroy/go-fs"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	hasStaticDir bool
	PublicDir    string
)

func init() {
	if cwd, err := os.Getwd(); err != nil {
		// ignore error
		hasStaticDir = false
	} else {
		PublicDir = path.Join(cwd, "public")
		if fs.PathExists(PublicDir) {
			hasStaticDir = true
		}
	}
}

func Get(context *gin.Context) {
	filename := context.Param("filename")
	filePath := path.Join(PublicDir, filename)

	// do not allow dot prefix
	if strings.HasPrefix(".", filename) {
		http.NotFound(context.Writer, context.Request)
		return
	}

	if fs.PathExists(filePath) == false {
		// if the path not found
		http.NotFound(context.Writer, context.Request)
		return
	}
	http.ServeFile(context.Writer, context.Request, filePath)
}
