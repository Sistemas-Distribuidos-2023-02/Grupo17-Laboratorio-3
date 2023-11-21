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
	mensaje = "_ Vanguardia " + mensaje
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

	// Imprimir el mensaje recibido
	mensajeRecibido := resp.GetReply()
	fmt.Printf("Mensaje recibido de vuelta: %s\n", mensajeRecibido)

	comandos := strings.Split(mensaje," ")
	log := ""
	for i := 2; i < len(comandos); i++ {
		if i + 1 == len(comandos) {
			log += comandos[i]
		} else {
			log += comandos[i] + " "
		}
	}
	escribirEnArchivo := fmt.Sprintf("%s %s\n", log, mensajeRecibido)
	if err := escribirEnLog("log_vanguardia.txt", escribirEnArchivo); err != nil {
		return fmt.Errorf("Error al escribir en el archivo de registro: %v", err)
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
