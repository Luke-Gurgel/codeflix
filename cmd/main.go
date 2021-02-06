package main

import (
	"os"

	"github.com/Luke-Gurgel/codeflix/application/grpc"
	"github.com/Luke-Gurgel/codeflix/infra/db"
	"github.com/jinzhu/gorm"
)

var database *gorm.DB

func main() {
	database = db.ConnectDB(os.Getenv("env"))
	grpc.StartGrpcServer(database, 50051)
}
