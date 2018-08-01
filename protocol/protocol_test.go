package protocol

import (
	"testing"
	"reflect"
)

func Test_check__sig(t *testing.T) {
	cases := map[string]bool{
		"mg": true,
		"m":  false,
		"g":  false,
		"":   false,
		"dt": false,
	}
	for k, v := range cases {
		if check__sig([]byte(k)) != v {
			t.Fatal("case fails", k)
		}
	}
}

func TestMUGSOFT_Parse(t *testing.T) {
	cases := map[ERRCODE]MUGSOFT{
		ERR_SUCCES: MUGSOFT{
			signature: [2]byte{'m', 'g'},
			GameType:  2,
			CMD:       CMD_NEWGAME,
			DataLen:   1,
			Data:      []byte{1},
		},
		ERR_DATA_LEN: MUGSOFT{
			signature: [2]byte{'m', 'g'},
			GameType:  2,
			CMD:       CMD_NEWGAME,
			DataLen:   1,
			Data:      []byte{1, 2},
		},
	}

	for k, v := range cases {
		var found = new(MUGSOFT)
		err__code := found.Parse(v.Bytes())
		if k != err__code {
			t.Errorf("expected error code: %v, found: %v ", k, err__code)
		}
		if !reflect.DeepEqual(found.Bytes(), v.Bytes()) {
			t.Errorf("expected MUGSOFT: %v, found: %v ", found.Bytes(), v.Bytes())
		}
	}

	expexted := MUGSOFT{
		signature: [2]byte{'o', 'g'},
		GameType:  2,
		CMD:       CMD_NEWGAME,
		DataLen:   1,
		Data:      []byte{1},
	}
	var found = new(MUGSOFT)
	err__code := found.Parse(expexted.Bytes())
	if ERR_INVALID_SIG != err__code {
		t.Errorf("expected error code: %v, found: %v ", ERR_INVALID_SIG, err__code)
	}
}

func TestMUGSOFT_Bytes(t *testing.T) {
	type cases struct {
		input    MUGSOFT
		expected []byte
	}
	casesMap := map[bool]cases{
		true: cases{
			input: MUGSOFT{
				signature: [2]byte{'m', 'g'},
				GameType:  2,
				CMD:       1,
				DataLen:   1,
				Data:      []byte{1},
			},
			expected: []byte{'m', 'g', 2, 0, 1, 0, 1, 0, 0, 0, 1},
		},
		false: cases{
			input: MUGSOFT{
				signature: [2]byte{'m', 'g'},
				GameType:  2,
				CMD:       1,
				DataLen:   1,
				Data:      []byte{1},
			},
			expected: []byte{'m', 'g', 0, 2, 0, 1, 0, 1, 0, 0, 0, 1},
		},
	}
	for k, v := range casesMap {
		found := v.input.Bytes()
		if k == (!reflect.DeepEqual(v.expected, found)) {
			t.Errorf("expected: %v, found: %v", v.expected, found)
		}
	}
}
