package main

import (
	"fmt"

	"github.com/MrSergey98/cloud_technologies_and_backend_development/mymodule"
)

func main() {
	message := mymodule.Hello("World")
	fmt.Println(message)
}
