package sugar

import (
	"fmt"
	"testing"
)

func Test_Error(t *testing.T) {
	err := execute()
	if err != nil {
		fmt.Println("error")
	}
}

func execute() error {
	var err FileErrors

	return err
}
