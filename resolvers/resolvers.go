package resolvers

import (
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds/cb"
	fdc "github.com/littlebunch/fdc-api/model"
	"github.com/littlebunch/fdc-graphql/utils"
	"gopkg.in/couchbase/gocb.v1"
)

//Resolver type for resolving queries
type Resolver struct {
	Ds *cb.Cb
	Cs fdc.Config
}

//Food queries for a single Food by fdcId
func (r *Resolver) Food(p graphql.ResolveParams) (interface{}, error) {
	var food fdc.Food
	food.FdcID = p.Args["id"].(string)
	err := r.Ds.Get(food.FdcID, &food)
	if err != nil {
		return nil, err
	}
	return food, nil
}

//FoodSearchCount finds the number of hits for a proposed SearchRequest
func (r *Resolver) FoodSearchCount(p graphql.ResolveParams) (interface{}, error) {
	var (
		sr        fdc.SearchRequest
		err, errs error
		c         int
		rs        []interface{}
	)
	sr, errs = utils.Searchquery(p)
	sr.IndexName = r.Cs.CouchDb.Fts
	sr.Max = 1
	if c, err = r.Ds.Search(sr, &rs); err != nil {
		return nil, err
	}
	return c, errs
}

//Foods queries a list of Food objects by a list of fdcIds
func (r *Resolver) Foods(p graphql.ResolveParams) (interface{}, error) {
	var (
		dt   *fdc.DocType
		s    string
		errs error
		fIDs string
	)
	where := fmt.Sprintf("type=\"%s\" ", dt.ToString(fdc.FOOD))
	if p.Args["fdcids"] != nil {
		fIDs, s = utils.Fdcids(p.Args["fdcids"].([]interface{}))
		if s != "" {
			errs = utils.Seterror(&errs, s)
		}
		if fIDs != "" {
			where += fmt.Sprintf("AND fdcId in [%s]", fIDs)
		}
	}

	rs, _ := r.Ds.Browse(r.Cs.CouchDb.Bucket, where, int64(0), int64(2), "fdcId", "desc")
	return rs, errs
}

//FoodSearch query for a SearchRequest
func (r *Resolver) FoodSearch(p graphql.ResolveParams) (interface{}, error) {
	var (
		sr        fdc.SearchRequest
		err, errs error
		rs        []interface{}
	)
	sr, errs = utils.Searchquery(p)

	sr.IndexName = r.Cs.CouchDb.Fts
	if _, err = r.Ds.Search(sr, &rs); err != nil {
		return nil, err
	}
	return rs, errs
}

//FoodsBrowse queries a list of foods based on a Browse object
func (r *Resolver) FoodsBrowse(p graphql.ResolveParams) (interface{}, error) {
	var (
		dt                  *fdc.DocType
		max, page           int
		sort, order, source string
		errs                error
	)
	b := p.Args["browse"].(map[string]interface{})
	if b["max"] == nil {
		max = 50
	} else {
		max = b["max"].(int)
	}
	if max > 150 {
		errs = utils.Seterror(&errs, "cannot return more than 150 items")
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
	if max > utils.MAXPAGE {
		errs = utils.Seterror(&errs, fmt.Sprintf("max parameter cannot exceed %d", utils.MAXPAGE))
		max = utils.MAXPAGE
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
		errs = utils.Seterror(&errs, "unrecognized sort parameter.  Must be 'company', 'foodDescription' or 'fdcId'")
		sort = "fdcId"
	}
	offset := page * max
	where := fmt.Sprintf("type=\"%s\" ", dt.ToString(fdc.FOOD))

	if source != "" {
		where = where + fmt.Sprintf(" AND dataSource = '%s'", source)
	}
	rs, _ := r.Ds.Browse(r.Cs.CouchDb.Bucket, where, int64(offset), int64(max), sort, order)
	return rs, errs
}

//Nutrientdata queries a list of Nutrientdata based on a list of fdcIds and nutrientIds
func (r *Resolver) Nutrientdata(p graphql.ResolveParams) (interface{}, error) {

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
	fIDs, _ = utils.Fdcids(p.Args["fdcids"].([]interface{}))

	// build an int array of nutrient numbers
	for _, gnid := range p.Args["nutids"].([]interface{}) {
		nIDs = append(nIDs, gnid.(int))
	}
	// put the nutrientno array into a string for the query
	nstr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(nIDs)), ","), "[]")

	if nstr != "" {
		q = fmt.Sprintf("select nutrientdata.* from %s as nutrientdata where type=\"NUTDATA\" and fdcId in [%s] and nutrientNumber in [%s] order by fdcId,nutrientNumber", r.Cs.CouchDb.Bucket, fIDs, nstr)
	} else {
		q = fmt.Sprintf("select nutrientdata.* from %s as nutrientdata where type=\"NUTDATA\" and fdcId in [%s] order by fdcId,nutrientNumber", r.Cs.CouchDb.Bucket, fIDs)
	}
	rows, err = r.Ds.Conn.ExecuteN1qlQuery(gocb.NewN1qlQuery(q), nil)
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
}

//Nutrients queries a list of nutrients
func (r *Resolver) Nutrients(p graphql.ResolveParams) (interface{}, error) {
	var dt *fdc.DocType
	return r.Ds.GetDictionary(r.Cs.CouchDb.Bucket, dt.ToString(fdc.NUT), 0, 300)
}
