package utils

import (
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
