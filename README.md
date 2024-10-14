# MNB QR code standard implementation in Go

[![Go](https://github.com/Gerifield/mnb-qr-go/actions/workflows/go.yml/badge.svg)](https://github.com/Gerifield/mnb-qr-go/actions/workflows/go.yml)

Standard: https://www.mnb.hu/penzforgalom/azonnalifizetes/utmutatok

PDF:
- https://www.mnb.hu/letoltes/qr-kod-utmutato-20190712.pdf
- https://www.mnb.hu/letoltes/qr-kod-utmutato-20190712-en.pdf

This code is mainly a lib, but there's a binary to test and show how you could use it.

Install:
```
go get -u github.com/gerifield/mnb-qr-go/...
```

## Using the server

```
$ go run src/cmd/qr-server/qr-server.go
```

Different terminal:
```
$ curl -X POST "http://127.0.0.1:8080" -d '{"pngSize":128,"kind":"RTP","bic":"abcdefghijk","name":"Test User","iban":"HU00123456789012345678901234","expire":360}' --output test.png
$ open test.png
```

### Possible JSON fields (and types)

Check the MNB docs for more details.

Reqired:
- `kind` - string (`RTP` or `HCT`)
- `bic` - string (`8` or `11` character, the `8` char long will get a `XXX` postfix)
- `name` - string (70 chars max, recipient or sender name)
- `iban` - string (28 chars)
- `expire` - int (seconds added to the current time)
- `pngSize` - int (generated image size in pixels `128` or `256` should be fine)

Optional:
- `amount` - int (amount in HUF, optional)
- `purpose` - string (4 char, from a fixed set, check the `purposeCodes` variable in the code)
- `message` - string (70 chars max, message added to the code)
- `shopID` - string (35 chars max)
- `shopID` - string (35 chars max)
- `merchDevID` - string (35 chars max)
- `invoiceID` - string (35 chars max)
- `customerID` - string (35 chars max)
- `credTranID` - string (35 chars max)
- `loyaltyID` - string (35 chars max)
- `navCheckID` - string (35 chars max)

### Build using docker

```
$ docker build -t mnb-qr-server .
```

It'll produce a container named `mnb-qr-server`. You could run it:
```
$ docker run -d -p8080:8080 mnb-qr-server
```

### Using the prebuild image

```
$ docker run -d -p8080:8080 docker.io/gerifield/mnb-qr-server:latest
```

## Command line tool usage
```
$ mnb-qr-gen -bic CIBHHUHB -name "Test Name" -iban HU90107001234567890123456789 -amount 5 -message "Hello\!"
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


## Docker usage

For this you only need a docker on your machine to build and run the server.
Use the following commands, the first will build the binary (and run the tests) and the second will start the server
and allow you to connect to it on your localhost's port 8080: 

```
$ docker build -t mnb-qr .
$ docker run -p8080:8080 mnb-qr:latest
```

# Development ideas:

- It looks like the EPC standard QR looks very similar to this format, we could add support for that too! https://en.wikipedia.org/wiki/EPC_QR_code - https://www.europeanpaymentscouncil.eu/sites/default/files/kb/file/2018-05/EPC069-12%20v2.1%20Quick%20Response%20Code%20-%20Guidelines%20to%20Enable%20the%20Data%20Capture%20for%20the%20Initiation%20of%20a%20SCT.pdf
