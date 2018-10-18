package telemetry

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/heptio/go-telemetry/pkg/logging"
	pb "github.com/heptio/go-telemetry/pkg/proto/telemetry"
)

const (
	// DefaultAddress is the default address of the telemetry service
	DefaultAddress = "telemetry.heptio.com:80"
)

// Interface defines telemetry
type Interface interface {
	SendEvent(string, Measurements) error
	With(Labels) Interface
	Close()
}

// Labels hold metadata about a telemetry event
type Labels map[string]string

// Measurements hold metrics describing a telemetry event
type Measurements map[string]int64

// NilClient enables a noop telemetry interface
type NilClient struct{}

// Client is a telemetry client
type Client struct {
	address    string
	connection *grpc.ClientConn
	timeout    time.Duration
	closed     bool
	pbClient   pb.TelemetryClient
	logger     logging.Logger
	labels     Labels
}

// SendEvent for NilClient is a noop
func (*NilClient) SendEvent(string, Measurements) error { return nil }

// With for NilClient returns itself
func (n *NilClient) With(Labels) Interface { return n }

// Close for NilClient is a noop
func (*NilClient) Close() {}

// Ensure clients implements Interface at compile time
var _ Interface = (*Client)(nil)
var _ Interface = (*NilClient)(nil)

// NewClient constructs a new teller client
func NewClient(address string, timeout time.Duration, logger interface{}) (*Client, error) {
	adaptedLogger, err := logging.Adapt(logger)
	if err != nil {
		return nil, errors.Wrap(err, "adapting to logger")
	}
	adaptedLogger.Debugf("Creating new telemetry client")

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to telemetry server")
	}

	return &Client{
		address:    address,
		connection: conn,
		timeout:    timeout,
		closed:     false,
		pbClient:   pb.NewTelemetryClient(conn),
		logger:     adaptedLogger,
	}, nil
}

// Close shuts down the underlying connection
func (t *Client) Close() {
	if !t.closed {
		t.connection.Close()
	}
	t.closed = true
	t.logger.Debugf("Telemetry connection closed")
}

// With creates a new client with labels which are additive on top of existing labels
func (t *Client) With(labels Labels) Interface {
	newLabels := Labels{}

	for k, v := range t.labels {
		newLabels[k] = v
	}

	for k, v := range labels {
		newLabels[k] = v
	}

	c := *t
	c.labels = newLabels
	return &c
}

// SendEvent reports an event to telemetry server
func (t *Client) SendEvent(name string, measurements Measurements) error {
	err := t.sendEvents([]pb.Event{{
		Name:         name,
		Timestamp:    ptypes.TimestampNow(),
		Labels:       t.labels,
		Measurements: measurements,
	}})

	if err != nil {
		t.logger.Debugf("failure to send event: %s", err)
		return err
	}

	return nil
}

// sendEvents allows sending a batch of telemetry events
func (t *Client) sendEvents(events []pb.Event) error {
	if t.closed {
		return errors.New("used closed/unready client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)

	// cancelling context cleans up grpc on early return from error
	defer cancel()

	stream, err := t.pbClient.Record(ctx)
	if err != nil {
		return errors.New("create stream")
	}

	for _, event := range events {
		if err := stream.Send(&event); err != nil {
			return errors.Wrap(err, "streaming event")
		}
	}

	if _, err := stream.CloseAndRecv(); err != nil {
		return errors.New("failed ack")
	}

	return nil
}
