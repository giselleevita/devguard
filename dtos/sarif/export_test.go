package sarif

import (
	"encoding/json"
	"testing"
)

func TestExportRuleProperties(t *testing.T) {
	var rule ReportingDescriptor
	if err := json.Unmarshal([]byte(`{"id":"example","properties":{"tags":["security"],"precision":"high"}}`), &rule); err != nil {
		t.Fatal(err)
	}
	for name, properties := range map[string]PropertyBag{
		"decoded":       *rule.Properties,
		"after storage": {AdditionalProperties: rule.Properties.Fields()},
		"stored":        {AdditionalProperties: map[string]any{"tags": []any{"security"}, "precision": "high"}},
	} {
		t.Run(name, func(t *testing.T) {
			data, err := json.Marshal(properties)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != `{"precision":"high","tags":["security"]}` {
				t.Errorf("unexpected property bag: %s", data)
			}
		})
	}
	id := "example"
	data, err := json.Marshal(Result{RuleID: &id})
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	if string(fields["ruleId"]) != `"example"` || fields["ruleID"] != nil {
		t.Errorf("unexpected result: %s", data)
	}
}
