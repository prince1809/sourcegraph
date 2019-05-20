// +build generate

package main

import (
	"github.com/prometheus/common/log"
	"github.com/shurcooL/vfsgen"
	"net/http"
)

func main() {
	dir := "../../../../../ui/assets"
	err := vfsgen.Generate(http.Dir(dir), vfsgen.Options{
		PackageName:  "assets",
		BuildTags:    "dist",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
