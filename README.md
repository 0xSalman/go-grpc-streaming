# gRPC Bi-Directional Streaming

A Go gRPC service that implements a bi-directional streaming method to
find the maximum number.

## Problem

Implement a gRPC bi-directional streaming client and server system:

- The client will stream numbers and sign each number with a private key
- The server will receive numbers, verify request signature with a matching public key
- When the signature is valid, the server will stream back maximum number it has received so far

## Solution

- Created a simple gRPC bi-directional streaming client and server using protobuf
- The request object from the client sends a `int64` number and signature of the signed number
- The response object from server send back a `int64` number
- The client uses go routines to send & receive numbers in parallel
- A simple interface `crypto/Key` is created to allow different implementations of
how to read public & private keys. The default implementation is `crypto/FileKey` ,
which as name suggests reads private & public keys from a filesystem
- Similarly, the `crypto/PublicKey` & `crypto/PrivateKey` interfaces allow different implementations
of how to sign data and verify its signature. The default implementations
`crypto/RSAPublicKey` & `crypto/RSAPrivateKey` uses sha256 and RSA X.509 PEM formatted keys
to sign and verify
- `config/config.go` makes it easier to change the values of the most important variables
either by updating the default values or using the environment vars.
Please see environment variables section for more details
- The client and server integration tests build the `server/server.go` executable
and run it va `exec.Command`. It is done this way to have control over the
server process and kill it at the end of the tests to free the port

## Possible Improvements

- Maybe use sha-512 in RSA key implementations to save bandwidth
- On the server side, maybe use chain of responsibility pattern to verify
signature and determine the new maximum number
- Extract out the shared code between tests to a common package

## Requirements

1. make (http://www.gnu.org/software/make/)
2. Go 1.11+
3. protobuf tool (https://github.com/google/protobuf/releases)

## Getting Started

1. Generate a new RSA key in `$HOME/.ssh`
    1. Go to `$HOME/.ssh` via terminal
    2. Run command `openssl genrsa -out maxnumber_rsa_private.pem 2048`
    3. Run command `openssl rsa -in maxnumber_rsa_private.pem -pubout -out maxnumber_rsa_public.pem`
    4. Run command `chmod 600 maxnumber_rsa*` to change file permissions
2. `make install`
3. `make gen-rpc`
4. `make run-server` - by default server runs on port `7000`
5. `make run-client`

To run tests, do: `make test`

## Environment Variables

There are number of variables that can be configured via environment vars
to override their default values (i.e., gRPC server port etc).
The following variables are configurable:

- `GRPC_PORT`, default value is `7000`
- `GRPC_PRIVATE_KEY`, default value is `$HOME/.ssh/maxnumber_rsa_private.pem`
- `GRPC_PUBLIC_KEY`, default value is `$HOME/.ssh/maxnumber_rsa_public.pem`
- `GRPC_TOTAL_NUMBERS`, total numbers to send; default value is `15`
- `GRPC_NUMBER_MULTIPLIER`, random number multiplier; default value is `100`

