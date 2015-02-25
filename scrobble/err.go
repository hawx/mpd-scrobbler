package scrobble

import "fmt"

type Err struct {
	name string
	err  error
}

func (e *Err) Error() string {
	return fmt.Sprintf("%v: %v", e.name, e.err.Error())
}
