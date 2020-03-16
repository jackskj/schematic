package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/jackskj/schematic/schematic/pb"
)

const (
	schematicPort = 15010
)

var (
	err error
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", schematicPort))
	if err != nil {
		log.Fatal("%s", err)
	}
	grpcServer := grpc.NewServer()
	schematicSrv, err := newSchemticServer()
	if err != nil {
		log.Fatal("%s", err)
	}
	pb.RegisterSchematicServer(grpcServer, schematicSrv)
	grpcServer.Serve(lis)
}
