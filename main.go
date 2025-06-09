package main

import (
	"context"
	"log"
	"net"
	"time"

	"legoas/legoas/proto"
	pb "legoas/legoas/proto"
	"legoas/services/account"
	"legoas/services/menu"
	"legoas/services/office"
	"legoas/services/role"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Mongo connection error: %v", err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	roleRepo := role.NewRoleRepository(db)
	roleService := role.NewRoleServiceServer(roleRepo)
	proto.RegisterRoleServiceServer(grpcServer, roleService)

	officeRepo := office.NewOfficeRepository(db)
	officeService := office.NewOfficeServiceServer(officeRepo)
	proto.RegisterOfficeServiceServer(grpcServer, officeService)

	menuRepo := menu.NewMenuRepository(db)
	menuService := menu.NewMenuServiceServer(menuRepo)
	proto.RegisterMenuServiceServer(grpcServer, menuService)

	grpcServer := grpc.NewServer()
	pb.RegisterAccountServiceServer(grpcServer, account.NewAccountServiceServer(mongoClient))

	reflection.Register(grpcServer)

	log.Println("MongoDB connected.")
	log.Println("gRPC server running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
