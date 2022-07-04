// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"bytes"
)

func main() {
	b := bytes.NewBuffer(make([]byte,0,255))
	str:="12" + "34"
	b.Grow(int(str))
	fmt.Println(b)
}
      