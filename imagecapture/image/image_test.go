package image

import (
	"os"
	"testing"
)

func TestIsImageBlack(t *testing.T) {
	type args struct {
		src *os.File
	}
	tests := []struct {
		name    string
		args    args
		wantR   bool
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				src: nil,
			},
			wantR:   true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			tt.args.src, err = os.Open("/Users/mikhailcherniaev/Documents/GitHub/homeWebCamera/imagecapture/image/img.png")
			if err != nil {
				t.Errorf("IsImageBlack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer func() {
				tt.args.src.Close()
			}()
			gotR, err := IsImageBlack(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsImageBlack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("IsImageBlack() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
