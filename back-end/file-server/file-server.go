package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

type rjFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func NewRjFileSystem(root string) *rjFileSystem {
	return &rjFileSystem{
		FileSystem: gin.Dir(root, false),
		root:       root,
		indexes:    true,
	}
}

func (l *rjFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		_, err := os.Stat(name)
		if err != nil {
			return false
		}

		return true
	}
	return false
}

// RjServe returns a middleware handler that serves static files in the given directory.
func RjServe(urlPrefix string, fs static.ServeFileSystem) gin.HandlerFunc {
	fileserver := http.FileServer(fs)

	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}

	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else if c.Request.Method == http.MethodGet {
			http.ServeFile(c.Writer, c.Request, "./static/index.html")
			c.Abort()
		}
	}
}

func main() {
	debug := flag.Bool("d", false, "Sets debugging mode, Cross-Origin Resource Sharing policy won't discriminate against the request origin (\"Access-Control-Allow-Origin\" header is \"*\")")
	port := ":80"

	flag.Parse()

	router := gin.Default()

	if *debug {
		port = ":9001"

		router.Use(func(c *gin.Context) {
			headers := c.Writer.Header()

			headers.Set("Access-Control-Allow-Origin", "*")
			headers.Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, auth, token")
			headers.Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")

			c.Next()

			if c.Request.Method == http.MethodOptions {
				c.Status(200)
			} else {
				headers.Set("Content-Type", "application/json")
			}
		})
	}

	apiServerURL, err := url.Parse("http://api-server")

	if err != nil {
		panic(err)
	}

	apiHandler := httputil.NewSingleHostReverseProxy(apiServerURL)

	router.Any("/api/*path", func(c *gin.Context) {
		apiHandler.ServeHTTP(c.Writer, c.Request)
	})

	if err != nil {
		panic(err)
	}

	router.NoRoute(RjServe("/", NewRjFileSystem("./static/")))

	// router.Use()

	router.GET("/hello-world", func(c *gin.Context) {
		c.Writer.Write([]byte("HELLO WORLD!"))
	})

	log.Fatal(router.Run(port))
}
