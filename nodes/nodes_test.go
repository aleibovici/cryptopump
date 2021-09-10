package nodes

import (
	"testing"

	"github.com/aleibovici/cryptopump/types"
)

func TestNode_GetRole(t *testing.T) {
	type args struct {
		configData  *types.Config
		sessionData *types.Session
	}
	tests := []struct {
		name string
		n    Node
		args args
	}{
		{
			name: "success",
			n:    Node{},
			args: args{
				configData: &types.Config{
					TestNet: true,
				},
				sessionData: &types.Session{
					MasterNode: true,
				},
			},
		},
		{
			name: "success",
			n:    Node{},
			args: args{
				configData: &types.Config{
					TestNet: false,
				},
				sessionData: &types.Session{
					MasterNode: true,
				},
			},
		},
		{
			name: "success",
			n:    Node{},
			args: args{
				configData: &types.Config{},
				sessionData: &types.Session{
					MasterNode: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Node{}
			n.GetRole(tt.args.configData, tt.args.sessionData)
		})
	}
}

func TestNode_ReleaseMasterRole(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name string
		n    Node
		args args
	}{
		{
			name: "success",
			n:    Node{},
			args: args{
				sessionData: &types.Session{
					MasterNode: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Node{}
			n.ReleaseMasterRole(tt.args.sessionData)
		})
	}
}
