package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type HexColour int

func (h HexColour) Int() int {
    return int(h)
}

func (h HexColour) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"#%06x"`, h)), nil
}

func (h *HexColour) UnmarshalJSON(data []byte) error {
	str := strings.TrimPrefix(string(data), `"`)
	str = strings.TrimPrefix(str, "#")
	str = strings.TrimSuffix(str, `"`)

	i, err := strconv.ParseInt(str, 16, 32)
	if err != nil {
        return err
    }

	if i < 0 || i > 0xFFFFFF {
        return fmt.Errorf("invalid hex colour: %s", str)
    }

	*h = HexColour(i)
	return nil
}