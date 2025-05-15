package utils

import (
	"log"
	"net/url"
	"sort"
	"time"

	"github.com/rs/xid"
)

func MakeCategoryId() string {
	return "cat" + xid.New().String()
}

func MakeProductId() string {
	return "pro" + xid.New().String()
}

func MakeProductTypeId() string {
	return "pty" + xid.New().String()
}

func MakeBannerId() string {
	return "ban" + xid.New().String()
}

func MakeOrderId() string {
	return "ord" + xid.New().String()
}

func ConvertUnixToDateTime(format string, t int64) (string, error) {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Println("load location err:", err)
		return "", err
	}
	formattedDate := time.Unix(t, 0).In(location).Format(format)
	return formattedDate, nil
}

func SortParams(params url.Values) url.Values {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedParams := make(url.Values)
	for _, k := range keys {
		sortedParams[k] = params[k]
	}

	return sortedParams
}
