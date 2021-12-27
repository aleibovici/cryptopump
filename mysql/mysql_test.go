package mysql

import (
	"testing"

	"github.com/aleibovici/cryptopump/types"
	_ "github.com/go-sql-driver/mysql"
)

func TestGetThreadCount(t *testing.T) {

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
			wantCount: 0,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := GetThreadCount(tt.args.sessionData)
			if (err != nil) != tt.wantErr && (gotCount > tt.wantCount) {
				return
			}
		})
	}
}

func TestGetThreadAmount(t *testing.T) {

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

func TestGetSessionStatus(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name       string
		args       args
		wantStatus string
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					Db: DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetSessionStatus(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetGlobal(t *testing.T) {
	type args struct {
		sessiondata *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessiondata: &types.Session{
					Db: DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, err := GetGlobal(tt.args.sessiondata)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGlobal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetProfit(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					Db: DBInit(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := GetProfit(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetProfitByThreadID(t *testing.T) {
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
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetProfitByThreadID(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProfitByThreadID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetThreadTransactionByThreadID(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetThreadTransactionByThreadID(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadTransactionByThreadID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetOrderTransactionCount(t *testing.T) {
	type args struct {
		sessionData *types.Session
		side        string
	}
	tests := []struct {
		name      string
		args      args
		wantCount float64
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
				side: "SELL",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetOrderTransactionCount(tt.args.sessionData, tt.args.side)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderTransactionCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetThreadTransactiontUpmarketPriceCount(t *testing.T) {
	type args struct {
		sessionData *types.Session
		price       float64
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
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
				price: 0.0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := GetThreadTransactiontUpmarketPriceCount(tt.args.sessionData, tt.args.price)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadTransactiontUpmarketPriceCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("GetThreadTransactiontUpmarketPriceCount() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestGetOrderByOrderID(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID:         "c683ok5mk1u1120gnmmg",
					Db:               DBInit(),
					ForceSellOrderID: 8551815,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetOrderByOrderID(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderByOrderID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetThreadLastTransaction(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetThreadLastTransaction(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadLastTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetThreadTransactionByPriceHigher(t *testing.T) {
	type args struct {
		marketData  *types.Market
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
				marketData: &types.Market{
					Price: 0.0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetThreadTransactionByPriceHigher(tt.args.marketData, tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadTransactionByPriceHigher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetThreadTransactionByPrice(t *testing.T) {
	type args struct {
		marketData  *types.Market
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
				marketData: &types.Market{
					Price: 0.0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetThreadTransactionByPrice(tt.args.marketData, tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadTransactionByPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetOrderTransactionPending(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetOrderTransactionPending(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderTransactionPending() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetThreadTransactionDistinct(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetThreadTransactionDistinct(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadTransactionDistinct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetOrderSymbol(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name       string
		args       args
		wantSymbol string
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
			wantSymbol: "BTCUSDT",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSymbol, err := GetOrderSymbol(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSymbol != tt.wantSymbol {
				t.Errorf("GetOrderSymbol() = %v, want %v", gotSymbol, tt.wantSymbol)
			}
		})
	}
}

func TestGetOrderTransactionSideLastTwo(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetOrderTransactionSideLastTwo(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderTransactionSideLastTwo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetLastOrderTransactionSide(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetLastOrderTransactionSide(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastOrderTransactionSide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetLastOrderTransactionPrice(t *testing.T) {
	type args struct {
		sessionData *types.Session
		Side        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
				Side: "SELL",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetLastOrderTransactionPrice(tt.args.sessionData, tt.args.Side)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastOrderTransactionPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetThreadTransactionCount(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID: "c683ok5mk1u1120gnmmg",
					Db:       DBInit(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetThreadTransactionCount(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetThreadTransactionCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDeleteThreadTransactionByOrderID(t *testing.T) {
	type args struct {
		sessionData *types.Session
		orderID     int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					Db: DBInit(),
				},
				orderID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteThreadTransactionByOrderID(tt.args.sessionData, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteThreadTransactionByOrderID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
