package promptparser

import (
	"fmt"
	"strings"
)

// GeneratePrompt converts a slice of PromptNodes into a string representation.
func GeneratePrompt(nodes []PromptNode) string {
	var sb strings.Builder
	for i, node := range nodes {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(generateNode(node))
	}
	return sb.String()
}

func generateNode(node PromptNode) string {
	switch node.Type {
	case "plain":
		return node.Value.(string)
	case "extra_networks":
		return fmt.Sprintf("<%s:%s>", node.Value, node.Args[0].Value)
	case "positive", "negative":
		content := generateNode(node.Args[0])
		depth := node.Value.(int)
		bracket := "("
		closeBracket := ")"
		if node.Type == "negative" {
			bracket = "["
			closeBracket = "]"
		}
		return strings.Repeat(bracket, depth) + content + strings.Repeat(closeBracket, depth)
	case "weighted":
		return fmt.Sprintf("(%s:%g)", generateNode(node.Args[0]), node.Value.(float64))
	case "alternate":
		var parts []string
		for _, arg := range node.Args {
			parts = append(parts, generateNode(arg))
		}
		return "[" + strings.Join(parts, "|") + "]"
	case "scheduled_full":
		return fmt.Sprintf("[%s:%s:%g]", generateNode(node.Args[0]), generateNode(node.Args[1]), node.Value.(float64))
	case "scheduled_from":
		return fmt.Sprintf("[%s::%g]", generateNode(node.Args[0]), node.Value.(float64))
	case "scheduled_to":
		return fmt.Sprintf("[:%s:%g]", generateNode(node.Args[0]), node.Value.(float64))
	default:
		return ""
	}
}
