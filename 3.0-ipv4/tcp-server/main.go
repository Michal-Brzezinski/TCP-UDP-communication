package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
    // 1. Tworzymy gniazdo TCP nasłuchujące na porcie 8080
    listener, err := net.Listen("tcp4", ":8080")
    if err != nil {
        log.Fatalf("Błąd Listen: %v", err)
    }
    fmt.Println("Serwer TCP IPv4 działa na porcie 8080...")

    for {
        // 2. Czekamy na klienta
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Błąd Accept: %v", err)
            continue
        }

        // 3. Obsługa klienta w osobnej gorutynie
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()

    fmt.Printf("Połączono z %s\n", conn.RemoteAddr())

    reader := bufio.NewScanner(conn)

    for reader.Scan() {
        msg := reader.Text()
        fmt.Printf("Odebrano: %s\n", msg)

        // 4. Echo — odsyłamy to samo
        conn.Write([]byte("Serwer: " + msg + "\n"))
    }

    fmt.Printf("Rozłączono: %s\n", conn.RemoteAddr())
}

/*	===================================================

CO DZIEJE SIĘ W TYM KODZIE?:

- net.Listen("tcp4", ":8080") — wymusza IPv4.

- Accept() — blokuje do momentu połączenia.

- go handleClient() — każdy klient w osobnej gorutynie.

- bufio.Scanner — wygodne czytanie linii.

- Serwer odsyła echo — idealne do testów.

=======================================================*/