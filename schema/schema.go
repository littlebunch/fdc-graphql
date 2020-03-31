package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds/cb"
	fdc "github.com/littlebunch/fdc-api/model"
	"github.com/littlebunch/fdc-graphql/resolvers"
	"github.com/littlebunch/fdc-graphql/types"
)

// InitSchema -- Create and return the FDC schema which is based on the fdc.Foods package
func InitSchema(cb cb.Cb, cs fdc.Config) (graphql.Schema, error) {
	var t types.Types
	r := resolvers.Resolver{Ds: &cb, Cs: cs}
	t.InitTypes()
	// Define the queries
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"foods": &graphql.Field{
				Type: graphql.NewList(t.Food),
				Args: graphql.FieldConfigArgument{
					"fdcids": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewList(graphql.String)),
					},
				},
				Description: "Returns a list of foods.  Parameters sent in the browse input object.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return r.Foods(p)
				},
			},
			"foodsBrowse": &graphql.Field{
				Type: graphql.NewList(t.Food),
				Args: graphql.FieldConfigArgument{
					"browse": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(t.BrowseRequest),
					},
				},
				Description: "Returns a list of foods.  Parameters sent in the browse input object.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return r.FoodsBrowse(p)
				},
			},
			"food": &graphql.Field{
				Type: t.Food,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Description: "Returns a food for a given fdcId.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return r.Food(p)
				},
			},

			"nutrients": &graphql.Field{
				Type:        graphql.NewList(t.Nutrient),
				Description: "Returns a list of nutrients used in the database",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return r.Nutrients(p)
				},
			},
			"nutrientdata": &graphql.Field{
				Type: graphql.NewList(t.NutrientData),
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
					return r.Nutrientdata(p)
				},
			},
			"foodsSearchCount": &graphql.Field{
				Type: graphql.Int,
				Args: graphql.FieldConfigArgument{
					"search": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(t.SearchRequest),
					},
				},
				Description: "Returns a count of items returned by a search",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return r.FoodSearchCount(p)
				},
			},
			"foodsSearch": &graphql.Field{
				Type: graphql.NewList(t.FoodSearch),
				Args: graphql.FieldConfigArgument{
					"search": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(t.SearchRequest),
					},
				},
				Description: "Returns a list of foods.  Parameters sent in the browse input object.",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return r.FoodSearch(p)
				},
			},
		},
	})
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
}
