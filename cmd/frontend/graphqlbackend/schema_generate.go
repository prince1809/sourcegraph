// +build generate

package main

import (
	"log"
	"io/ioutil"
)

func main() {
	out, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		log.Fatal(err)
	}

	out, err = graphqlfile.str
}
