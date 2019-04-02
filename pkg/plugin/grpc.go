package plugin

import (
	"context"
	"encoding/json"

	"github.com/heptio/developer-dash/pkg/plugin/proto"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// GRPCClient is the dashboard GRPC client.
type GRPCClient struct {
	broker Broker
	client proto.PluginClient
}

var _ Service = (*GRPCClient)(nil)

// NewGRPCClient creates an instance of GRPCClient.
func NewGRPCClient(broker Broker, client proto.PluginClient) *GRPCClient {
	return &GRPCClient{
		client: client,
		broker: broker,
	}
}

func (c *GRPCClient) run(fn func() error) error {
	if fn == nil {
		return errors.New("client function is nil")
	}

	var s *grpc.Server
	defer func() {
		if s != nil {
			s.Stop()
		}
	}()

	serverFunc := func(options []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(options...)
		return s
	}

	brokerID := c.broker.NextId()
	go c.broker.AcceptAndServe(brokerID, serverFunc)

	return fn()
}

// Register register a plugin.
func (c *GRPCClient) Register(dashboardAPIAddress string) (Metadata, error) {
	var m Metadata

	err := c.run(func() error {
		registerRequest := &proto.RegisterRequest{
			DashboardAPIAddress: dashboardAPIAddress,
		}

		resp, err := c.client.Register(context.Background(), registerRequest)
		if err != nil {
			return errors.Wrap(err, "grpc client register")
		}

		capabilities := convertToCapabilities(resp.Capabilities)

		m = Metadata{
			Name:         resp.PluginName,
			Description:  resp.Description,
			Capabilities: capabilities,
		}

		return nil
	})

	if err != nil {
		return Metadata{}, err
	}

	return m, nil
}

// Print prints an object.
func (c *GRPCClient) Print(object runtime.Object) (PrintResponse, error) {
	var pr PrintResponse

	err := c.run(func() error {
		in, err := createObjectRequest(object)
		if err != nil {
			return err
		}

		resp, err := c.client.Print(context.Background(), in)
		if err != nil {
			return errors.Wrap(err, "grpc client print")
		}

		var items []component.FlexLayoutItem
		if len(resp.Items) > 0 {
			if err := json.Unmarshal(resp.Items, &items); err != nil {
				return err
			}
		}

		config, err := convertToSummarySections(resp.Config)
		if err != nil {
			return errors.Wrap(err, "convert config sections")
		}

		status, err := convertToSummarySections(resp.Status)
		if err != nil {
			return errors.Wrap(err, "convert status sections")
		}

		pr = PrintResponse{
			Config: config,
			Status: status,
			Items:  items,
		}

		return nil
	})

	if err != nil {
		return PrintResponse{}, err
	}

	return pr, nil
}

func createObjectRequest(object runtime.Object) (*proto.ObjectRequest, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	or := &proto.ObjectRequest{
		Object: data,
	}

	return or, err
}

// PrintTab creates a tab for an object.
func (c *GRPCClient) PrintTab(object runtime.Object) (*component.Tab, error) {
	var tab component.Tab

	err := c.run(func() error {
		in, err := createObjectRequest(object)
		if err != nil {
			return err
		}

		resp, err := c.client.PrintTab(context.Background(), in)
		if err != nil {
			return errors.Wrap(err, "grpc client print tab")
		}

		var to component.TypedObject
		if err := json.Unmarshal(resp.Layout, &to); err != nil {
			return err
		}

		c, err := to.ToComponent()
		if err != nil {
			return err
		}

		layout, ok := c.(*component.FlexLayout)
		if !ok {
			return errors.Errorf("expected to be flexlayout was: %T", c)
		}

		tab.Name = resp.Name
		tab.Contents = *layout

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &tab, nil
}

// GRPCServer is the grpc server the dashboard will use to communicate with the
// the plugin.
type GRPCServer struct {
	Impl   Service
	broker Broker
}

var _ proto.PluginServer = (*GRPCServer)(nil)

// Register register a plugin.
func (s *GRPCServer) Register(ctx context.Context, registerRequest *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	m, err := s.Impl.Register(registerRequest.DashboardAPIAddress)
	if err != nil {
		return nil, err
	}

	capabilities := convertFromCapabilities(m.Capabilities)

	return &proto.RegisterResponse{
		PluginName:   m.Name,
		Description:  m.Description,
		Capabilities: &capabilities,
	}, nil
}

// Print prints an object.
func (s *GRPCServer) Print(ctx context.Context, objectRequest *proto.ObjectRequest) (*proto.PrintResponse, error) {
	u, err := decodeObjectRequest(objectRequest)
	if err != nil {
		return nil, err
	}

	pr, err := s.Impl.Print(u)
	if err != nil {
		return nil, errors.Wrap(err, "grpc server print")
	}

	itemBytes, err := json.Marshal(pr.Items)
	if err != nil {
		return nil, err
	}

	config, err := convertFromSummarySections(pr.Config)
	if err != nil {
		return nil, err
	}

	status, err := convertFromSummarySections(pr.Status)
	if err != nil {
		return nil, err
	}

	out := &proto.PrintResponse{
		Config: config,
		Status: status,
		Items:  itemBytes,
	}

	return out, nil
}

// ObjectStatus generates status for an object.
func (s *GRPCServer) ObjectStatus(context.Context, *proto.ObjectRequest) (*proto.ObjectStatusResponse, error) {
	panic("not implemented")
}

func decodeObjectRequest(req *proto.ObjectRequest) (*unstructured.Unstructured, error) {
	m := map[string]interface{}{}

	err := json.Unmarshal(req.Object, &m)
	if err != nil {
		return nil, err
	}

	u := &unstructured.Unstructured{Object: m}
	return u, nil
}

// PrintTab prints a tab for an object.
func (s *GRPCServer) PrintTab(ctx context.Context, objectRequest *proto.ObjectRequest) (*proto.PrintTabResponse, error) {
	u, err := decodeObjectRequest(objectRequest)
	if err != nil {
		return nil, err
	}

	tab, err := s.Impl.PrintTab(u)
	if err != nil {
		return nil, errors.Wrap(err, "grpc server print tab")
	}

	layoutBytes, err := json.Marshal(tab.Contents)
	if err != nil {
		return nil, err
	}

	out := &proto.PrintTabResponse{
		Name:   tab.Name,
		Layout: layoutBytes,
	}

	return out, nil
}

// WatchAdd is called when a watched GVK has a new object added.
func (s *GRPCServer) WatchAdd(context.Context, *proto.WatchRequest) (*proto.Empty, error) {
	panic("not implemented")
}

// WatchUpdate is caleld when a watched GVK has an object updated.
func (s *GRPCServer) WatchUpdate(context.Context, *proto.WatchRequest) (*proto.Empty, error) {
	panic("not implemented")
}

// WatchDelete is called when a watched GVK has an object deleted.
func (s *GRPCServer) WatchDelete(context.Context, *proto.WatchRequest) (*proto.Empty, error) {
	panic("not implemented")
}
