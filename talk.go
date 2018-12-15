
package main

import (
  "bufio"
  "fmt"
  "net"
  "os"
  "os/signal"
)

func main() {

  fmt.Println("Received arguments", os.Args)

  go server()

  go client(os.Args[1])

  interruptSignal := make(chan os.Signal)
  signal.Notify(interruptSignal, os.Interrupt)
  <-interruptSignal
  fmt.Println("Ending...")
}

func server() {

  addr := net.TCPAddr{Port: 28216}
  listener, err := net.ListenTCP("tcp", &addr)

  if err == nil {

    defer listener.Close()

    for {

      conn, err := listener.AcceptTCP()

      if err != nil {

        go incomingHandler(conn)

      }

    }

  }

}

func incomingHandler(conn *net.TCPConn) {

  defer conn.Close()

  in := make([]byte, 8)

  sender := conn.RemoteAddr().String()

  for {

    n, err := conn.Read(in)

    if err != nil { return }

    msg := string(in[:n])

    // TODO get a timestamp
    fmt.Println(sender, msg)

  }

}

func client(host string) {

  reader := bufio.NewReader(os.Stdin)

  for {

    msg, err := reader.ReadString('\n')

    if err != nil { return }

    fmt.Println("You", ":", msg)

    // TODO send the string

  }

}

