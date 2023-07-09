package grpc

import (
	"log"

	"github.com/romanchechyotkin/car_booking_service/proto/pb"
	"google.golang.org/grpc"
)

func NewCarsManagementClient(host, port string) pb.CarsManagementClient {
	conn, err := grpc.Dial("localhost:5500", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewCarsManagementClient(conn)

	return client
}
