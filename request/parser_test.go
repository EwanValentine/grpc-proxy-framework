package request

import (
	"testing"

	"github.com/matryer/is"
)

func TestFlattenInputs(t *testing.T) {
	is := is.New(t)
	args, err := FlattenInputs(
		[]byte(`{"name": "Ewan"}`),
		map[string]string{"phone": "12345"},
		[]string{"phone"},
	)
	is.NoErr(err)
	is.Equal(map[string]interface{}{
		"name":  "Ewan",
		"phone": "12345",
	}, args)
}
