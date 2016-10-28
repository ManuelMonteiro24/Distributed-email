package main

import "fmt"
import "os"


func main() {
    fmt.Printf("Hello welcome to the Distributed email\n")

    //Authentication????

        for {
          fmt.Printf("Menu\n 1.New email 2.Check your inbox 3. Exit\n")
            menu()
        }
}
//Menu inicial ?? adaptar ao metodo de autenticação do utilizador
func menu(){
  var input int
  n, err := fmt.Scanln(&input)
  if n < 1 || n > 2 || err != nil {
    fmt.Println("Invalid input\n")
    return
  }

  switch input {
        case 1: save_info := Dest_information()
                fmt.Println(save_info) 
        case 2:
                fmt.Println("Chamar funcao check inbox")
        case 3:
                fmt.Println("Goodbye thank you for choosing us")
                os.Exit(0)
        default:
                fmt.Println("Invalid Input")
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
//Meter aqui funcao do amado que recebe o texto que o utilizador quer enviar
//e introduz num ficheiro para enviar???
