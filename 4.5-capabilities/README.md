## Co to jest setcap:

    Raw socket wymaga uprawnień uprzywilejowanych, bo pozwala omijać normalny stos TCP/UDP. Zamiast uruchamiać program jako root, nadaję mu capability cap_net_raw+ep przez setcap. Dzięki temu tylko ten program ma prawo tworzyć raw socket, ale nie ma pełnych uprawnień roota.

## Polecenie `sudo setcap cap_net_raw+ep ./raw-client`

### `cap_net_raw`:

Pozwala procesowi:

- tworzyć raw sockets,
- wysyłać/odbierać pakiety na poziomie IP,
- bez potrzeby bycia rootem.

#### `+ep` - znaczenie:

- e — effective (uprawnienie jest aktywne w czasie działania programu),
- p — permitted (program może to uprawnienie mieć).

## Polecenie `sudo setcap cap_net_bind_service+ep ./udp-server`

`cap_net_bind_service` - znaczenie:

Pozwala procesowi bindować porty <1024 bez bycia rootem. Normalnie porty 1–1023 są „uprzywilejowane” i wymagają roota.
Dzięki temu mój serwer UDP może nasłuchiwać np. na porcie 53 czy 80 bez uruchamiania go jako root.

#### Żeby sprawdzić: 
1. Zmiana w kodzie serwera port z 9090 na np. 53.
2. Budowanie: go build -o udp-server.
3. Nadanie capability: sudo setcap cap_net_bind_service+ep ./udp-server.
4. Uruchominie: ./udp-server (bez sudo) → działa na porcie 53.