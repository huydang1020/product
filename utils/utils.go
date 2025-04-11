package utils

import (
	"github.com/rs/xid"
)

func MakePartnerId() string {
	return "part" + xid.New().String()
}
