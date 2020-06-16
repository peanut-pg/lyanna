package util

import (
	"os"
	"strings"
)

const (
	LYANNA_ENV  = "lyanna_env"
	PRODUCT_ENV = "product"
	TEST_ENV    = "test"
)

var (
	curLyannaEnv string = TEST_ENV
)

func init() {
	curLyannaEnv = strings.ToLower(os.Getenv(LYANNA_ENV))
	curLyannaEnv = strings.TrimSpace(curLyannaEnv)

	if len(curLyannaEnv) == 0 {
		curLyannaEnv = TEST_ENV
	}
}

func IsProduct() bool {
	return curLyannaEnv == PRODUCT_ENV
}

func IsTest() bool {
	return curLyannaEnv == TEST_ENV
}

func GetEnv() string {
	return curLyannaEnv
}
