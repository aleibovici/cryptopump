package mysql

import (
	"database/sql"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aleibovici/cryptopump/types"
	_ "github.com/go-sql-driver/mysql"
)

// NewMock returns a new mock database and sqlmock.Sqlmock
func NewMock() (*sql.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestGetThreadCount(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
				},
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	columns := []string{"count"}
	mock.ExpectBegin()                                                      /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadCount()")). /* call procedure */
										WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
				},
			},
			wantAmount: 0,
			wantErr:    false,
		},
	}

	columns := []string{"amountNullFloat64"}
	mock.ExpectBegin()                                                                  /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadTransactionAmount()")). /* call procedure */
												WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"threadID"}
	mock.ExpectBegin()                                                        /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetSessionStatus()")). /* call procedure */
											WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"profit", "profitNet", "profitPct", "transactTime"}
	mock.ExpectBegin()                                                 /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetGlobal()")). /* call procedure */
										WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"profit", "profitNet", "percentage"}
	mock.ExpectBegin()                                                 /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetProfit()")). /* call procedure */
										WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"fiatNullFloat64", "percentageNullFloat64"}
	mock.ExpectBegin()                                                            /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetProfitByThreadID(?)")). /* call procedure */
											WithArgs(tests[0].args.sessionData.ThreadID). /* with args */
											WillReturnRows(sqlmock.NewRows(columns))      /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"orderID", "cumulativeQuoteQty", "price", "executedQuantity"}
	mock.ExpectBegin()                                                                       /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadTransactionByThreadID(?)")). /* call procedure */
													WithArgs(tests[0].args.sessionData.ThreadID). /* with args */
													WillReturnRows(sqlmock.NewRows(columns))      /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
				side: "SELL",
			},
			wantErr: false,
		},
	}

	columns := []string{"count"}
	mock.ExpectBegin()                                                                     /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetOrderTransactionCount(?,?,?)")). /* call procedure */
												WithArgs(
								tests[0].args.sessionData.ThreadID,
								tests[0].args.side,
								-60). /* with args */
		WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
				price: 0.0,
			},
			wantErr: false,
		},
	}

	columns := []string{"count"}
	mock.ExpectBegin()                                                                                  /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadTransactiontUpmarketPriceCount(?,?)")). /* call procedure */
														WithArgs(
								tests[0].args.sessionData.ThreadID,
								tests[0].args.price). /* with args */
		WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:               db,
					ForceSellOrderID: 8551815,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"OrderID", "Price", "ExecutedQuantity", "CummulativeQuoteQty", "TransactTime"}
	mock.ExpectBegin()                                                            /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetOrderByOrderID(?,?)")). /* call procedure */
											WithArgs(
								tests[0].args.sessionData.ForceSellOrderID,
								tests[0].args.sessionData.ThreadID). /* with args */
		WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"CumulativeQuoteQuantity", "OrderID", "Price", "ExecutedQuantity", "TransactTime"}
	mock.ExpectBegin()                                                                 /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadLastTransaction(?)")). /* call procedure */
												WithArgs(tests[0].args.sessionData.ThreadID). /* with args */
												WillReturnRows(sqlmock.NewRows(columns))      /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
				marketData: &types.Market{
					Price: 0.0,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"CumulativeQuoteQuantity", "OrderID", "Price", "ExecutedQuantity", "TransactTime"}
	mock.ExpectBegin()                                                                            /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadTransactionByPriceHigher(?,?)")). /* call procedure */
													WithArgs(
								tests[0].args.sessionData.ThreadID,
								tests[0].args.marketData.Price). /* with args */
		WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
				marketData: &types.Market{
					Price: 0.0,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"CumulativeQuoteQuantity", "OrderID", "Price", "ExecutedQuantity", "TransactTime"}
	mock.ExpectBegin()                                                                      /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadTransactionByPrice(?,?)")). /* call procedure */
												WithArgs(
								tests[0].args.sessionData.ThreadID,
								tests[0].args.marketData.Price). /* with args */
		WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"OrderID", "Symbol"}
	mock.ExpectBegin()                                                                   /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetOrderTransactionPending(?)")). /* call procedure */
												WithArgs(tests[0].args.sessionData.ThreadID). /* with args */
												WillReturnRows(sqlmock.NewRows(columns))      /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
				},
			},
		},
	}

	columns := []string{"threadID", "threadIDSession"}
	mock.ExpectBegin()                                                                    /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadTransactionDistinct()")). /* call procedure */
												WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"symbol"}
	mock.ExpectBegin()                                                       /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetOrderSymbol(?)")). /* call procedure */
											WithArgs(tests[0].args.sessionData.ThreadID). /* with args */
											WillReturnRows(sqlmock.NewRows(columns))      /* return 1 row */

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetOrderSymbol(tt.args.sessionData)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrderSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetOrderTransactionSideLastTwo(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"side1", "side2"}
	mock.ExpectBegin()                                                                       /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetOrderTransactionSideLastTwo(?)")). /* call procedure */
													WithArgs(tests[0].args.sessionData.ThreadID). /* with args */
													WillReturnRows(sqlmock.NewRows(columns))      /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"side"}
	mock.ExpectBegin()                                                                    /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetLastOrderTransactionSide(?)")). /* call procedure */
												WithArgs(tests[0].args.sessionData.ThreadID). /* with args */
												WillReturnRows(sqlmock.NewRows(columns))      /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
				Side: "SELL",
			},
			wantErr: false,
		},
	}

	columns := []string{"price"}
	mock.ExpectBegin()                                                                       /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetLastOrderTransactionPrice(?,?)")). /* call procedure */
													WithArgs(
								tests[0].args.sessionData.ThreadID,
								tests[0].args.Side). /* with args */
		WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
				},
			},
			wantErr: false,
		},
	}

	columns := []string{"count"}
	mock.ExpectBegin()                                                                  /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetThreadTransactionCount(?)")). /* call procedure */
												WithArgs(tests[0].args.sessionData.ThreadID). /* with args */
												WillReturnRows(sqlmock.NewRows(columns))      /* return 1 row */

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

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
				},
				orderID: 1,
			},
			wantErr: false,
		},
	}

	columns := []string{"count"}
	mock.ExpectBegin()                                                                         /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.DeleteThreadTransactionByOrderID(?)")). /* call procedure */
													WithArgs(tests[0].args.orderID).         /* with args */
													WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteThreadTransactionByOrderID(tt.args.sessionData, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteThreadTransactionByOrderID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveOrder(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

	type args struct {
		sessionData   *types.Session
		order         *types.Order
		orderIDSource int64
		orderPrice    float64
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
					Db:       db,
				},
				order: &types.Order{
					ClientOrderID:           "0",
					CumulativeQuoteQuantity: 0,
					ExecutedQuantity:        0,
					OrderID:                 0,
					Price:                   0,
					Side:                    "0",
					Status:                  "0",
					Symbol:                  "0",
					TransactTime:            0,
					ThreadID:                0,
					ThreadIDSession:         0,
					OrderIDSource:           0,
				},
				orderIDSource: 0,
				orderPrice:    0,
			},
			wantErr: false,
		},
	}

	mock.ExpectBegin()                                                                        /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.SaveOrder(?,?,?,?,?,?,?,?,?,?,?,?)")). /* call procedure */
													WithArgs( /* with args */
			tests[0].args.order.ClientOrderID,
			tests[0].args.order.CumulativeQuoteQuantity,
			tests[0].args.order.ExecutedQuantity,
			tests[0].args.order.OrderID,
			tests[0].args.orderIDSource,
			tests[0].args.order.Price,
			tests[0].args.order.Side,
			tests[0].args.order.Status,
			tests[0].args.order.Symbol,
			tests[0].args.order.TransactTime,
			tests[0].args.sessionData.ThreadID,
			tests[0].args.sessionData.ThreadIDSession).
		WillReturnRows(sqlmock.NewRows([]string{""}))
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveOrder(tt.args.sessionData, tt.args.order, tt.args.orderIDSource, tt.args.orderPrice); (err != nil) != tt.wantErr {
				t.Errorf("SaveOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateOrder(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

	type args struct {
		sessionData             *types.Session
		OrderID                 int64
		CumulativeQuoteQuantity float64
		ExecutedQuantity        float64
		Price                   float64
		Status                  string
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
					Db:       db,
				},
				OrderID:                 0,
				CumulativeQuoteQuantity: 0,
				ExecutedQuantity:        0,
				Price:                   0,
				Status:                  "",
			},
			wantErr: false,
		},
	}

	mock.ExpectBegin()                                                            /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.UpdateOrder(?,?,?,?,?)")). /* call procedure */
											WithArgs( /* with args */
								tests[0].args.OrderID,
								tests[0].args.CumulativeQuoteQuantity,
								tests[0].args.ExecutedQuantity,
								tests[0].args.Price,
								tests[0].args.Status).
		WillReturnRows(sqlmock.NewRows([]string{""})) /* return empty row */
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateOrder(tt.args.sessionData, tt.args.OrderID, tt.args.CumulativeQuoteQuantity, tt.args.ExecutedQuantity, tt.args.Price, tt.args.Status); (err != nil) != tt.wantErr {
				t.Errorf("UpdateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateSession(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

	type args struct {
		configData  *types.Config
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
				configData: &types.Config{
					ExchangeName: "binance",
				},
				sessionData: &types.Session{
					Db:              db,
					ThreadID:        "c683ok5mk1u1120gnmmg",
					ThreadIDSession: "c683ok5mk1u1120gnmmg",
					SymbolFiat:      "",
					SymbolFiatFunds: 0,
					DiffTotal:       0,
					Status:          false,
				},
			},
			wantErr: false,
		},
	}

	mock.ExpectBegin()                                                                  /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.UpdateSession(?,?,?,?,?,?,?)")). /* call procedure */
												WithArgs( /* with args */
								tests[0].args.sessionData.ThreadID,
								tests[0].args.sessionData.ThreadIDSession,
								tests[0].args.configData.ExchangeName,
								tests[0].args.sessionData.SymbolFiat,
								tests[0].args.sessionData.SymbolFiatFunds,
								tests[0].args.sessionData.DiffTotal,
								tests[0].args.sessionData.Status).
		WillReturnRows(sqlmock.NewRows([]string{""})) /* return empty row */
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateSession(tt.args.configData, tt.args.sessionData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateGlobal(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
					Global: &types.Global{
						Profit:    0,
						ProfitNet: 0,
						ProfitPct: 0,
					},
				},
			},
			wantErr: false,
		},
	}

	mock.ExpectBegin()                                                           /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.UpdateGlobal(?,?,?,?)")). /* call procedure */
											WithArgs( /* with args */
								tests[0].args.sessionData.Global.Profit,
								tests[0].args.sessionData.Global.ProfitNet,
								tests[0].args.sessionData.Global.ProfitPct,
								time.Now().Unix()).
		WillReturnRows(sqlmock.NewRows([]string{""})) /* return empty row */
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateGlobal(tt.args.sessionData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateGlobal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveGlobal(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

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
					Db: db,
					Global: &types.Global{
						Profit:    0,
						ProfitNet: 0,
						ProfitPct: 0,
					},
				},
			},
			wantErr: false,
		},
	}

	mock.ExpectBegin()                                                         /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.SaveGlobal(?,?,?,?)")). /* call procedure */
											WithArgs( /* with args */
								tests[0].args.sessionData.Global.Profit,
								tests[0].args.sessionData.Global.ProfitNet,
								tests[0].args.sessionData.Global.ProfitPct,
								time.Now().Unix()).
		WillReturnRows(sqlmock.NewRows([]string{""})) /* return empty row */
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveGlobal(tt.args.sessionData); (err != nil) != tt.wantErr {
				t.Errorf("SaveGlobal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveSession(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

	type args struct {
		configData  *types.Config
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
				configData: &types.Config{
					ExchangeName: "binance",
				},
				sessionData: &types.Session{
					Db:              db,
					ThreadID:        "c683ok5mk1u1120gnmmg",
					ThreadIDSession: "c683ok5mk1u1120gnmmg",
					SymbolFiat:      "BTCUSDT",
					SymbolFiatFunds: 0,
					DiffTotal:       0,
					Status:          false,
				},
			},
			wantErr: false,
		},
	}

	mock.ExpectBegin()                                                                /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.SaveSession(?,?,?,?,?,?,?)")). /* call procedure */
												WithArgs( /* with args */
								tests[0].args.sessionData.ThreadID,
								tests[0].args.sessionData.ThreadIDSession,
								tests[0].args.configData.ExchangeName,
								tests[0].args.sessionData.SymbolFiat,
								tests[0].args.sessionData.SymbolFiatFunds,
								tests[0].args.sessionData.DiffTotal,
								tests[0].args.sessionData.Status).
		WillReturnRows(sqlmock.NewRows([]string{""})) /* return empty row */
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveSession(tt.args.configData, tt.args.sessionData); (err != nil) != tt.wantErr {
				t.Errorf("SaveSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

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
					Db:       db,
					ThreadID: "c683ok5mk1u1120gnmmg",
				},
			},
			wantErr: false,
		},
	}

	mock.ExpectBegin()                                                      /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.DeleteSession(?)")). /* call procedure */
										WithArgs( /* with args */
								tests[0].args.sessionData.ThreadID).
		WillReturnRows(sqlmock.NewRows([]string{""})) /* return empty row */
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteSession(tt.args.sessionData); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveThreadTransaction(t *testing.T) {

	db, mock := NewMock()
	defer db.Close()

	type args struct {
		sessionData             *types.Session
		OrderID                 int64
		CumulativeQuoteQuantity float64
		Price                   float64
		ExecutedQuantity        float64
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
					Db:       db,
					ThreadID: "c683ok5mk1u1120gnmmg",
				},
				OrderID:                 0,
				CumulativeQuoteQuantity: 0,
				Price:                   0,
				ExecutedQuantity:        0,
			},
			wantErr: false,
		},
	}

	mock.ExpectBegin()                                                                       /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.SaveThreadTransaction(?,?,?,?,?,?")). /* call procedure */
													WithArgs( /* with args */
								tests[0].args.sessionData.ThreadID,
								tests[0].args.sessionData.ThreadIDSession,
								tests[0].args.OrderID,
								tests[0].args.CumulativeQuoteQuantity,
								tests[0].args.Price,
								tests[0].args.ExecutedQuantity).
		WillReturnRows(sqlmock.NewRows([]string{""})) /* return empty row */
	mock.ExpectCommit()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveThreadTransaction(tt.args.sessionData, tt.args.OrderID, tt.args.CumulativeQuoteQuantity, tt.args.Price, tt.args.ExecutedQuantity); (err != nil) != tt.wantErr {
				t.Errorf("SaveThreadTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
