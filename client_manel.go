package main

import (
  "crypto/rand"
  "crypto/rsa"
  "fmt"
  "os"
  "encoding/pem"
  "io/ioutil"
  "crypto/x509"
)


func main() {
    fmt.Printf("Hello welcome to the Distributed email\n")
     user_name := auth_user()
    fmt.Println("User: ", user_name)
        for {
          fmt.Printf("Main Menu\n 1.New email 2.Check your inbox 3. Exit\n")
            menu()
        }
}

func auth_user()(string){
  var user_name = data_input("Insert Username: ")
  // detect if key file associated to user exists
	 _, err := os.Stat(user_name)

	  if os.IsNotExist(err) {
    //create file
		var user_file, err = os.Create(user_name)
		checkError(err)
    defer user_file.Close()

    // generate user key
    userKey, err := rsa.GenerateKey(rand.Reader, 2048)
    checkError(err)

    fmt.Println("User Key : ", userKey)

    //encode key to file
    pemdata := pem.EncodeToMemory(
        &pem.Block{
            Type: "RSA PRIVATE KEY",
            Bytes: x509.MarshalPKCS1PrivateKey(userKey),
        },
    )
    ioutil.WriteFile(user_name, pemdata, 0644)
		fmt.Printf("User keys generated!\n")

  }else{
    //
    fmt.Printf("User login successful!\n")
    file_data, err := ioutil.ReadFile(user_name)
    checkError(err)
    pemdata, _ := pem.Decode(file_data)
    userKey, err := x509.ParsePKCS1PrivateKey(pemdata.Bytes)
    checkError(err)

    fmt.Println("User Key : ", userKey)
  }
  return user_name
}

func menu(){
  var input int
  n, err := fmt.Scanln(&input)
  if n < 1 || n > 2 || err != nil {
    fmt.Println("Invalid input\n")
    return
  }
  switch input {
        case 1:
                fmt.Println("New email option chosen\n")
        case 2:
                fmt.Println("Check inbox option chosen\n")
        case 3:
                fmt.Println("Goodbye thank you for choosing us\n")
                os.Exit(0)
        default:
                fmt.Println("Invalid Input\n")
  }
}

func rcv_dest() {
    //recieve contact
    var recp_name = data_input("Insert Recipient: ")

    //check if contact exists
    var _, err = os.Stat(/contacts/recp_name)
  	if os.IsNotExist(err) {
      fmt.Println(err.Error())
  		os.Exit(3)
    } else{
      
    }
}

func data_input(msg string) (string){
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

func Dest_information() (dest string){
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

func checkError(err error) {
  	if err != nil {
  		fmt.Println(err.Error())
  		os.Exit(0)
  	}
}
