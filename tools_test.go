package tools

import (
	"testing"
	"reflect"
)

func TestInt2LE(t *testing.T) {
	type cases struct {
		expected [4]byte
		pass     bool
	}
	casesMap := map[uint]cases{
		1: cases{
			expected: [4]byte{1, 0, 0, 0},
			pass:     true,
		},
		257: cases{
			expected: [4]byte{1, 1, 0, 0},
			pass:     true,
		},
		1234: cases{
			expected: [4]byte{4, 3, 2, 1},
			pass:     false,
		},
	}
	for input, v := range casesMap {
		found := Int2LE(input)
		if v.pass == (!reflect.DeepEqual(v.expected, found)) {
			t.Errorf("expected: %v, found: %v", v.expected, found)
		}
	}
}

func TestLE2Int(t *testing.T) {
	type cases struct {
		input []byte
		pass  bool
	}
	casesMap := map[uint]cases{
		1: cases{
			input: []byte{1},
			pass:  true,
		},
		256: cases{
			input: []byte{0, 1},
			pass:  true,
		},
		261: cases{
			input: []byte{5, 1},
			pass:  true,
		},
		512: cases{
			input: []byte{1, 2},
			pass:  false,
		},
	}
	for expected, v := range casesMap {
		found := LE2Int(v.input)
		if v.pass == (expected != found) {
			t.Errorf("expected: %v, found: %v", expected, found)
		}
	}
}
