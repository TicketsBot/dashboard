package types

import (
	"fmt"
	"strconv"
	"strings"
)

type Colour uint32

func (c Colour) Uint32() uint32 {
	return uint32(c)
}

func (c Colour) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"#%06x"`, c)), nil
}

func (c *Colour) UnmarshalJSON(b []byte) error {
	// Try to parse as int first
	if parsed, err := strconv.ParseUint(string(b), 10, 32); err == nil {
		*c = Colour(parsed)
		return nil
	}

	if len(b) < 2 {
		return fmt.Errorf("invalid colour")
	}

	s := string(b)[1 : len(b)-1]   // Trim quotes
	s = strings.TrimPrefix(s, `#`) // Trim # if present

	tmp, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return err
	}

	if tmp > 0xFFFFFF {
		return fmt.Errorf("colour value %v is out of range", tmp)
	}

	*c = Colour(tmp)
	return nil
}
