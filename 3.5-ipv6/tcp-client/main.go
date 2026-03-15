package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

    // Łączymy się z serwerem
    conn, err := net.Dial("tcp6", "[::1]:8080")
    // "tcp6" oznacza TCP dla IPv6.
    // ::1 to odpowiednik 127.0.0.1 w IPv6 (loopback).
    // Znowu: [::1]:8080 — nawiasy są obowiązkowe.

    if err != nil {
        log.Fatalf("[TCP IPv6: Klient]: Błąd Dial: %v", err)
    }
    defer conn.Close()

    // Typowy schemat w Go: 
    // 1. Najpierw tworzysz zasób (połączenie, plik, itp.).
    // 2. Sprawdzasz, czy się udało.
    // 3. Jeśli tak – od razu deklarujesz, że na końcu funkcji chcesz go zamknąć (defer).

    // ============================================
    //  Wypisanie portu źródłowego klienta
    // ============================================
    fmt.Printf("[TCP IPv6: Klient]: Mój adres lokalny (IP:port): %s\n", conn.LocalAddr())
    // LocalAddr() zwraca adres gniazda klienta, np. [::1]:56322
    // Port po prawej stronie to port źródłowy klienta TCP
    // Jest on przydzielany automatycznie przez system operacyjny
    // z zakresu tzw. portów efemerycznych (tymczasowych)
    // ============================================

    fmt.Println("[TCP IPv6: Klient]: Połączono z serwerem TCP.")
    fmt.Println("Wpisz wiadomość (lub 'quit' aby zakończyć):")

    scanner := bufio.NewScanner(os.Stdin)

    for {
        // Pętla nieskończona

        fmt.Print("> ")
        if !scanner.Scan() {
            // EOF na stdin (np. Ctrl+D) albo błąd
            fmt.Println("\n[TCP IPv6: Klient]: Koniec wejścia. Zamykanie połączenia.")
            break
        }
        text := scanner.Text()

        if text == "quit" {
            fmt.Println("[TCP IPv6: Klient]: Zamykanie połączenia na życzenie użytkownika.")
            break
        }

        // Wysyłamy wiadomość
        _, err := conn.Write([]byte(text + "\n"))
        if err != nil {
            fmt.Printf("[TCP IPv6: Klient]: Błąd przy wysyłaniu: %v\n", err)
            break
        }

        // Odbieramy odpowiedź
        reply := bufio.NewReader(conn)
        resp, err := reply.ReadString('\n')
        if err != nil {
            fmt.Printf("[TCP IPv6: Klient]: Błąd przy odbiorze odpowiedzi (serwer mógł się rozłączyć): %v\n", err)
            // Jeśli serwer został zamknięty, klient kończy działanie
            break
        }

        fmt.Printf("[TCP IPv6: Klient - Odpowiedź serwera]: %s", resp)
    }

    // defer conn.Close() zamknie połączenie tutaj
    fmt.Println("[TCP IPv6: Klient]: Połączenie zakończone.")
}
