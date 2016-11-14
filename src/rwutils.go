package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

/*Retorna o conteudo de um ficheiro de texto para uma string
string nof - nome ou path do ficheiro .txt
string message - conteudo do ficheito; mensagem de texto a ser enviada no email
 */
func ReadFile(nof string) (message string){
	messageB, err := ioutil.ReadFile(nof)
	Check(err)
	message = string(messageB)
	return message
}

/*Escreve (JSON) para ficheiro mail_ready.json a
estrutura Mail
 */
func WriteJSON(mail Mail) {
	mailB, _ := json.MarshalIndent(mail, "", "    ")
	f, err := os.Create("mail_ready.json")//Pode mudar-se o nome do ficheiro mais tarde
	Check(err)
	defer f.Close()
	f.Write(mailB)
}

/*Faz o inverso da função acima.
Pode ser útil para posterior implementação
 */
func ReadJSON() (m Mail) {
	var content string = ReadFile("mail_ready.json")
	json.Unmarshal([]byte(content), &m)
	return m
}