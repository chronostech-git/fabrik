package serialize

import "encoding/gob"

func Register(value any) {
	gob.Register(value)
}
