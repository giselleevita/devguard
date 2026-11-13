package controllers

import (
	"encoding/json"
	"testing"

	"github.com/l3montree-dev/devguard/database/models"
	"github.com/stretchr/testify/require"
)

func TestFirstPartySARIFExport(t *testing.T) {
	report := firstPartyVulnsToSARIF("example", []models.FirstPartyVuln{{
		RuleID: "example.rule",
		RuleProperties: map[string]any{
			"tags": []any{"security"}, "precision": "high",
		},
	}})
	result, err := json.Marshal(report.Runs[0].Results[0])
	require.NoError(t, err)
	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(result, &fields))
	require.JSONEq(t, `"example.rule"`, string(fields["ruleId"]))
	require.NotContains(t, fields, "ruleID")

	properties, err := json.Marshal(report.Runs[0].Tool.Driver.Rules[0].Properties)
	require.NoError(t, err)
	require.JSONEq(t, `{"tags":["security"],"precision":"high"}`, string(properties))
}
