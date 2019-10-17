package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/pkg/plugin"
)

func TestNewPlugin(t *testing.T) {
	name := "plugin-name"
	description := "description"
	capabilities := &plugin.Capabilities{}

	type args struct {
		name         string
		description  string
		capabilities *plugin.Capabilities
	}

	tests := []struct {
		name  string
		args  args
		isErr bool
	}{
		{
			name: "with correct parameters",
			args: args{
				name:         name,
				description:  description,
				capabilities: capabilities,
			},
		},
		{
			name: "blank name",
			args: args{
				name:         "",
				description:  description,
				capabilities: capabilities,
			},
			isErr: true,
		},
		{
			name: "blank description",
			args: args{
				name:         description,
				description:  "",
				capabilities: capabilities,
			},
			isErr: true,
		},
		{
			name: "nil capabilities",
			args: args{
				name:         name,
				description:  description,
				capabilities: nil,
			},
			isErr: true,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			args := test.args
			_, err := Register(args.name, args.description, args.capabilities)
			if test.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestNewPlugin_multiple_errors(t *testing.T) {
	_, err := Register("", "", nil)
	require.Error(t, err)
	require.Equal(t, "validation errors: requires name, requires description, requires capabilities", err.Error())
}

func TestPlugin_Serve(t *testing.T) {
	capabilities := &plugin.Capabilities{}

	ran := false

	serveOpt := func(p *Plugin) {
		p.serverFactory = func(service plugin.Service) {
			ran = true
		}
	}

	p, err := Register("name", "description", capabilities, serveOpt)
	require.NoError(t, err)

	p.Serve()

	assert.True(t, ran)
}
