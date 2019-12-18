package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds/cb"
	fdc "github.com/littlebunch/fdc-api/model"
	"gopkg.in/couchbase/gocb.v1"
)

// InitSchema -- Create and return the FDC schema which is based on the fdc.Foods package
func InitSchema(cb cb.Cb, cs fdc.Config) (graphql.Schema, error) {
	ds := &cb
	// food servings
	servingType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Serving",
		Fields: graphql.Fields{
			"nutrientBasis": &graphql.Field{
				Type:        graphql.String,
				Description: "Unit of measure which weight is reported -- either g or ml.",
			},
			"servingUnit": &graphql.Field{
				Type:        graphql.String,
				Description: "The household description of the serving",
			},
			"servingState": &graphql.Field{
				Type: graphql.String,
			},
			"weight": &graphql.Field{
				Type:        graphql.Float,
				Description: "unit of measure equilavent weight",
			},
			"value": &graphql.Field{
				Type:        graphql.Float,
				Description: "Portion size",
			},
			"dataPoints": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of data points used in calculating the serving",
			},
		},
	})
	// food categories
	foodGroupType := graphql.NewObject(graphql.ObjectConfig{
		Name: "foodGroup",
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
		},
	})
	// food meta data
	foodType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Food",
		Fields: graphql.Fields{
			"fdcId": &graphql.Field{
				Type:        graphql.String,
				Description: "Food Data Central ID assigned to the food",
			},
			"upc": &graphql.Field{
				Type:        graphql.String,
				Description: "UPC or GTIN number assigned to the food. Applies to Branded Food Products only",
			},
			"foodDescription": &graphql.Field{
				Type:        graphql.String,
				Description: "Name of the food",
			},
			"ingredients": &graphql.Field{
				Type:        graphql.String,
				Description: "The list of ingredients (as it appears on the product label).  Only available for Branded Food Products items",
			},
			"dataSource": &graphql.Field{
				Type:        graphql.String,
				Description: "Source of the food data.  SR = Standard Reference Legacy; FNDDS = Food Survey; GDSN = Global Food; LI = Label Insight",
			},
			"company": &graphql.Field{
				Type:        graphql.String,
				Description: "Manufacturer of the food",
			},
			"foodGroup": &graphql.Field{
				Type:        foodGroupType,
				Description: "Category assigned to the food.  Differs by dataSource",
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"servingSizes": &graphql.Field{
				Type:        graphql.NewList(servingType),
				Description: "Portion information.  A food may have several.",
			},
		},
	})
	// nutrient
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
	// nutrient derivation
	derivationType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Derivation",
		Description: "Procedure indicating how a food nutrient value was obtained",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"code": &graphql.Field{
				Type:        graphql.String,
				Description: "Code used for the derivation (e.g. A means analytical)",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "Description of the derivation",
			},
			"type": &graphql.Field{
				Type: graphql.Float,
			},
		},
	})
	// food nutrient data
	nutrientDataType := graphql.NewObject(graphql.ObjectConfig{
		Name: "NutrientData",
		Fields: graphql.Fields{
			"fdcId": &graphql.Field{
				Type:        graphql.String,
				Description: "Food Data Central id of the food to which the data belongs",
			},
			"source": &graphql.Field{
				Type: graphql.String,
			},
			"value": &graphql.Field{
				Type:        graphql.Float,
				Description: "Amount of the nutrient per 100g of food. Specified in unit defined in the unit field.",
			},
			"unit": &graphql.Field{
				Type:        graphql.String,
				Description: "The standard unit of measure for the nutrient (g or ml)",
			},
			"nutrientno": &graphql.Field{
				Type:        graphql.Int,
				Description: "ID of the nutrient to which the food nutrient pertains",
			},
			"nutrient": &graphql.Field{
				Type:        graphql.String,
				Description: "Name of the nutrient",
			},
			"datapoints": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of observations on which the value is based",
			},
			"min": &graphql.Field{
				Type:        graphql.Float,
				Description: "Minimum number of observations",
			},
			"max": &graphql.Field{
				Type:        graphql.Float,
				Description: "Maximum number of observations",
			},
			"derivation": &graphql.Field{
				Type:        derivationType,
				Description: "Derivation information",
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// browse request
	browseRequestType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "browse",
		Description: "Describes parameters for browse queries",
		Fields: graphql.InputObjectConfigFieldMap{
			"max": &graphql.InputObjectFieldConfig{
				Type:        graphql.Int,
				Description: "Maximum number of items to be returned.",
			},
			"page": &graphql.InputObjectFieldConfig{
				Type:        graphql.Int,
				Description: "Page as defined by the max parameter to start the list.  This is zero based and used to determine offsets in the list.",
			},
			"sort": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Field on which browse results are to be sorted.",
			},
			"order": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Sort order -- ASC or DESC.",
			},
			"searchTerms": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Terms to filter results list. (optional)",
			},
			"searchField": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Field to limit term searches (optional)",
			},
			"searchType": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Type of search to perform: MATCH, PHRASE, WILDCARD or REGEX. (optional)  ",
			},
		},
	})
	// Define the schema
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"foods": &graphql.Field{
				Type: graphql.NewList(foodType),
				Args: graphql.FieldConfigArgument{
					"browse": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(browseRequestType),
					},
				},
				Description: "Returns a list of foods.  Parameters sent in the browse input object.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var (
						dt                  *fdc.DocType
						max, page           int
						sort, order, source string
					)
					b := p.Args["browse"].(map[string]interface{})
					if b["max"] == nil {
						max = 50
					} else {
						max = b["max"].(int)
					}
					if b["page"] == nil {
						page = 0
					} else {
						page = b["page"].(int)
					}
					if b["sort"] == nil {
						sort = "fdcId"
					} else {
						sort = b["sort"].(string)
					}
					if b["order"] == nil {
						order = "ASC"
					} else {
						order = b["order"].(string)
					}
					if b["source"] != nil {
						source = b["source"].(string)
					}
					if max == 0 {
						max = 50
					}
					if max > 150 {
						return nil, errors.New("max parameter cannot exceed 150")
					}
					if page < 0 {
						page = 0
					}
					if sort == "" {
						sort = "fdcId"
					}
					if order == "" {
						order = "ASC"
					}

					offset := page * max
					where := fmt.Sprintf("type=\"%s\" ", dt.ToString(fdc.FOOD))
					if source != "" {
						where = where + fmt.Sprintf(" AND dataSource = '%s'", source)
					}

					return ds.Browse(cs.CouchDb.Bucket, where, int64(offset), int64(max), sort, order)
				},
			},
			"food": &graphql.Field{
				Type: foodType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Description: "Returns a food for a given fdcId.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var food fdc.Food
					food.FdcID = p.Args["id"].(string)
					err := ds.Get(food.FdcID, &food)
					if err != nil {
						return nil, err
					}
					return food, nil
				},
			},
			"nutrients": &graphql.Field{
				Type:        graphql.NewList(nutrientType),
				Description: "Returns a list of nutrients used in the database",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var dt *fdc.DocType
					return ds.GetDictionary(cs.CouchDb.Bucket, dt.ToString(fdc.NUT), 0, 300)
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
				Description: "Returns one or more nutrient values for a food.",
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
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
}
