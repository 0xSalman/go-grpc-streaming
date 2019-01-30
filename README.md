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

