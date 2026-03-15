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
# Podstawowe różnice między TCP a UDP:

1. TCP ma połączenie — UDP nie
2. TCP gwarantuje dostarczenie danych — UDP nie 
3. TCP jest strumieniowe — UDP jest pakietowe
4. TCP ma kontrolę przepływu i przeciążenia — UDP nie
5. TCP ma „połączenie” klient–serwer — UDP ma tylko adresy