package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	pb "main/proto"
	"fmt"
	"strings"
	"os"
	"context"
	"math/rand"
)

var fulcrumServers []string // Lista de direcciones IP de servidores Fulcrum

func RecibirMensajeInformantes(stream pb.BrokerLuna_RecibirMensajeInformantesServer) error {
	// Recibe mensaje de Informantes y elige al azar un servidor Fulcrum para enviarle la IP de este al informante
	request, err := stream.Recv()
	if err != nil {
		return err
	}

	// Selecciona un servidor Fulcrum al azar
	fulcrumIP := obtenerServidorFulcrumAlAzar()

	// Responde al informante con la IP del servidor Fulcrum
	respuesta := &pb.RespuestaInformante{FulcrumIp: fulcrumIP}
	if err := stream.Send(respuesta); err != nil {
		return err
	}

	return nil
}

func RecibirMensajeVanguardia(stream pb.BrokerLuna_RecibirMensajeVanguardiaServer) error {
	// Recibe mensaje de Vanguardia
	request, err := stream.Recv()
	if err != nil {
		return err
	}

	// Selecciona un servidor Fulcrum al azar
	fulcrumIP := obtenerServidorFulcrumAlAzar()

	// Crea una conexión gRPC con el servidor Fulcrum seleccionado
	conn, err := grpc.Dial(fulcrumIP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor Fulcrum: %v", err)
		return err
	}
	defer conn.Close()

	// Crea un cliente gRPC para comunicarse con el servidor Fulcrum
	fulcrumClient := pb.NewFulcrumClient(conn)

	// Envía el mensaje de la Vanguardia al servidor Fulcrum seleccionado
	respuestaFulcrum, err := fulcrumClient.ProcesarMensaje(context.Background(), &pb.MensajeVanguardia{Mensaje: request.Mensaje})
	if err != nil {
		log.Fatalf("Error al enviar mensaje a Fulcrum: %v", err)
		return err
	}

	// Reenvía la respuesta de Fulcrum al servidor Vanguardia
	respuestaVanguardia := &pb.RespuestaVanguardia{Respuesta: respuestaFulcrum.Respuesta}
	if err := stream.Send(respuestaVanguardia); err != nil {
		log.Fatalf("Error al enviar respuesta de Fulcrum a Vanguardia: %v", err)
		return err
	}

	return nil
}


func obtenerServidorFulcrumAlAzar() string {
	// Elige al azar un servidor Fulcrum de la lista
	index := rand.Intn(len(fulcrumServers))
	return fulcrumServers[index]
}

func main() {
	// Obtener las direcciones IP de los servidores Fulcrum desde las variables de entorno
	fulcrumServers := []string{
		os.Getenv("fulcrum1_server") + ":" + os.Getenv("fulcrum1_port"),
		os.Getenv("fulcrum2_server") + ":" + os.Getenv("fulcrum2_port"),
		os.Getenv("fulcrum3_server") + ":" + os.Getenv("fulcrum3_port"),
	}

	informante1Server := os.Getenv("informante1_server") + ":" + os.Getenv("informante1_port")
	informante2Server := os.Getenv("informante2_server") + ":" + os.Getenv("informante2_port")
	vanguardiaServer := os.Getenv("vanguardia_server") + ":" + os.Getenv("vanguardia_port")

	// Inicia el servidor BrokerLuna
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	pb.RegisterBrokerLunaServer(serv, &server{})

	log.Printf("BrokerLuna listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
