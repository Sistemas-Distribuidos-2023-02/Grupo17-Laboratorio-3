package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "main/proto"
)

func EnviarMensajeABrokerLuna(mensaje string, conn *grpc.ClientConn) error {
	client := pb.NewOMSClient(conn)

	// Construir el mensaje
	mensaje = "_ Informante " + mensaje
	request := &pb.Request{Message: mensaje}

	// Enviar el mensaje a BrokerLuna
	stream, err := client.NotifyBidirectional(context.Background())
	if err != nil {
		return fmt.Errorf("Error al abrir el flujo bidireccional: %v", err)
	}

	// Enviar el mensaje al servidor
	if err := stream.Send(request); err != nil {
		return fmt.Errorf("Error al enviar mensaje a BrokerLuna: %v", err)
	}

	// Recibir la respuesta del servidor
	resp, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("Error al recibir mensaje de BrokerLuna: %v", err)
	}

	// Imprimir la dirección IP recibida
	fmt.Printf("Dirección IP recibida: %s\n", resp.GetReply())

	// Enviar automáticamente el mensaje a la dirección IP recibida
	err = EnviarMensajeAFulcrum(mensaje, resp.GetReply())
	if err != nil {
		return err
	}

	// Recibir el mensaje de vuelta
	resp, err = stream.Recv()
	if err != nil {
		return fmt.Errorf("Error al recibir mensaje de vuelta: %v", err)
	}

	// Imprimir el mensaje de vuelta
	fmt.Printf("Mensaje recibido de vuelta: %s\n", resp.GetReply())

	return nil
}

func EnviarMensajeAFulcrum(mensaje, direccionIP string) error {
	// Implementa la lógica para enviar el mensaje a la dirección IP proporcionada.
	// Puedes usar las funciones estándar de Go para realizar operaciones de red, como 'net.Dial'.
	// En este ejemplo, simplemente imprimimos un mensaje simulando el envío a la dirección IP.
	fmt.Printf("Enviando mensaje '%s' a la dirección IP: %s\n", mensaje, direccionIP)
	return nil
}

func main() {
	conn, err := grpc.Dial("localhost:50070", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor BrokerLuna: %v", err)
	}
	defer conn.Close()

	lector := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Ingrese función (o 'exit' para salir): ")
		mensaje, _ := lector.ReadString('\n')
		mensaje = strings.TrimSuffix(mensaje, "\n")

		if mensaje == "exit" {
			break
		}

		err = EnviarMensajeABrokerLuna(mensaje, conn)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	fmt.Println("Programa finalizado.")
}
