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
				Type:        graphql.String,
				Description: "Unit of measure which weight is reported -- either g or ml.",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "Description of the serving",
			},
			"servingstate": &graphql.Field{
				Type: graphql.String,
			},
			"weight": &graphql.Field{
				Type:        graphql.Float,
				Description: "unit of measure equilavent weight",
			},
			"servingamount": &graphql.Field{
				Type:        graphql.Float,
				Description: "Portion size",
			},
			"datapoints": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of data points used in calculating the serving",
			},
		},
	})
	/*type FoodGroup struct {
		ID          int32  `json:"id" binding:"required"`
		Code        string `json:"code,omitempty"`
		Description string `json:"description" binding:"required"`
		LastUpdate  string `json:"lastUpdate,omitempty"`
		Type        string `json:"type" binding:"required"`
	}*/
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
				Description: "List of ingredients only available for Branded Food Products items",
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
			"servings": &graphql.Field{
				Type:        graphql.NewList(servingType),
				Description: "Portion information.  A food may have several.",
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

	browseRequestType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "browse",
		Fields: graphql.InputObjectConfigFieldMap{
			"max": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"page": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"sort": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"order": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"source": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
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
				/*Args: graphql.FieldConfigArgument{
					"max": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"sort": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"order": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"source": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},*/
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
						where = where + sourceFilter(source)
						fmt.Printf("WHERE=%s", where)
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
				Type: graphql.NewList(nutrientType),
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
func useIndex(sort string, order string) string {
	useindex := ""
	switch sort {
	case "foodDescription":
		useindex = "idx_fd"
	case "company":
		useindex = "idx_company"
	case "fdcid":
	default:
		useindex = "idx_fdcId"
	}
	if order == "desc" {
		useindex = useindex + "_desc"
	} else {
		useindex = useindex + "_asc"
	}
	return useindex
}
func sourceFilter(s string) string {
	w := ""
	if s != "" {
		if s == "BFPD" {
			w = fmt.Sprintf(" AND ( dataSource = '%s' OR dataSource='%s' )", "LI", "GDSN")
		} else {
			w = fmt.Sprintf(" AND dataSource = '%s'", s)
		}
	}
	return w
}
