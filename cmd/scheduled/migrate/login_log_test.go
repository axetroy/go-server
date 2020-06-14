// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package migrate

import (
	"testing"
	"time"
)

func TestLoginLog_GetTimeInterval(t *testing.T) {
	{
		now := time.Date(2000, 6, 1, 0, 0, 0, 0, time.Now().Location())
		tests := []struct {
			name string
			want time.Duration
		}{
			{
				name: "basic",
				want: time.Hour * 24 * 30,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := LoginLog{}
				if got := c.GetTimeInterval(now); got != tt.want {
					t.Errorf("GetTimeInterval() = %v, want %v", got, tt.want)
				}
			})
		}
	}

	{
		now := time.Date(2000, 5, 1, 0, 0, 0, 0, time.Now().Location())
		tests := []struct {
			name string
			want time.Duration
		}{
			{
				name: "basic",
				want: time.Hour * 24 * 31,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := LoginLog{}
				if got := c.GetTimeInterval(now); got != tt.want {
					t.Errorf("GetTimeInterval() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}

func TestLoginLog_Next(t *testing.T) {
	tests := []struct {
		name             string
		wantShouldGoNext bool
		wantErr          bool
	}{
		{
			name:             "basic",
			wantShouldGoNext: false,
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := LoginLog{}
			gotShouldGoNext, err := c.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotShouldGoNext != tt.wantShouldGoNext {
				t.Errorf("Next() gotShouldGoNext = %v, want %v", gotShouldGoNext, tt.wantShouldGoNext)
			}
		})
	}
}
