package main

import (
	"fmt"
	"time"

	tarantool "github.com/tarantool/go-tarantool"
)

type Tuple struct {
	/* instruct msgpack to pack this struct as array,
	 * so no custom packer is needed */
	_msgpack struct{} `msgpack:",asArray"`
	ID       uint
}

func trntTest() {
	//spaceNo := uint32(512)
	//indexNo := uint32(0)

	server := "127.0.0.1:3309"
	opts := tarantool.Opts{
		Timeout:       50 * time.Millisecond,
		Reconnect:     100 * time.Millisecond,
		MaxReconnects: 3,
		User:          "user",
		Pass:          "resu",
	}
	client, err := tarantool.Connect(server, opts)
	if err != nil {
		fmt.Printf("Failed to connect: %s\n", err.Error())
		return
	}

	resp, err := client.Ping()
	fmt.Println("Ping Code", resp.Code)
	fmt.Println("Ping Data", resp.Data)
	fmt.Println("Ping Error", err)

	resp, err = client.Insert(2, []interface{}{"SessID", "userID", 1, "en", 2323}) //"SessID", "userID", authorized, preferredlang, startTime
	fmt.Println("Insert Error", err)
	fmt.Println("Insert Code", resp.Code)
	fmt.Println("Insert Data", resp.Data)

}
