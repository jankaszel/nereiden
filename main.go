package main

import (
	"github.com/falafeljan/gin-simple-token-middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"fmt"
	"log"
	"strings"
	"time"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins: args.AllowedOrigins,
		AllowHeaders: []string{"Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	router.ForwardedByClientIP = true
	router.Use(limiterMiddleware(args.RateLimit))

	group := router.Group("/", tokenmiddleware.NewHandler(token))
	group.POST("/graphql", createGraphQLHandler(args.LetsEncryptEmail))

	log.Fatal(router.Run(
		strings.Join([]string{":", args.HTTPPort}, "")))
}
