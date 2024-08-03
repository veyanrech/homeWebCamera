package utils

import "testing"

func TestGenerateFilename(t *testing.T) {
	type args struct {
		additional string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"test1", args{"test"}, "2021_09_01_15_04_05_000_test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateFilename(tt.args.additional); got != tt.want {
				t.Errorf("GenerateFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
