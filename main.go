package main

import (
	"Peony/Peony_backend/routers"
)

func main() {
	router := routers.InitRoute()
	router.Run()
}
