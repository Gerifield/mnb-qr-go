# MNB QR code standard implementation in Go

Standard: https://www.mnb.hu/penzforgalom/azonnalifizetes/utmutatok

PDF:
- https://www.mnb.hu/letoltes/qr-kod-utmutato-20190712.pdf
- https://www.mnb.hu/letoltes/qr-kod-utmutato-20190712-en.pdf


Install:
```
$ go get -u https://github.com/gerifield/mnb-qr-go/...
```


Example usage:
```
$ mnb-qr-go -bic CIBHHUHB -name "Test Name" -iban HU90107001234567890123456789 -amount 5
RTP
001
1
CIBHHUHBXXX
Test Name
HU90107001234567890123456789
HUF5
20200520003312+2

Hello!






```

It'll generate an `out.png` and try to open it on the system.
