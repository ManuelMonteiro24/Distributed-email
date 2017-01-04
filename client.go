package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"distmail/kademlia"
	"encoding/pem"
	"fmt"
	b58 "github.com/jbenet/go-base58"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const COUNTER_INTERVAL = 1000000

//node ID is the SHA1 of its Public Key
func userID(pubkey *rsa.PublicKey) [20]byte {
	return sha1.Sum(kademlia.SerializePublicKey(pubkey))
}

func main() {

	fmt.Printf("\nHello welcome to the Distributed email\n")
	user_name, userKey := AuthUser()
	ID := userID(&userKey.PublicKey)

    fmt.Println("Insert IP and Port (type nothing to use STUN to find external IP and port)")
    IP := DataInput("IP: ")
    port := DataInput("Port: ")

    fmt.Println("Insert bootstrap IP and Port (type nothing to skip bootstrap step)")
    bIP := DataInput("bootstrapIP: ")
    bPort := DataInput("bootstrapPort: ")

	dht, _ := kademlia.InitDHT(ID[:], bIP, bPort, IP, port, userKey, ExtractToIDfromMail)

	for {
		fmt.Printf("Main Menu\n 1.New email 2.Check your inbox 3. Exit\n")
		Menu(user_name, userKey, dht)
	}
}

func ExtractToIDfromMail(data []byte) string {
	mail := ReadJSON(data)

	var header Header
	header.StringToHeader(mail.Header)

	return header.Resource
}

/*Creates new Mail struct and fills its fields
with the provided information; ready to be sent*/
func NewMail(userKey *rsa.PrivateKey, toPublicKey *rsa.PublicKey, message string, to string, from string) (mail Mail) {
	//create symmetric encryption key for AES encryption of payload
	var h Header
	var pow []byte
	t := time.Now()

	/*Fill header with mail specific info*/
	h.ZeroCount = 20
	h.Date = CreateDate(t)
	rsrc := userID(toPublicKey)
	h.Resource = b58.Encode(rsrc[:])
	h.RandString = RndStr(12) //random 12 byte string
	h.Counter = RndInt(COUNTER_INTERVAL)

	/*Hashcash to create valid proof_of_work*/
	header := h.HeaderToString() //header in string format
	pow = HashDigest(header)

	for !CheckZeroBits(pow, h.ZeroCount) {
		h.IncCounter() //increment counter to create new digest
		header = h.HeaderToString()
		pow = HashDigest(header)
	}
	mail.Header = header
	mail.Proof_of_Work = Encode64(pow)

	/*Generate payload string with format "PoW//\\TO//\\FROM//\\MESSAGE//\\SIGNATURE" */
	payloadArray := []string{mail.Proof_of_Work, to, from, message}
	payload := strings.Join(payloadArray, "//\\\\")

	/*Sign payload*/
	payloadSignature := Sign(HashDigest(payload), userKey)

	/*Encrypt (payload + signature) and from field with symKey*/
	symKey := GenSymKey()
	mail.From = SymEncrypt(symKey, from)
	mail.Payload = SymEncrypt(symKey, payload+"//\\\\"+payloadSignature) //Payload encrypted with symKey, encoded in Base64

	/*Encrypt symKey with recipient toPublicKey*/
	encSymKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, toPublicKey, symKey, []byte(""))
	checkError(err)
	mail.SymKey = Encode64(encSymKey)

	return mail
}

/*Reads Mail using user PrivateKey and sender PublicKey*/
func ReadMail(userKey *rsa.PrivateKey, mail Mail, dht *kademlia.DHT) {

	/*Integrity check*/
	if mail.Proof_of_Work != Encode64(HashDigest(mail.Header)) {
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
	from := SymDecrypt(symKey, mail.From)

	senderKey := dht.GetPubKeyByID(strings.Replace(from, "\n", "", -1))

	plParts := strings.Split(payload, "//\\\\")

	/*Check for valid sender signature*/
	d := HashDigest(strings.Replace(payload, "//\\\\"+plParts[4], "", -1))
	if !CheckSign(d, []byte(Decode64(plParts[4])), senderKey) {
		fmt.Print("Mail illicitly altered or corrupt!! (signature mismatch)")
		return
	}

	fmt.Print("\nMESSAGE\nFROM: ", from, "\nTO: ", plParts[1], "\nCONTENT: ", plParts[3], "\n")

}

/*Signs digest with priv private key
Returns signature in string format, encoded in Base64
*/
func Sign(digest []byte, priv *rsa.PrivateKey) string {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto //for simplicity
	signature, err := rsa.SignPSS(rand.Reader, priv, crypto.SHA1, digest, &opts)
	checkError(err)
	return Encode64(signature)
}

/*Check signature with senderKey and digest
 *Returns true if signature is valid and false otherwise
 */
func CheckSign(digest []byte, signature []byte, senderKey *rsa.PublicKey) bool {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	if rsa.VerifyPSS(senderKey, crypto.SHA1, digest, signature, &opts) == nil {
		return true
	} else {
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

	//check if Users dir exists if not creates one
	rootpath, _ := os.Getwd()
	subpath := filepath.Join(rootpath, "Users")
	err := os.MkdirAll(subpath, os.ModePerm)
	checkError(err)

	user_name := DataInput("Insert Username: \n")
	userPriv := filepath.Join(subpath, user_name+"_PrivateKey")
	_, err = os.Stat(userPriv)

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
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(userKey),
			},
		)

		ioutil.WriteFile(userPriv, pemdata, 0644)

		//create user PublicKey contact file
		PublicKey, err := x509.MarshalPKIXPublicKey(&userKey.PublicKey)
		checkError(err)

		pemdata = pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: PublicKey,
			},
		)

		userPub := filepath.Join(subpath, user_name+"_PublicKey")
		contact_file, err := os.Create(userPub)
		checkError(err)
		defer contact_file.Close()

		ioutil.WriteFile(userPub, pemdata, 0644)

		fmt.Printf("User keys generated!")
		return user_name, userKey

	} else {
		//load user key from file
		fmt.Println("User login successful!")
		file_data, err := ioutil.ReadFile(userPriv)
		checkError(err)
		pemdata, _ := pem.Decode(file_data)
		userKey, err := x509.ParsePKCS1PrivateKey(pemdata.Bytes)
		checkError(err)

		return user_name, userKey
	}
}

/*
func Pass_gen_or_check(user_name string){
	rcv_pass := DataInput("Insert Password: ")
}*/

func Sendit(user_name string, dest_name string, userKey *rsa.PrivateKey, dest_PublicKey *rsa.PublicKey, dht *kademlia.DHT) {
	file_name := FileToSend()
	fmt.Println("Sender:", user_name, "\nDest :", dest_name, "\nFile: ", file_name)
	message := ReadFile(file_name)

	fromID := userID(&userKey.PublicKey)
	mail := NewMail(userKey, dest_PublicKey, message, dest_name, b58.Encode(fromID[:]))

	fmt.Println()
	mailBytes := WriteJSON(mail)

	dht.SendEmail(mailBytes)
}

func Menu(user_name string, userKey *rsa.PrivateKey, dht *kademlia.DHT) {

	var input int
	n, err := fmt.Scanln(&input)
	if n < 1 || n > 2 || err != nil {
		fmt.Println("Invalid input")
		return
	}
	switch input {
	case 1:
		fmt.Println("New email option chosen:")
		dest_name, dest_PublicKey, exists := RcvDest()
		if exists {
			Sendit(user_name, dest_name, userKey, dest_PublicKey, dht)
		} else {
			fmt.Println("Contact not stored.")
			destID := DataInput("Insert dest ID:")
			if destID != "" {
				dest_PublicKey = dht.GetPubKeyByID(destID)
				if dest_PublicKey != nil {
					Sendit(user_name, dest_name, userKey, dest_PublicKey, dht)
				}

			}
		}

	case 2:
		fmt.Println("Check inbox option chosen")

		user_id := userID(&userKey.PublicKey)
		results, _ := dht.Lookup(b58.Encode(user_id[:]), 10)

		for _, result := range results {
			rcvMail := ReadJSON(result)
			ReadMail(userKey, rcvMail, dht)
		}

	case 3:
		fmt.Println("Goodbye thank you for choosing us")
		os.Exit(0)
	default:
		fmt.Println("Invalid Input")
	}
}

func RcvDest() (string, *rsa.PublicKey, bool) {

	//check if Contacts dir exists if not creates one
	rootpath, _ := os.Getwd()
	subpath := filepath.Join(rootpath, "Contacts")
	os.MkdirAll(subpath, os.ModePerm)

	//receive name of contact
	var recp_name = DataInput("Insert Recipient: ")

	//Check if contact exists
	subpath = filepath.Join(subpath, recp_name+"_PublicKey")
	fmt.Println(subpath)
	_, err := os.Stat(subpath)

	if os.IsNotExist(err) {
		fmt.Printf("No such contact in Contacts list!")
		return recp_name, nil, false
	}

	fmt.Println("Contact found!")
	file_data, err := ioutil.ReadFile(subpath)
	checkError(err)
	pemdata, _ := pem.Decode(file_data)
	pub, err := x509.ParsePKIXPublicKey(pemdata.Bytes)
	checkError(err)

	return recp_name, pub.(*rsa.PublicKey), true
}

func FileToSend() string {
	//recieve name of file
	var file_name = DataInput("Insert name of file to send: ")

	//Check if file exists
	_, err := os.Stat(file_name)
	if os.IsNotExist(err) {
		fmt.Println(err.Error())
		os.Exit(3)
	} else {
		fmt.Println("File found!")
	}
	return file_name
}

func DataInput(msg string) string {
	for {
		var data string
		fmt.Printf(msg)
		n, err := fmt.Scanln(&data)
		if n < 1 || n > 2 || err != nil {
			if n == 0 {
				return ""
			}
			fmt.Println("Invalid input")
		} else {
			return data
		}
	}
}

/*
func pwd(){
	//check if Users list file, with the passwords exists if not creates one
	_, err = os.Stat("users_list.txt")

	if os.IsNotExist(err){
		users_file, err := os.Create(users_list)
		checkError(err)
		defer users_list.Close()
	}
	flag := DataInput("Press 1 for: Login Or 2 for: Sign up\n")
	if(flag = 1 || flag =2){
		username := DataInput("Insert Username:\n")
		pwd := DataInput("Insert Password:\n")

		if(flag = 1){
			 //= HashDigest(pwd + salt)

		}
		if(flag = 2){
			salt := make([]byte, 32)
		  _, err := io.ReadFull(rand.Reader, salt)
		  checkError(err)

    	scanner := bufio.NewScanner(users_file)
    	for scanner.Scan() {
				plParts := strings.Split(scanner.Text(), "//\\\\")

    	}

    	if err := scanner.Err(); err != nil {
      	log.Fatal(err)
    	}
		}

	}else{

	}
}*/
