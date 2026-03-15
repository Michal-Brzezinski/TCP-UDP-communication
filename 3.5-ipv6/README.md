# Różnice w wykorzystaniu IPv6 względem wersji IPv4

## Adresy i protokoły:

- tcp4 -> tcp6

- udp4 -> udp6

- "0.0.0.0" -> "::"

- "127.0.0.1" -> "::1"

- "127.0.0.1:8080" -> "[::1]:8080"

- ":8080" -> "[::]:8080"

## Reszta:

Logika serwera/klienta, obsługa sygnałów, gorutyny, mapy klientów, echo–protokół jest identyczna.

API Go (net.Listen, net.Dial, net.ListenUDP, net.DialUDP, ReadFromUDP, WriteToUDP) działa tak samo, tylko z innym „rodzajem” (tcp6/udp6) i innymi adresami.