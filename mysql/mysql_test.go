package mysql

import (
	"os"
	"testing"

	"github.com/aleibovici/cryptopump/types"
	_ "github.com/go-sql-driver/mysql"
)

func Test_UsingEnvvar(t *testing.T) {

	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "swatch!12")
	os.Setenv("DB_TCP_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_NAME", "cryptopump")

}

func TestGetThreadCount(t *testing.T) {

	Test_UsingEnvvar(t)

	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					Db: DBInit(),
				},
			},
			wantCount: 1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := GetThreadCount(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("GetThreadCount() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestGetThreadAmount(t *testing.T) {

	Test_UsingEnvvar(t)

	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name       string
		args       args
		wantAmount float64
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					Db: DBInit(),
				},
			},
			wantAmount: 0,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAmount, err := GetThreadAmount(tt.args.sessionData)
			if (err == nil) && gotAmount > 0 {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAmount != tt.wantAmount {
				t.Errorf("GetThreadAmount() = %v, want %v", gotAmount, tt.wantAmount)
			}
		})
	}
}

func TestGetProfit(t *testing.T) {

	Test_UsingEnvvar(t)

	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name           string
		args           args
		wantFiat       float64
		wantPercentage float64
		wantErr        bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					Db: DBInit(),
				},
			},
			wantFiat:       0,
			wantPercentage: 0,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiat, gotPercentage, err := GetProfit(tt.args.sessionData)
			if (err == nil) && gotFiat > 0 {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFiat != tt.wantFiat {
				t.Errorf("GetProfit() gotFiat = %v, want %v", gotFiat, tt.wantFiat)
			}
			if gotPercentage != tt.wantPercentage {
				t.Errorf("GetProfit() gotPercentage = %v, want %v", gotPercentage, tt.wantPercentage)
			}
		})
	}
}
