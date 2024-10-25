package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.Start(":3000")
}
