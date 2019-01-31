package crypto

import (
  "crypto"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha256"
  "crypto/x509"
  "encoding/pem"
  "errors"
  "fmt"
  "reflect"
)

func parseRSAKey(key []byte) (interface{}, error) {
  block, _ := pem.Decode(key)
  if block == nil {
    return nil, errors.New("key not found")
  }
  
  switch block.Type {
  case "PUBLIC KEY":
    return x509.ParsePKIXPublicKey(block.Bytes)
  case "RSA PRIVATE KEY":
    return x509.ParsePKCS1PrivateKey(block.Bytes)
  default:
    return nil, fmt.Errorf("unsupported key type %s", block.Type)
  }
}

func verifyKeyType(actual, expected interface{}) (interface{}, error) {
  if reflect.TypeOf(actual) == reflect.TypeOf(expected) {
    return actual, nil
  }
  return nil, fmt.Errorf("unsupported key type %T", actual)
}

type RSAPrivateKey struct {
  *rsa.PrivateKey
}

func NewRSAPrivateKey(key []byte) (PrivateKey, error) {
  rawKey, err := parseRSAKey(key)
  if err != nil {
    return nil, err
  }
  
  rawKey, err = verifyKeyType(rawKey, &rsa.PrivateKey{})
  if err != nil {
    return nil, err
  }
  
  rsaKey := rawKey.(*rsa.PrivateKey)
  return &RSAPrivateKey{rsaKey}, nil
}

func (r RSAPrivateKey) Sign(data []byte) ([]byte, error) {
  h := sha256.New()
  h.Write(data)
  d := h.Sum(nil)
  return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}

// TODO to be implemented
func (r RSAPrivateKey) Decrypt(data []byte) ([]byte, error) {
  return nil, errors.New("function is not implemented yet")
}

type RSAPublicKey struct {
  *rsa.PublicKey
}

func NewRSAPublicKey(key []byte) (PublicKey, error) {
  rawKey, err := parseRSAKey(key)
  if err != nil {
    return nil, err
  }
  
  rawKey, err = verifyKeyType(rawKey, &rsa.PublicKey{})
  if err != nil {
    return nil, err
  }
  
  rsaKey := rawKey.(*rsa.PublicKey)
  return &RSAPublicKey{rsaKey}, nil
}

func (r RSAPublicKey) Verify(data, signature []byte) (bool, error) {
  h := sha256.New()
  h.Write(data)
  d := h.Sum(nil)
  
  err := rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, signature)
  if err == nil {
    return true, nil
  }
  return false, err
}

func (r RSAPublicKey) VerifyString(data []byte, signature string) (bool, error) {
  return r.Verify(data, []byte(signature))
}

// TODO to be implemented
func (r RSAPublicKey) Encrypt(data []byte) ([]byte, error) {
  return nil, errors.New("function is not implemented yet")
}
