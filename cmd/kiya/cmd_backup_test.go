package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"testing/fstest"
)

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

func TestBackupWithoutEncryption(t *testing.T) {
	input := map[string]interface{}{
		"bar": "bar string",
	}

	buf, err := json.Marshal(input)

	require.NoError(t, err)

	backup := Backup{
		Secret:    "",
		Encrypted: false,
		Data:      buf,
	}

	mockFS := fstest.MapFS{
		"backup_test": {
			Data: []byte(backup.String()),
		},
	}

	backBuf, err := mockFS.ReadFile("backup_test")
	require.NoError(t, err)

	backup2 := Backup{}
	backup2.FromString(string(backBuf))

	require.Equal(t, len(backup.Data), len(backup2.Data))

	backupData := make(map[string]interface{})

	err = json.Unmarshal(backup2.Data, &backupData)
	require.NoError(t, err)

	bar := input["bar"]
	bar2 := backupData["bar"]

	require.Equal(t, bar, bar2)
}
