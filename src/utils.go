package main

import (
	"time"
	//"fmt"
	"math/rand"
	"strconv"
	"strings"
	"encoding/base64"
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
func BuildHeader(to string) string {
	var zeroBits float64 = 20
	var currentTime time.Time = time.Now()
	var rndString string = RandString(6)//6 bytes para poder ser codificada em base64 sem precisar de padding
	var headerArray = make([]string, 5)

	y, m, d := currentTime.Date()//ano, mês, dia

	headerArray[0] = "1" //SHA version
	headerArray[1] = strconv.Itoa(int(zeroBits)) //número de pre-image 0 bits para calculo da hash
	headerArray[2] =  strconv.Itoa(d)  + strconv.Itoa(int(m)) + strings.TrimPrefix(strconv.Itoa(y), "20")
	headerArray[3] = to
	headerArray[4] = Encode64([]byte(rndString))
	headerArray[5] = strconv.Itoa(RndInt())

	return strings.Join(headerArray, ":")
}

/*Retorna uma string gerada aleatoriamente de n bytes*/
const letterArray = "absdefghijklmnopqrstuvyxwzABCDEFGHIJKLMNOPQRSTUVYXWZ"
func RandString(n int) string {
	rndBytes := make([]byte, n)
	for i := range rndBytes {
		rndBytes[i] = letterArray[RndInt(len(letterArray))]
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