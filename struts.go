package main

import (
	"reflect"
)

type Mail struct{
	Header string `json:"header"` //
	Proof_of_Work string `json:"proof_of_work"` //SHA1 hash value of 'header'
	From string `json:"from"`//From field
	Payload string `json:"payload"` //encrypted with symmetric key; format = {PoW//\\FROM//\\TO//\\MESSAGE//\\SIGNATURE}
	SymKey string `json:"symmkey"`//symmetric key encrypted with "to" public key; SymKey in base64
}

type Header struct{
	ZeroCount int //num of most significant zero bits required for hash digest
	Date string //Date
	Resource string //"To" field from Mail
	RandString string //randomly generated string
	Counter int //counter
}

/*Gives a value to the field of Mail
specified by index = {0:Header, 1:Payload, 2:Proof_of_Work, 3:SymmKey}
 */
func ( m *Mail) AddField(index int, content string) {
	mValue := reflect.ValueOf(m).Elem()//reflect.ValueOf(m)returns new Value inited to the concrete value stored in the interface &m
	fieldValue := mValue.Field(index)// Field(index) returns Value of field with index "index"
	fieldValue.SetString(content)//set actual field value
}
