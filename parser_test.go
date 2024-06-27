// Package promptparser provides functionality to parse and generate prompts
// for AI image generation models.
package promptparser

import (
	"reflect"
	"testing"
)

var testCases = []struct {
	name   string
	prompt string
	nodes  []PromptNode
}{
	{
		name:   "Simple prompt",
		prompt: "masterpiece, 1girl, blonde hair, <lora:Zelda_v1:0.5>, (chromatic aberration:0.7), sharp focus, hyper detailed, (fog:0.7), <hypernet:sxz-bloom:0.5>, [real photo], [highlight:dark:0.9], (((good anatomy)))",
		nodes: []PromptNode{
			{Type: "plain", Value: "masterpiece"},
			{Type: "plain", Value: "1girl"},
			{Type: "plain", Value: "blonde hair"},
			{Type: "extra_networks", Value: "lora", Args: []PromptNode{{Type: "plain", Value: "Zelda_v1:0.5"}}},
			{Type: "weighted", Value: 0.7, Args: []PromptNode{{Type: "plain", Value: "chromatic aberration"}}},
			{Type: "plain", Value: "sharp focus"},
			{Type: "plain", Value: "hyper detailed"},
			{Type: "weighted", Value: 0.7, Args: []PromptNode{{Type: "plain", Value: "fog"}}},
			{Type: "extra_networks", Value: "hypernet", Args: []PromptNode{{Type: "plain", Value: "sxz-bloom:0.5"}}},
			{Type: "negative", Value: 1, Args: []PromptNode{{Type: "plain", Value: "real photo"}}},
			{Type: "scheduled_full", Value: 0.9, Args: []PromptNode{{Type: "plain", Value: "highlight"}, {Type: "plain", Value: "dark"}}},
			{Type: "positive", Value: 3, Args: []PromptNode{{Type: "plain", Value: "good anatomy"}}},
		},
	},
	{
		name:   "Multiple lora prompt",
		prompt: "<lora:detailed_eye:0.5>, <lora:LORA MODEL:0.2>, kpop, asian 18 years old girl with pony hair, sbg, breast grab, breast lift, nipples, <lora:SelfBreastGrab:0.8>, couple, hetero, nude girl standing behind girl, make up, <lora:aiKorea_sd15:1>",
		nodes: []PromptNode{
			{Type: "extra_networks", Value: "lora", Args: []PromptNode{{Type: "plain", Value: "detailed_eye:0.5"}}},
			{Type: "extra_networks", Value: "lora", Args: []PromptNode{{Type: "plain", Value: "LORA MODEL:0.2"}}},
			{Type: "plain", Value: "kpop"},
			{Type: "plain", Value: "asian 18 years old girl with pony hair"},
			{Type: "plain", Value: "sbg"},
			{Type: "plain", Value: "breast grab"},
			{Type: "plain", Value: "breast lift"},
			{Type: "plain", Value: "nipples"},
			{Type: "extra_networks", Value: "lora", Args: []PromptNode{{Type: "plain", Value: "SelfBreastGrab:0.8"}}},
			{Type: "plain", Value: "couple"},
			{Type: "plain", Value: "hetero"},
			{Type: "plain", Value: "nude girl standing behind girl"},
			{Type: "plain", Value: "make up"},
			{Type: "extra_networks", Value: "lora", Args: []PromptNode{{Type: "plain", Value: "aiKorea_sd15:1"}}},
		},
	},
}

func TestPromptParser(t *testing.T) {
	parser := NewPromptParser()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parser.Parse(tc.prompt)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}

			if !reflect.DeepEqual(result, tc.nodes) {
				t.Errorf("Parse result mismatch.\nExpected: %+v\nGot: %+v", tc.nodes, result)
			}

			// Test generation
			generatedPrompt := GeneratePrompt(result)
			if generatedPrompt != tc.prompt {
				t.Errorf("Generated prompt mismatch.\nExpected: %s\nGot: %s", tc.prompt, generatedPrompt)
			}
		})
	}
}

func TestGeneratePrompt(t *testing.T) {

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GeneratePrompt(tc.nodes)
			if result != tc.prompt {
				t.Errorf("Generated prompt mismatch.\nExpected: %s\nGot: %s", tc.prompt, result)
			}
		})
	}
}
