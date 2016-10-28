package main

type Mail struct{
	To string `json:"to"`
	Proof_of_Work string `json:"proof_of_work"`
	Payload string `json:"payload"`
	/*Assinatura (dentro ou não da payload??)
	 *Timestamp
	 */
}
