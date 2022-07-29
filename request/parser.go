package request

import (
	"encoding/json"
)

type Args map[string]interface{}

// FlattenInputs -
func FlattenInputs(body []byte, params map[string]string, pathArgs []string) (map[string]interface{}, error) {
	args := make(Args)

	if len(body) > 0 {
		var a map[string]interface{}
		if err := json.Unmarshal(body, &a); err != nil {
			return args, err
		}

		for k, v := range a {
			args[k] = v
		}
	}

	for _, arg := range pathArgs {
		args[arg] = params[arg]
	}

	return args, nil
}
