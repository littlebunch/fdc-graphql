// Package main creates and starts a web server providing a REST endpoint for
// GraphQL queries
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds"
	"github.com/littlebunch/fdc-api/ds/cb"
	fdc "github.com/littlebunch/fdc-api/model"
	"github.com/littlebunch/fdc-graphql/schema"
)

const (
	maxListSize    = 150
	defaultListMax = 50
	apiVersion     = "1.0.0 Beta"
)

var (
	d   = flag.Bool("d", false, "Debug")
	i   = flag.Bool("i", false, "Initialize the authentication store")
	c   = flag.String("c", "config.yml", "YAML Config file")
	l   = flag.String("l", "/tmp/fdcgql.out", "send log output to this file -- defaults to /tmp/fdcgcl.out")
	p   = flag.String("p", "8000", "TCP port to used")
	r   = flag.String("r", "graphql", "root path to deploy -- defaults to 'v1'")
	cs  fdc.Config
	err error
	dc  ds.DataSource
)

// process cli flags; build the config and init an Mongo client and a logger
func init() {
	var (
		lfile *os.File
	)
	lfile, err = os.OpenFile(*l, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", *l, ":", err)
	}
	m := io.MultiWriter(lfile, os.Stdout)
	log.SetOutput(m)
}

func main() {

	var cb cb.Cb
	flag.Parse()
	// get configuration
	cs.GetConfig(c)
	// Create a datastore and connect to it
	dc = &cb
	err = dc.ConnectDs(cs)
	if err != nil {
		log.Fatal("Cannot get datastore connection %v.", err)
	}
	defer dc.CloseDs()
	// initialize our jwt authentication
	//var u *auth.User
	//if *i {
	//	u.BootstrapUsers(session, cs.MongoDb.Collection)
	//}
	//authMiddleware := u.AuthMiddleware(session, cs.MongoDb.Collection)
	//router := gin.Default()
	schema, err := schema.InitSchema(cb, cs)
	if err != nil {
		log.Fatal("Cannot create the schema %v\n", err)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group(fmt.Sprintf("%s", *r))
	{
		//v1.POST("/login", authMiddleware.LoginHandler)
		v1.GET("/", gin.WrapH(handler.Playground("GraphQL playground", "/graphql")))
		v1.GET("", func(c *gin.Context) {
			result := graphql.Do(graphql.Params{
				Schema:        schema,
				RequestString: c.Query("query"),
			})
			c.JSON(http.StatusOK, result)
		})
		v1.POST("", func(c *gin.Context) {
			type Q struct {
				Query string `json:"query"`
			}
			var q Q
			err := c.BindJSON(&q)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":  "json decoding : " + err.Error(),
					"status": http.StatusBadRequest,
				})
				return
			}

			result := graphql.Do(graphql.Params{
				Schema:        schema,
				RequestString: q.Query,
			})
			c.JSON(http.StatusOK, result)
		})

	}
	endless.ListenAndServe(":"+*p, router)

}
