package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
    addr := net.UDPAddr{
        Port: 9090,
        IP:   net.ParseIP("0.0.0.0"),
    }

    // 1. Tworzymy socket UDP
    conn, err := net.ListenUDP("udp4", &addr)
    if err != nil {
        log.Fatalf("[UDP IPv4: Serwer]: Błąd ListenUDP: %v", err)
    }
    defer conn.Close()

    fmt.Println("[UDP IPv4: Serwer]: Serwer UDP IPv4 działa na porcie 9090...")

    // ============================================================
    //  Mechanizm eleganckiego zamknięcia serwera (Ctrl+C)
    // ============================================================
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    // Tworzony jest kanał (chan), do którego będą wysyłane sygnały z systemu operacyjnego.

    // - os.Signal oznacza typ danych reprezentujący sygnał systemowy.
    // - 1 oznacza bufor kanału – może przechować jeden sygnał.

    
    go func() {
        <-stop
        fmt.Println("\n[UDP IPv4: Serwer]: Otrzymano sygnał zakończenia. Zamykam socket UDP...")
        conn.Close()
        os.Exit(0)
    }()
    // ============================================================

    // Mapa przechowująca adresy klientów, którzy już wysłali datagram
    // Map w Go to w rzeczywistości hash table.
    knownClients := make(map[string]bool)

    buf := make([]byte, 1024)

    for {
        // 2. Odbieramy datagram
        n, remoteAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Printf("[UDP IPv4: Serwer]: Błąd ReadFromUDP: %v", err)
            continue
        }

        // remoteAddr.String() → "IP:port"
        addrStr := remoteAddr.String()

        // Jeśli klient pierwszy raz wysłał datagram — informujemy o tym
        if !knownClients[addrStr] {
            fmt.Printf("[UDP IPv4: Serwer]: Nowy klient: %s\n", addrStr)
            knownClients[addrStr] = true
        }

        msg := string(buf[:n])
        fmt.Printf("[UDP IPv4: Serwer]: Odebrano od %s: %s\n", remoteAddr, msg)

        // 3. Odsyłamy odpowiedź
        conn.WriteToUDP([]byte("[UDP IPv4: Serwer]: Serwer odsyła odpowiedź: "+msg), remoteAddr)
    }
}
