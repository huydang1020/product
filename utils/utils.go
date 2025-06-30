package utils

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/rs/xid"
	"golang.org/x/text/unicode/norm"
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

func MakeOrderShipId() string {
	return "osh" + xid.New().String()
}

func MakeOrderDetailId() string {
	return "odt" + xid.New().String()
}

func MakeReviewId() string {
	return "rev" + xid.New().String()
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

func Include(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Chuyển tiếng Việt có dấu sang không dấu và thành slug
func ToSlug(input string) string {
	// Normalize để tách dấu ra
	t := norm.NFD.String(input)
	slug := strings.Builder{}
	for _, r := range t {
		switch {
		case unicode.Is(unicode.Mn, r):
			continue // bỏ dấu
		case r == 'đ':
			slug.WriteRune('d')
		case r == 'Đ':
			slug.WriteRune('d')
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			slug.WriteRune(unicode.ToLower(r))
		default:
			slug.WriteRune(' ')
		}
	}
	// Thay nhiều dấu cách thành dấu gạch ngang
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(strings.TrimSpace(slug.String()), "-")
}

func SendReqPost(url string, headers map[string]string, body []byte) (int, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}
	if headers != nil {
		for k, val := range headers {
			req.Header.Set(k, val)
		}
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		req.Close = true
		resp.Body.Close()
	}()
	body, _ = io.ReadAll(resp.Body)
	return resp.StatusCode, body, nil
}
