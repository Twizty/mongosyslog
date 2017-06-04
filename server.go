package main

import (
  "strconv"

  "gopkg.in/mgo.v2"
  "gopkg.in/mcuadros/go-syslog.v2"
)


type Server struct {
  workersCount int
  port         int
  useUDP       bool
  mongoURL     string
}

func NewServer(workersCount int, port int, useUDP bool, mongoURL string) *Server {
  return &Server{workersCount, port, useUDP, mongoURL}
}

func (self *Server) Serve() {
  hostport := "127.0.0.1:" + strconv.Itoa(self.port)

  session, err := mgo.Dial(self.mongoURL)
  if err != nil {
    panic(err)
  }
  db := session.DB("logger")
  collection := db.C("logs")

  channel := make(syslog.LogPartsChannel, self.workersCount * 2)
  handler := syslog.NewChannelHandler(channel)

  server := syslog.NewServer()
  server.SetFormat(syslog.RFC5424)
  server.SetHandler(handler)

  if (self.useUDP) {
    server.ListenUDP(hostport)
  } else {
    server.ListenTCP(hostport)
  }
  server.Boot()

  for i := 0; i < self.workersCount; i++ {
    go startWorker(channel, collection)
  }

  server.Wait()
}

func startWorker(c syslog.LogPartsChannel, collection *mgo.Collection) {
  for logParts := range c {
    collection.Insert(logParts)
  }
}
