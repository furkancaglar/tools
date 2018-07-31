package protocol

import "testing"

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
func Test_Unmarshal(t *testing.T) {
	// t.Fail()
}
