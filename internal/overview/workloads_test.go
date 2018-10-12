package overview

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
)

func Test_printCronJobTable(t *testing.T) {
	tbl := &metav1beta1.Table{
		ColumnDefinitions: []metav1beta1.TableColumnDefinition{
			{
				Name:        "Name",
				Type:        "string",
				Format:      "name",
				Description: "Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
				Priority:    0,
			},
			{
				Name:        "Schedule",
				Type:        "string",
				Format:      "",
				Description: "The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.",
				Priority:    0,
			},
			{
				Name:        "Suspend",
				Type:        "boolean",
				Format:      "",
				Description: "This flag tells the controller to suspend subsequent executions, it does not apply to already started executions.  Defaults to false.",
				Priority:    0,
			},
			{
				Name:        "Active",
				Type:        "integer",
				Format:      "",
				Description: "A list of pointers to currently running jobs.",
				Priority:    0,
			},
			{
				Name:        "Last Schedule",
				Type:        "string",
				Format:      "",
				Description: "Information when was the last time the job was successfully scheduled.",
				Priority:    0,
			},
			{
				Name:        "Age",
				Type:        "string",
				Format:      "",
				Description: "CreationTimestamp is a timestamp representing the server time when this object was created. Itis not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.\n\nPopulated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata",
				Priority:    0,
			},
			{
				Name:        "Containers",
				Type:        "string",
				Format:      "",
				Description: "Names of each container in the template.",
				Priority:    1,
			},
			{
				Name:        "Images",
				Type:        "string",
				Format:      "",
				Description: "Images referenced by each container in the template.",
				Priority:    1,
			},
			{
				Name:        "Selector",
				Type:        "string",
				Format:      "",
				Description: "A label query over pods that should match the pod count. Normally, the system sets this field for you. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors",
				Priority:    1,
			},
			{
				Name:        "Labels",
				Type:        "string",
				Format:      "",
				Description: "",
				Priority:    0,
			},
		},
		Rows: []metav1beta1.TableRow{
			{
				Cells: []interface{}{
					"hello",
					"*/1 * * * *",
					"False",
					int64(0),
					"30s",
					"<unknown>",
					"hello",
					"busybox",
					"<none>",
					"<none>",
				},
			},
		},
	}

	contentTable, err := printContentTable("Title", "default", "/prefix", tbl, cronJobTransforms)
	require.NoError(t, err)

	expected := newTable("Title")
	expected.Columns = []tableColumn{
		newCol("Name"),
		newCol("Schedule"),
		newCol("Suspend"),
		newCol("Active"),
		newCol("Last Schedule"),
		newCol("Age"),
		newCol("Containers"),
		newCol("Images"),
		newCol("Selector"),
		newCol("Labels"),
	}

	expected.AddRow(tableRow{
		"Name":          newLinkText("hello", "/prefix/workloads/cron-jobs/hello?namespace=default"),
		"Schedule":      newStringText("*/1 * * * *"),
		"Suspend":       newStringText("False"),
		"Active":        newStringText("0"),
		"Last Schedule": newStringText("30s"),
		"Age":           newStringText("<unknown>"),
		"Containers":    newStringText("hello"),
		"Images":        newStringText("busybox"),
		"Selector":      newStringText("<none>"),
		"Labels":        newStringText("<none>"),
	})

	assert.Equal(t, expected, *contentTable)
}

func Test_printDeploymentTable(t *testing.T) {
	tbl := &metav1beta1.Table{
		ColumnDefinitions: []metav1beta1.TableColumnDefinition{
			{
				Name:        "Name",
				Type:        "string",
				Format:      "name",
				Description: "Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
				Priority:    0,
			},
			{
				Name:        "Desired",
				Type:        "string",
				Format:      "",
				Description: "Number of desired pods. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1.",
				Priority:    0,
			},
			{
				Name:        "Current",
				Type:        "string",
				Format:      "",
				Description: "Total number of non-terminated pods targeted by this deployment (their labels match the selector).",
				Priority:    0,
			},
			{
				Name:        "Up-to-date",
				Type:        "string",
				Format:      "",
				Description: "Total number of non-terminated pods targeted by this deployment that have the desired templatespec.",
				Priority:    0,
			},
			{
				Name:        "Available",
				Type:        "string",
				Format:      "",
				Description: "Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.",
				Priority:    0,
			},
			{
				Name:        "Age",
				Type:        "string",
				Format:      "",
				Description: "CreationTimestamp is a timestamp representing the server time when this object was created. Itis not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.\n\nPopulated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata",
				Priority:    0,
			},
			{
				Name:        "Containers",
				Type:        "string",
				Format:      "",
				Description: "Names of each container in the template.",
				Priority:    1,
			},
			{
				Name:        "Images",
				Type:        "string",
				Format:      "",
				Description: "Images referenced by each container in the template.",
				Priority:    1,
			},
			{
				Name:        "Selector",
				Type:        "string",
				Format:      "",
				Description: "Label selector for pods. Existing ReplicaSets whose pods are selected by this will be the onesaffected by this deployment.",
				Priority:    1,
			},
			{
				Name:        "Labels",
				Type:        "string",
				Format:      "",
				Description: "",
				Priority:    0,
			},
		},
		Rows: []metav1beta1.TableRow{
			{
				Cells: []interface{}{
					"nginx-deployment",
					3,
					3,
					3,
					3,
					"<unknown>",
					"nginx",
					"nginx:1.7.9",
					"app=nginx",
					"<none>",
				},
			},
			{
				Cells: []interface{}{
					"krex-debug-pod",
					1,
					1,
					1,
					1,
					"<unknown>",
					"krex-debug-pod",
					"ubuntu:latest",
					"run=krex-debug-pod",
					"<none>",
				},
			},
		},
	}

	contentTable, err := printContentTable("Title", "default", "/prefix", tbl, deploymentTransforms)
	require.NoError(t, err)

	expected := newTable("Title")
	expected.Columns = []tableColumn{
		newCol("Name"),
		newCol("Desired"),
		newCol("Current"),
		newCol("Up-to-date"),
		newCol("Available"),
		newCol("Age"),
		newCol("Containers"),
		newCol("Images"),
		newCol("Selector"),
		newCol("Labels"),
	}

	expected.AddRow(tableRow{
		"Name":       newLinkText("nginx-deployment", "/prefix/workloads/deployments/nginx-deployment?namespace=default"),
		"Desired":    newStringText("3"),
		"Current":    newStringText("3"),
		"Up-to-date": newStringText("3"),
		"Available":  newStringText("3"),
		"Age":        newStringText("<unknown>"),
		"Containers": newStringText("nginx"),
		"Images":     newStringText("nginx:1.7.9"),
		"Selector":   newStringText("app=nginx"),
		"Labels":     newStringText("<none>"),
	})

	expected.AddRow(tableRow{
		"Name":       newLinkText("krex-debug-pod", "/prefix/workloads/deployments/krex-debug-pod?namespace=default"),
		"Desired":    newStringText("1"),
		"Current":    newStringText("1"),
		"Up-to-date": newStringText("1"),
		"Available":  newStringText("1"),
		"Age":        newStringText("<unknown>"),
		"Containers": newStringText("krex-debug-pod"),
		"Images":     newStringText("ubuntu:latest"),
		"Selector":   newStringText("run=krex-debug-pod"),
		"Labels":     newStringText("<none>"),
	})

	assert.Equal(t, expected, *contentTable)
}

func newCol(name string) tableColumn {
	return tableColumn{
		Name:     name,
		Accessor: name,
	}
}
