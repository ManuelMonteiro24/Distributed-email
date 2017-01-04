package main

import (
	"encoding/json"
	"io/ioutil"
//	"os"
)

/*Retorna o conteudo de um ficheiro de texto para uma string
string nof - nome ou path do ficheiro .txt
string message - conteudo do ficheito; mensagem de texto a ser enviada no email
*/
func ReadFile(nof string) (message string) {
	messageB, err := ioutil.ReadFile(nof)
	Check(err)
	message = string(messageB)
	return message
}

/*Escreve (JSON) para ficheiro mail_ready.json a
estrutura Mail
*/
func WriteJSON(mail Mail) []byte {
	res, _ := json.MarshalIndent(mail, "", "    ")
	return res
}

/*Faz o inverso da função acima.
Pode ser útil para posterior implementação
*/
func ReadJSON(data []byte) (m Mail) {
	json.Unmarshal(data, &m)
	return m
}
