package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	fdc "github.com/littlebunch/fdc-api/model"
)

// Maximum number of FDC Id's that may be requested per query
const (
	MAXIDS  = 100
	MAXPAGE = 150
)

//Fdcids creates a csv string representation of an array of fdcids for use in a query
func Fdcids(fids []interface{}) (string, string) {
	i := 0
	fIDs := ""
	var err string
	for _, fid := range fids {
		fIDs += fmt.Sprintf("\"%s\",", fid.(string))
		i++
		if i > MAXIDS {
			err = fmt.Sprintf("number of fdcId's should not exceed %d", MAXIDS)
			break
		}
	}
	fIDs = strings.Trim(fIDs, ",")
	return fIDs, err
}

//Seterror adds an error to an error array
func Seterror(err *error, msg string) error {

	if *err == nil {
		*err = errors.New(msg)
	} else {
		*err = fmt.Errorf("%w;"+msg, *err)
	}
	return *err

}

//Searchquery builds a SearchRequest from query parameters
func Searchquery(p graphql.ResolveParams) (fdc.SearchRequest, error) {
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
		errs = Seterror(&errs, "cannot return more than 150 items")
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
		errs = Seterror(&errs, fmt.Sprintf("max parameter cannot exceed %d", MAXPAGE))
		max = MAXPAGE
	}
	if page < 0 {
		page = 0
	}

	sr.Max = max

	if b["type"] != nil {
		t := b["type"].(string)
		if t != fdc.PHRASE && t != fdc.WILDCARD && t != fdc.REGEX {
			errs = Seterror(&errs, fmt.Sprintf("Search type must be %s, %s, or %s ", fdc.PHRASE, fdc.WILDCARD, fdc.REGEX))
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
