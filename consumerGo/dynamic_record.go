// this package determins the structure of the received message by retreiving file descriptor from the schematic layer
package main

import (
	"context"
	"google.golang.org/grpc"
	"log"

	// "github.com/golang/protobuf/proto"
	desc "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jackskj/schematic/consumerGo/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

var (
	schematicTarget = "localhost:15010"
	ctx             = context.Background()
	schema          *pb.Schema

	// biometric message descriptors, the message schma is completly known
	fd protoreflect.FileDescriptor
	md protoreflect.MessageDescriptor
)

// retrevieng bometric schema from schematic layer
func init() {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(schematicTarget, grpc.WithInsecure()); err != nil {
		log.Fatalf("failed to dial schematic: %v", err)
	}
	scClient := pb.NewSchematicClient(conn)

	// getting apple watch biometric data
	schemaRequest := pb.SchemaReq{
		Topic:   "apple_watch",
		Version: "v1",
	}
	for {
		schema, err = scClient.GetSchema(ctx, &schemaRequest)
		if err == nil {
			break
		}
	}

	// deserializing file descriptor
	descriptor := desc.FileDescriptorProto{}
	err = proto.Unmarshal(schema.GetFileDescriptor(), &descriptor)
	if err != nil {
		log.Fatalf("failed to unmarshal file descriptor: %v", err)
	}
	files := protoregistry.Files{}
	fd, err = protodesc.NewFile(&descriptor, &files)
	if err != nil {
		log.Fatalf("failed to retreive file descriptor: %v", err)
	}
	md = fd.Messages().ByName("Biometrics")

	// Print retreived file descriptor
	marshal := protojson.MarshalOptions{
		Multiline: true,
	}
	buff, _ := marshal.Marshal(&descriptor)
	log.Println("\nSuccessfully fetched schema:\n %s", string(buff))
}

// Marshals dynamic message whoes type was retreved from schematic layer
func MarshalBiometrics(buff []byte) *dynamicpb.Message {
	//instantiatine new dynamic message
	bio := dynamicpb.NewMessage(md)
	//marshaling byres received from kafka to the message
	err = proto.Unmarshal(buff, bio)
	if err != nil {
		log.Fatalf("error deserializeing message: %v", err)
	}
	return bio
}
