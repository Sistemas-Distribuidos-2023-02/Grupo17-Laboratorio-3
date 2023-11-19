package main

import (
	"net"
	"log"
	pb "main/proto"
	"google.golang.org/grpc"
	"strings"
	"fmt"
	"strconv"
	"os"
	"context"
	"sync"
)

type RelojVector struct {
    X int
    Y int
    Z int
}

func agregarRegistro(servidor *ServidorFulcrum, sector, base string, cantidad, oscuridad int) {
    // Implementar lógica para agregar un registro y actualizar reloj y log
}

func borrarRegistro(servidor *ServidorFulcrum, sector, base string) {
    // Implementar lógica para borrar un registro y actualizar reloj y log
}

func renombrarBase(servidor *ServidorFulcrum, sector, base, nuevoNombre string) {
    // Implementar lógica para renombrar una base y actualizar reloj y log
}

func actualizarValor(servidor *ServidorFulcrum, sector, base string, nuevoValor int) {
    // Implementar lógica para actualizar el valor de una base y actualizar reloj y log
}

func propagarCambios(servidor *ServidorFulcrum) {
    // Implementar lógica para propagar cambios y mantener la consistencia eventual
}

// Función para manejar problemas de consistencia y realizar un merge.
func manejarConsistencia(servidores ...*ServidorFulcrum) {
    // Implementar lógica para manejar la consistencia y realizar un merge
}

func main(){

	//Coneccion con el servidor.
	//SE DEBE DIFERENCIAR LOS DATANODES POR EL PUERTO 50052 SERA NODO 1 Y 50053 NODO 2
	listener, err := net.Listen("tcp", "localhost:" + os.Getenv("data_node_port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	archivo ,_  = os.Create("DATA.txt")
	serv := grpc.NewServer()
	personas = make(map[int]string)
	pb.RegisterOMSServer(serv, &server{})
	log.Printf("server listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}