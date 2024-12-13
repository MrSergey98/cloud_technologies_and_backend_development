package mymodule

import "fmt"

// Hello возвращает приветственное сообщение.
func Hello(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
