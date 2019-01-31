package config

import "testing"

func TestAbsolutePath_WithPrefix(t *testing.T) {
  path := "~/logs/max-numbers"
  userHome, _ := UserHomeDir()
  expected := userHome + "/logs/max-numbers"
  actual, _ := AbsolutePath(path)
  if expected != actual {
    t.Errorf("Got: %s, wanted: %s\n", actual, expected)
  }
}

func TestAbsolutePath_WithNoPrefix(t *testing.T) {
  path := "/var/logs/nginx"
  expected := "/var/logs/nginx"
  actual, _ := AbsolutePath(path)
  if expected != actual {
    t.Errorf("Got: %s, wanted: %s\n", actual, expected)
  }
}
