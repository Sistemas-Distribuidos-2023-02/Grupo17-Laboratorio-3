package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	pb "main/proto"
	"strings"
	"os"
	"context"
	"math/rand"
)

var fulcrumServers []string // Lista de direcciones IP de servidores Fulcrum

type server struct {
	pb.UnimplementedOMSServer
}

func (s *server) NotifyBidirectional(stream pb.OMS_NotifyBidirectionalServer) error {
	request, err := stream.Recv()
	if err != nil {
		return err
	}
	comandos := strings.Split(request.Message," ")
	if len(comandos) != 6 {
		return nil
	}

	if comandos[1] == "Informante" {
		ip := fulcrumServers[rand.Intn(len(fulcrumServers))]
		mensaje := &pb.Response{Reply: ip}
		if err := stream.Send(mensaje); err != nil {
			return err
		}

	}
	if comandos[1] == "Vanguardia" {
		return RedireccionMensaje(request.Message,stream)
	}

	return nil
}

func RedireccionMensaje(msg string,stream pb.OMS_NotifyBidirectionalServer) error {
	ip := fulcrumServers[rand.Intn(len(fulcrumServers))]
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("No se pudo conectar al servidor: %v", err)
		return err
	}
	client := pb.NewOMSClient(conn)
	stream_fulcrum, err := client.NotifyBidirectional(context.Background())
	if err != nil {
        log.Fatalf("Error al abrir el flujo bidireccional: %v", err)
		return err
	}
	mensaje_send := &pb.Request{Message: msg}
	if err := stream_fulcrum.Send(mensaje_send); err != nil {
        log.Fatalf("Error al enviar mensaje: %v", err)
		return err
	}

	mensaje_recv, err := stream_fulcrum.Recv()
	if err != nil {
        log.Fatalf("Error al recibir el mensaje: %v", err)
		return err
	}

	mensaje := &pb.Response{Reply: mensaje_recv.Reply}
	if err := stream.Send(mensaje); err != nil {
		return err
	}
	return nil
}

func main() {
	// Obtener las direcciones IP de los servidores Fulcrum desde las variables de entorno
	fulcrumServers = []string{
		os.Getenv("fulcrum1_server") + ":" + os.Getenv("fulcrum1_port"),
		os.Getenv("fulcrum2_server") + ":" + os.Getenv("fulcrum2_port"),
		os.Getenv("fulcrum3_server") + ":" + os.Getenv("fulcrum3_port"),
	}

	//informante1Server := os.Getenv("informante1_server") + ":" + os.Getenv("informante1_port")
	//informante2Server := os.Getenv("informante2_server") + ":" + os.Getenv("informante2_port")
	//vanguardiaServer := os.Getenv("vanguardia_server") + ":" + os.Getenv("vanguardia_port")

	// Inicia el servidor BrokerLuna
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	pb.RegisterOMSServer(serv, &server{})

	log.Printf("BrokerLuna listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
