package main

import (
  "io"
  "log"
  "math/rand"
  "time"
  
  "github.com/salman-ahmad/grpc-streaming/config"
  "github.com/salman-ahmad/grpc-streaming/crypto"
  pb "github.com/salman-ahmad/grpc-streaming/proto"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
)

func main() {
  
  conf := loadConfig()
  conn := startClient(conf.Port)
  defer stopClient(conn)
  
  // generate random numbers between 0 and
  // conf.NumbersToSend * conf.NumberMultiplier
  numbers := make([]int64, conf.NumbersToSend)
  for i := 0; i < conf.NumbersToSend; i++ {
    numbers[i] = int64(rand.Intn((i + 1) * conf.NumberMultiplier))
  }
  
  client := pb.NewSimpleClient(conn)
  rsaPrivateKey := rsaPrivateKey(conf.PrivateKey)
  maxNumber := findMaxNumber(client, rsaPrivateKey, numbers)
  log.Printf("finished with maxNumber %d\n", maxNumber)
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

func startClient(port string) *grpc.ClientConn {
  log.Println("startClient()")
  conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("failed to connect to server: %v", err)
  }
  log.Printf("connected to server localhost:%s\n", port)
  return conn
}

func stopClient(conn *grpc.ClientConn) {
  log.Println("stopClient()")
  if err := conn.Close(); err != nil {
    log.Fatalf("failed to stop client: %v\n", err)
  }
}

func rsaPrivateKey(key string) crypto.PrivateKey {
  log.Println("rsaPrivateKey()")
  privKeyPath, err := config.AbsolutePath(key)
  if err != nil {
    log.Fatalf("failed to calculate private key's absloute path :%v\n", err)
  }
  
  privateKey, err := crypto.NewFileKey(privKeyPath)
  if err != nil {
    log.Fatalf("failed to load private key: %v\n", err)
  }
  log.Printf("using private key %s\n", privKeyPath)
  
  rsaPrivateKey, err := crypto.NewRSAPrivateKey(privateKey.Bytes())
  if err != nil {
    log.Fatalf("failed to read private key: %v\n", err)
  }
  log.Println("parsed RSA private key")
  return rsaPrivateKey
}

// invoke server to find the maximum number
func findMaxNumber(
  client pb.SimpleClient,
  privateKey crypto.PrivateKey,
  numbers []int64) int64 {
  
  log.Println("findMaxNumber()")
  stream, err := client.FindMaxNumber(context.Background())
  if err != nil {
    log.Fatalf("failed to open stream: %v", err)
  }
  
  // go routine to stream numbers to server
  go sendNumbers(stream, privateKey, numbers)
  
  // go routine to receive maximum number from server
  maxNumReceiver := make(chan int64)
  go getMaxNumber(stream, maxNumReceiver)
  
  var maxNumber int64
  for num := range maxNumReceiver {
    maxNumber = num
    log.Printf("received new maxNumber %d\n", maxNumber)
  }
  return maxNumber
}

// send the given numbers and sleep between each send
func sendNumbers(
  stream pb.Simple_FindMaxNumberClient,
  privateKey crypto.PrivateKey,
  numbers []int64) {
  
  log.Println("sendNumbers()")
  for _, number := range numbers {
    signature, err := privateKey.Sign(crypto.Int64ToBytes(number))
    if err != nil {
      log.Fatalf("failed to sign the request: %v\n", err)
    }
    
    request := &pb.MaxNumberRequest{Number: number, Signature: signature}
    if err := stream.Send(request); err != nil {
      log.Fatalf("failed to send the request: %v\n", err)
    }
    log.Printf("sent new number %d\n", request.Number)
    time.Sleep(time.Millisecond * 200)
  }
  
  if err := stream.CloseSend(); err != nil {
    log.Fatalf("failed to close the stream: %v\n", err)
  }
}

// receive max number from server and
// close the channel when stream is finished
func getMaxNumber(
  stream pb.Simple_FindMaxNumberClient,
  maxNumReceiver chan int64) {
  
  log.Println("getMaxNumber()")
  for {
    response, err := stream.Recv()
    if err == io.EOF {
      close(maxNumReceiver)
      return
    }
    if err != nil {
      log.Fatalf("failed to receive stream response: %v\n", err)
    }
    
    maxNumReceiver <- response.Number
  }
}
