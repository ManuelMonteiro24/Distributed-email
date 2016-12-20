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
            menu(user_name)
        }
}

func auth_user()(string){
  var user_name = data_input("Insert Username: ")
  // detect if key file associated to user exists
	 _, err := os.Stat(user_name + "_PrivateKey")

	  if os.IsNotExist(err) {
    //create user key file
		var user_file, err = os.Create(user_name + "_PrivateKey")
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
    ioutil.WriteFile(user_name + "_PrivateKey", pemdata, 0644)

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

  }else{
    //load user key from file
    fmt.Printf("User login successful!\n")
    file_data, err := ioutil.ReadFile(user_name + "_PrivateKey")
    checkError(err)
    pemdata, _ := pem.Decode(file_data)
    userKey, err := x509.ParsePKCS1PrivateKey(pemdata.Bytes)
    checkError(err)

    fmt.Println("User Key : ", userKey)
  }
  return user_name
}

func menu(user_name string){
  var input int
  n, err := fmt.Scanln(&input)
  if n < 1 || n > 2 || err != nil {
    fmt.Println("Invalid input\n")
    return
  }
  switch input {
        case 1:
                fmt.Println("New email option chosen\n")
                dest_name := rcv_dest()
                file_name := file_to_send()
                fmt.Println("Sender:", user_name, "\nDest :", dest_name, "\nFile: ", file_name,"\n")

        case 2:
                fmt.Println("Check inbox option chosen\n")
        case 3:
                fmt.Println("Goodbye thank you for choosing us\n")
                os.Exit(0)
        default:
                fmt.Println("Invalid Input\n")
  }
}

func rcv_dest()(string){
    //recieve name of contact
    var recp_name = data_input("Insert Recipient: ")

    //Check if contact exists
    _, err := os.Stat(recp_name + "_PublicKey")
  	if os.IsNotExist(err) {
      fmt.Println(err.Error())
  		os.Exit(3)
    } else{
      fmt.Println("Contact found!\n")
    }
    return recp_name
}

func file_to_send()(string){
    //recieve name of file
    var file_name = data_input("Insert name of file to send: ")

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
