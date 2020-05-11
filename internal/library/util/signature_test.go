package util_test

import (
	"github.com/axetroy/go-server/internal/library/util"
	"testing"
)

func TestSignature(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "common",
			args: args{
				input: "123",
			},
			want: "aa154db7952aa6a2656fed90d0d88f2e87560a6ca7d7ed180ac76705fdc1639b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.Signature(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Signature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Signature() got = %v, want %v", got, tt.want)
			}
		})
	}
}
