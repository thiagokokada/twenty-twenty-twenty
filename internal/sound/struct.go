package sound

import "github.com/gopxl/beep/v2"

type sound struct {
	name   string
	buffer *beep.Buffer
}

func (s sound) String() string {
	return s.name
}
