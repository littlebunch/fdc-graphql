package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/littlebunch/fdc-api/ds/cb"
	fdc "github.com/littlebunch/fdc-api/model"
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
		if i > MAXIDS {
			err = fmt.Sprintf("number of fdcId's should not exceed %d", MAXIDS)
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
