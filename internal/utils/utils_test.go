// Copyright (c) 2024 Christopher Watson
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import "testing"

func TestEqualRunes(t *testing.T) {
	tests := []struct {
		name string
		a    []rune
		b    []rune
		want bool
	}{
		{
			name: "Equal empty slices",
			a:    []rune{},
			b:    []rune{},
			want: true,
		},
		{
			name: "Equal non-empty slices",
			a:    []rune{'a', 'b', 'c'},
			b:    []rune{'a', 'b', 'c'},
			want: true,
		},
		{
			name: "Different lengths",
			a:    []rune{'a', 'b'},
			b:    []rune{'a', 'b', 'c'},
			want: false,
		},
		{
			name: "Same length different runes",
			a:    []rune{'a', 'b', 'c'},
			b:    []rune{'a', 'b', 'd'},
			want: false,
		},
		{
			name: "Unicode runes",
			a:    []rune{'世', '界'},
			b:    []rune{'世', '界'},
			want: true,
		},
		{
			name: "Different Unicode runes",
			a:    []rune{'世', '界'},
			b:    []rune{'你', '好'},
			want: false,
		},
		{
			name: "Nil slices",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "One nil one empty",
			a:    nil,
			b:    []rune{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualRunes(tt.a, tt.b); got != tt.want {
				t.Errorf("EqualRunes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name string
		f    float32
		low  float32
		high float32
		want float32
	}{
		{
			name: "Within range",
			f:    5.0,
			low:  0.0,
			high: 10.0,
			want: 5.0,
		},
		{
			name: "Below range",
			f:    -5.0,
			low:  0.0,
			high: 10.0,
			want: 0.0,
		},
		{
			name: "Above range",
			f:    15.0,
			low:  0.0,
			high: 10.0,
			want: 10.0,
		},
		{
			name: "Equal to low",
			f:    0.0,
			low:  0.0,
			high: 10.0,
			want: 0.0,
		},
		{
			name: "Equal to high",
			f:    10.0,
			low:  0.0,
			high: 10.0,
			want: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.f, tt.low, tt.high); got != tt.want {
				t.Errorf("Clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClampInt(t *testing.T) {
	tests := []struct {
		name string
		i    int
		low  int
		high int
		want int
	}{
		{
			name: "Within range",
			i:    5,
			low:  0,
			high: 10,
			want: 5,
		},
		{
			name: "Below range",
			i:    -5,
			low:  0,
			high: 10,
			want: 0,
		},
		{
			name: "Above range",
			i:    15,
			low:  0,
			high: 10,
			want: 10,
		},
		{
			name: "Equal to low",
			i:    0,
			low:  0,
			high: 10,
			want: 0,
		},
		{
			name: "Equal to high",
			i:    10,
			low:  0,
			high: 10,
			want: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClampInt(tt.i, tt.low, tt.high); got != tt.want {
				t.Errorf("ClampInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
