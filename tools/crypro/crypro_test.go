package crypro

import (
	"fmt"
	"testing"
)

func TestEncryptString(t *testing.T) {
	fmt.Println(Md5WithSalt("ms-co!@#$12398mmx", "ZfnXw"))
}
