package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/uuid"
	pb "github.com/scottypate/grpc-rest/codegen/go/v1/vehicle"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	DatabaseConnection()
}

func main() {
	log.Println("gRPC server is running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterVehicleServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}
}

var (
	port = flag.Int("port", 50051, "gRPC server port")
	err  error
	DB   *gorm.DB
)

type server struct {
	pb.UnimplementedVehicleServiceServer
}

type Vehicle struct {
	Id    string `json:"id"`
	Vin   string `json:"vin"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int64  `json:"year"`
	Trim  string `json:"trim"`
}

func DatabaseConnection() {
	host := "postgres"
	port := "5432"
	dbName := "postgres"
	dbUser := "postgres"
	password := "LocalDevelopmentOnly"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(Vehicle{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	log.Println("Database connection successful...")
}

func CreateId(v *pb.Vehicle) string {
	return uuid.NewSHA1(
		uuid.MustParse("f2f82775-c1cb-4c45-b536-a79fce94c18b"),
		[]byte(
			strings.ToUpper(
				strings.ToUpper(v.GetVin())+
					strings.ToUpper(v.GetMake()),
			),
		),
	).String()
}

func (*server) CreateVehicle(ctx context.Context, req *pb.CreateVehicleRequest) (*pb.CreateVehicleResponse, error) {
	log.Println("Create Vehicle")
	vehicle := req.GetVehicle()
	vehicle.Id = CreateId(vehicle)

	data := Vehicle{
		Id:    vehicle.GetId(),
		Vin:   vehicle.GetVin(),
		Make:  vehicle.GetMake(),
		Model: vehicle.GetModel(),
		Year:  vehicle.GetYear(),
		Trim:  vehicle.GetTrim(),
	}

	if res := DB.Create(&data); res.Error != nil {
		return nil, res.Error
	}

	return &pb.CreateVehicleResponse{
		Vehicle: &pb.Vehicle{
			Id:    vehicle.GetId(),
			Vin:   vehicle.GetVin(),
			Make:  vehicle.GetMake(),
			Model: vehicle.GetModel(),
			Year:  vehicle.GetYear(),
			Trim:  vehicle.GetTrim(),
		},
	}, nil
}

func (*server) GetVehicle(ctx context.Context, req *pb.GetVehicleRequest) (*pb.GetVehicleResponse, error) {
	log.Println("Read Vehicle", req.GetVin())
	var vehicle pb.Vehicle
	res := DB.Find(&vehicle, "vin = ?", req.GetVin())
	if res.RowsAffected == 0 {
		return nil, res.Error
	}
	return &pb.GetVehicleResponse{
		Vehicle: &vehicle,
	}, nil
}

func (*server) GetVehicles(ctx context.Context, req *pb.GetVehiclesRequest) (*pb.GetVehiclesResponse, error) {
	fmt.Println("Read vehicles")
	vehicles := []*pb.Vehicle{}
	res := DB.Find(&vehicles)
	if res.RowsAffected == 0 {
		return nil, res.Error
	}
	return &pb.GetVehiclesResponse{
		Vehicles: vehicles,
	}, nil
}
