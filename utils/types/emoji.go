package types

import (
	"encoding/json"
	"fmt"
	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/objects/guild/emoji"
)

type Emoji struct {
	IsCustomEmoji bool
	Name          string
	Id            *uint64
}

func NewEmoji(emojiName *string, emojiId *uint64) Emoji {
	if emojiName == nil || *emojiName == "" {
		return Emoji{
			IsCustomEmoji: false,
			Name:          "",
			Id:            nil,
		}
	}

	return Emoji{
		IsCustomEmoji: emojiId != nil,
		Name:          *emojiName,
		Id:            emojiId,
	}
}

func (e Emoji) IntoGdl() *emoji.Emoji {
	if e.IsCustomEmoji {
		return &emoji.Emoji{
			Id:   objects.NewNullableSnowflake(*e.Id),
			Name: e.Name,
		}
	} else {
		if e.Name == "" {
			return nil
		} else {
			return &emoji.Emoji{
				Name: e.Name,
			}
		}
	}
}

type customEmoji struct {
	Name string  `json:"name"`
	Id   *uint64 `json:"id,string"`
}

func (e *Emoji) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw == nil {
		return fmt.Errorf("emoji data was nil")
	}

	switch v := raw.(type) {
	case string:
		e.IsCustomEmoji = false
		e.Name = v
	case map[string]interface{}:
		var decoded customEmoji
		if err := json.Unmarshal(data, &decoded); err != nil {
			return err
		}

		e.IsCustomEmoji = true
		e.Name = decoded.Name
		e.Id = decoded.Id
	default:
		return fmt.Errorf("unknown type")
	}

	return nil
}

func (e Emoji) MarshalJSON() ([]byte, error) {
	if e.IsCustomEmoji {
		return json.Marshal(customEmoji{
			Name: e.Name,
			Id:   e.Id,
		})
	} else {
		return json.Marshal(e.Name)
	}
}
