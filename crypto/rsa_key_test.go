package crypto

import (
  "crypto/rsa"
  "errors"
  "testing"
)

func mockPublicKey() string {
  return `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxv11IUrwEaWPR6/r+gKe
PG/zv4SBqMmwAj6eBu4OOl4XnS228NClf16iSwMl59/maKFGsg1nnfBClaZv4zKA
snDYl/QVJJBsVZbVNY29c0gaDhM4ZnUXUgY5mG6tvDCh3TitVPxGmaogSkgxIPN4
g7OoxXb8nl45Rno3MkPLIwL8t3UcMqwe5WI837EtO5pUKQB+KjgV1OAaQy7ptyHv
4hHHefxNoU0adVPYMipNqIzsPuuaE38DcCr6utmZetmedmMesGCZUPGKjTOetxQ+
wbsfwTT/kTyMaTTB6/aF7un64NzWeqogqDzvfwemSq1qXO7mAvpRHti80Nx47fJC
TQIDAQAB
-----END PUBLIC KEY-----`
}

func mockWrongPublicKey() string {
  return `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyJusHBvW5+kJW4y33ET0
CP3p3PzfVDSu3U9gf26U557NLULuY9Qo8hUizzcS4U6W8A+/H+zGNY2iuT/fqtmV
AhD3uKvro/u7BoC4T5KQtMMuQI/HATlzkAoQVpQZn5eMniHSmksbyjU+v09B8p4W
uJg6rGoC3/nZVc3bNmDO3BepM0mNG8wEF8jWIIq2ntAVndduNCAtkjRnjyqHKrpo
m5033+c9+Oz9RAqu8C4/Wi/20t1geJWcb+uulf/t2xVWRNxk8hqpwXzwOWKL2zwR
yo31T2njACfSUMfZmwX93jBE+Z7g+pNRZRWRHeCJuFcrwmNR+vPQxm1mlIh+JLZW
aQIDAQAB
-----END PUBLIC KEY-----`
}

func mockPrivateKey() string {
  return `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAxv11IUrwEaWPR6/r+gKePG/zv4SBqMmwAj6eBu4OOl4XnS22
8NClf16iSwMl59/maKFGsg1nnfBClaZv4zKAsnDYl/QVJJBsVZbVNY29c0gaDhM4
ZnUXUgY5mG6tvDCh3TitVPxGmaogSkgxIPN4g7OoxXb8nl45Rno3MkPLIwL8t3Uc
Mqwe5WI837EtO5pUKQB+KjgV1OAaQy7ptyHv4hHHefxNoU0adVPYMipNqIzsPuua
E38DcCr6utmZetmedmMesGCZUPGKjTOetxQ+wbsfwTT/kTyMaTTB6/aF7un64NzW
eqogqDzvfwemSq1qXO7mAvpRHti80Nx47fJCTQIDAQABAoIBAAyGVYID8npZ3lvX
wdWZppYNQd1THMof77kkcdPj1fdshrX486PSrigHL9Xi29bta9Y4GHgKifQR9E7x
C+fT/O++VJOz5ETJ5le4x7C4PC1uY11xbkJcqlwaUjO6+6p1sSp4b8iCnHr9j0y9
oIH/cR1xCHVtWNcq/RXniWPbioSohduvzup2AFkdCigyW6Nja6okb6L/THVNXHtq
rSuGDG3YlotxdftoSoW+94QIY1A89Af1002Z+RNsuc+blfxpe9XQL8PiEXeJF/q5
PprMzx9iYlXgyN6iVJOFN12ib3cyz7u8YJTeMu1S8HgKUKBCHefdlBj62618pPg1
uKFesg0CgYEA/m8ehjm5gtBqH/fGXvSbhkImItmsUDXVmeklzovv3ungF41ehm1M
FvbXBQLf2o+rwffK0iV8MOYvSLk1Db4WbDN1tsmZS/iRfGYIT3CV5I7orHCp9Li2
qewZq/6yqwJjiEimpUvT/zrjAwtoQz6qfEWXnsolGXef16cYaSDtYbcCgYEAyDb7
YaN+0QMyH2+Mu204lFYu9YIi7F8xHdAwaogmF24rlL1jSmhN2VpGoTnY23G72Zq7
izU9TbRiBMZJxsHSgoynbXuMUqUcAUF4Dupul/+2+Zj13URvGG26RzZSno/KGRG4
MHwtdmrwbvdkA/8eaaep8fpnFMD1GL6bPFrOrBsCgYEA/BfRIKD1I52oaMAw9kha
CC5mZsVRq6+LUhHleb7BDhagB/X0IDEO4Pn1lWuBrKYJQghoFss5P6HyW5XV8SXU
RaS/Dzqz/sfsLltSBJPCkFDgTGrcmjKiGb5quTWEhVe6kn+ZTdHR3OLVpmCZD3d5
p+O0FIqpM5CI+T0APLl5OgUCgYBLuVfcduzY+p9zeko8/TNAD1SVcJHq2poGD56w
PCxEAlwjVnn+Q3LmORmrkuhtHxgQVlCGdy1nfUjxS1nN/bKzw6TzaJ4LB/2OkAdr
hMktXf8DahHbjS2DjMS+eFJJPFMQpj4GwIClYA7tuU2voUcMaOiC59Ui6VQJ9tVZ
v3KZbwKBgEOuPRLTqs5CrY1WyoJ7dsVC7EbBD6Lwcoc6UvWgkFNVayekawZq5S43
qBpk19ErRswsGDYFuBpqYH2xLMdTBhFFo/2xUh656POXd2ReeyV01+VWorD9aBUL
ZB4zMMx0qTlJOFF7DPOYWOYBsNdlfkH2INafdAxNYY7pYDAUjJmJ
-----END RSA PRIVATE KEY-----`
}

func TestParseRSAKey_Public(t *testing.T) {
  keyBytes := []byte(mockPublicKey())
  _, actual := parseRSAKey(keyBytes)
  if actual != nil {
    t.Errorf("Got: %v, wanted: %v\n", actual, nil)
  }
}

func TestParseRSAKey_Private(t *testing.T) {
  keyBytes := []byte(mockPrivateKey())
  _, actual := parseRSAKey(keyBytes)
  if actual != nil {
    t.Errorf("Got: %v, wanted: %v\n", actual, nil)
  }
}

func TestParseRSAKey_Error(t *testing.T) {
  expected := errors.New("key not found")
  keyBytes := []byte("Luke, I am your father!")
  _, actual := parseRSAKey(keyBytes)
  if actual == nil || actual.Error() != expected.Error() {
    t.Errorf("Got: %v, wanted: %v\n", actual, expected)
  }
}

func TestVerifyKeyType_Private(t *testing.T) {
  parseKey, _ := parseRSAKey([]byte(mockPrivateKey()))
  _, actual := verifyKeyType(parseKey, &rsa.PrivateKey{})
  if actual != nil {
    t.Errorf("Got: %v, wanted: %v\n", actual, nil)
  }
}

func TestVerifyKeyType_Error(t *testing.T) {
  expected := errors.New("unsupported key type *rsa.PrivateKey")
  keyBytes := []byte(mockPrivateKey())
  parseKey, _ := parseRSAKey(keyBytes)
  _, actual := verifyKeyType(parseKey, rsa.PrivateKey{})
  if actual == nil || actual.Error() != expected.Error() {
    t.Errorf("Got: %v, wanted: %v\n", actual, expected)
  }
}

func TestRSAPrivateKey_Sign(t *testing.T) {
  keyBytes := []byte(mockPrivateKey())
  rsaKey, _ := NewRSAPrivateKey(keyBytes)
  data := []byte("The force is strong with this one")
  
  sig, err := rsaKey.Sign(data)
  if err != nil {
    t.Errorf("Got: %v, wanted: %v\n", err, nil)
  }
  if len(sig) <= 0 {
    t.Errorf("Got: %s, wanted: %s\n", "valid signature", "empty signature")
  }
}

func TestRSAPublicKey_Verify(t *testing.T) {
  rsaPrivateKey, err := NewRSAPrivateKey([]byte(mockPrivateKey()))
  data := []byte("The force is strong with this one")
  sig, _ := rsaPrivateKey.Sign(data)
  
  rsaPublicKey, _ := NewRSAPublicKey([]byte(mockPublicKey()))
  validSig, err := rsaPublicKey.Verify(data, sig)
  if err != nil {
    t.Errorf("Got: %v, wanted: %v\n", err, nil)
  }
  if validSig == false {
    t.Errorf("Got: %v, wanted: %v\n", validSig, true)
  }
}

func TestRSAPublicKey_Verify_WrongKey(t *testing.T) {
  rsaPrivateKey, err := NewRSAPrivateKey([]byte(mockPrivateKey()))
  data := []byte("The force is strong with this one")
  sig, _ := rsaPrivateKey.Sign(data)
  
  rsaPublicKey, _ := NewRSAPublicKey([]byte(mockWrongPublicKey()))
  validSig, err := rsaPublicKey.Verify(data, sig)
  if err == nil {
    t.Errorf("Got: %s, wanted: %v\n", "crypto/rsa: verification error", nil)
  }
  if validSig {
    t.Errorf("Got: %v, wanted: %v\n", validSig, false)
  }
}
