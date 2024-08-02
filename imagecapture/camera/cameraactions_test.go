package camera

import (
	"reflect"
	"testing"

	"github.com/veyanrech/homeWebCamera/config"
	"github.com/veyanrech/homeWebCamera/utils"
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

	lggr := utils.NewConsoleLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCameraByOS(config.NewConfig(), lggr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCameraByOS() = %v, want %v", got, tt.want)
			}
		})
	}
}
