package sound

import "github.com/gopxl/beep"

type sound struct {
	name   string
	buffer *beep.Buffer
}

func (s sound) String() string {
	return s.name
}
