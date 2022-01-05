package loader

import (
	"database/sql"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aleibovici/cryptopump/types"
	"github.com/paulbellamy/ratecounter"
)

// NewMock returns a new mock database and sqlmock.Sqlmock
func NewMock() (*sql.DB, sqlmock.Sqlmock) {

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestLoadSessionDataAdditionalComponents(t *testing.T) {
	type args struct {
		sessionData *types.Session
		marketData  *types.Market
		configData  *types.Config
	}

	db, mock := NewMock()
	defer db.Close()

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					ThreadID:    "c683ok5mk1u1120gnmmg",
					Symbol:      "BTCUSD",
					Db:          db,
					RateCounter: ratecounter.NewRateCounter(5 * time.Second),
					Global:      &types.Global{Profit: 0},
				},
				marketData: &types.Market{},
				configData: &types.Config{},
			},
			want:    []byte{},
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
			_, err := LoadSessionDataAdditionalComponents(tt.args.sessionData, tt.args.marketData, tt.args.configData)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSessionDataAdditionalComponents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestLoadSessionDataAdditionalComponentsAsync(t *testing.T) {
	type args struct {
		sessionData *types.Session
	}

	db, mock := NewMock()
	defer db.Close()

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				sessionData: &types.Session{
					Db:     db,
					Global: &types.Global{},
				},
			},
		},
	}

	columns := []string{"profit", "profitNet", "profitPct", "transactTime"}
	mock.ExpectBegin()                                                 /* begin transaction */
	mock.ExpectQuery(regexp.QuoteMeta("call cryptopump.GetGlobal()")). /* call procedure */
										WillReturnRows(sqlmock.NewRows(columns)) /* return 1 row */

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadSessionDataAdditionalComponentsAsync(tt.args.sessionData)
		})
	}
}
