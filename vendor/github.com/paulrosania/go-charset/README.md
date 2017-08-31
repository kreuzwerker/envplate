# Character set conversion for Go

[![Build Status](https://travis-ci.org/paulrosania/go-charset.svg?branch=master)](https://travis-ci.org/paulrosania/go-charset)

A fork of [go-charset](https://code.google.com/p/go-charset/), which itself is a
port of inferno's convcs for Go, which supports conversion to and from utf-8 for
the following character sets:

* big5
* ibm437
* ibm850
* ibm866
* iso-8859-1
* iso-8859-10
* iso-8859-15
* iso-8859-2
* iso-8859-3
* iso-8859-4
* iso-8859-5
* iso-8859-6
* iso-8859-7
* iso-8859-8
* iso-8859-9
* koi8-r
* us-ascii
* utf-16
* utf-16be
* utf-16le
* utf-8
* windows-1250
* windows-1251
* windows-1252

This project also includes an extra package which links to the GNU iconv library
and adds all the character sets available from it.

## Documentation

Full API documentation is available here:

http://godoc.org/github.com/paulrosania/go-charset/charset

## Contributors

* Roger Peppe <rogpeppe@gmail.com> (original author)
* Paul Rosania <paul@rosania.org>

## License

go-charset is available under the BSD 3-clause license. Details in `LICENSE`.
