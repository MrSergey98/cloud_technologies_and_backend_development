package main

import (
    "fmt"
    "github.com/MrOganes/GO/mymodule"
)

func main() {
    message := mymodule.Hello("World")
    fmt.Println(message)
}
