// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package migrate

import (
	"testing"
)

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
