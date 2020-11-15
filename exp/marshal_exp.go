package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
)

type userInfo struct {
	UserName string `json:"user"`
	Password string `json:"passwd"`
	Count    int    `json:"count"`
}

func main() {
	u1 := &userInfo{
		UserName: "harry",
		Password: "12345678",
		Count:    100,
	}
	var ub []byte
	var err error
	if ub, err = json.Marshal(u1); err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(ub))
	var u2 userInfo
	if err = json.Unmarshal(ub, &u2); err != nil {
		log.Fatal(err)
	}
	fmt.Println(u2)

	fmt.Println("===============")

	var buf bytes.Buffer
	if err = gob.NewEncoder(&buf).Encode(u1); err != nil {
		log.Fatal(err)
	}
	ub = buf.Bytes()

	if err = gob.NewDecoder(bytes.NewBuffer(ub)).Decode(&u2); err != nil {
		log.Fatal(err)
	}
	fmt.Println(u2)

}
