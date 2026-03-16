# Co to jest raw socket?

Normalnie, gdy używam:

`net.DialUDP`, `net.ListenUDP`, `net.Dial("tcp")`

to system operacyjny:

- sam tworzy nagłówek IP,

- sam tworzy nagłówek UDP/TCP,

- sam liczy checksumy,

- sam ustawia pola TTL, ID, flags, fragmenty,

- sam ustawia porty źródłowe,

- sam ustawia długości pakietów.

Raw socket to tryb, w którym:

System operacyjny `NIE tworzy nagłówków IP ani UDP`, zatem należy je zbudować ręcznie.

# Co należy zrobić? Plan działania

## 1. Utworzyć raw socket:

```go
syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
```

## 2. Zbudować nagłówek IPv4 (20 bajtów)

```
0               7 8             15 16            23 24            31
+----------------+----------------+----------------+----------------+
|Version|  IHL   |Type of Service |        Total Length             |
+----------------+----------------+----------------+----------------+
|        Identification           |Flags|   Fragment Offset         |
+----------------+----------------+----------------+----------------+
|  TTL  |Protocol|    Header Checksum                               |
+----------------+----------------+----------------+----------------+
|                      Source IP Address                            |
+----------------+----------------+----------------+----------------+
|                   Destination IP Address                          |
+----------------+----------------+----------------+----------------+
```

Najważniejsze pola:

- Version = 4

- IHL = 5 (czyli 20 bajtów nagłówka)

- Protocol = 17 (UDP)

- Total Length = 20 + 8 + payload

- Source IP = Mój adres

- Destination IP = adres serwera UDP

- Checksum -> trzeba policzyć


## 3. Zbudować nagłówek UDP (8 bajtów)
```
0               7 8             15 16            23 24            31
+----------------+----------------+----------------+----------------+
|     Src Port   |   Dest Port   |     Length     |   Checksum      |
+----------------+----------------+----------------+----------------+

```

Najważniejsze pola:

- Src Port = dowolny (np. 54321)

- Dest Port = port mojego serwera UDP (np. 9090)

- Length = 8 + payload

- Checksum = trzeba policzyć (pseudo-header):

    Checksum UDP liczy się z:

        - Nagłówka UDP

        - Payloadu

        - Pseudo-headera (to jest najgorsza część)

    Pseudo-header:
    ```
    Source IP (4 bajty)
    Destination IP (4 bajty)
    Zero (1 bajt)
    Protocol (1 bajt = 17)
    UDP Length (2 bajty)
    ```


## 4. Zbudować payload (np. "Hello from raw socket")

## 5. Policz checksum UDP (i ewentualnie IP)

## 6. Połączyć to wszystko w jeden pakiet:

        [IPv4 header][UDP header][payload]

## 7. Wysłać pakiet przez syscall.Sendto

## 8. Twój serwer UDP z zadania 3.0 powinien go odebrać.
