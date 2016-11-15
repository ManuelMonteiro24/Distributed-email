package main

import "reflect"

type Mail struct{
	Header string `json:"header"` //
	Payload string `json:"payload"` //encrypted with symmetric key
	SymmKey string `json:"symmkey"`//symmetric key encrypted with "to" public key
	Proof_of_Work string `json:"proof_of_work"` //SHA1 hash value of 'header'
	/*
	 *randString could be initialization vector (IV) if Cipher Block Chaining is used to encrypt
	 *Timestamp - string proof of work dentro de um certo intervalo de tempo
	 */
}

/*Gives a value to the field of Mail
specified by index = {0:Header, 1:Payload, 2:SymmKey, 3:Proof_of_Work}
 */
func ( m *Mail) AddField(index int, content string) {
	mValue := reflect.ValueOf(m).Elem()//reflect.ValueOf(m)returns new Value inited to the concrete value stored in the interface &m
	fieldValue := mValue.Field(index)// Field(index) returns Value of field with index "index"
	fieldValue.SetString(content)//set actual field value
}