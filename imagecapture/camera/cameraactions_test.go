package camera

import (
	"reflect"
	"testing"
)

func TestNewCameraByOS(t *testing.T) {
	type args struct {
		dn []string
	}
	tests := []struct {
		name string
		args args
		want Camera
	}{
		{
			name: "Test NewCameraByOS",
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCameraByOS(tt.args.dn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCameraByOS() = %v, want %v", got, tt.want)
			}
		})
	}
}
