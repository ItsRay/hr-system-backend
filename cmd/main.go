package cmd

import (
	"log"
	"net/http"
)

func main() {
	db := database.InitDB() // 初始化 DB
	router := http.NewServeMux()

	employees.RegisterRoutes(router, db)
	leaves.RegisterRoutes(router, db)
	positionlevels.RegisterRoutes(router, db)

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", router)
}
