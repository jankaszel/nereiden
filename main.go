package main

import (
	"fmt"
	"github.com/falafeljan/gin-simple-token-middleware"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func main() {
	args := getArgs()

	if args.InProduction {
		gin.SetMode("release")
	}

	token, err := generateToken()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The security token for this instance is:\n\t%s\n"+
		"Please use this token with the 'Authorization' HTTP header as `Token <token>`.\n"+
		"For more information, please refer to the documentation.\n\n", token)

	router := gin.Default()
	router.Use(createCORSMiddleware(args.AllowedOrigins))
	router.Use(limiterMiddleware(args.RateLimit))

	group := router.Group("/", tokenmiddleware.NewHandler(token))
	group.POST("/graphql", createGraphQLHandler(args.LetsEncryptEmail))

	log.Fatal(router.Run(
		strings.Join([]string{":", args.HTTPPort}, "")))
}
