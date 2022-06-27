package types

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils"
	"time"
)

type DateTimeLocal time.Time

var format = "2006-01-02T15:04"

func NewDateTimeLocalFromPtr(t *time.Time) *DateTimeLocal {
	if t == nil {
		return nil
	}

	return utils.Ptr(DateTimeLocal(*t))
}

func TimeOrNil(d *DateTimeLocal) *time.Time {
	if d == nil {
		return nil
	}

	return utils.Ptr(d.Time())
}

func (d DateTimeLocal) Time() time.Time {
	return time.Time(d)
}

func (d DateTimeLocal) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(d).Format(format))), nil
}

func (d *DateTimeLocal) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid DateTimeLocal")
	}

	tmp, err := time.Parse(format, string(b[1:len(b)-1]))
	if err != nil {
		return err
	}

	*d = DateTimeLocal(tmp)
	return nil
}
