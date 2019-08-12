package main

import (
	"encoding/json"
	"fmt"

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
			"Nutrientno": &graphql.Field{
				Type: graphql.Int,
			},
			"Tagname": &graphql.Field{
				Type: graphql.String,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Unit": &graphql.Field{
				Type: graphql.String,
			},
			"Type": &graphql.Field{
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
			"ID": &graphql.Field{
				Type: graphql.Int,
			},
			"Code": &graphql.Field{
				Type: graphql.String,
			},
			"Description": &graphql.Field{
				Type: graphql.String,
			},
			"Type": &graphql.Field{
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
			"Source": &graphql.Field{
				Type: graphql.String,
			},
			"Value": &graphql.Field{
				Type: graphql.String,
			},
			"Unit": &graphql.Field{
				Type: graphql.String,
			},
			"Nutrientno": &graphql.Field{
				Type: graphql.Int,
			},
			"Nutrient": &graphql.Field{
				Type: graphql.String,
			},
			"Datapoints": &graphql.Field{
				Type: graphql.Int,
			},
			"Min": &graphql.Field{
				Type: graphql.Float,
			},
			"Max": &graphql.Field{
				Type: graphql.Float,
			},
			"Derivation": &graphql.Field{
				Type: derivationType,
			},
			"Type": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	nutrientIDType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "nutids",
		Fields: graphql.InputObjectConfigFieldMap{
			"ids": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
	})
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
						Type: nutrientIDType,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					type P struct {
						nids []int
					}
					var (
						nut     fdc.NutrientData
						nutdata []fdc.NutrientData
						nids    P
					)

					nut.FdcID = p.Args["fdcid"].(string)
					nutids := p.Args["nutids"].(map[string]interface{})
					fmt.Println("NUTIDS ", nutids)
					b, _ := json.Marshal(nutids)
					json.Unmarshal(b, &nids)
					fmt.Println("NIDS=", nids.nids)
					if true {
						q := fmt.Sprintf("select nutrientdata.* from %s as nutrientdata where type=\"NUTDATA\" and fdcId=\"%s\"", cs.CouchDb.Bucket, nut.FdcID)
						query := gocb.NewN1qlQuery(q)
						rows, err := cb.Conn.ExecuteN1qlQuery(query, nil)
						if err != nil {
							return nil, err
						}
						for rows.Next(&nut) {
							nutdata = append(nutdata, nut)
						}
					} else {
						err := dc.Get(fmt.Sprintf("%s_%d", nut.FdcID, nut.Nutrientno), &nut)
						if err != nil {
							return nil, err
						}
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
