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
			want: "707bcdd6f7e49913fd0be83f3ac3046b5b445492efe331189f4210ff7d114f45",
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
