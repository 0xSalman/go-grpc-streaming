package crypto

type Key interface {
  Bytes() []byte
  String() string
}

type PrivateKey interface {
  Sign(data []byte) ([]byte, error)
  Decrypt(data []byte) ([]byte, error)
}

type PublicKey interface {
  Verify(data, signature []byte) (bool, error)
  VerifyString(data []byte, signature string) (bool, error)
  Encrypt(data []byte) ([]byte, error)
}
