package model

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cometbft/cometbft/abci/types"
)

func TestSetGet(t *testing.T) {
	os.RemoveAll(dbPath)
	NewDB()
	defer DBINS.Close()

	Set("test", []byte("test"))
	v, err := Get("test")
	if err != nil {
		t.Error(err)
	}
	if string(v) != "test" {
		t.Error("value not equal")
	}

	// defer os.RemoveAll(dbPath)
}

func TestGetAbciMessage(t *testing.T) {
	os.RemoveAll(dbPath)
	NewDB()
	defer DBINS.Close()

	initv := types.ValidatorUpdate{Power: 10000}

	key := "validator" + fmt.Sprint(time.Now().Unix())
	err := SetAbciMessage("", key, &initv)
	if err != nil {
		t.Error(err)
	}

	returnValue, err := GetAbciMessage[types.ValidatorUpdate]("", key)
	if err != nil {
		t.Error(err)
	}

	if returnValue.Power != 10000 {
		t.Error("value not equal")
	}

	list, err := GetAbciMessageList[types.ValidatorUpdate]("", "validator")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(list)
}
