// Package main creates and starts a web server
package main

// @APITitle Brand Foods Product Database

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
	fdc "github.com/littlebunch/fdc-api/model"
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
	l   = flag.String("l", "/tmp/bfpd.out", "send log output to this file -- defaults to /tmp/bfpd.out")
	p   = flag.String("p", "8000", "TCP port to used")
	r   = flag.String("r", "v1", "root path to deploy -- defaults to 'v1'")
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

	//var cb cb.Cb
	flag.Parse()
	// get configuration
	/*cs.GetConfig(c)
	// Create a datastore and connect to it
	dc = &cb
	err = dc.ConnectDs(cs)
	if err != nil {
		log.Fatal("Cannot get datastore connection %v.", err)
	}
	defer dc.CloseDs()*/
	// initialize our jwt authentication
	//var u *auth.User
	//if *i {
	//	u.BootstrapUsers(session, cs.MongoDb.Collection)
	//}
	//authMiddleware := u.AuthMiddleware(session, cs.MongoDb.Collection)
	//router := gin.Default()

	foodType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Food",
		Fields: graphql.Fields{
			"fdcId": &graphql.Field{
				Type: graphql.String,
			},
			"upc": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"ingredients": &graphql.Field{
				Type: graphql.String,
			},
			"dataSource": &graphql.Field{
				Type: graphql.String,
			},
			"company": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"foods": &graphql.Field{
				Type: graphql.NewList(foodType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					/*query := gocb.NewN1qlQuery("SELECT META(account).id, account.* FROM example AS account WHERE account.type = 'account'")
					rows, err := bucket.ExecuteN1qlQuery(query, nil)
					if err != nil {
						return nil, err
					}
					var accounts []Account
					var row Account
					for rows.Next(&row) {
						accounts = append(accounts, row)
					}*/
					return nil, nil
				},
			},
		},
	})
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group(fmt.Sprintf("%s", *r))
	{
		//v1.POST("/login", authMiddleware.LoginHandler)
		v1.GET("/", gin.WrapH(handler.Playground("GraphQL playground", "/api")))
		v1.GET("/graphql", func(c *gin.Context) {
			result := graphql.Do(graphql.Params{
				Schema:        schema,
				RequestString: c.Param("query"),
			})
			c.JSON(http.StatusOK, result)
		})

	}
	endless.ListenAndServe(":"+*p, router)

}
