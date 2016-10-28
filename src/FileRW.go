package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

/*Verificar se existe algum erro*/
func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*Retorna o conteudo de um ficheiro de texto para uma string
string nof - nome ou path do ficheiro .txt
string message - conteudo do ficheito; mensagem de texto a ser enviada no email
 */
func readFile(nof string) (message string){
	messageB, err := ioutil.ReadFile(nof)
	check(err)
	message = string(messageB)
	return message
}

/*Escreve (JSON) para ficheiro mail_ready.json a
estrutura Mail
 */
func writeJSON(mail Mail) {
	f, err := os.Create("mail_ready.json")//Pode mudar-se o nome do ficheiro mais tarde
	check(err)
	defer f.Close()
	mailB, _ := json.MarshalIndent(mail, "", "    ")
	f.Write(mailB)
}

/*Faz o inverso da função acima.
Pode ser útil para posterior implementação
 */
func readJSON() (mail Mail) {
	var content string = readFile("mail_ready.json")
	var m Mail
	json.Unmarshal([]byte(content), &m)
	return m
}