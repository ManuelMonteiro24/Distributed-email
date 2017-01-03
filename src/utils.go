package main


import (
	"time"
	"math/rand"
	"fmt"
	"os"
	"encoding/base64"
)

const letterArray = "absdefghijklmnopqrstuvyxwzABCDEFGHIJKLMNOPQRSTUVYXWZ"


/*Verificar se existe algum erro*/
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

/*Função para construir a string do field Header
da estrutura Mail. Layout em strut.go
 */
func buildHeader(to string) (header string) {
	return header

}

/*Retorna uma string gerada aleatoriamente de n bytes*/
func RndStr(n int) string {
	rndBytes := make([]byte, n)
	var a int
	for i := range rndBytes{
		a = RndInt(len(letterArray))
		rndBytes[i] = letterArray[a]
	}
	return string(rndBytes)
}

/*Retorna um inteiro aleatório*/
func RndInt(interval int) int {
	rSource := rand.NewSource(time.Now().UnixNano()) //Nova seed para que rand.Intn() não gere sempre a mesma sequência de números
	var r *rand.Rand = rand.New(rSource)

	return r.Intn(interval)
}

/*Retorna a conversão de source para um string codificada em base64*/
func Encode64(source []byte) string {
	return base64.StdEncoding.EncodeToString(source)
}

/*Retorna []byte correspondente à source descodidicada em base64*/
func Decode64(source string) (b []byte) {
	b, _ = base64.StdEncoding.DecodeString(source)
	return b
}
