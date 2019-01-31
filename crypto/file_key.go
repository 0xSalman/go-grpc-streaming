package crypto

import "io/ioutil"

type FileKey struct {
  content []byte
}

func NewFileKey(path string) (Key, error) {
  contentBytes, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  return &FileKey{content: contentBytes}, nil
}

func (f FileKey) Bytes() []byte {
  return f.content
}

func (f FileKey) String() string {
  return string(f.content)
}
