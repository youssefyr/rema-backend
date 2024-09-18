// dbinit/dbinit.go
package dbinit

import (
	"log"
	"xira/db"
)

var Client *db.PrismaClient

func DbInit() {
	Client = db.NewClient()
	if err := Client.Prisma.Connect(); err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
}

func DbDisconnect() {
	if err := Client.Prisma.Disconnect(); err != nil {
		log.Fatalf("failed to disconnect Prisma client: %v", err)
	}
}
