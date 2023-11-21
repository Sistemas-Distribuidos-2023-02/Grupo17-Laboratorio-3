package main

import (
	// "context"
	"log"
	"net"
	// "os"
	"strings"

	"google.golang.org/grpc"
	pb "main/proto"
)

type fulcrumServer struct {
	pb.UnimplementedOMSServer
}

func (s *fulcrumServer) NotifyBidirectional(stream pb.OMS_NotifyBidirectionalServer) error {
	// Recibir el mensaje del BrokerLuna
	request, err := stream.Recv()
	if err != nil {
		return err
	}

	// Dividir el mensaje en palabras
	comandos := strings.Split(request.Message, " ")
	if len(comandos) != 6 {
		return nil
	}

	// Ejecutar función correspondiente según el tipo de mensaje
	if comandos[1] == "Vanguardia" {
		return FuncionVanguardia(request.Message, stream)
	} else if comandos[1] == "Informante" {
		return FuncionInformante(request.Message, stream)
	}

	return nil
}

func FuncionVanguardia(msg string, stream pb.OMS_NotifyBidirectionalServer) error {
	// Lógica específica para mensajes de Vanguardia
	// ...

	// Enviar respuesta al BrokerLuna
	respuesta := &pb.Response{Reply: "Respuesta desde Vanguardia"}
	return stream.Send(respuesta)
}

func FuncionInformante(msg string, stream pb.OMS_NotifyBidirectionalServer) error {
	// Lógica específica para mensajes de Informante
	// ...

	// Enviar respuesta al BrokerLuna
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
	pb.RegisterOMSServer(serv, &fulcrumServer{})

	log.Printf("FulcrumServer listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
