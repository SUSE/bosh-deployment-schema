package main

import (
	"fmt"
	"io/ioutil"
	"os"

	wraperr "github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// func init() {
//         loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
// }

func loadSpec(fpath string) (*loads.Document, error) {
	document, err := loads.Spec(fpath)
	if err != nil {
		return nil, wraperr.Wrap(err, "Failed to load spec")
	}

	document, err = document.Expanded(&spec.ExpandOptions{RelativeBase: fpath})
	if err != nil {
		return nil, wraperr.Wrap(err, "Failed to expand spec")
	}

	if err := validate.Spec(document, strfmt.Default); err != nil {
		return nil, wraperr.Wrap(err, "Spec is invalid")
	}

	return document, nil
}

func loadYAML(fpath string) (interface{}, error) {
	var data interface{}

	manifestContents, err := ioutil.ReadFile(fpath)
	if err != nil {
		return data, err
	}

	if err := yaml.Unmarshal(manifestContents, data); err != nil {
		return data, err
	}

	return data, nil
}

func main() {
	fpath := os.Args[1]

	document, err := loadSpec(fpath)
	if err != nil {
		panic(err)
	}

	data, err := loadYAML(os.Args[2])
	if err != nil {
		panic(err)
	}

	sch := document.Schema()
	err = validate.AgainstSchema(sch, data, strfmt.Default)
	ve, ok := err.(*errors.CompositeError)
	if ok {
		fmt.Printf("validation ok\n")
		os.Exit(0)
	}
	fmt.Printf("validation failed\n")
	fmt.Printf("%#v\n", ve)
	os.Exit(1)
}
