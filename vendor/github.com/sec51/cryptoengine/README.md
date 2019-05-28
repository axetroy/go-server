### Build status

[![Build Status](https://travis-ci.org/sec51/cryptoengine.svg?branch=master)](https://travis-ci.org/sec51/cryptoengine)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/sec51/cryptoengine/)

### CryptoEngine package

This simplifies even further the usage of the NaCl crypto primitives,
by taking care of the `nonce` part.
It uses a KDF, specifically HKDF to compute the nonces.

### Big Picture

The encryption and decryption phases are the following:

```

Message -> Encrypt -> EncryptedMessage -> ToBytes() -> < = NETWORK = >  <- FromBytes() -> EncryptedMessage -> Decrypt -> Message

```

### Usage

1- Import the library

```
import github.com/sec51/cryptoengine
```

2- Instanciate the `CryptoEngine` object via:

```
	engine, err := cryptoengine.InitCryptoEngine("Sec51")
	if err != nil {
		return err
	}
```
See the godoc for more info about the InitCryptoEngine parameter

3- Encrypt a message using symmetric encryption

```
    message := "the quick brown fox jumps over the lazy dog"
	engine.NewMessage(message)
	if err != nil {
		return err
	}
```

4- Serialize the message to a byte slice, so that it can be safely sent to the network

```
	messageBytes, err := tcp.ToBytes()
	if err != nil {
		t.Fatal(err)
	}	
```

5- Parse the byte slice back to a message

```
	message, err := MessageFromBytes(messageBytes)
	if err != nil {
		t.Fatal(err)
	}
```

### License

Copyright (c) 2015 Sec51.com <info@sec51.com>

Permission to use, copy, modify, and distribute this software for any
purpose with or without fee is hereby granted, provided that the above 
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE. 

