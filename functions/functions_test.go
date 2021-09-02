package functions

import (
	"testing"
)

func TestFloat64ToStr(t *testing.T) {
	type args struct {
		value float64
		prec  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: args{
				value: 1.11111111,
				prec:  8,
			},
			want: "1.11111111",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64ToStr(tt.args.value, tt.args.prec); got != tt.want {
				t.Errorf("Float64ToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToFloat64(t *testing.T) {
	type args struct {
		value int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "success",
			args: args{
				value: 1,
			},
			want: 1.00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToFloat64(tt.args.value); got != tt.want {
				t.Errorf("IntToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToInt(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name  string
		args  args
		wantR int
	}{
		{
			name: "success",
			args: args{
				value: "1",
			},
			wantR: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := StrToInt(tt.args.value); gotR != tt.wantR {
				t.Errorf("StrToInt() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestStrToFloat64(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name  string
		args  args
		wantR float64
	}{
		{
			name: "success",
			args: args{
				value: "11.11",
			},
			wantR: 11.11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotR := StrToFloat64(tt.args.value); gotR != tt.wantR {
				t.Errorf("StrToFloat64() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestGetThreadID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "success",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetThreadID(); got == "" {
				t.Errorf("GetThreadID() = %v, want %v", got, tt.want)
			}
		})
	}
}
