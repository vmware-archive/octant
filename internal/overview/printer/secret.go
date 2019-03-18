package printer

import (
	"context"
	"fmt"

	"github.com/heptio/developer-dash/internal/overview/link"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

var (
	secretTableCols = component.NewTableCols("Name", "Labels", "Type", "Data", "Age")
	secretDataCols  = component.NewTableCols("Key")
)

// SecretListHandler is a printFunc that lists secrets.
func SecretListHandler(ctx context.Context, list *corev1.SecretList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("list of secrets is nil")
	}

	table := component.NewTable("Secrets", secretTableCols)

	for _, secret := range list.Items {
		row := component.TableRow{}

		row["Name"] = link.ForObject(&secret, secret.Name)
		row["Labels"] = component.NewLabels(secret.ObjectMeta.Labels)
		row["Type"] = component.NewText(string(secret.Type))
		row["Data"] = component.NewText(fmt.Sprintf("%d", len(secret.Data)))
		row["Age"] = component.NewTimestamp(secret.ObjectMeta.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

// SecretHandler is a printFunc for printing a secret summary.
func SecretHandler(ctx context.Context, secret *corev1.Secret, options Options) (component.ViewComponent, error) {
	o := NewObject(secret)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return secretConfiguration(*secret)
	}, 16)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return secretData(*secret)
		},
		Width: component.WidthFull,
	})

	return o.ToComponent(ctx, options)
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
