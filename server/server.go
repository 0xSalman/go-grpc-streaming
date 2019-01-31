package main

import (
  "io"
  "log"
  "net"
  
  "github.com/salman-ahmad/grpc-streaming/config"
  "github.com/salman-ahmad/grpc-streaming/crypto"
  pb "github.com/salman-ahmad/grpc-streaming/proto"
  "google.golang.org/grpc"
)

type server struct {
  publicKey crypto.PublicKey
}

// TODO maybe use Chain of Responsibility pattern to verify
//  signature and determine new maximum number
func (s server) FindMaxNumber(stream pb.Simple_FindMaxNumberServer) error {
  log.Println("FindMaxNumber()")
  var maxNumber int64
  
  for {
    // receive new request from stream
    request, err := stream.Recv()
    if err == io.EOF {
      log.Println("end of stream")
      return nil
    }
    if err != nil {
      log.Printf("failed to receive stream request: %v\n", err)
      return err
    }
    log.Printf("received new number %d\n", request.Number)
    
    // TODO maybe extract this block to a separate method?
    numberBytes := crypto.Int64ToBytes(request.Number)
    verified, err := s.publicKey.Verify(numberBytes, request.Signature)
    if err != nil {
      log.Printf("failed to verify signature: %v\n", err)
      return err
    }
    
    // when signature is verified and new number is larger
    // then update new max number and send it to stream
    if verified && crypto.IsNewInt64Max(maxNumber, request.Number) {
      maxNumber = request.Number
      resp := &pb.MaxNumberResponse{Number: maxNumber}
      if err := stream.Send(resp); err != nil {
        log.Printf("failed to send stream response: %v\n", err)
        return err
      }
      log.Printf("sent new maxNumber %d\n", maxNumber)
    }
  }
}

func main() {
  
  conf := loadConfig()
  rsaPublicKey := rsaPublicKey(conf.PublicKey)
  server := &server{publicKey: rsaPublicKey}
  grpcServer := grpc.NewServer()
  
  pb.RegisterSimpleServer(grpcServer, server)
  listener := startListener(conf.Port)
  err := grpcServer.Serve(listener)
  if err != nil {
    log.Fatalf("failed to start server: %v\n", err)
  }
}

func loadConfig() *config.Config {
  log.Println("loadConfig()")
  conf, err := config.LoadConfig()
  if err != nil {
    log.Fatalf("failed to read configuration :%v\n", err)
  }
  log.Println("processed configuration")
  return conf
}

func rsaPublicKey(key string) crypto.PublicKey {
  log.Println("rsaPublicKey()")
  pubKeyPath, err := config.AbsolutePath(key)
  if err != nil {
    log.Fatalf("failed to calculate private key's absloute path :%v\n", err)
  }
  
  publicKey, err := crypto.NewFileKey(pubKeyPath)
  if err != nil {
    log.Fatalf("failed to load public key: %v\n", err)
  }
  log.Printf("using public key %s\n", pubKeyPath)
  rsaPublicKey, err := crypto.NewRSAPublicKey(publicKey.Bytes())
  
  if err != nil {
    log.Fatalf("failed to read public key: %v\n", err)
  }
  log.Println("parsed RSA public key")
  return rsaPublicKey
}

func startListener(port string) net.Listener {
  log.Println("startListener()")
  lis, err := net.Listen("tcp", ":"+port)
  if err != nil {
    log.Fatalf("failed to listen to port: %v\n", err)
  }
  log.Printf("starting server on port %s\n", port)
  return lis
}
