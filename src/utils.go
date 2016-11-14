package main

import (
	"time"
	"fmt"
	"math/rand"
	//"reflect"
)

/*Verificar se existe algum erro*/
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

/*Função para construir a string do field Header
da estrutura Mail. Layout em strut.go
 */
func BuildHeader(to string) (header string) {
	//var zeroBits int = 20
	var currentTime time.Time  = time.Now()


	fmt.Print(currentTime.Date())
	return header
}

/*Retorna uma string de n bytes*/
const letterArray = "absdefghijklmnopqrstuvyxwzABCDEFGHIJKLMNOPQRSTUVYXWZ"
func RandString(n int) string {
	rSource := rand.NewSource(time.Now().UnixNano()) //Nova seed para que rand.Intn() não gere sempre a mesma sequência de números
	var r *rand.Rand = rand.New(rSource)

	rndBytes := make([]byte, n)
	for i := range rndBytes {
		rndBytes[i] = letterArray[r.Intn(len(letterArray))]
	}
	return string(rndBytes)
}