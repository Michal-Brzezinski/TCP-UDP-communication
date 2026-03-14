package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
    addr := net.UDPAddr{
        Port: 9090,
        IP:   net.ParseIP("0.0.0.0"),
    }

    // 1. Tworzymy socket UDP
    conn, err := net.ListenUDP("udp4", &addr)
    if err != nil {
        log.Fatalf("Błąd ListenUDP: %v", err)
    }
    defer conn.Close()

    fmt.Println("Serwer UDP IPv4 działa na porcie 9090...")

    buf := make([]byte, 1024)

    for {
        // 2. Odbieramy datagram
        n, remoteAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            log.Printf("Błąd ReadFromUDP: %v", err)
            continue
        }

        msg := string(buf[:n])
        fmt.Printf("Odebrano od %s: %s\n", remoteAddr, msg)

        // 3. Odsyłamy odpowiedź
        conn.WriteToUDP([]byte("Serwer: "+msg), remoteAddr)
    }
}
