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
        IP:   net.ParseIP("::1"),
    }

    // 1. Tworzymy socket UDP
    // DialUDP w UDP NIE tworzy połączenia (bo UDP jest bezpołączeniowe),
    // ale ustawia domyślny adres docelowy, dzięki czemu można używać conn.Write().
    // Drugi argument (laddr) to adres lokalny — jeśli nil, system przydzieli port efemeryczny.

    conn, err := net.DialUDP("udp6", nil, &serverAddr)
    if err != nil {
        log.Fatalf("[UDP IPv6: Klient]: Błąd DialUDP: %v", err)
    }
    defer conn.Close()

    fmt.Println("[UDP IPv6: Klient]: Połączono z serwerem UDP.")

    // ============================================================
    //  Wypisanie portu źródłowego klienta UDP
    // ============================================================
    fmt.Printf("[UDP IPv6: Klient]: Mój adres lokalny (IP:port): %s\n", conn.LocalAddr())
    // LocalAddr() zwraca adres gniazda klienta, np. 127.0.0.1:56322
    // Port po prawej stronie to port źródłowy klienta UDP
    // Jest on przydzielany automatycznie przez system operacyjny
    // z zakresu portów efemerycznych (tymczasowych)
    // ============================================================


    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Wpisz wiadomość: ")
        scanner.Scan()
        text := scanner.Text()

        if text == "quit" {
            fmt.Println("[UDP IPv6: Klient]: Zamykanie klienta UDP.")
            break
        }

        // 2. Wysyłamy datagram
        // W UDP nie ma połączenia — to tylko wysłanie pakietu do serwera
        _, err := conn.Write([]byte(text))
        if err != nil {
            fmt.Printf("[UDP IPv6: Klient]: Błąd przy wysyłaniu: %v\n", err)
            break
        }

        // 3. Odbieramy odpowiedź
        // ReadFromUDP zwraca:
        // - liczbę bajtów
        // - adres nadawcy (serwera)
        // - błąd
        buf := make([]byte, 1024)
        n, _, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Printf("[UDP IPv6: Klient]: Błąd przy odbiorze odpowiedzi: %v\n", err)
            break
        }

        fmt.Printf("[UDP IPv6: Klient - Odpowiedź serwera]: %s\n", string(buf[:n]))
    }

    fmt.Println("[UDP IPv6: Klient]: Połączenie zakończone.")
}
