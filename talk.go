
package main

import (
  "bufio"
  "fmt"
  "io"
  "net"
  "os"
  "os/signal"
  "strconv"
)

func main() {

  fmt.Println("Received arguments", os.Args)

  out := make(chan string)

  go server(out, os.Args[1])

  go client(out, os.Args[2])

  interruptSignal := make(chan os.Signal)
  signal.Notify(interruptSignal, os.Interrupt)

  defer fmt.Println("Ending...")

  for {
    select {
      case <-interruptSignal:
        return
      case s := <-out:
        fmt.Print(s)
    }
  }
}

func server(out chan string, listenPort string) {

  port, _ := strconv.Atoi(listenPort)
  addr := net.TCPAddr{IP: net.IP{127, 0, 0, 1}, Port: port}
  listener, err := net.ListenTCP("tcp", &addr)

  if err == nil {

    defer listener.Close()

    for {

      conn, err := listener.AcceptTCP()

      if err == nil {

        go incomingHandler(out, conn)

      } else {

        out <- fmt.Sprintf("Some error\n(%v)\n", err)

      }

    }

  }

}

func incomingHandler(out chan string, conn *net.TCPConn) {

  defer conn.Close()

  reader := bufio.NewReader(conn)

  sender := conn.RemoteAddr().String()

  for {

    msg, err := reader.ReadString('\n')

    if err != nil {
      if err != io.EOF {
        out <- fmt.Sprintf("Error reading input\n(%v)\n", err)
      }
      return
    }

    // TODO get a timestamp
    out <- fmt.Sprintf("%v: %v", sender, msg)

  }

}

func client(out chan string, host string) {

  reader := bufio.NewReader(os.Stdin)

  for {

    msg, err := reader.ReadString('\n')

    if err != nil { return }

    sendMessage(out, host, msg)

  }

}

func sendMessage(out chan string, host, msg string) {

  addr, err := net.ResolveTCPAddr("tcp", host)

  if err != nil {
    out <- fmt.Sprintf("Unable to resolve host: %v\n(%v)\n", host, err)
    return
  }

  conn, err := net.DialTCP("tcp", nil, addr)

  if err != nil {
    out <- fmt.Sprintf("Unabled to connect to host: %v\n(%v)\n", host, err)
    return
  }

  conn.SetKeepAlive(false)
  defer func() {
    err := conn.Close()
    if err != nil {
      out <- fmt.Sprintf("Error closing socket\n(%v)\n", err)
    }
  }()

  n, err := conn.Write([]byte(msg))

  if err != nil {
    out <- fmt.Sprintf("Failed to write message to server. Wrote %v\n(%v)\n", n, err)
    return
  }
}

