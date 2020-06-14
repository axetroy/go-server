// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package migrate

import (
	"testing"
	"time"
)

func Test_generateTableName(t *testing.T) {
	type args struct {
		tableName string
		date      time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic",
			args: args{
				tableName: "user",
				date:      time.Date(2000, 6, 1, 0, 0, 0, 0, time.Now().Location()),
			},
			want: "user_200006",
		},
		{
			name: "basic2",
			args: args{
				tableName: "login_log",
				date:      time.Date(2020, 11, 20, 0, 0, 0, 0, time.Now().Location()),
			},
			want: "login_log_202011",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateTableName(tt.args.tableName, tt.args.date); got != tt.want {
				t.Errorf("generateTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
