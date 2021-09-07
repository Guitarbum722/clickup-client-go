package main

import (
	"reflect"
	"testing"
)

func Test_chunkSlice(t *testing.T) {
	type args struct {
		slice     []string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "TestSuccess Chunks of 2",
			args: args{
				slice:     []string{"one", "two", "three", "four", "five"},
				chunkSize: 2,
			},
			want: [][]string{
				{
					"one",
					"two",
				},
				{
					"three",
					"four",
				},
				{
					"five",
				},
			},
		},
		{
			name: "TestSuccess Chunks of 3",
			args: args{
				slice:     []string{"one", "two", "three", "four", "five"},
				chunkSize: 3,
			},
			want: [][]string{
				{
					"one",
					"two",
					"three",
				},
				{
					"four",
					"five",
				},
			},
		},
		{
			name: "TestSuccess Chunks of 5",
			args: args{
				slice:     []string{"one", "two", "three", "four", "five"},
				chunkSize: 5,
			},
			want: [][]string{
				{
					"one",
					"two",
					"three",
					"four",
					"five",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := chunkSlice(tt.args.slice, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("chunkSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
