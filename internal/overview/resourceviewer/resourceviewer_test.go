package resourceviewer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	"github.com/heptio/developer-dash/internal/overview/objectvisitor"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type stubbedVisitor struct{ visitErr error }

var _ objectvisitor.Visitor = (*stubbedVisitor)(nil)

func (v stubbedVisitor) Visit(context.Context, objectvisitor.ClusterObject) error {
	return v.visitErr
}

func stubVisitor(fail bool) ViewerOpt {
	return func(rv *ResourceViewer) error {
		sv := &stubbedVisitor{}
		if fail {
			sv.visitErr = errors.Errorf("fail")
		}

		rv.visitor = sv
		return nil
	}
}

func Test_ResourceViewer(t *testing.T) {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment",
			UID:  types.UID("deployment"),
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()
	c := cachefake.NewMockCache(controller)

	rv, err := New(nil, c, stubVisitor(false))
	require.NoError(t, err)

	ctx := context.Background()

	vc, err := rv.Visit(ctx, deployment)
	require.NoError(t, err)
	assert.NotNil(t, vc)
}

func Test_ResourceViewer_visitor_fails(t *testing.T) {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment",
			UID:  types.UID("deployment"),
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()
	c := cachefake.NewMockCache(controller)

	rv, err := New(nil, c, stubVisitor(true))
	require.NoError(t, err)

	ctx := context.Background()

	_, err = rv.Visit(ctx, deployment)
	require.Error(t, err)
}
