# ip2asn

## original readme

This package provides the ability to search for ASN / Country Code data given
an IP address.

It depends on the awesome work by @jedisct1 (Frank Denis).  If you use this,
consider sponsoring him!

https://github.com/sponsors/jedisct1

## updated readme

This CLI tool could use offline ASN data to check the info of specified IP or IP list.

Use Makefile to download ASN info file.

Then compile the command tool:

```
go build ./cmd/ip2asn/main.go
```

And we got an executable file named `main` or `main.exe`. Make sure the ASN data file is in the same directory of compiled executable file. Provide our IP list file `test.txt` and check them by:

```
./main test.txt
```

We can get line-by-line results.