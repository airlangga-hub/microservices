package catalog

import "os"

func main() {
	port := os.Getenv("CATALOG_PORT")
	
	NewRepository()
}