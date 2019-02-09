package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/manno/kin-openapi/openapi3"
	yaml "gopkg.in/yaml.v2"
)

func loadSpec(path string) (*openapi3.Swagger, error) {
	fmt.Printf("loading spec from: %s\n", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	loader := openapi3.NewSwaggerLoader()
	doc, err := loader.LoadSwaggerFromYAMLData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec: %s", err)
	}

	fmt.Printf("loaded: %s\n", doc.Info.Title)

	err = doc.Validate(loader.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to validate spec: %s", err)
	}

	return doc, nil
}

func loadYAML(path string) (data []byte, err error) {
	fmt.Printf("loading doc from: %s\n", path)

	data, err = ioutil.ReadFile(path)
	return
}

// kin-openapi does not support map[interface{}]interface{}
// https://stackoverflow.com/questions/40737122/convert-yaml-to-json-without-struct?answertab=votes#tab-top
func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func main() {
	fpath := os.Args[1]

	spec, err := loadSpec(fpath)
	if err != nil {
		panic(err)
	}

	schema := spec.Components.Schemas["DeploymentManifest"]
	if schema == nil {
		for k, _ := range spec.Components.Schemas {
			fmt.Printf("found schema: %s\n", k)
		}
		panic("failed to find BOSH schema: DeploymentManifest")
	}

	data, err := loadYAML(os.Args[2])
	if err != nil {
		panic(fmt.Errorf("failed to load yaml: %s", err))
	}

	var manifest interface{}
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		panic(fmt.Errorf("failed to jsonize manifest: %s", err))
	}

	err = schema.Value.VisitJSON(convert(manifest))
	if v, ok := err.(*openapi3.SchemaError); ok {
		text, _ := yaml.Marshal(v.Value)
		fmt.Println(string(text))
		fmt.Println(v)
		fmt.Printf("failed on field '%s': %s\n", v.SchemaField, v.Reason)
		os.Exit(1)
	}
	fmt.Println("validation: ok")
}
