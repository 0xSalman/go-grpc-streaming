package main

import (
  "log"
  "net"
  "os"
  "os/exec"
  "testing"
  "time"
  
  "github.com/salman-ahmad/grpc-streaming/config"
  pb "github.com/salman-ahmad/grpc-streaming/proto"
)

var simpleClient pb.SimpleClient
var conf *config.Config

func buildServer() {
  log.Println("buildServer()")
  if _, err := os.Stat("../server/server"); os.IsNotExist(err) {
    log.Println("file does not exist; build it")
    cmd := exec.Command("go", "build", "-o", "server/server", "server/server.go")
    cmd.Dir = ".."
    err := cmd.Start()
    if err != nil {
      log.Fatalf("Server file failed to build: %v\n", err)
    }
  }
}

func startServer() *exec.Cmd {
  log.Println("startServer()")
  cmdStr := "server/server"
  serverCmd := exec.Command(cmdStr)
  serverCmd.Dir = ".."
  
  err := serverCmd.Start()
  if err != nil {
    log.Fatalf("Server failed to start: %v\n", err)
  }
  
  // wait for server to start and open the port
  log.Println("waiting for server to start...")
  time.Sleep(1 * time.Second)
  _, err = net.Listen("tcp", ":"+conf.Port)
  if err != nil {
    log.Println("server started")
  }
  
  return serverCmd
}

func stopServer(serverCmd *exec.Cmd) {
  log.Println("stopServer()")
  if err := serverCmd.Process.Kill(); err != nil {
    log.Fatalf("failed to kill process: %v\n", err)
  }
}

func TestMain(m *testing.M) {
  conf = loadConfig()
  buildServer()
  serverCmd := startServer()
  clientConn := startClient(conf.Port)
  simpleClient = pb.NewSimpleClient(clientConn)
  
  returnCode := m.Run()
  
  stopClient(clientConn)
  stopServer(serverCmd)
  os.Exit(returnCode)
}

func TestRunFindMaxNumber(t *testing.T) {
  numbersToSend := []int64{-100, 1, 4, 100, 30, 50, 203, 1111, 1301, 2004}
  expectedMaxNumber := int64(2004)
  privateKey := rsaPrivateKey(conf.PrivateKey)
  actualMaxNumber := findMaxNumber(simpleClient, privateKey, numbersToSend)
  if actualMaxNumber != expectedMaxNumber {
    t.Errorf("Got: %d, wanted: %d\n", actualMaxNumber, expectedMaxNumber)
  }
}
