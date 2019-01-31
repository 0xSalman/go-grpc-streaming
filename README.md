# Go gRPC Bi-Directional Streaming

A Go gRPC service that implements a bi-directional streaming method to
find the maximum number.

The client streams numbers to the server and signs each number with a private key.
The server verifies the signature with a public key. When the signature is valid,
it streams back the maximum number it has received so far.

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
3. `make generate-rpc`
4. `make run-server` - by default server runs on port 7000
5. `make run-client`

To run tests, do: `make test`

## Environment Variables

There are number of variables that can be configured via environment vars
to override their default values (i.e., gRPC server port etc).
The following variables are configurable:

- `GRPC_PORT`, default value is `7000`
- `GRPC_PRIVATE_KEY`, default value is `$HOME/.ssh/maxnumber_rsa_private.pem`
- `GRPC_PUBLIC_KEY`, default value is `$HOME/.ssh/maxnumber_rsa_public.pem`
- `GRPC_TOTAL_NUMBERS`, default value is `15`
- `GRPC_NUMBER_MULTIPLIER`, default value is `100`

