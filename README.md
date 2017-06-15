[![Build Status](https://travis-ci.org/dedis/kyber.svg?branch=master)](https://travis-ci.org/dedis/kyber)

DeDiS Advanced Crypto Library for Go
====================================

This package provides a toolbox of advanced cryptographic primitives for Go,
targeting applications like [Dissent](http://dedis.cs.yale.edu/dissent/)
that need more than straightforward signing and encryption.
Please see the
[GoDoc documentation for this package](http://godoc.org/github.com/DeDiS/crypto)
for details on the library's purpose and API functionality.

Installing
----------

First make sure you have [Go](https://golang.org)
version 1.3 or newer installed.

The basic crypto library requires only Go and a few
third-party Go-language dependencies that can be installed automatically
as follows:

	go get github.com/dedis/kyber
	cd $GOPATH/src/github.com/dedis/kyber
	go get ./... # install 3rd-party dependencies

You should then be able to test its basic function as follows:

	cd $GOPATH/src/github.com/dedis/kyber
	go test -v

You can recursively test all the packages in the library as follows,
keeping in mind that some sub-packages will only build
if certain dependencies are satisfied as described below:

	go test -v ./...

Issues
------

- Traditionally, ECDH (Elliptic curve Diffie-Hellman) derives the shared secret
from the x point only. In this framework, you can either manually retrieve the
value or use the MarshalBinary method to take the combined (x, y) value as the
shared secret. We recommend the latter process for new softare/protocols using
this framework as it is cleaner and generalizes across different types of
groups (e.g., both integer and elliptic curves), although it will likely be
incompatible with other implementations of ECDH.
http://en.wikipedia.org/wiki/Elliptic_curve_Diffie%E2%80%93Hellman

