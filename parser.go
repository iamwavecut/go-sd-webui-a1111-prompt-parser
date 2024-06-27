package promptparser

// PromptNode represents a node in the prompt structure.
type PromptNode struct {
	Type  string
	Value interface{}
	Args  []PromptNode
}

// PromptParser is responsible for parsing prompts.
type PromptParser struct{}

// NewPromptParser creates a new PromptParser instance.
func NewPromptParser() *PromptParser {
	return &PromptParser{}
}

