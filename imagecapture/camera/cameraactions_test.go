package camera

import (
	"os"
	"reflect"
	"testing"

	"github.com/veyanrech/homeWebCamera/imagecapture/config"
	"github.com/veyanrech/homeWebCamera/imagecapture/utils"
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

	lggr := utils.NewConsoleLogger(utils.LogLevel(1))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var filename string

			switch opsys := utils.GetOS(); opsys {
			case "darwin":
				filename = "." + string(os.PathSeparator) + "macos.config.json"
			case "windows":
				filename = "." + string(os.PathSeparator) + "win.config.json"
			}

			if got := NewCameraByOS(config.NewConfig(filename), lggr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCameraByOS() = %v, want %v", got, tt.want)
			}
		})
	}
}
