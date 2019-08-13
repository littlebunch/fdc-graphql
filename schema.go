package main

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds/cb"
	fdc "github.com/littlebunch/fdc-api/model"
	"gopkg.in/couchbase/gocb.v1"
)

// Create and return the FDC schema which is based on the fdc.Foods package
func initSchema(cb cb.Cb) (graphql.Schema, error) {
	/*
		type Serving struct {
			Nutrientbasis string  `json:"100UnitNutrientBasis,omitempty"`
			Description   string  `json:"householdServingUom"`
			Servingstate  string  `json:"servingState,omitempty"`
			Weight        float32 `json:"weightInGmOrMl"`
			Servingamount float32 `json:"householdServingValue,omitempty"`
			Datapoints    int32   `json:"datapoints,omitempty"`
			}
	*/
	servingType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Serving",
		Fields: graphql.Fields{
			"nutrientbasis": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"servingstate": &graphql.Field{
				Type: graphql.String,
			},
			"weight": &graphql.Field{
				Type: graphql.Float,
			},
			"servingamount": &graphql.Field{
				Type: graphql.Float,
			},
			"datapoints": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})
	/*
		type Food struct {
			UpdatedAt       time.Time  `json:"lastChangeDateTime,omitempty"`
			FdcID           string     `json:"fdcId" binding:"required"`
			NdbNo           string     `json:"ndbno,omitempty"`
			Upc             string     `json:"upc,omitempty"`
			Description     string     `json:"foodDescription" binding:"required"`
			Source          string     `json:"dataSource"`
			PublicationDate time.Time  `json:"publicationDateTime"`
			Ingredients     string     `json:"ingredients,omitempty"`
			Manufacturer    string     `json:"company,omitempty"`
			Group           *FoodGroup `json:"foodGroup,omitempty"`
			Servings        []Serving  `json:"servingSizes,omitempty"``
			Type       string      `json:"type" binding:"required"`
			InputFoods []InputFood `json:"inputfoods,omitempty"`
		}
	*/
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
	/*
		type Nutrient struct {
			NutrientID uint   `json:"id" binding:"required"`
			Nutrientno uint   `json:"nutrientno" binding:"required"`
			Tagname    string `json:"tagname,omitempty"`
			Name       string `json:"name"  binding:"required"`
			Unit       string `json:"unit"  binding:"required"`
			Type       string `json:"type"  binding:"required"`
		}
	*/
	nutrientType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Nutrient",
		Fields: graphql.Fields{
			"nutrientno": &graphql.Field{
				Type: graphql.Int,
			},
			"tagname": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"unit": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	/*
		type Derivation struct {
		ID          int32  `json:"id" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
		Type        string `json:"type"`
	}*/
	derivationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Derivation",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"code": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.Float,
			},
		},
	})
	/*
		type NutrientData struct {
			FdcID      string      `json:"fdcId" binding:"required"`
			Source     string      `json:"Datasource"`
			Type       string      `json:"type"`
			Value      float32     `json:"valuePer100UnitServing"`
			Unit       string      `json:"unit"  binding:"required"`
			Derivation *Derivation `json:"derivation,omitempty"`
			Nutrientno uint        `json:"nutrientNumber"`
			Nutrient   string      `json:"nutrientName"`
			Datapoints int         `json:"datapoints,omitempty"`
			Min        float32     `json:"min,omitempty"`
			Max        float32     `json:"max,omitempty"`
		}
	*/
	nutrientDataType := graphql.NewObject(graphql.ObjectConfig{
		Name: "NutrientData",
		Fields: graphql.Fields{
			"fdcId": &graphql.Field{
				Type: graphql.String,
			},
			"source": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type: graphql.Float,
			},
			"unit": &graphql.Field{
				Type: graphql.String,
			},
			"nutrientno": &graphql.Field{
				Type: graphql.Int,
			},
			"nutrient": &graphql.Field{
				Type: graphql.String,
			},
			"datapoints": &graphql.Field{
				Type: graphql.Int,
			},
			"min": &graphql.Field{
				Type: graphql.Float,
			},
			"max": &graphql.Field{
				Type: graphql.Float,
			},
			"derivation": &graphql.Field{
				Type: derivationType,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	/*nutrientIDType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "nutids",
		Fields: graphql.InputObjectConfigFieldMap{
			"ids": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})*/
	// Define the schema
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
			"nutrients": &graphql.Field{
				Type: graphql.NewList(nutrientType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					q := fmt.Sprintf("select nutrient.* from %s as nutrient where type=\"NUT\"", cs.CouchDb.Bucket)
					query := gocb.NewN1qlQuery(q)
					rows, err := cb.Conn.ExecuteN1qlQuery(query, nil)
					if err != nil {
						return nil, err

					}
					var nutrients []fdc.Nutrient
					var row fdc.Nutrient
					for rows.Next(&row) {
						nutrients = append(nutrients, row)
					}
					return nutrients, nil
				},
			},
			"nutrientdata": &graphql.Field{
				Type: graphql.NewList(nutrientDataType),
				Args: graphql.FieldConfigArgument{
					"fdcid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},

					"nutids": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					var (
						nut     fdc.NutrientData
						nutdata []fdc.NutrientData
						nIDs    []int
						q       string
					)

					nut.FdcID = p.Args["fdcid"].(string)
					// build an int array of nutrient numbers
					for _, gnid := range p.Args["nutids"].([]interface{}) {
						nIDs = append(nIDs, gnid.(int))
					}
					// put the nutrientno array into a string for the query
					nstr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(nIDs)), ","), "[]")

					if nstr != "" {
						q = fmt.Sprintf("select nutrientdata.* from %s as nutrientdata where type=\"NUTDATA\" and fdcId=\"%s\" and nutrientNumber in [%s]", cs.CouchDb.Bucket, nut.FdcID, nstr)
					} else {
						q = fmt.Sprintf("select nutrientdata.* from %s as nutrientdata where type=\"NUTDATA\" and fdcId=\"%s\"", cs.CouchDb.Bucket, nut.FdcID)
					}
					query := gocb.NewN1qlQuery(q)
					rows, err := cb.Conn.ExecuteN1qlQuery(query, nil)
					if err != nil {
						return nil, err
					}
					// put the query results into the nutrientdata array
					for rows.Next(&nut) {
						nutdata = append(nutdata, nut)
					}
					return nutdata, nil
				},
			},
		},
	})
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	return s, err
}
