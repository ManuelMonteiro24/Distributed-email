package main

import "fmt"
import "os"

func main() {
    fmt.Printf("Hello welcome to the Distributed email\n")
        for {
          fmt.Printf("Menu\n 1.New email 2.Check your inbox 3. Exit\n Enter Number:\n")
            menu()
        }
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
      mail_information()
  case 2:
        fmt.Println("Chamar funcao check inbox\n")
  case 3:
      fmt.Println("Goodbye thank you for choosing us\n")
      os.Exit(0)
  default:
      fmt.Println("Invalid Input\n")
  }

}

func mail_information(){
  fmt.Println("Email information\n")
  var i = 1
  for i > 0{
    fmt.Println("From: \n")
    var input string
    n, err := fmt.Scanln(&input)
    if n < 1 || n > 2 || err != nil {
      fmt.Println("Invalid input\n")
      return
    } else {
      i = 0
      fmt.Println(input)
    }
  }
}
