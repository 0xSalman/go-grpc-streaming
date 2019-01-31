package main

import (
  "context"
  "io"
  "log"
  "net"
  "os"
  "os/exec"
  "testing"
  "time"
  
  "github.com/salman-ahmad/grpc-streaming/config"
  "github.com/salman-ahmad/grpc-streaming/crypto"
  pb "github.com/salman-ahmad/grpc-streaming/proto"
  "google.golang.org/grpc"
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

func startClient(port string) *grpc.ClientConn {
  log.Println("startClient()")
  var opts []grpc.DialOption
  opts = append(opts, grpc.WithInsecure())
  
  conn, err := grpc.Dial("localhost:"+port, opts...)
  if err != nil {
    log.Fatalf("failed to connect to server: %v\n", err)
  }
  return conn
}

func stopClient(conn *grpc.ClientConn) {
  log.Println("stopClient()")
  if err := conn.Close(); err != nil {
    log.Fatalf("failed to stop client: %v\n", err)
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

func rsaPrivateKey() crypto.PrivateKey {
  privKeyPath, err := config.AbsolutePath(conf.PrivateKey)
  if err != nil {
    log.Fatalf("failed to calculate private key's absloute path :%v\n", err)
  }
  
  privateKey, err := crypto.NewFileKey(privKeyPath)
  if err != nil {
    log.Fatalf("failed to load private key: %v\n", err)
  }
  
  rsaPrivateKey, err := crypto.NewRSAPrivateKey(privateKey.Bytes())
  if err != nil {
    log.Fatalf("failed to read private key: %v\n", err)
  }
  return rsaPrivateKey
}

func TestFindMaxNumber(t *testing.T) {
  numbersToSend := []int64{1, 4, 100, 30, 50, 203}
  expectedMaxNumber := int64(203)
  done := make(chan struct{})
  
  stream, err := simpleClient.FindMaxNumber(context.Background())
  if err != nil {
    log.Fatalf("%vFindMaxNumber(_) = _ %v", simpleClient, err)
  }
  
  // send numbers
  rsaPrivateKey := rsaPrivateKey()
  go func() {
    for _, number := range numbersToSend {
      signature, err := rsaPrivateKey.Sign(crypto.Int64ToBytes(number))
      if err != nil {
        log.Fatalf("failed to sign the request: %v\n", err)
      }
      req := &pb.MaxNumberRequest{Number: number, Signature: signature}
      if err := stream.Send(req); err != nil {
        log.Fatalf("failed to send the request: %v\n", err)
      }
    }
    
    if err := stream.CloseSend(); err != nil {
      log.Fatalf("failed to close the stream: %v\n", err)
    }
  }()
  
  // receive max numbers
  var actualMaxNumber int64
  go func() {
    for {
      response, err := stream.Recv()
      if err == io.EOF {
        close(done)
        return
      }
      if err != nil {
        log.Fatalf("failed to receive stream response: %v\n", err)
      }
      actualMaxNumber = response.Number
    }
  }()
  
  <-done
  if actualMaxNumber != expectedMaxNumber {
    t.Errorf("Got: %d, wanted: %d\n", actualMaxNumber, expectedMaxNumber)
  }
}
