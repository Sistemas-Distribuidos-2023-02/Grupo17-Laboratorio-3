package main

import (
	// "context"
	"log"
	"net"
	"os"
	"strings"
	"fmt"
	"google.golang.org/grpc"
	pb "main/proto"
)

type server struct {
	pb.UnimplementedOMSServer
}
func (s *server) NotifyBidirectional(stream pb.OMS_NotifyBidirectionalServer) error {
	// Recibir el mensaje del BrokerLuna
	request, err := stream.Recv()
	if err != nil {
		return err
	}
	comandos := strings.Split(request.Message," ")

	log := ""
	for i := 2; i < len(comandos); i++ {
		if i + 1 == len(comandos) {
			log += comandos[i]
		} else {
			log += comandos[i] + " "
		}
	}
	log += "\n"
	escribirEnLog("log.txt",log)

	// Ejecutar función correspondiente según el tipo de mensaje
	if comandos[1] == "Vanguardia" {
		return FuncionVanguardia(request.Message, stream)
	} else if comandos[1] == "Informante" {
		return FuncionInformante(request.Message, stream)
	}

	return nil
}

func escribirEnLog(nombreArchivo string, mensaje string) error {
	archivo, err := os.OpenFile(nombreArchivo, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer archivo.Close()

	_, err = archivo.WriteString(mensaje)
	return err
}

func FuncionVanguardia(msg string, stream pb.OMS_NotifyBidirectionalServer) error {
	fmt.Println("Mensaje recibido: ",msg)
	respuesta := &pb.Response{Reply: "Respuesta desde Vanguardia"}
	return stream.Send(respuesta)
}

func FuncionInformante(msg string, stream pb.OMS_NotifyBidirectionalServer) error {
	fmt.Println("Mensaje recibido: ",msg)
	comandos := strings.Split(msg," ")
	// if comandos[2] == "GetSoldados"
	respuesta := &pb.Response{Reply: "Respuesta desde Informante"}
	return stream.Send(respuesta)
}

func main() {
	// Inicia el servidor Fulcrum
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	pb.RegisterOMSServer(serv, &server{})

	log.Printf("FulcrumServer listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
