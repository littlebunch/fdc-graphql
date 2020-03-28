package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds/cb"
	fdc "github.com/littlebunch/fdc-api/model"
	"github.com/littlebunch/fdc-graphql/schema"
)

func foodResolver(p graphql.ResolveParams, ds *cb.Cb) (interface{}, error) {
	var food fdc.Food
	food.FdcID = p.Args["id"].(string)
	err := ds.Get(food.FdcID, &food)
	if err != nil {
		return nil, err
	}
	return food, nil
}
func fdcids(fids []interface{}) (string, string) {
	i := 0
	fIDs := ""
	var err string
	for _, fid := range fids {
		fIDs += fmt.Sprintf("\"%s\",", fid.(string))
		i++
		if i > schema.MAXIDS {
			err = fmt.Sprintf("number of fdcId's should not exceed %d", schema.MAXIDS)
			break
		}
	}
	fIDs = strings.Trim(fIDs, ",")
	return fIDs, err
}
func setError(err *error, msg string) error {

	if *err == nil {
		*err = errors.New(msg)
	} else {
		*err = fmt.Errorf("%w;"+msg, *err)
	}
	return *err

}
func searchquery(p graphql.ResolveParams) (fdc.SearchRequest, error) {
	var (
		sr        fdc.SearchRequest
		max, page int
		errs      error
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
	if max > schema.MAXPAGE {
		errs = setError(&errs, fmt.Sprintf("max parameter cannot exceed %d", schema.MAXPAGE))
		max = schema.MAXPAGE
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

	return sr, errs
}
