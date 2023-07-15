package defaults

import (
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNil(t *testing.T) {
	var myString *string
	ApplyDefaults(NewDefaultApplicator[*string](NilCheck[string], &myString, utils.Ptr("hello")))
	assert.Equal(t, "hello", *myString)
}
