package main

import (
	"fmt"
)

/*Gera um ficheiro .json com os parametros da estrutura Mail definida em struts.
Lê o ficheiro .json criado de volta para a estrutura Mail
 */
func main() {
	//var nof string = "message.txt"
	//payload := readFile(nof)
	var mail Mail

	(&mail).AddField(1,"1")
	(&mail).AddField(0,"0")
	WriteJSON(mail)
	mail2 := ReadJSON()
	fmt.Print(mail2.Header + " " + mail2.Proof_of_Work + " " + mail2.Payload)//só para debug
}

/*func main() {

	/*dat, err := ioutil.ReadFile("codeGO.txt")
	check(err)
	fmt.Println(string(dat))

	f, err := os.Open("codeGo.txt")
	/*b1 := make([]byte, 3)
	n1, err := f.Read(b1)
	check(err)
	fmt.Println("%d bytes = %s", n1, string(dat))

	//Read Files from byte 6 as set by Seek()
	o2, err := f.Seek(6,0)
	check(err)
	b2 := make([]byte, 5)
	n2, err := f.Read(b2)
	check(err)
	fmt.Printf("%d bytes = @ %d: %s", n2, o2, string(b2))
}*/