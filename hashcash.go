package main

import(
	"strings"
	"strconv"
	"crypto/sha1"
	"math"
//	"fmt"
	//"bytes"
)

/*Returns header hashcode obtained with SHA1*/
func HashDigest(header string) []byte {
	hasher := sha1.New()
	hasher.Write([]byte(header))
	return hasher.Sum(nil)
}

/*Increments header counter*/
func (h *Header) IncCounter() {
	h.Counter++
}

/*Checks if first zeroBits of b are zero
 *Returns true if condition is checked and false otherwise.
*/
func CheckZeroBits(b []byte, zeroBits int) bool {
	var counter int = 0
	var n int = 8 - zeroBits

	for i := range b {
		if b[i] == 0x00 {
			counter = counter + 8
		} else if float64(b[i]) < math.Pow(float64(2), float64(n)) {
			counter = counter + 8-n
		}
		//Break condition
		if counter >= zeroBits {
			break
		}
		n = 8 - (zeroBits - counter) //número de bits dentro do byte
	}
	if counter >= zeroBits {
		return true
	} else {
		return false
	}

}

/*Returns a string with header information.
 *Header layout in struts.go
 */
func (h Header) HeaderToString() (string) {
	hElem := make([]string, 6)

	hElem[0] = "1" //SHA version
	hElem[1] = strconv.Itoa(h.ZeroCount)
	hElem[2] = h.Date
	hElem[3] = h.Resource
	hElem[4] = h.RandString
	hElem[5] = Encode64([]byte(strconv.Itoa(h.Counter)))

	return strings.Join(hElem, ":")
}

/*Fills in h *Header struct with header string information*/
func (h *Header) StringToHeader(header string) {
	hElem := strings.Split(header, ":")

	h.ZeroCount, _ = strconv.Atoi(hElem[1])
	h.Date = hElem[2]
	h.Resource = hElem[3]
	h.RandString = string(Decode64(hElem[4]))
	h.Counter, _ = strconv.Atoi(string(Decode64(hElem[5])))
}
