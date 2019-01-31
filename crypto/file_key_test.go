package crypto

import (
  "io/ioutil"
  "log"
  "reflect"
  "testing"
)

func writeToFile(file, content string) {
  contentBytes := []byte(content)
  err := ioutil.WriteFile(file, contentBytes, 0644)
  if err != nil {
    log.Fatal(err)
  }
}

func TestFileKey_String(t *testing.T) {
  expected := "this is a file reading test"
  file := "./file_test_string"
  writeToFile(file, expected)
  
  key, _ := NewFileKey(file)
  actual := key.String()
  
  if actual != expected {
    t.Errorf("Got: %s, wanted: %s\n", actual, expected)
  }
}

func TestFileKey_Bytes(t *testing.T) {
  content := "this is another file reading test"
  file := "./file_test_bytes"
  writeToFile(file, content)
  
  expected := []byte(content)
  key, _ := NewFileKey(file)
  actual := key.Bytes()
  
  if !reflect.DeepEqual(actual, expected) {
    t.Errorf("Got: %v, wanted: %v\n", actual, expected)
  }
}

func TestFileKey_Error(t *testing.T) {
  file := "./random_file"
  _, err := NewFileKey(file)
  if err == nil {
    t.Errorf("Got %v, wanted: %s", nil, "error")
  }
}
