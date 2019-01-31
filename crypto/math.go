package crypto

import "encoding/binary"

func Int64ToBytes(number int64) []byte {
  byteData := make([]byte, 8)
  binary.LittleEndian.PutUint64(byteData, uint64(number))
  return byteData
}

func IsNewInt64Max(max, newNumber int64) bool {
  if newNumber > max {
    return true
  }
  return false
}
