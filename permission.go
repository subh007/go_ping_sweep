// permission
package go_ping_sweep

import (
	"fmt"
	"os"
)

// This funciton return the admin priviledge for the
// executing.
func IsAdmin() bool {
	if os.Getuid() == 0 {
		return true
	} else {
		fmt.Println("must be run with the root priviledge.")
		os.Exit(-1)
	}
	return false
}
