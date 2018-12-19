package overview

import (
	"context"
	"sort"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type ConfigMapSummary struct{}

var _ View = (*ConfigMapSummary)(nil)

func NewConfigMapSummary(prefix, namespace string, c clock.Clock) View {
	return &ConfigMapSummary{}
}

func (cms *ConfigMapSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	configMap, err := retrieveConfigMap(object)
	if err != nil {
		return nil, err
	}

	detail, err := printConfigMapSummary(configMap)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

// ConfigMapDetails describe the details of a kubernetes corev1.ConfigMap
type ConfigMapDetails struct{}

// NewConfigMapDetails constructs a new ConfigMapDetails object
func NewConfigMapDetails(prefix, namespace string, c clock.Clock) View {
	return &ConfigMapDetails{}
}

// Content describes human readable object details
func (cm *ConfigMapDetails) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	configMap, ok := object.(*corev1.ConfigMap)
	if !ok {
		return nil, errors.Errorf("expected object to be a ConfigMap, it was %T", object)
	}

	emptyMessage := "ConfigMap does not contain any data"
	table := content.NewTable("ConfigMap Data", emptyMessage)
	table.Columns = []content.TableColumn{
		tableCol("Key"),
		tableCol("Value"),
	}

	var keys []string
	for k := range configMap.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := configMap.Data[k]
		row := content.TableRow{
			"Key":   content.NewStringText(k),
			"Value": content.NewStringText(v),
		}
		table.AddRow(row)
	}

	return []content.Content{&table}, nil
}

func retrieveConfigMap(object runtime.Object) (*corev1.ConfigMap, error) {
	rc, ok := object.(*corev1.ConfigMap)
	if !ok {
		return nil, errors.Errorf("expected object to be a ConfigMap, it was %T", object)
	}

	return rc, nil
}
