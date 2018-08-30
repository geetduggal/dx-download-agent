package main

import (
	"fmt"
	// The dxda package should contain all core functionality
	"github.com/geetduggal/dxda"
)

// The CLI is simply a wrapper around the dxda package
func main() {
	token, method := dxda.GetToken()
	fmt.Printf("Obtained token using %s\n", method)
	fmt.Println(token)
	fmt.Println(dxda.WhoAmI(token))
}
