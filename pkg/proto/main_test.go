package proto_test

import (
	"github.com/axetroy/go-server/pkg/proto"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	s := proto.NewProto("auth", map[string]interface{}{"hello": "world"})

	type args struct {
		link string
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.Proto
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				link: "auth://eyJoZWxsbyI6IndvcmxkIn0=",
			},
			wantErr: false,
			want:    &s,
		},
		{
			name: "basic",
			args: args{
				link: "auth://123123",
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := proto.Parse(tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
