package crypto

import (
  "encoding/binary"
  "testing"
)

func TestIsNewInt64Max_WithLarge(t *testing.T) {
  expected := true
  actual := IsNewInt64Max(100, 250)
  if expected != actual {
    t.Errorf("Got: %v, wanted: %v\n", actual, expected)
  }
}

func TestIsNewInt64Max_WithSmall(t *testing.T) {
  var expected bool
  actual := IsNewInt64Max(100, 50)
  if expected != actual {
    t.Errorf("Got: %v, wanted: %v\n", actual, expected)
  }
}

func TestIsNewInt64Max_WithEqual(t *testing.T) {
  var expected bool
  actual := IsNewInt64Max(100, 100)
  if expected != actual {
    t.Errorf("Got: %v, wanted: %v\n", actual, expected)
  }
}

func TestInt64ToBytes(t *testing.T) {
  expected := int64(73)
  numberBytes := Int64ToBytes(expected)
  actual := binary.LittleEndian.Uint64(numberBytes)
  if expected != int64(actual) {
    t.Errorf("Got: %d, wanted: %d\n", actual, expected)
  }
}
