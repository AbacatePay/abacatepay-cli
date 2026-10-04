package webhook

import (
	"encoding/json"
	"strings"
)

// resourceAliases maps an event prefix to the data key holding its resource
// when the two differ (e.g. "payout.completed" carries data.withdraw).
var resourceAliases = map[string]string{
	"payout": "withdraw",
}

// parseEnvelope extracts the event type and the ID worth showing for a
// webhook message.
//
// API v2 payloads nest the resource under a key named after the event
// prefix ("transparent.completed" -> data.transparent.id), so that's tried
// first. Older payloads with a flat data.id and the top-level id (present on
// subscription events) are kept as fallbacks. The ID is "" when none apply.
func parseEnvelope(message []byte) (event, id string, err error) {
	var raw struct {
		Event string                     `json:"event"`
		ID    string                     `json:"id"`
		Data  map[string]json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(message, &raw); err != nil {
		return "", "", err
	}

	prefix, _, _ := strings.Cut(raw.Event, ".")

	candidates := []string{prefix}
	if alias, ok := resourceAliases[prefix]; ok {
		candidates = append(candidates, alias)
	}

	for _, key := range candidates {
		if id := nestedID(raw.Data[key]); id != "" {
			return raw.Event, id, nil
		}
	}

	if id := stringValue(raw.Data["id"]); id != "" {
		return raw.Event, id, nil
	}

	return raw.Event, raw.ID, nil
}

func nestedID(obj json.RawMessage) string {
	if len(obj) == 0 {
		return ""
	}

	var v struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(obj, &v); err != nil {
		return ""
	}

	return v.ID
}

func stringValue(v json.RawMessage) string {
	if len(v) == 0 {
		return ""
	}

	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return ""
	}

	return s
}
