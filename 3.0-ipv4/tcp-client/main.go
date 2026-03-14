package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
    // 1. Łączymy się z serwerem
    conn, err := net.Dial("tcp4", "127.0.0.1:8080")
    if err != nil {
        log.Fatalf("Błąd Dial: %v", err)
    }
    defer conn.Close()

    fmt.Println("Połączono z serwerem TCP.")

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Wpisz wiadomość: ")
        scanner.Scan()
        text := scanner.Text()

        // 2. Wysyłamy wiadomość
        conn.Write([]byte(text + "\n"))

        // 3. Odbieramy odpowiedź
        reply := bufio.NewReader(conn)
        resp, _ := reply.ReadString('\n')
        fmt.Printf("Odpowiedź: %s", resp)
    }
}