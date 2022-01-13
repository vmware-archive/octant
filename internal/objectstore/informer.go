package objectstore

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
)

type interruptibleInformer struct {
	stopCh   chan struct{}
	informer informers.GenericInformer
	gvr      schema.GroupVersionResource
}

func (i interruptibleInformer) Stop() {
	close(i.stopCh)
}
