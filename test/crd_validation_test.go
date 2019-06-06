package test

import (
	"github.com/RHsyseng/operator-utils/pkg/validation"
	atlasmapv1alpha1 "github.com/atlasmap/atlasmap-operator/pkg/apis/atlasmap/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var crdTypeMap = map[string]interface{}{
	"atlasmap_v1alpha1_atlasmap_crd.yaml": &atlasmapv1alpha1.AtlasMap{},
}

func TestCRDSchemas(t *testing.T) {
	for crdFileName, amqType := range crdTypeMap {
		schema := getSchema(t, crdFileName)
		missingEntries := schema.GetMissingEntries(amqType)
		for _, missing := range missingEntries {
			assert.Fail(t, "Discrepancy between CRD and Struct",
				"Missing or incorrect schema validation at %v, expected type %v  in CRD file %v", missing.Path, missing.Type, crdFileName)
		}
	}
}

func TestSampleCustomResources(t *testing.T) {

	var crFileName, crdFileName string = "atlasmap_v1alpha1_atlasmap_cr.yaml", "atlasmap_v1alpha1_atlasmap_crd.yaml"
	assert.NotEmpty(t, crdFileName, "No matching CRD file found for CR suffixed: %s", crFileName)

	schema := getSchema(t, crdFileName)
	yamlString, err := ioutil.ReadFile("../deploy/crds/" + crFileName)
	assert.NoError(t, err, "Error reading %v CR yaml", crFileName)
	var input map[string]interface{}
	assert.NoError(t, yaml.Unmarshal([]byte(yamlString), &input))
	assert.NoError(t, schema.Validate(input), "File %v does not validate against the CRD schema", crFileName)
}

func getSchema(t *testing.T, crdFile string) validation.Schema {

	yamlString, err := ioutil.ReadFile("../deploy/crds/" + crdFile)
	assert.NoError(t, err, "Error reading CRD yaml %v", yamlString)

	schema, err := validation.New([]byte(yamlString))
	assert.NoError(t, err)

	return schema
}
