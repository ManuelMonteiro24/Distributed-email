package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"time"
	"strconv"
)

const COUNTER_INTERVAL = 1000000

func main() {
    fmt.Printf("\nHello welcome to the Distributed email\n")
    user_name, userKey:= AuthUser()
    fmt.Println("User_name: ", user_name)
    fmt.Println("User_key: ", userKey)

    for {
        fmt.Printf("Main Menu\n 1.New email 2.Check your inbox 3. Exit\n")
        Menu(user_name, userKey)
    }


}

/*Creates new Mail struct and fills its fields
with the provided information; ready to be sent*/
func NewMail(userKey *rsa.PrivateKey, toPublicKey *rsa.PublicKey, message string, to string, from string) (mail Mail){
	//create symmetric encryption key for AES encryption of payload
	var h Header
	var pow []byte
	t := time.Now()

	/*Fill header with mail specific info*/
	h.ZeroCount = 20
	h.Date = CreateDate(t)
	h.From = from
	rsrc, err := x509.MarshalPKIXPublicKey(toPublicKey)
	checkError(err)
	h.Resource = Encode64(rsrc)//Para ler de volta, usar -- bytes, err, := x509.ParsePKIXPublicKey(Decode64(h.Resource))
	h.RandString = RndStr(12) //random 12 byte string
	h.Counter = RndInt(COUNTER_INTERVAL)

	/*Hashcash to create valid proof_of_work*/
	header := h.HeaderToString() //header in string format
	pow = HashDigest(header)

	for !CheckZeroBits(pow, h.ZeroCount)  {
		h.IncCounter()//increment counter to create new digest
		header = h.HeaderToString()
		pow = HashDigest(header)
	}
	mail.Header = header
	mail.Proof_of_Work = Encode64(pow)

	/*Generate payload string with format "PoW//\\TO//\\FROM//\\MESSAGE//\\SIGNATURE" */
	payloadArray := []string{mail.Proof_of_Work, to, h.From, message}
	payload := strings.Join(payloadArray, "//\\\\")

	/*Sign payload*/
	payloadSignature := Sign(HashDigest(payload), userKey)

	/*Encrypt (payload + signature) with symKey*/
	symKey := GenSymKey()
	mail.Payload = SymEncrypt(symKey, payload + "//\\\\" + payloadSignature) //Payload encrypted with symKey, encoded in Base64

	/*Encrypt symKey with recipient toPublicKey*/
	encSymKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, toPublicKey, symKey, []byte(""))
	checkError(err)
	mail.SymKey = Encode64(encSymKey)

	return mail
}

/*Reads Mail using user PrivateKey and sender PublicKey*/
func ReadMail(userKey *rsa.PrivateKey, senderKey *rsa.PublicKey, mail Mail) {

	/*Integrity check*/
	if mail.Proof_of_Work != Encode64(HashDigest(mail.Header)){
		fmt.Print("Mail may be spam or was illicitly altered!! (pow test fail)\n")
		return
	}
	var header Header
	header.StringToHeader(mail.Header)

	/*Check if user PublicKey is equal to header.Resource(TO) */
	pub, err := x509.MarshalPKIXPublicKey(&userKey.PublicKey)
	checkError(err)
	if header.Resource == Encode64(pub) {
		fmt.Print("Header.Resource == userKey.PublicKey\n")
	}

	/*Retrieve symKey*/
	symKey, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, userKey, Decode64(mail.SymKey), []byte(""))
	checkError(err)

	/*Retrieve payload*/
	payload := SymDecrypt(symKey, mail.Payload)
	plParts := strings.Split(payload, "//\\\\")

	/*Check for valid sender signature*/
	d := HashDigest(strings.Replace(payload, "//\\\\" + plParts[4], "", -1))
	if !CheckSign(d, []byte(Decode64(plParts[4])), senderKey) {
		fmt.Print("Mail illicitly altered or corrupt!! (signature mismatch)")
		return
	}

	fmt.Print("\nMESSAGE\nFROM: ", plParts[2], "\nTO: ", plParts[1], "\nCONTENT: ", plParts[3], "\n")

}

/*Signs digest with priv private key
Returns signature in string format, encoded in Base64
 */
func Sign(digest []byte, priv *rsa.PrivateKey) string{
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto //for simplicity
	signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA1, digest, &opts)
	checkError(err)
	return Encode64(signature)
}

/*Check signature with senderKey and digest
 *Returns true if signature is valid and false otherwise
 */
func CheckSign(digest []byte, signature []byte, senderKey *rsa.PublicKey) bool{
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	if rsa.VerifyPSS(senderKey, crypto.SHA1, digest, signature, &opts) == nil {
		return true
	}else {
		return false
	}
}

/*Create a date string with current date of the form DDMMYYHHmm */
func CreateDate(t time.Time) (date string) {

	DD := t.Day()
	MM := int(t.Month())

	day := strconv.Itoa(DD)
	month := strconv.Itoa(MM)
	if DD < 10 {
		day = "0" + day
	}
	if MM < 10 {
		month = "0" + month
	}

	hh, mm, _ := t.Clock()
	hours := strconv.Itoa(hh)
	minutes := strconv.Itoa(mm)
	if hh < 10 {
		hours = "0" + hours
	}
	if mm < 10 {
		minutes = "0" + minutes
	}
	date = day + month + strconv.Itoa(t.Year()-2000) + hours + minutes
	return date
}

func AuthUser() (string, *rsa.PrivateKey) {
        var user_name = DataInput("Insert Username: ")

        var userPriv = user_name + "_PrivateKey"
	 _, err := os.Stat(userPriv)

        // detect if key file associated to user exists
        if os.IsNotExist(err) {
	//create user key file
		var user_file, err = os.Create(userPriv)
		checkError(err)
		defer user_file.Close()

                // generate user key
                userKey, err := rsa.GenerateKey(rand.Reader, 2048)
                checkError(err)

                //encode key to file
                pemdata := pem.EncodeToMemory(
                &pem.Block{
                        Type: "RSA PRIVATE KEY",
                        Bytes: x509.MarshalPKCS1PrivateKey(userKey),
                },
                )
                ioutil.WriteFile(userPriv, pemdata, 0644)

                //create user PublicKey contact file
                PublicKey, err := x509.MarshalPKIXPublicKey(&userKey.PublicKey)
                checkError(err)

                pemdata = pem.EncodeToMemory(
                &pem.Block{
                        Type: "RSA PUBLIC KEY",
                        Bytes: PublicKey,
                },
                )

                var contact_name string = user_name + "_PublicKey"
                contact_file, err := os.Create(contact_name)
		checkError(err)
                defer contact_file.Close()

                ioutil.WriteFile(contact_name, pemdata, 0644)

		fmt.Printf("User keys generated!\n")
                return user_name, userKey
        }else{
        //load user key from file
                fmt.Printf("User login successful!\n")
                file_data, err := ioutil.ReadFile(userPriv)
                checkError(err)
                pemdata, _ := pem.Decode(file_data)
                userKey, err := x509.ParsePKCS1PrivateKey(pemdata.Bytes)
                checkError(err)

                return user_name, userKey
        }

}

func Menu(user_name string, userKey *rsa.PrivateKey){

	var input int
	n, err := fmt.Scanln(&input)
	if n < 1 || n > 2 || err != nil {
		fmt.Println("Invalid input\n")
		return
	}
	switch input {
	case 1:
		fmt.Println("New email option chosen\n")
		dest_name, dest_PublicKey := RcvDest()
		fmt.Println("Dest PublicKey : ", dest_PublicKey)
		file_name := FileToSend()
		fmt.Println("Sender:", user_name, "\nDest :", dest_name, "\nFile: ", file_name,"\n")
		message := ReadFile(file_name)
		mail := NewMail(userKey, dest_PublicKey, message, dest_name, user_name)

		fmt.Println()

		WriteJSON(mail)

	case 2:
		fmt.Println("Check inbox option chosen\n")
		rcvMail := ReadJSON("mail_ready.json")
		_, pub := RcvDest()
		ReadMail(userKey, pub, rcvMail)
	case 3:
		fmt.Println("Goodbye thank you for choosing us\n")
		os.Exit(0)
	default:
		fmt.Println("Invalid Input\n")
	}
}

func RcvDest() (string, *rsa.PublicKey){
    //receive name of contact
    var recp_name = DataInput("Insert Recipient: ")

    //Check if contact exists
	_, err := os.Stat(recp_name + "_PublicKey")
	if os.IsNotExist(err) {
		fmt.Printf("No such contact in Contacts list!\n")
		os.Exit(3)
	}

	fmt.Println("Contact found!\n")
	file_data, err := ioutil.ReadFile(recp_name + "_PublicKey")
	checkError(err)
	pemdata, _ := pem.Decode(file_data)
	pub, err := x509.ParsePKIXPublicKey(pemdata.Bytes)
	checkError(err)

	return recp_name, pub.(*rsa.PublicKey)
}

func FileToSend() (string){
	//recieve name of file
	var file_name = DataInput("Insert name of file to send: ")

	//Check if file exists
	_, err := os.Stat(file_name)
	if os.IsNotExist(err) {
		fmt.Println(err.Error())
		os.Exit(3)
	} else{
		fmt.Println("File found!\n")
	}
	return file_name
}

func DataInput(msg string) (string){
  for{
    var data string
      fmt.Println(msg)
    n, err := fmt.Scanln(&data)
      if n < 1 || n > 2 || err != nil {
        fmt.Println("Invalid input\n")
      }else{
        return data
      }
  }
}

func DestInformation() (dest string){
  fmt.Println("Email information")
  for {
    fmt.Println("To: ")
    var input string
    n, err := fmt.Scanln(&input)
    if n < 1 || n > 2 || err != nil {
      fmt.Println("Invalid input")
    } else {
      return input
      fmt.Println(input)
      }
    }
}
