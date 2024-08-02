package client

import (
	"fmt"
	"testing"
)

func TestRoundBufferQueue_Add(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				s: []string{"1", "2", "3", "4", "5"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Test 1")
			rbq := NewRoundBufferQueue(3)
			for _, s := range tt.args.s {
				rbq.Add(s)
				fmt.Println(rbq.Get())
			}
			fmt.Println(len(rbq.queue))
		})
	}
	for _, tt := range tests {
		fmt.Println("Test 2")
		t.Run(tt.name, func(t *testing.T) {
			rbq := NewRoundBufferQueue(3)
			fmt.Println(rbq.Get())
			fmt.Println(len(rbq.queue))
		})
	}
	for _, tt := range tests {
		fmt.Println("Test 3")
		t.Run(tt.name, func(t *testing.T) {
			rbq := NewRoundBufferQueue(3)
			rbq.Add(tt.args.s[0])
			fmt.Println(rbq.Get())
			fmt.Println(rbq.Get())
			fmt.Println(len(rbq.queue))
		})
	}
}
