// +build generate
package main

import "io/ioutil"

func main()  {
	out, err := ioutil.ReadFile("schema.graphql")
}
