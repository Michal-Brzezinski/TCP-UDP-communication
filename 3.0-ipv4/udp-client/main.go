package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
    serverAddr := net.UDPAddr{
        Port: 9090,
        IP:   net.ParseIP("127.0.0.1"),
    }

    // 1. Tworzymy socket UDP
    conn, err := net.DialUDP("udp4", nil, &serverAddr)
    if err != nil {
        log.Fatalf("Błąd DialUDP: %v", err)
    }
    defer conn.Close()

    fmt.Println("Połączono z serwerem UDP.")

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Wpisz wiadomość: ")
        scanner.Scan()
        text := scanner.Text()

        // 2. Wysyłamy datagram
        conn.Write([]byte(text))

        // 3. Odbieramy odpowiedź
        buf := make([]byte, 1024)
        n, _, _ := conn.ReadFromUDP(buf)
        fmt.Printf("Odpowiedź: %s\n", string(buf[:n]))
    }
}
