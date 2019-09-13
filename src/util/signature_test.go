package util_test

import (
	"github.com/axetroy/go-server/src/util"
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
			want: "38130d15223d921022b83673aaf26dde49bdb103e9471a1f52a801919364924a",
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
