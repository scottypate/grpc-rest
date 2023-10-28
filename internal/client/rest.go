package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/scottypate/grpc-rest/codegen/go/v1/vehicle"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "server:50051", "the address to connect to")
)

type Vehicle struct {
	Id    string `json:"id"`
	Vin   string `json:"vin"`
	Make  string `json:"make"`
	Model string `json:"model"`
	Year  int64  `json:"year"`
	Trim  string `json:"trim"`
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect to gRPC server %v", err)
	}

	defer conn.Close()
	client := pb.NewVehicleServiceClient(conn)

	r := gin.Default()

	r.GET("/vehicles", func(ctx *gin.Context) {
		res, err := client.GetVehicles(ctx, &pb.GetVehiclesRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"vehicles": res.Vehicles,
		})
	})
	r.GET("/vehicles/:vin", func(ctx *gin.Context) {
		vin := ctx.Param("vin")
		res, err := client.GetVehicle(ctx, &pb.GetVehicleRequest{Vin: vin})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"vehicle": res.Vehicle,
		})
	})
	r.POST("/vehicles", func(ctx *gin.Context) {
		var vehicle Vehicle

		err := ctx.ShouldBind(&vehicle)

		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		data := &pb.Vehicle{
			Id:    vehicle.Id,
			Vin:   vehicle.Vin,
			Make:  vehicle.Make,
			Model: vehicle.Model,
			Year:  vehicle.Year,
			Trim:  vehicle.Trim,
		}
		res, err := client.CreateVehicle(ctx, &pb.CreateVehicleRequest{
			Vehicle: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"vehicle": res.Vehicle,
		})
	})

	r.Run(":50001")
}
