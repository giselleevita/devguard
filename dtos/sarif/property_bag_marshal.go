package sarif

import (
	"encoding/json"
)

// MarshalJSON emits SARIF property-bag extensions alongside standard fields.
func (j PropertyBag) MarshalJSON() ([]byte, error) {
	properties := make(map[string]any, len(j.AdditionalProperties)+1)
	if len(j.Tags) > 0 {
		properties["tags"] = j.Tags
	}
	for name, value := range j.AdditionalProperties {
		if name == "tags" {
			continue
		}
		properties[name] = value
	}
	return json.Marshal(properties)
}

// MarshalJSON corrects the generated schema's non-standard ruleID capitalization.
func (j Result) MarshalJSON() ([]byte, error) {
	type resultAlias Result
	encoded, err := json.Marshal(resultAlias(j))
	if err != nil {
		return nil, err
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		return nil, err
	}
	if ruleID, ok := fields["ruleID"]; ok {
		delete(fields, "ruleID")
		fields["ruleId"] = ruleID
	}
	return json.Marshal(fields)
}
