package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_createConditionsTableErrs(t *testing.T) {
	// Nil object
	table, ok, err := createConditionsTable(nil, "", nil)
	assert.EqualError(t, err, "object is nil")
	assert.False(t, ok)
	assert.Nil(t, table)

	// No status found
	noStatus := &unstructured.Unstructured{Object: map[string]interface{}{"noStatus": nil}}
	table, ok, err = createConditionsTable(noStatus, "", nil)
	assert.Nil(t, err)
	assert.False(t, ok)
	assert.NotNil(t, table)

	// Bad status, not a map[string]interface{}
	badStatus := &unstructured.Unstructured{Object: map[string]interface{}{"status": 1}}
	table, ok, err = createConditionsTable(badStatus, "", nil)
	assert.EqualError(t, err, ".status accessor error: 1 is of the type int, expected map[string]interface{}")
	assert.False(t, ok)
	assert.Nil(t, table)

	// No conditions found
	noConditions := &unstructured.Unstructured{Object: map[string]interface{}{"status": map[string]interface{}{"noConditions": nil}}}
	table, ok, err = createConditionsTable(noConditions, "", nil)
	assert.Nil(t, err)
	assert.False(t, ok)
	assert.NotNil(t, table)
}
