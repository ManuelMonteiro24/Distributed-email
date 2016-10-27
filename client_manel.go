package main

import "fmt"
import "os"

func main() {
    fmt.Printf("Hello welcome to the Distributed email\n")
        for {
            fmt.Printf("Menu\n 1. New email\n 2. Check your inbox\n 3. Exit\n Enter Number:\n")
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
      fmt.Println("Chamar funcao criar mail")
  case 2:
        fmt.Println("Chamar funcao check inbox")
  case 3:
      fmt.Printf("Hello welcome to the Distributed email\n")
      os.Exit(0)
  default:
      fmt.Println("Invalid Input\n")
  }

}
