package schema

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds/cb"
	fdc "github.com/littlebunch/fdc-api/model"
	"gopkg.in/couchbase/gocb.v1"
)

// Maximum number of FDC Id's that may be requested per query
const (
	MAXIDS  = 100
	MAXPAGE = 150
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
	type FoodMeta struct {
		FdcID        string `json:"fdcId" binding:"required"`
		Upc          string `json:"upc"`
		Description  string `json:"foodDescription" binding:"required"`
		Ingredients  string `json:"ingredients,omitempty"`
		Source       string `json:"dataSource"`
		Manufacturer string `json:"company,omitempty"`
		Type         string `json:"type"`
		Category     string `json:"foodgroup.description,omitempty"`
	}
	foodSearchType := graphql.NewObject(graphql.ObjectConfig{
		Name: "FoodSearch",
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
			"category": &graphql.Field{
				Type:        graphql.String,
				Description: "Category assigned to the food.  Differs by dataSource",
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

	searchRequestType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:        "query",
		Description: "Describes parameters for search queries",
		Fields: graphql.InputObjectConfigFieldMap{
			"terms": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Terms to include in the search",
			},
			"field": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Limit search terms to a particular field ",
			},
			"type": &graphql.InputObjectFieldConfig{
				Type:        graphql.String,
				Description: "Type of search to run",
			},
			"page": &graphql.InputObjectFieldConfig{
				Type:        graphql.Int,
				Description: "Page as defined by the max parameter to start the list.  This is zero based and used to determine offsets in the list.",
			},
			"max": &graphql.InputObjectFieldConfig{
				Type:        graphql.Int,
				Description: "Maximum number of items to return. ",
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
		},
	})

	// Define the schema
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"foods": &graphql.Field{
				Type: graphql.NewList(foodType),
				Args: graphql.FieldConfigArgument{
					"fdcids": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
				},
				Description: "Returns a list of foods.  Parameters sent in the browse input object.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var (
						dt   *fdc.DocType
						s    string
						errs error
						r    interface{}
						fIDs string
					)
					where := fmt.Sprintf("type=\"%s\" ", dt.ToString(fdc.FOOD))
					if p.Args["fdcids"] != nil {
						fIDs, s = fdcids(p.Args["fdcids"].([]interface{}))
						if s != "" {
							errs = setError(&errs, s)
						}
						if fIDs != "" {
							where += fmt.Sprintf("AND fdcId in [%s]", fIDs)
						}
					}

					r, _ = ds.Browse(cs.CouchDb.Bucket, where, int64(0), int64(2), "fdcId", "desc")
					return r, errs
				},
			},
			"foodsBrowse": &graphql.Field{
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
						errs                error
						r                   interface{}
					)
					b := p.Args["browse"].(map[string]interface{})
					if b["max"] == nil {
						max = 50
					} else {
						max = b["max"].(int)
					}
					if max > 150 {
						errs = setError(&errs, "cannot return more than 150 items")
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
					if max > MAXPAGE {
						errs = setError(&errs, fmt.Sprintf("max parameter cannot exceed %d", MAXPAGE))
						max = MAXPAGE
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
					if sort != "foodDescription" && sort != "company" && sort != "fdcId" {
						errs = setError(&errs, "unrecognized sort parameter.  Must be 'company', 'foodDescription' or 'fdcId'")
						sort = "fdcId"
					}
					offset := page * max
					where := fmt.Sprintf("type=\"%s\" ", dt.ToString(fdc.FOOD))

					if source != "" {
						where = where + fmt.Sprintf(" AND dataSource = '%s'", source)
					}
					r, _ = ds.Browse(cs.CouchDb.Bucket, where, int64(offset), int64(max), sort, order)
					return r, errs
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
					"fdcids": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},

					"nutids": &graphql.ArgumentConfig{
						Type: graphql.NewList(graphql.Int),
					},
				},
				Description: "Returns one or more nutrient values for a food.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					var (
						nut       fdc.NutrientData
						nutdata   []fdc.NutrientData
						rows      gocb.QueryResults
						nIDs      []int
						fIDs      string
						q         string
						err, errs error
					)

					// build a string array of FDC id's
					fIDs, _ = fdcids(p.Args["fdcids"].([]interface{}))

					// build an int array of nutrient numbers
					for _, gnid := range p.Args["nutids"].([]interface{}) {
						nIDs = append(nIDs, gnid.(int))
					}
					// put the nutrientno array into a string for the query
					nstr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(nIDs)), ","), "[]")

					if nstr != "" {
						q = fmt.Sprintf("select nutrientdata.* from %s as nutrientdata where type=\"NUTDATA\" and fdcId in [%s] and nutrientNumber in [%s] order by fdcId,nutrientNumber", cs.CouchDb.Bucket, fIDs, nstr)
					} else {
						q = fmt.Sprintf("select nutrientdata.* from %s as nutrientdata where type=\"NUTDATA\" and fdcId in [%s] order by fdcId,nutrientNumber", cs.CouchDb.Bucket, fIDs)
					}
					rows, err = cb.Conn.ExecuteN1qlQuery(gocb.NewN1qlQuery(q), nil)
					if err != nil {
						return nil, err
					}
					// put the query results into the nutrientdata array
					for rows.Next(&nut) {
						nutdata = append(nutdata, nut)
					}
					if err == nil {
						err = errs
					}
					return nutdata, err
				},
			},
			"foodsSearch": &graphql.Field{
				Type: graphql.NewList(foodSearchType),
				Args: graphql.FieldConfigArgument{
					"search": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(searchRequestType),
					},
				},
				Description: "Returns a list of foods.  Parameters sent in the browse input object.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var (
						sr        fdc.SearchRequest
						max, page int
						err, errs error
						r         []interface{}
					)
					b := p.Args["search"].(map[string]interface{})
					if b["max"] == nil {
						max = 50
					} else {
						max = b["max"].(int)
					}
					if max > 150 {
						errs = setError(&errs, "cannot return more than 150 items")
					}
					if b["page"] == nil {
						page = 0
					} else {
						page = b["page"].(int)
					}

					if max <= 0 {
						max = 50
					}
					if max > MAXPAGE {
						errs = setError(&errs, fmt.Sprintf("max parameter cannot exceed %d", MAXPAGE))
						max = MAXPAGE
					}
					if page < 0 {
						page = 0
					}

					sr.Max = max

					if b["type"] != nil {
						t := b["type"].(string)
						if t != fdc.PHRASE && t != fdc.WILDCARD && t != fdc.REGEX {
							errs = setError(&errs, fmt.Sprintf("Search type must be %s, %s, or %s ", fdc.PHRASE, fdc.WILDCARD, fdc.REGEX))
							sr.SearchType = ""
						} else {
							sr.SearchType = t
						}
					}
					if b["field"] != nil {
						sr.SearchField = b["field"].(string)
						if strings.ToLower(sr.SearchField) == "category" {
							sr.SearchField = "foodGroup.description"
						}
					}
					if b["terms"] != nil {
						sr.Query = b["terms"].(string)
					}
					if sr.SearchType == fdc.REGEX {
						sr.SearchField += "_kw"
					}
					sr.Page = page * max
					sr.IndexName = cs.CouchDb.Fts
					if _, err = ds.Search(sr, &r); err != nil {
						return nil, err
					}

					return r, errs
				},
			},
		},
	})
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
}
