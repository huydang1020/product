package utils

import (
	"log"
	"testing"
)

func Test_slug(t *testing.T) {
	key := "Xịt/lăn khử mùi Rexona 72H kháng khuẩn thể thao dành cho nam 135/45ml"
	a := ToSlug(key)
	log.Println("a: ", a)
}
