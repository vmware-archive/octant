package printer

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/gridlayout"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

var (
	secretTableCols = component.NewTableCols("Name", "Labels", "Type", "Data", "Age")
	secretDataCols  = component.NewTableCols("Key")
)

// SecretListHandler is a printFunc that lists secrets.
func SecretListHandler(list *corev1.SecretList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("list of secrets is nil")
	}

	table := component.NewTable("Secrets", secretTableCols)

	for _, secret := range list.Items {
		row := component.TableRow{}
		secretPath, err := gvkPathFromObject(&secret)
		if err != nil {
			return nil, errors.Wrapf(err, "build path for secret %s", secret.Name)
		}
		row["Name"] = component.NewLink("", secret.Name, secretPath)
		row["Labels"] = component.NewLabels(secret.ObjectMeta.Labels)
		row["Type"] = component.NewText(string(secret.Type))
		row["Data"] = component.NewText(fmt.Sprintf("%d", len(secret.Data)))
		row["Age"] = component.NewTimestamp(secret.ObjectMeta.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

// SecretHandler is a printFunc for printing a secret summary.
func SecretHandler(secret *corev1.Secret, options Options) (component.ViewComponent, error) {
	if secret == nil {
		return nil, errors.New("secret is nil")
	}
	gl := gridlayout.New()

	configSection := gl.CreateSection(8)
	configView, err := secretConfiguration(*secret)
	if err != nil {
		return nil, errors.Wrapf(err, "summarize configuration for secret %s", secret.Name)
	}
	configSection.Add(configView, 12)

	dataSection := gl.CreateSection(8)
	dataView, err := secretData(*secret)
	if err != nil {
		return nil, errors.Wrapf(err, "summary data for secret %s", secret.Name)
	}
	dataSection.Add(dataView, 24)

	return gl.ToGrid(), nil
}

func secretConfiguration(secret corev1.Secret) (*component.Summary, error) {
	var sections []component.SummarySection

	sections = append(sections, component.SummarySection{
		Header:  "Type",
		Content: component.NewText(string(secret.Type)),
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func secretData(secret corev1.Secret) (*component.Table, error) {
	table := component.NewTable("Data", secretDataCols)

	for key := range secret.Data {
		row := component.TableRow{}
		row["Key"] = component.NewText(key)

		table.Add(row)
	}

	return table, nil
}
