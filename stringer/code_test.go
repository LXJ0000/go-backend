package stringer

import (
	"fmt"
	"log"
	"testing"
)

func TestCode(t *testing.T) {
	code := CODE_OK
	fmt.Println("-- ", code, " --")
	log.Fatal("-- ", code.String(), " --")

}
