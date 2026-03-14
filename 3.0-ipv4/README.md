## NA 3.0

DIAGRAM PRZEPŁYWU DANYCH dla:

### TCP — połączeniowy, strumieniowy:

```
+---------+        SYN        +---------+
| Klient  | ----------------> | Serwer  |
+---------+ <---------------- +---------+
            SYN/ACK
+---------+ ----------------> +---------+
| Klient  |       ACK        | Serwer  |
+---------+ <---------------- +---------+
            Połączenie OK

Klient:  Write() -----> Serwer: Read()
Serwer:  Write() -----> Klient: Read()

Zamykanie:
FIN ---> FIN/ACK ---> ACK
```

### UDP — bezpołączeniowy, datagramowy

```
+---------+      Datagram      +---------+
| Klient  | -----------------> | Serwer  |
+---------+                    +---------+

Serwer nie potwierdza odbioru.
Brak połączenia, brak handshake.
Każdy pakiet jest niezależny.

```
