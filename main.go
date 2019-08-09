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

	gocb "gopkg.in/couchbase/gocb.v1"

	"github.com/99designs/gqlgen/handler"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds"
	"github.com/littlebunch/fdc-api/ds/cb"
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
	/*Nutrientbasis string  `json:"100UnitNutrientBasis,omitempty"`
	Description   string  `json:"householdServingUom"`
	Servingstate  string  `json:"servingState,omitempty"`
	Weight        float32 `json:"weightInGmOrMl"`
	Servingamount float32 `json:"householdServingValue,omitempty"`
	Datapoints    int32   `json:"datapoints,omitempty"`*/

	servingType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Serving",
		Fields: graphql.Fields{
			"Nutrientbasis": &graphql.Field{
				Type: graphql.String,
			},
			"Description": &graphql.Field{
				Type: graphql.String,
			},
			"Servingstate": &graphql.Field{
				Type: graphql.String,
			},
			"Weight": &graphql.Field{
				Type: graphql.Float,
			},
			"Servingamount": &graphql.Field{
				Type: graphql.Float,
			},
			"Datapoints": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})
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
			"servings": &graphql.Field{
				Type: graphql.NewList(servingType),
			},
		},
	})
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"foods": &graphql.Field{
				Type: graphql.NewList(foodType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					query := gocb.NewN1qlQuery("select food.* from gnutdata as food where type=\"FOOD\" AND dataSource=\"LI\" LIMIT 50")
					rows, err := cb.Conn.ExecuteN1qlQuery(query, nil)
					if err != nil {
						return nil, err

					}
					var foods []fdc.Food
					var row fdc.Food
					for rows.Next(&row) {
						foods = append(foods, row)
					}
					return foods, nil
				},
			},
			"food": &graphql.Field{
				Type: foodType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var food fdc.Food
					food.FdcID = p.Args["id"].(string)
					err := dc.Get(food.FdcID, &food)
					if err != nil {
						return nil, err
					}
					return food, nil
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
		v1.GET("/", gin.WrapH(handler.Playground("GraphQL playground", "")))
		v1.GET("", func(c *gin.Context) {
			fmt.Printf("Query=%s\n", c.Query("query"))
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
