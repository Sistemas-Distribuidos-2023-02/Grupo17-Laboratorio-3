package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"math/rand"
	"time"
	"log"
	"google.golang.org/grpc"
	"context"
	pb "main/proto"
	
)

func EnviarMensajeABrokerLuna(mensaje string, conn *grpc.ClientConn) error {
	client := pb.NewBrokerLunaClient(conn)

	// Envía el mensaje a BrokerLuna
	resp, err := client.EnviarMensaje(context.Background(), &pb.Request{Message: mensaje})
	if err != nil {
		return fmt.Errorf("Error al enviar mensaje a BrokerLuna: %v", err)
	}

	// Imprime la dirección IP recibida
	fmt.Printf("Dirección IP recibida: %s\n", resp.GetIpAddress())

	// Envía automáticamente el mensaje a la dirección IP recibida
	err = EnviarMensajeAFulcrum(mensaje, resp.GetIpAddress())
	if err != nil {
		return err
	}

	// Recibe el mensaje de vuelta
	resp, err = client.RecibirMensajeDeDireccionIP(context.Background(), &pb.Request{Message: mensaje, IpAddress: resp.GetIpAddress()})
	if err != nil {
		return fmt.Errorf("Error al recibir mensaje de vuelta: %v", err)
	}

	// Imprime el mensaje de vuelta
	fmt.Printf("Mensaje recibido de vuelta: %s\n", resp.GetMessage())

	return nil
}

func EnviarMensajeAFulcrum(mensaje, direccionIP string) error {
	// Aquí puedes implementar la lógica para enviar el mensaje a la dirección IP proporcionada.
	// Puedes usar las funciones estándar de Go para realizar operaciones de red, como 'net.Dial'.
	// En este ejemplo, simplemente imprimimos un mensaje simulando el envío a la dirección IP.
	fmt.Printf("Enviando mensaje '%s' a la dirección IP: %s\n", mensaje, direccionIP)
	return nil
}

func main() {
	conn, err := grpc.Dial(os.Getenv("boker_server") + ":" + os.Getenv("broker_port"), grpc.WithInsecure()) 
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
