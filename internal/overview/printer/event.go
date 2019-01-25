package printer

import (
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/internal/view/component"
)

// EventListHandler is a printFunc that lists events.
func EventListHandler(list *corev1.EventList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Kind", "Message", "Reason", "Type",
		"First Seen", "Last Seen")
	table := component.NewTable("Events", cols)

	for _, event := range list.Items {
		row := component.TableRow{}

		objectPath, err := ObjectReferencePath(event.InvolvedObject)
		if err != nil {
			return nil, err
		}

		infoItems := []component.ViewComponent{
			component.NewLink("", event.InvolvedObject.Name, objectPath),
			component.NewText("", fmt.Sprintf("%d", event.Count)),
		}
		info := component.NewList("", infoItems)

		row["Kind"] = info
		row["Message"] = component.NewText("", event.Message)
		row["Reason"] = component.NewText("", event.Reason)
		row["Type"] = component.NewText("", event.Type)
		row["First Seen"] = component.NewTimestamp(event.FirstTimestamp.Time)
		row["Last Seen"] = component.NewTimestamp(event.LastTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}
