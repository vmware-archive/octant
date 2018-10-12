package printers

import (
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/kubernetes/pkg/printers"
)

type Handler struct{}

var _ printers.PrintHandler = (*Handler)(nil)

func (h *Handler) Handler(columns []string, columnsWithWide []string, printFunc interface{}) error {
	panic("not implemented")
}

func (h *Handler) TableHandler(columns []metav1beta1.TableColumnDefinition, printFunc interface{}) error {
	panic("not implemented")
}

func (h *Handler) DefaultTableHandler(columns []metav1beta1.TableColumnDefinition, printFunc interface{}) error {
	panic("not implemented")
}
