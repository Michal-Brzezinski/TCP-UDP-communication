## Polecenie: `sudo iptables -A INPUT -p udp --dport 9090 -j DROP`

Poszczególne składowe:

- iptables - narzędzie do konfiguracji firewalla w jądrze Linux (stare, ale wciąż popularne).

- -A INPUT - dodaj (Append) regułę do łańcucha INPUT:

    INPUT - ruch przychodzący do lokalnej maszyny.

- -p udp - dotyczy tylko pakietów UDP.

- --dport 9090 - dotyczy pakietów, których port docelowy to 9090.

- -j DROP - akcja: wyrzuć pakiet, nie przepuszczaj, nie odpowiadaj.

Efekt:

    Serwer UDP na porcie 9090 nie dostanie żadnego pakietu, nawet jeśli raw‑client je wysyła.


## Polecenie: `sudo iptables -D INPUT -p udp --dport 9090 -j DROP`

- -D - usuwa regułę z łańcucha INPUT.

- Parametry muszą być identyczne jak przy -A.

Efekt:

    Ruch UDP na port 9090 znowu jest przepuszczany, a serwer UDP znów odbiera pakiety.