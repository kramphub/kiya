package main

import "testing"

func TestBackup_String(t *testing.T) {
	bak := Backup{
		Secret: "foo bar",
		Data:   []byte{0x0, 0x1f, 0x2, 0x1e},
	}

	bakStr := bak.String()

	bak2 := Backup{}
	bak2.FromString(bakStr)

	if bak.Secret != bak2.Secret {
		t.Fail()
	}

	if len(bak.Data) != len(bak2.Data) {
		t.Fail()
	}

	if bak.Data[0] != bak2.Data[0] || bak.Data[1] != bak2.Data[1] || bak.Data[2] != bak2.Data[2] || bak.Data[3] != bak2.Data[3] {
		t.Fail()
	}
}
