package main

type Mail struct{
	To string `json:"to"`
	Proof_of_Work string `json:"proof_of_work"` //SHA1 hash value of 'header'
	Payload string `json:"payload"`
	/*
	 *Key - symmetric key encrypted with public key
	 *Payload encrypted with symmetric key
	 *Date
	 *Timestamp - string proof of work dentro de um certo intervalo de tempo
 	 *Header - version?:bits:date:to:randmString:counter
	 */
}
