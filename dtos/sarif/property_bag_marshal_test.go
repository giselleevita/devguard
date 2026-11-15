package sarif

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPropertyBagMarshalsAdditionalPropertiesInline(t *testing.T) {
	properties := PropertyBag{
		Tags: []string{"security"},
		AdditionalProperties: map[string]any{
			"precision": "high",
		},
	}

	encoded, err := json.Marshal(properties)
	if err != nil {
		t.Fatalf("marshal property bag: %v", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatalf("unmarshal property bag JSON: %v", err)
	}
	if _, ok := fields["AdditionalProperties"]; ok {
		t.Fatalf("extension properties were nested: %s", encoded)
	}

	var tags []string
	if err := json.Unmarshal(fields["tags"], &tags); err != nil {
		t.Fatalf("decode tags: %v", err)
	}
	if len(tags) != 1 || tags[0] != "security" {
		t.Fatalf("unexpected tags: %v", tags)
	}
	var precision string
	if err := json.Unmarshal(fields["precision"], &precision); err != nil {
		t.Fatalf("decode precision: %v", err)
	}
	if precision != "high" {
		t.Fatalf("unexpected precision: %q", precision)
	}
}

func TestResultMarshalsRuleIdWithSarifCasing(t *testing.T) {
	ruleID := "example-rule"
	encoded, err := json.Marshal(Result{RuleID: &ruleID})
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	if !strings.Contains(string(encoded), `"ruleId":"example-rule"`) {
		t.Fatalf("missing SARIF ruleId field: %s", encoded)
	}
	if strings.Contains(string(encoded), `"ruleID"`) {
		t.Fatalf("non-standard ruleID field was emitted: %s", encoded)
	}
}
