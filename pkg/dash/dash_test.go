/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

package dash

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/event"

	pkglog "github.com/vmware-tanzu/octant/pkg/log"
)

func TestRunner_ValidateKubeconfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "/test1", []byte(""), 0755)
	afero.WriteFile(fs, "/test2", []byte(""), 0755)

	separator := string(filepath.ListSeparator)

	tests := []struct {
		name     string
		fileList string
		expected string
		isErr    bool
	}{
		{
			name:     "single path",
			fileList: "/test1",
			expected: "/test1",
			isErr:    false,
		},
		{
			name:     "multiple paths",
			fileList: "/test1" + separator + "/test2",
			expected: "/test1" + separator + "/test2",
			isErr:    false,
		},
		{
			name:     "single path not found",
			fileList: "/unknown",
			expected: "",
			isErr:    true,
		},
		{
			name:     "multiple paths not found",
			fileList: "/unknown" + separator + "/unknown2",
			expected: "",
			isErr:    true,
		},
		{
			name:     "multiple file path; missing a config",
			fileList: "/test1" + separator + "/unknown",
			expected: "/test1",
			isErr:    false,
		},
		{
			name:     "invalid path",
			fileList: "not a filepath",
			expected: "",
			isErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := log.NopLogger()
			path, err := ValidateKubeConfig(logger, test.fileList, fs)
			if test.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, path, test.expected)
		})
	}
}

type inMemoryListener struct {
	conns chan net.Conn
}

func NewInMemoryListener() *inMemoryListener {
	return &inMemoryListener{conns: make(chan net.Conn)}
}

func (iml *inMemoryListener) Accept() (net.Conn, error) {
	return <-iml.conns, nil
}
func (iml *inMemoryListener) Close() error {
	return nil
}
func (iml *inMemoryListener) Dial(network, addr string) (net.Conn, error) {
	server, client := net.Pipe()
	iml.conns <- AddrOverriddingConn{server}
	return AddrOverriddingConn{client}, nil
}
func (iml *inMemoryListener) Addr() net.Addr {
	return localhostAddr{}
}

type AddrOverriddingConn struct {
	net.Conn
}

func (AddrOverriddingConn) LocalAddr() net.Addr {
	return localhostAddr{}
}
func (AddrOverriddingConn) RemoteAddr() net.Addr {
	return localhostAddr{}
}

type localhostAddr struct{}

func (localhostAddr) Network() string {
	return "tcp"
}
func (localhostAddr) String() string {
	return "127.0.0.1:7777"
}

func TestNewRunnerLoadsValidKubeConfigFilteringNonexistent(t *testing.T) {
	srv := fakeK8sAPIThatForbidsWatchingCRDs()
	defer srv.Close()
	stubRiceBox("dist/dash-frontend")
	kubeConfig := tempFile(makeKubeConfig("test-context", srv.URL))
	defer os.Remove(kubeConfig.Name())

	listener := NewInMemoryListener()
	cancel, _ := makeRunner(
		Options{
			KubeConfig: strings.Join(
				[]string{
					"/non/existent/kubeconfig",
					kubeConfig.Name(),
				},
				string(filepath.ListSeparator),
			),
			Listener: listener,
		},
		log.NopLogger(),
	)
	defer cancel()
	kubeConfigEvent := waitForKubeConfigEvent(listener)

	require.Equal(t, "test-context", kubeConfigEvent.Data.CurrentContext)
}

func TestNewRunnerRunsLoadingAPIWhenStartedWithoutKubeConfig(t *testing.T) {
	srv := fakeK8sAPIThatForbidsWatchingCRDs()
	defer srv.Close()
	stubRiceBox("dist/dash-frontend")

	listener := NewInMemoryListener()
	cancel, _ := makeRunner(Options{Listener: listener}, log.NopLogger())
	defer cancel()
	kubeConfig := makeKubeConfig("test-context", srv.URL)
	websocketWrite(
		fmt.Sprintf(`{
	"type": "action.octant.dev/uploadKubeConfig",
	"payload": {"kubeConfig": "%s"}
}`, base64.StdEncoding.EncodeToString(kubeConfig)),
		listener,
	)
	// wait for API to reload
	for {
		if websocketReadTimeout(listener, 10*time.Millisecond) {
			break
		}
	}
	kubeConfigEvent := waitForKubeConfigEvent(listener)

	require.Equal(t, "test-context", kubeConfigEvent.Data.CurrentContext)
}

func TestNewRunnerShutsDownPluginsWhenStoppedBeforeReceivingKubeConfig(t *testing.T) {
	stubRiceBox("dist/dash-frontend")
	listener := NewInMemoryListener()
	shutdownCh := make(chan bool)
	options := Options{
		Listener: listener,
	}
	logger := log.NopLogger()
	ctx, cancel := context.WithCancel(context.Background())
	runner, err := NewRunner(ctx, logger, options)
	require.NoError(t, err)

	go runner.Start(options, make(chan bool), shutdownCh)
	cancel()

	select {
	case <-time.After(10 * time.Millisecond):
		require.Fail(t, "failed to shut down within 10 ms")
	case <-shutdownCh:
	}
}

func websocketWrite(message string, listener *inMemoryListener) error {
	dialer := websocket.DefaultDialer
	dialer.NetDial = listener.Dial
	wsConn, _, err := dialer.Dial("ws://127.0.0.1:7777/api/v1/stream", nil)
	if err != nil {
		return err
	}
	w, err := wsConn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}
	w.Close()
	wsConn.Close()
	return nil
}

func websocketReadTimeout(listener *inMemoryListener, timeout time.Duration) bool {
	dialer := websocket.DefaultDialer
	dialer.NetDial = listener.Dial
	wsConn, _, _ := dialer.Dial("ws://127.0.0.1:7777/api/v1/stream", nil)
	defer wsConn.Close()
	reader := make(chan interface{}, 1)
	go func() {
		wsConn.NextReader()
		reader <- nil
	}()
	select {
	case <-reader:
		return true
	case <-time.After(timeout):
		return false
	}
}

func makeKubeConfig(currentContext, serverAddr string) []byte {
	return []byte(fmt.Sprintf(`contexts:
- context: {cluster: cluster}
  name: %s
clusters:
- cluster: {server: %s}
  name: cluster
current-context: %s
`, currentContext, serverAddr, currentContext))
}

func waitForKubeConfigEvent(listener *inMemoryListener) kubeConfigEvent {
	var message kubeConfigEvent
	dialer := websocket.DefaultDialer
	dialer.NetDial = listener.Dial
	wsConn, resp, err := dialer.Dial("ws://127.0.0.1:7777/api/v1/stream", nil)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}
	defer wsConn.Close()
	for {
		msgBytes, _ := readNextMessage(wsConn)
		json.Unmarshal(msgBytes, &message)
		if message.Type == event.EventTypeKubeConfig {
			break
		}
	}
	return message
}

type kubeConfigEvent struct {
	Type event.EventType `json:"type"`
	Data struct {
		CurrentContext string `json:"currentContext"`
	} `json:"data"`
}

func tempFile(contents []byte) *os.File {
	tmpFile, _ := ioutil.TempFile("", "")
	tmpFile.Write(contents)
	tmpFile.Close()
	return tmpFile
}

func makeRunner(options Options, logger pkglog.Logger) (context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	runner, err := NewRunner(ctx, logger, options)
	if err != nil {
		return cancel, err
	}
	go runner.Start(options, make(chan bool), make(chan bool))
	return cancel, nil
}

func fakeK8sAPIThatForbidsWatchingCRDs() *httptest.Server {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	srv := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api":
				w.Write([]byte(fmt.Sprintf(`{
	"kind": "APIVersions",
	"versions": ["v1"],
	"serverAddressByClientCIDRs": [
		{
			"clientCIDR": "0.0.0.0/0",
			"serverAddress": "%s"
		}
	]
}`, l.Addr().String())))
			case "/apis":
				w.Write([]byte(`{
	"kind": "APIGroupList",
	"apiVersion": "v1",
	"groups": [
		{
			"name": "apiextensions.k8s.io",
			"versions": [
				{
					"groupVersion": "apiextensions.k8s.io/v1beta1",
					"version": "v1beta1"
				}
			],
			"preferredVersion": {
				"groupVersion": "apiextensions.k8s.io/v1beta1",
				"version": "v1beta1"
			}
		}
	]
}`))
			case "/apis/apiextensions.k8s.io/v1beta1":
				w.Write([]byte(`{
	"kind": "APIResourceList",
	"apiVersion": "v1",
	"groupVersion": "apiextensions.k8s.io/v1beta1",
	"resources": [
		{
			"name": "customresourcedefinitions",
			"singularName": "",
			"namespaced": false,
			"kind": "CustomResourceDefinition",
			"verbs": [
				"create",
				"delete",
				"deletecollection",
				"get",
				"list",
				"patch",
				"update",
				"watch"
			],
			"shortNames": [
				"crd",
				"crds"
			],
			"storageVersionHash": "jfWCUB31mvA="
		},
		{
			"name": "customresourcedefinitions/status",
			"singularName": "",
			"namespaced": false,
			"kind": "CustomResourceDefinition",
			"verbs": [
				"get",
				"patch",
				"update"
			]
		}
	]
}`))
			case "/apis/authorization.k8s.io/v1/selfsubjectaccessreviews":
				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{
	"kind": "SelfSubjectAccessReview",
	"apiVersion": "authorization.k8s.io/v1",
	"metadata": {
		"creationTimestamp": null
	},
	"spec": {
		"resourceAttributes": {
			"verb": "watch"
		}
	},
	"status": {
		"allowed": false
	}
}`))
			}
		}),
	)
	srv.Listener.Close()
	srv.Listener = l
	srv.Start()
	return srv
}

func stubRiceBox(name string) {
	_, callingGoFile, _, _ := runtime.Caller(0)
	pkgDir := filepath.Dir(callingGoFile)
	os.MkdirAll(filepath.Join(pkgDir, "../../web", name), 0755)
}

func readNextMessage(conn *websocket.Conn) ([]byte, error) {
	_, reader, err := conn.NextReader()
	if err != nil {
		return nil, err
	}
	msgBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}
