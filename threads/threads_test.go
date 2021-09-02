package threads

import (
	"testing"

	"github.com/aleibovici/cryptopump/types"
)

var sessionData = &types.Session{
	ThreadID:   "c2q3mt84a8024t1f6590",
	Symbol:     "BTCUSDT",
	SymbolFiat: "USDT",
	MasterNode: false,
}

func TestThread_Lock(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name string
		tr   Thread
		args args
		want bool
	}{
		{
			name: "success",
			tr:   Thread{},
			args: args{
				sessionData: sessionData,
			},
			want: true,
		},
		{
			name: "success",
			tr:   Thread{},
			args: args{
				sessionData: &types.Session{
					ThreadID: "",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Thread{}
			if got := tr.Lock(tt.args.sessionData); got != tt.want {
				t.Errorf("Thread.Lock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestThread_Unlock(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name string
		tr   Thread
		args args
	}{
		{
			name: "success",
			tr:   Thread{},
			args: args{
				sessionData: sessionData,
			},
		},
		{
			name: "success",
			tr:   Thread{},
			args: args{
				sessionData: &types.Session{
					ThreadID: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Thread{}
			tr.Unlock(tt.args.sessionData)
		})
	}
}
