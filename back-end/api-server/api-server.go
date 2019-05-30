package main

import (
	"log"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	var err error
	var hashedPassword []byte

	for hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14); err != nil; {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 14)
	}

	return string(hashedPassword)
}

func getUUID() string {
	var err error
	var uid uuid.UUID

	for uid, err = uuid.NewV4(); err != nil; {
		uid, err = uuid.NewV4()
	}

	return uid.String()
}

func main() {
	port := ":80"
	router := gin.Default()

	apiNoAuth := router.Group("/api")

	apiNoAuth.GET("/hello-world", func(c *gin.Context) {
		c.Writer.Write([]byte("HELLO WORLD!"))
	})

	// apiAuth := router.Group("/api")

	log.Fatal(router.Run(port))
}
