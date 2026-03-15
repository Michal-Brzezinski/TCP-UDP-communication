package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
    // 1. Tworzymy gniazdo TCP nasłuchujące na porcie 8080
    listener, err := net.Listen("tcp4", ":8080")

    // net.Listen("tcp4", ...") — wymusza IPv4.
    // Drugi argument, ":8080", to tzw. adres sieciowy w formacie "host:port"
    // jeżeli host jest pusty (tak jak tutaj) to Go interpretuje to jako 
    // „wszystkie interfejsy sieciowe” co znaczy, że
    // będzie nasłuchiwał na: 0.0.0.0:8080 (symbolicznie: „wszędzie”),

    // czyli będzie dostępny jako:

    //     127.0.0.1:8080 (z tej samej maszyny),

    //     192.168.0.5:8080 (z innego komputera w sieci lokalnej),

    //     itd., zależnie od konfiguracji.

    if err != nil {
        log.Fatalf("[TCP IPv4: Serwer]: Błąd Listen: %v", err)
    }
    fmt.Println("[TCP IPv4: Serwer]: Serwer TCP IPv4 działa na porcie 8080...")

    // ============================================
    //  Mechanizm eleganckiego zamknięcia serwera
    // ============================================
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    // Tworzony jest kanał (chan), do którego będą wysyłane sygnały z systemu operacyjnego.

    // - os.Signal oznacza typ danych reprezentujący sygnał systemowy.
    // - 1 oznacza bufor kanału – może przechować jeden sygnał.

    // Funkcja make w Go służy do tworzenia i inicjalizacji trzech specjalnych typów danych:

    // slice – dynamiczna lista elementów
    // 
    // map – tablica asocjacyjna (klucz -> wartość)
    // 
    // chan(nel) – mechanizm komunikacji między goroutines
    // 
    // I tylko tych trzech. Nie można jej użyć do structów, intów itd.
    // Najpierw alokuje pamięć, inicjalizuje wewnętrzną strukturę i zwraca gotowy 
    // do użycia obiekt (nie wskaźnik tak jak new())


    go func() {
        <-stop
        fmt.Println("\n[TCP IPv4: Serwer]: Otrzymano sygnał zakończenia. Zamykam listener...")
        listener.Close() // zamyka nasłuch — Accept() zwróci błąd
        os.Exit(0)
    }()
    
    // ============================================

    for {
        // 2. Czekamy na klienta
        conn, err := listener.Accept()
        // Accept() blokuje, dopóki jakiś klient się nie połączy.

        if err != nil {
            log.Printf("[TCP IPv4: Serwer]: Błąd Accept: %v", err)
            continue
        }

        // 3. Obsługa klienta w osobnej gorutynie
        go handleClient(conn)
        // Uruchamia handleClient w osobnej gorutynie (lekki wątek Go).
        // Dzięki temu serwer może obsługiwać wielu klientów równocześnie — każdy klient w osobnej gorutynie.

        // Odróżnienie klientów odbywa się poprzez porty źródłowe
        // każdy klient TCP, kiedy wysyła pakiet, musi mieć przypisany port źródłowy
    }
}

func handleClient(conn net.Conn) {
    // defer oznacza: „wykonaj to na końcu funkcji”.
    defer conn.Close()
    // - Gdy handleClient się zakończy (np. klient się rozłączy, wystąpi błąd), gniazdo zostanie zamknięte.
    // - Zamknięcie gniazda wysyła do klienta sygnał, że nie będzie więcej danych (EOF).

    fmt.Printf("[TCP IPv4: Serwer]: Połączono z %s\n", conn.RemoteAddr())
    // RemoteAddr() – adres klienta (IP:port).

    reader := bufio.NewScanner(conn)
    // Tworzymy Scanner nad połączeniem:
    // Scanner czyta dane z conn linia po linii (domyślnie do \n).

    for reader.Scan() {
        /* reader.Scan():

            - Zwraca true, gdy udało się przeczytać kolejną linię.

            - Zwraca false, gdy:

                - nastąpił EOF (klient zamknął połączenie),

                - albo wystąpił błąd.
        */

        msg := reader.Text()
        // reader.Text() – zwraca tekst ostatnio wczytanej linii (bez \n).
        
        fmt.Printf("[TCP IPv4: Serwer]: Odebrano od %s: %s\n", conn.RemoteAddr(), msg)
        
        // 4. Echo — odsyłamy to samo
        conn.Write([]byte("[TCP IPv4: Serwer odsyła odebraną wiadomość]: " + msg + "\n"))
        // conn.Write(...) – wysyłamy echo z dopisanym tekstem i \n

        // Gdy klient zamknie połączenie, Scan() przestanie zwracać true → wychodzimy z pętli.
    }

    if err := reader.Err(); err != nil {
        fmt.Printf("[TCP IPv4: Serwer]: Błąd przy czytaniu od %s: %v\n", conn.RemoteAddr(), err)
    }

    fmt.Printf("[TCP IPv4: Serwer]: Rozłączono: %s\n", conn.RemoteAddr())

    /*
    Po wyjściu z pętli:

    1. Wypisujemy info o rozłączeniu - jak powyżej
    2. Funkcja się kończy → defer conn.Close() zamyka gniazdo (jeśli jeszcze nie było zamknięte).
    */
}
