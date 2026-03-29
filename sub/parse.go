package sub

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ParseSubscription parses a JSON array of subscription objects from raw byte data.
func ParseSubscription(data []byte) ([]Subscription, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("parse subscription: %w", err)
	}
	delim, ok := t.(json.Delim)
	if !ok || delim != '[' {
		return nil, fmt.Errorf("parse subscription: expected JSON array, got %v", t)
	}
	subs := make([]Subscription, 0)
	for dec.More() {
		var sub Subscription
		if err := dec.Decode(&sub); err != nil {
			return nil, fmt.Errorf("parse subscription: decode element: %w", err)
		}
		subs = append(subs, sub)
	}
	return subs, nil
}
