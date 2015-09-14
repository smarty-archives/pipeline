package persist

import (
	"io/ioutil"
	"log"
)

//go:generate go install github.com/smartystreets/gunit/gunit
//go:generate gunit

func init() {
	log.SetOutput(ioutil.Discard)
}
