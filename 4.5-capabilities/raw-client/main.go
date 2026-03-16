package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"syscall"
)

func main() {

    // ============================================
    // 1. Tworzymy RAW SOCKET
    // ============================================
    //
    // syscall.Socket(domain, type, protocol)
    //
    // domain   = AF_INET  -> IPv4
    // type     = SOCK_RAW -> gniazdo surowe
    // protocol = IPPROTO_RAW -> pozwala nam samodzielnie tworzyć nagłówki IP
    //
    // UWAGA: to wymaga uprawnień roota (lub capabilities, o czym będzie w zadaniu 4.5)
    //
	// fd - deskryptor socketu
    fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
    if err != nil {
        log.Fatalf("Błąd tworzenia raw socket: %v", err)
    }
    defer syscall.Close(fd)

    fmt.Println("[RAW SOCKET]: Utworzono gniazdo surowe.")

    // ============================================
    // 2. Dane pakietu
    // ============================================

    sourceIP := net.ParseIP("127.0.0.1").To4()
    destIP := net.ParseIP("127.0.0.1").To4()

	// uint – liczba bez znaku, 16 bitów = 2 bajty
    sourcePort := uint16(54321)	// dowolny port źródłowy (nie musi być otwarty, bo to RAW SOCKET)
    destPort := uint16(9090) // port Twojego serwera UDP z zadania 3.0

    payload := []byte("Hello from RAW SOCKET!")	// zawartość payloadu można podać dowolną, ale niech będzie czytelna

    // ============================================
    // 3. Budujemy nagłówek UDP (8 bajtów)
    // ============================================

    udpHeader := make([]byte, 8)	// UDP header ma stałą długość 8 bajtów

    // Source Port
    binary.BigEndian.PutUint16(udpHeader[0:2], sourcePort)

    // Destination Port
    binary.BigEndian.PutUint16(udpHeader[2:4], destPort)

    // UDP Length = 8 bajtów nagłówka + długość payloadu
    udpLength := uint16(8 + len(payload))
    binary.BigEndian.PutUint16(udpHeader[4:6], udpLength)

    // Checksum na razie 0 — policzymy później
    binary.BigEndian.PutUint16(udpHeader[6:8], 0)

    // ============================================
    // 4. Budujemy nagłówek IPv4 (20 bajtów)
    // ============================================

    ipHeader := make([]byte, 20)

    // Version (4 bity) + IHL (4 bity)
    // Version = 4 (IPv4)
    // IHL = 5 (5 * 4 bajty = 20 bajtów nagłówka)
    ipHeader[0] = (4 << 4) | 5
	// 4 << 4 przesuwa 4 o 4 bity w lewo, czyli daje 0x40, a 5 to 0x05, więc razem 0x45

    // Type of Service (DSCP/ECN) — 0
	// Differentiated Services Code Point
	// Określa priorytet pakietu w sieci.
	// ECN Służy do informowania o przeciążeniu sieci bez gubienia pakietów.
	// W naszym przypadku 0 = brak priorytetu, normalny ruch

    ipHeader[1] = 0

    // Total Length = IP header (20) + UDP header (8) + payload
    totalLength := uint16(20 + 8 + len(payload))
    binary.BigEndian.PutUint16(ipHeader[2:4], totalLength)

    // Identification — dowolne
    binary.BigEndian.PutUint16(ipHeader[4:6], sourcePort)

    // Flags + Fragment Offset — 0 (nie fragmentujemy)
    binary.BigEndian.PutUint16(ipHeader[6:8], 0)

    // TTL — Time To Live = czas życia pakietu, ustawiamy na 64 (wystarczająco duża wartość, żeby dotarł do celu)
    ipHeader[8] = 64

    // Protocol — 17 = UDP
    ipHeader[9] = syscall.IPPROTO_UDP

    // Header checksum na razie 0 — policzymy później
    binary.BigEndian.PutUint16(ipHeader[10:12], 0)

    // Source IP
    copy(ipHeader[12:16], sourceIP)

    // Destination IP
    copy(ipHeader[16:20], destIP)

    // ============================================
    // 5. Liczymy checksum UDP (pseudo-header!)
    // ============================================

    udpChecksum := udpChecksumIPv4(sourceIP, destIP, udpHeader, payload)
    binary.BigEndian.PutUint16(udpHeader[6:8], udpChecksum)

    // ============================================
    // 6. Liczymy checksum IP
    // ============================================

    ipChecksum := checksum(ipHeader)
    binary.BigEndian.PutUint16(ipHeader[10:12], ipChecksum)

    // ============================================
    // 7. Składamy cały pakiet
    // ============================================

    packet := append(ipHeader, udpHeader...)
    packet = append(packet, payload...)

    // ============================================
    // 8. Wysyłamy pakiet przez raw socket
    // ============================================

    addr := syscall.SockaddrInet4{
        Port: int(destPort),
    }
    copy(addr.Addr[:], destIP)

    err = syscall.Sendto(fd, packet, 0, &addr)
    if err != nil {
        log.Fatalf("Błąd Sendto: %v", err)
    }

    fmt.Println("[RAW SOCKET]: Pakiet wysłany!")
    fmt.Println("Sprawdź Wireshark + serwer UDP z zadania 3.0.")
}

//
// ============================================================
// Funkcja checksum — klasyczna suma 16-bitowa z RFC 1071
// ============================================================
//
func checksum(data []byte) uint16 {
    sum := uint32(0)

    // Sumujemy po 16 bitów
    for i := 0; i < len(data)-1; i += 2 {
        sum += uint32(binary.BigEndian.Uint16(data[i : i+2]))
    }

    // Jeśli nieparzysta liczba bajtów
    if len(data)%2 == 1 {
        sum += uint32(data[len(data)-1]) << 8
    }

    // Dodajemy przeniesienia
    for (sum >> 16) > 0 {
        sum = (sum & 0xFFFF) + (sum >> 16)
    }
	// Czyli gdy powstanie przeniesienie z 16 bitów, dodaje się je z powrotem.

    // Negacja bitowa
    return ^uint16(sum)
}

//
// ============================================================
// Liczenie checksum UDP z pseudo-headerem (RFC 768)
// ============================================================
//
func udpChecksumIPv4(srcIP, dstIP net.IP, udpHeader, payload []byte) uint16 {

    // Pseudo-header:
    //  4 bajty Source IP
    //  4 bajty Destination IP
    //  1 bajt zero
    //  1 bajt protocol (17)
    //  2 bajty UDP length

    pseudo := []byte{}
    pseudo = append(pseudo, srcIP...)
    pseudo = append(pseudo, dstIP...)
    pseudo = append(pseudo, 0)                   // zero
    pseudo = append(pseudo, syscall.IPPROTO_UDP) // protocol
    udpLength := uint16(len(udpHeader) + len(payload))
    tmp := make([]byte, 2)
    binary.BigEndian.PutUint16(tmp, udpLength)
    pseudo = append(pseudo, tmp...)

    // Całość do checksum:
    data := append(pseudo, udpHeader...)
    data = append(data, payload...)

    return checksum(data)
}
