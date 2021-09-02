package logger

import (
	"testing"

	"github.com/aleibovici/cryptopump/types"
)

func TestLogEntry_Do(t *testing.T) {
	type fields struct {
		Config   *types.Config
		Market   *types.Market
		Session  *types.Session
		Order    *types.Order
		Message  string
		LogLevel string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "success",
			fields: fields{
				Config:   &types.Config{},
				Market:   &types.Market{},
				Session:  &types.Session{},
				Order:    &types.Order{},
				Message:  "UP",
				LogLevel: "infoLevel",
			},
		},
		{
			name: "success",
			fields: fields{
				Config:   &types.Config{},
				Market:   &types.Market{},
				Session:  &types.Session{},
				Order:    &types.Order{},
				Message:  "BUY",
				LogLevel: "infoLevel",
			},
		},
		{
			name: "success",
			fields: fields{
				Config:   &types.Config{},
				Market:   &types.Market{},
				Session:  &types.Session{},
				Order:    &types.Order{},
				Message:  "SELL",
				LogLevel: "infoLevel",
			},
		},
		{
			name: "success",
			fields: fields{
				Config: &types.Config{
					Debug: true,
				},
				Market:   &types.Market{},
				Session:  &types.Session{},
				Order:    &types.Order{},
				Message:  "CANCELED",
				LogLevel: "infoLevel",
			},
		},
		{
			name: "success",
			fields: fields{
				Config: &types.Config{
					Debug: true,
				},
				Market:   &types.Market{},
				Session:  &types.Session{},
				Order:    &types.Order{},
				Message:  "",
				LogLevel: "infoLevel",
			},
		},
		{
			name: "success",
			fields: fields{
				Config:   &types.Config{},
				Market:   &types.Market{},
				Session:  &types.Session{},
				Order:    &types.Order{},
				Message:  "",
				LogLevel: "debugLevel",
			},
		},
		{
			name: "success",
			fields: fields{
				Config:   &types.Config{},
				Market:   &types.Market{},
				Session:  &types.Session{},
				Order:    &types.Order{},
				Message:  "",
				LogLevel: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logEntry := LogEntry{
				Config:   tt.fields.Config,
				Market:   tt.fields.Market,
				Session:  tt.fields.Session,
				Order:    tt.fields.Order,
				Message:  tt.fields.Message,
				LogLevel: tt.fields.LogLevel,
			}
			logEntry.Do()
		})
	}
}
