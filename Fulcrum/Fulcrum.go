package main

import (
	"context"
	"bufio"
	"strconv"
	"time"
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

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (s *server) NotifyBidirectional(stream pb.OMS_NotifyBidirectionalServer) error {
	// Recibir el mensaje del BrokerLuna
	request, err := stream.Recv()
	if err != nil {
		return err
	}
	comandos := strings.Split(request.Message," ")
	reloj_de_vectores[fulcrum_id]++

	if len(comandos) >= 2 && comandos[0] == "Consistencia" {
		rv := strings.Split(comandos[1],",")
		for i := 0; i < len(rv); i++ {
			temp,err := strconv.Atoi(rv[i])
			temp = int(temp)
			if err != nil {
				continue
			}
			reloj_de_vectores[i] = Max(temp,reloj_de_vectores[i]);
		}
		return nil
	}

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
    fmt.Println("Mensaje recibido: ", msg)
    comandos := strings.Split(msg, " ")
    if comandos[2] == "GetSoldados" {
        archivoPath := comandos[3] + ".txt"
        if _, err := os.Stat(archivoPath); err == nil {
            archivo, err := os.Open(archivoPath)
            if err != nil {
                return err
            }
            defer archivo.Close()
            scanner := bufio.NewScanner(archivo)
            for scanner.Scan() {
                partes := strings.Split(scanner.Text(), " ")
                if len(partes) == 2 {
                    palabra := partes[0]
                    numero := partes[1]
                    if palabra == comandos[4] {
                        respuesta := &pb.Response{Reply: numero}
                        return stream.Send(respuesta)
                    }
                }
            }
            respuesta := &pb.Response{Reply: "No existe la base " + comandos[4] + " en el sector " + comandos[3]}
            return stream.Send(respuesta)
        }
        respuesta := &pb.Response{Reply: "No existe el sector " + comandos[3]}
        return stream.Send(respuesta)
    }
    respuesta := &pb.Response{Reply: "Comando Inválido"}
    return stream.Send(respuesta)
}



func FuncionInformante(msg string, stream pb.OMS_NotifyBidirectionalServer) error {
	fmt.Println("Mensaje recibido: ", msg)
	comandos := strings.Split(msg, " ")
	switch comandos[2] {
	case "AgregarBase":
		if _, err := os.Stat(comandos[3] + ".txt"); err == nil {
			// Archivo de sector existe, revisar si la base ya existe
			if baseYaExiste(comandos[3]+".txt", comandos[4]) {
				respuesta := &pb.Response{Reply: "Esa base ya existe"}
				return stream.Send(respuesta)
			}
			// Base no existe, agregarla
			if err := agregarBase(comandos[3]+".txt", comandos[4], comandos[5]); err != nil {
				respuesta := &pb.Response{Reply: fmt.Sprintf("Error al agregar la base: %v", err)}
				return stream.Send(respuesta)
			}
		} else {
			// Archivo de sector no existe, crearlo y agregar la base
			if err := crearSectorYAgregarBase(comandos[3]+".txt", comandos[4], comandos[5]); err != nil {
				respuesta := &pb.Response{Reply: fmt.Sprintf("Error al crear el sector y agregar la base: %v", err)}
				return stream.Send(respuesta)
			}
		}

	case "RenombrarBase":
		if _, err := os.Stat(comandos[3] + ".txt"); err == nil {
			// Archivo de sector existe, actualizar el nombre de la base
			if err := renombrarBase(comandos[3]+".txt", comandos[4], comandos[5]); err != nil {
				respuesta := &pb.Response{Reply: fmt.Sprintf("Error al renombrar la base: %v", err)}
				return stream.Send(respuesta)
			}
		}

	case "ActualizarValor":
		if _, err := os.Stat(comandos[3] + ".txt"); err == nil {
			// Archivo de sector existe, actualizar el valor de la base
			if err := actualizarValor(comandos[3]+".txt", comandos[4], comandos[5]); err != nil {
				respuesta := &pb.Response{Reply: fmt.Sprintf("Error al actualizar el valor de la base: %v", err)}
				return stream.Send(respuesta)
			}
		}

	case "BorrarBase":
		if _, err := os.Stat(comandos[3] + ".txt"); err == nil {
			// Archivo de sector existe, borrar la base
			if err := borrarBase(comandos[3]+".txt", comandos[4]); err != nil {
				respuesta := &pb.Response{Reply: fmt.Sprintf("Error al borrar la base: %v", err)}
				return stream.Send(respuesta)
			}
		}

	default:
		respuesta := &pb.Response{Reply: "Comándo Inválido"}
		return stream.Send(respuesta)
	}

	respuesta := &pb.Response{Reply: "Comando exitoso"}
	return stream.Send(respuesta)
	return nil
}

func baseYaExiste(filename, nombreBase string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linea := scanner.Text()
		palabras := strings.Fields(linea)
		if len(palabras) >= 1 && palabras[0] == nombreBase {
			// La base ya existe
			return true
		}
	}

	return false
}

func agregarBase(filename, nombreBase, numero string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s %s\n", nombreBase, numero)); err != nil {
		return err
	}

	return nil
}

func crearSectorYAgregarBase(filename, nombreBase, numero string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s %s\n", nombreBase, numero)); err != nil {
		return err
	}

	return nil
}

func renombrarBase(filename, antiguoNombre, nuevoNombre string) error {
	lines, err := leerArchivo(filename)
	if err != nil {
		return err
	}

	for i, line := range lines {
		palabras := strings.Fields(line)
		if len(palabras) >= 1 && palabras[0] == antiguoNombre {
			// Actualizar el nombre de la base
			lines[i] = fmt.Sprintf("%s %s", nuevoNombre, palabras[1])
			break
		}
	}

	if err := escribirArchivo(filename, lines); err != nil {
		return err
	}

	return nil
}

func actualizarValor(filename, nombreBase, nuevoNumero string) error {
	lines, err := leerArchivo(filename)
	if err != nil {
		return err
	}

	for i, line := range lines {
		palabras := strings.Fields(line)
		if len(palabras) >= 1 && palabras[0] == nombreBase {
			// Actualizar el número de la base
			lines[i] = fmt.Sprintf("%s %s", palabras[0], nuevoNumero)
			break
		}
	}

	if err := escribirArchivo(filename, lines); err != nil {
		return err
	}

	return nil
}

func borrarBase(filename, nombreBase string) error {
	lines, err := leerArchivo(filename)
	if err != nil {
		return err
	}

	var nuevasLineas []string
	borrar := false

	for _, line := range lines {
		palabras := strings.Fields(line)
		if len(palabras) >= 1 && palabras[0] == nombreBase {
			// No añadir la línea al slice para borrar la base
			borrar = true
		} else {
			nuevasLineas = append(nuevasLineas, line)
		}
	}

	if borrar {
		if err := escribirArchivo(filename, nuevasLineas); err != nil {
			return err
		}
	}

	return nil
}

func leerArchivo(filename string) ([]string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func escribirArchivo(filename string, lines []string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
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

func Consistencia() error {
	for i := 0; i < len(fulcrumServers); i++ {
		if i != fulcrum_id {
			ip := fulcrumServers[i]
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
					msg += strconv.Itoa(reloj_de_vectores[i]) + ",";
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
				temp,err := strconv.Atoi(rv[i])
				if err == nil {
					continue
				}
				temp = int(temp)
				reloj_de_vectores[i] = Max(temp,reloj_de_vectores[i]);
			}
		}
	}
	return nil
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
