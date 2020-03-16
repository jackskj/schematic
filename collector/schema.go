package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/proto"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jackskj/schematic/collector/pb"
	sc "github.com/jackskj/schematic/collector/schemas/pb"
)

var (
	schematicTarget = "localhost:15010"
	ctx             = context.Background()
)

// for demo purposes only, collector registers biometric schema.
func registerSchema() {
	// connecto to schematic layer
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(schematicTarget, grpc.WithInsecure()); err != nil {
		log.Fatalf("failed to dial schematic: %v", err)
	}

	scClient := pb.NewSchematicClient(conn)
	bio := &sc.Biometrics{}

	//registering schema with the the biometric file descriptor
	bioFd, _ := bio.Descriptor()
	reader, err := gzip.NewReader(bytes.NewReader(bioFd))
	if err != nil {
		log.Fatalf("Failed to decompress file descriptor: %s", err)
	}
	out, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Failed to read compressed file descriptor: %s", err)
	}

	descriptor := desc.FileDescriptorProto{}
	err = proto.Unmarshal(out, &descriptor)
	if err != nil {
		log.Fatalf("Failed unmarshal descriptor: %s", err)
	}

	fileDescriptor, _ := proto.Marshal(&descriptor)
	schema := pb.Schema{
		Topic:          "apple_watch",
		Version:        "v1",
		FileDescriptor: fileDescriptor,
	}
	if _, err := scClient.RegisterSchema(ctx, &schema); err != nil {
		log.Fatalf("failed to register schema: %v", err)
	}
}
