package main

import (
  "flag"
  "fmt"
)

func main() {
  workers := flag.Int("w", 1, "Number of workers")
  port := flag.Int("p", 1337, "Port")
  useUDP := flag.Bool("udp", false, "Use UDP as the transport protocol")
  url := flag.String("url", "mongodb://localhost:27017", "URL to your mongodb")

  flag.Parse()

  fmt.Println(*workers, *port, *useUDP, *url)

  NewServer(*workers, *port, *useUDP, *url).Serve()
}
