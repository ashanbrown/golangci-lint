//args: -Enolintlint -Elll
//config: linters-settings.nolintlint.allow-leading-space=false
package p

import "fmt"

func nolintlint() {
	fmt.Println() //nolint:bob // leading space should be dropped
	fmt.Println() //nolint:bob // leading spaces should be dropped
	// note that the next lines will retain trailing whitespace when fixed
	fmt.Println() 
	fmt.Println() 
	fmt.Println() //nolint:alice,bob // enabled linter should be dropped
	fmt.Println() //nolint:alice, bob // enabled linter should be dropped but whitespace preserved
	fmt.Println() //nolint:alice,lll,bob,nolintlint // enabled linter should not be dropped when nolintlint is nolinted
}
