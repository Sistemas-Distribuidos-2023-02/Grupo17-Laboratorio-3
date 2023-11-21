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

var fulcrumServers []string 
var reloj_de_vectores []int
var fulcrum_id int

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
<<<<<<< HEAD
	reloj_de_vectores[fulcrum_id]++

	if len(comandos) >= 2 && comandos[0] == "Consistencia" {
		rv := strings.Split(comandos[1],",")
		for i := 0; i < len(rv); i++ {
			temp := strconv.Atoi(rv[i])
			reloj_de_vectores[i] = math.Max(temp,reloj_de_vectores[i]);
		}
		return nil
	}
=======
>>>>>>> 0d48060e584e8d1b2a2629551550f09cb2bed121

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

func schedule(f func(), interval time.Duration) *time.Ticker {
    ticker := time.NewTicker(interval)
    go func() {
        for range ticker.C {
            f()
        }
    }()
    return ticker
}

func Consistencia(){
	for i := 0; i < len(fulcrumServers); i++ {
		if i != fulcrum_id {
			ip = fulcrumServers[i]
			conn, err := grpc.Dial(ip, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("No se pudo conectar al servidor: %v", err)
				continue
			}
			client := pb.NewOMSClient(conn)
			stream_fulcrum, err := client.NotifyBidirectional(context.Background())
			if err != nil {
				log.Fatalf("Error al abrir el flujo bidireccional: %v", err)
				continue
			}

			msg := "Consistencia "
			for i := 0; i < len(fulcrumServers); i++ {
				if i + 1 == len(fulcrumServers){
					msg += strconv.Itoa(reloj_de_vectores[i]);
				} else {
					mas += strconv.Itoa(reloj_de_vectores[i]) + ",";
				}
			}

			mensaje_send := &pb.Request{Message: msg}
			if err := stream_fulcrum.Send(mensaje_send); err != nil {
				log.Fatalf("Error al enviar mensaje: %v", err)
				continue
			}

			mensaje_recv, err := stream_fulcrum.Recv()
			if err != nil {
				log.Fatalf("Error al recibir el mensaje: %v", err)
				return err
			}

			
			rv := strings.Split(strings.Split(mensaje_recv.Reply," ")[1],",")
			for i := 0; i < len(rv); i++ {
				temp := strconv.Atoi(rv[i])
				reloj_de_vectores[i] = math.Max(temp,reloj_de_vectores[i]);
			}
		}
	}
}

func main() {
	fulcrumServers = []string{
		"localhost:50051",
		"localhost:50052",
		"localhost:50053",
	}

	reloj_de_vectores = []int{
		0,0,0,
	}

	for i := 0; i < len(fulcrumServers); i++ {
		if(fulcrumServers[i] == "localhost:50051"){
			fulcrum_id = i
		}
		
	}

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
