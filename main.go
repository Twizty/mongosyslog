package main

import (
  "flag"
  "fmt"

  "gopkg.in/mgo.v2"
  "gopkg.in/mcuadros/go-syslog.v2"
)

func startWorker(c syslog.LogPartsChannel, collection *mgo.Collection) {
  for logParts := range c {
    collection.Insert(logParts)
  }
}

func main() {
  workers := flag.Int("w", 1, "Number of workers")
  port := flag.String("p", "1337", "Port")
  useUDP := flag.Bool("udp", false, "Use UDP as the transport protocol")
  url := flag.String("url", "mongodb://localhost:27017", "URL to your mongodb")

  flag.Parse()

  fmt.Println(*workers, *port, *useUDP, *url)

  hostport := "127.0.0.1:" + *port

  session, err := mgo.Dial(*url)
  if err != nil {
    panic(err)
  }
  db := session.DB("test")
  collection := db.C("logs")

  channel := make(syslog.LogPartsChannel, *workers * 2)
  handler := syslog.NewChannelHandler(channel)

  server := syslog.NewServer()
  server.SetFormat(syslog.RFC5424)
  server.SetHandler(handler)

  if (*useUDP) {
    server.ListenUDP(hostport)
  } else {
    server.ListenTCP(hostport)
  }
  server.Boot()

  for i := 0; i < *workers; i++ {
    go startWorker(channel, collection)
  }

  server.Wait()
}
