package main

import (
	"backend-brevet/config"
	"backend-brevet/seed/master"
	"fmt"
	"log"
)

func main() {
	db := config.ConnectDB()

	fmt.Println("Seeding users...")
	if err := master.SeedUsers(db); err != nil {
		log.Fatalf("failed seeding users: %v", err)
	}

	fmt.Println("Seeding prices...")
	if err := master.SeedPrices(db); err != nil {
		log.Fatalf("failed seeding prices: %v", err)
	}

	fmt.Println("Seeding courses...")
	if err := master.SeedCourses(db); err != nil {
		log.Fatalf("failed seeding courses: %v", err)
	}

	fmt.Println("Seeding done!")
}
