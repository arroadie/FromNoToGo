package main

import (
  "fmt"             // standard library
  "bytes"           // Podría muy bien estar en fmt, pero bueno
  "os"              // funcionalidades del sistema operacional
  "bufio"           // Buffed io
  "net/http"        // THE INTERNET
)

func main() {

  clients, err := os.Open("/Users/thiago/tierra.clients.txt")
  if err != nil {
    panic(err)
  }
  defer clients.Close() // Buena práctica

  iterador := bufio.NewScanner(clients)
  url := "http://localhost:8000/geodetails/v1.0/clients?client_id=geoservice-test" // Que nadie pregunte porque solo inicialize eso aquí
  total := 0

  for iterador.Scan() {
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(iterador.Text())))
    if err != nil {
      panic(err)
    }
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")

    cliente := &http.Client{} // optimización de memoria, me encantó
    respuesta, err := cliente.Do(req)
    if err != nil {
        fmt.Println(err)
    } // La tercera vez que escribo eso, capaz está bueno hacer un metodo para lo mismo, no?
    defer respuesta.Body.Close()

    if respuesta.Status == "200 OK" {
      total++
    }

  }

  fmt.Println(total, " requests realizados con exito")
}
