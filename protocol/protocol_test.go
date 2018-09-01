package protocol

import (
	"testing"
	"reflect"
	"os"
	"fmt"
	"github.com/mugsoft/tools/bytesize"
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
			Type:      1,
			CMD:       CMD_NEWGAME,
			DataLen:   1,
			Data:      []byte{1},
		},
	}

	for _, v := range cases {
		var found = new(MUGSOFT)
		_ = found.Parse(v.Bytes())
		if !reflect.DeepEqual(found.Bytes(), v.Bytes()) {
			t.Errorf("expected MUGSOFT: %v, found: %v ", v.Bytes(), found.Bytes())
		}
	}

	cases2 := map[ERRCODE]MUGSOFT{
		ERR_SUCCES: MUGSOFT{
			signature: [2]byte{'m', 'g'},
			Type:      1,
			CMD:       CMD_NEWGAME,
			DataLen:   1,
			Data:      []byte{1},
		},
		ERR_DATA_LEN: MUGSOFT{
			signature: [2]byte{'m', 'g'},
			Type:      2,
			CMD:       CMD_NEWGAME,
			DataLen:   1,
			Data:      []byte{1, 2},
		},
	}

	for k, v := range cases2 {
		var found = new(MUGSOFT)
		err__code := found.Parse(v.Bytes())
		if k != err__code {
			t.Errorf("expected error code: %v, found: %v ", k, err__code)
		}
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
				Type:      2,
				CMD:       1,
				DataLen:   1,
				Data:      []byte{1},
			},
			expected: []byte{'m', 'g', 2, 0, 1, 0, 1, 0, 0, 0, 1},
		},
		false: cases{
			input: MUGSOFT{
				signature: [2]byte{'m', 'g'},
				Type:      2,
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

type fake__sock struct {
	data      []byte
	len__data int
	offset    int
	len__buf  int
}

func (sck *fake__sock) Read(b []byte) (n int, err error) {
	sck.len__buf = len(b)

	for k, _ := range b {
		if sck.len__data < sck.offset+k {

			sck.offset += k
			return k, nil
		}

		b[k] = sck.data[sck.offset+k]
		if k > 104728 {
			sck.offset += k
			return k, nil
		}
	}

	sck.offset += sck.len__buf

	return len(b), nil
}

func TestMUGSOFT_Scan(t *testing.T) {
	var fake__sck = new(fake__sock)
	var prot = new(MUGSOFT)

	var dt = make([]byte, (bytesize.MB*50)+m__meta_len)
	var file, err = os.Open("1KB__file.bytes")
	if nil != err {
		fmt.Errorf("os.Open error : %v", err)
		return
	}

	_, err = file.Read(dt)
	if nil != err {
		fmt.Errorf("file.Read error : %v", err)
		return
	}
	fake__sck.data = dt
	fake__sck.len__data = len(fake__sck.data)

	for i := 12; i > 0; i-- {
		err = prot.Scan(fake__sck)
		if nil != err {
			t.Errorf("mugsoft scan error : %v", err)
			return
		}
		fake__sck.offset = 0
	}
}
