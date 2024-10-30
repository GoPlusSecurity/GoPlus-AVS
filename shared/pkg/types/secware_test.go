package types

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestSecwareTask_JSON(t *testing.T) {
	task := SecwareTask{
		SecwareId:      1,
		SecwareVersion: 2,
		SignedTx:       []byte{0xab, 0xcd},
		StartTime:      0x12345678,
		EndTime:        0x87654321,
		Args:           `{"key":"value"}`,
	}
	jsonBytes, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	expect := `{"secware_id":1,"secware_version":2,"signed_tx":"0xabcd","start_time":"0x12345678","end_time":"0x87654321","args":"{\"key\":\"value\"}"}`
	actual := string(jsonBytes)
	if actual != expect {
		t.Fatalf("expect '%s', got '%s'", expect, actual)
	}

	var newTask SecwareTask
	if err := json.Unmarshal(jsonBytes, &newTask); err != nil {
		t.Error(err)
	}

	if task.SecwareId != newTask.SecwareId ||
		task.SecwareVersion != newTask.SecwareVersion ||
		!bytes.Equal(task.SignedTx, newTask.SignedTx) ||
		task.StartTime != newTask.StartTime ||
		task.EndTime != newTask.EndTime ||
		task.Args != newTask.Args {
		t.Fatalf("expect '%#v', got '%#v'", task, newTask)
	}
}
