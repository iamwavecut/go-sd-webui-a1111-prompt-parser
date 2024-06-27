package promptparser

import (
	"fmt"
	"strconv"
	"strings"
)

func (p *PromptParser) parseExtraNetworks(text string) (PromptNode, error) {
	content := text[1 : len(text)-1]
	parts := strings.SplitN(content, ":", 2)
	if len(parts) != 2 {
		return PromptNode{}, fmt.Errorf("invalid extra networks format: %q", text)
	}

	return PromptNode{
		Type:  "extra_networks",
		Value: parts[0],
		Args:  []PromptNode{{Type: "plain", Value: parts[1]}},
	}, nil
}

func (p *PromptParser) parsePositive(text string) (PromptNode, error) {
	content := text[1 : len(text)-1]
	depth := 1
	for strings.HasPrefix(content, "(") && strings.HasSuffix(content, ")") {
		content = content[1 : len(content)-1]
		depth++
	}

	return PromptNode{
		Type:  "positive",
		Value: depth,
		Args:  []PromptNode{{Type: "plain", Value: content}},
	}, nil
}

func (p *PromptParser) parseNegative(text string) (PromptNode, error) {
	depth := 0
	content := text

	for strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]") {
		content = content[1 : len(content)-1]
		depth++
	}

	return PromptNode{
		Type:  "negative",
		Value: depth,
		Args:  []PromptNode{{Type: "plain", Value: strings.TrimSpace(content)}},
	}, nil
}

func (p *PromptParser) parseWeighted(text string) (PromptNode, error) {
	// Remove surrounding parentheses if present
	text = strings.TrimPrefix(text, "(")
	text = strings.TrimSuffix(text, ")")

	parts := strings.SplitN(text, ":", 2)
	if len(parts) != 2 {
		return PromptNode{}, fmt.Errorf("invalid weighted format: %s", text)
	}

	weight, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return PromptNode{}, fmt.Errorf("invalid weight: %s", parts[1])
	}

	return PromptNode{
		Type:  "weighted",
		Value: weight,
		Args:  []PromptNode{{Type: "plain", Value: parts[0]}},
	}, nil
}

func (p *PromptParser) Parse(text string) ([]PromptNode, error) {
	parts := strings.Split(text, ",")
	result := make([]PromptNode, 0, len(parts))

	for _, part := range parts {
		node, err := p.parseNodeNonRecursive(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("error parsing node %q: %w", part, err)
		}
		result = append(result, node)
	}

	return result, nil
}

func (p *PromptParser) parseNodeNonRecursive(text string) (PromptNode, error) {
	text = strings.TrimSpace(text)

	switch {
	case strings.HasPrefix(text, "<") && strings.HasSuffix(text, ">"):
		return p.parseExtraNetworks(text)
	case strings.HasPrefix(text, "(") && strings.HasSuffix(text, ")"):
		if strings.Contains(text, ":") {
			return p.parseWeighted(text)
		}
		return p.parsePositive(text)
	case strings.HasPrefix(text, "[") && strings.HasSuffix(text, "]"):
		return p.parseSquareBrackets(text)
	default:
		return PromptNode{Type: "plain", Value: text}, nil
	}
}

func (p *PromptParser) parseSquareBrackets(text string) (PromptNode, error) {
	switch {
	case strings.Contains(text, "|"):
		return p.parseAlternate(text)
	case strings.Count(text, ":") == 2:
		return p.parseScheduledFull(text)
	case strings.HasPrefix(text[1:], ":"):
		return p.parseScheduledTo(text)
	case strings.Contains(text, "::"):
		return p.parseScheduledFrom(text)
	case strings.Contains(text, ":"):
		return p.parseScheduledTo(text)
	default:
		return p.parseNegative(text)
	}
}
func (p *PromptParser) parseAlternate(text string) (PromptNode, error) {
	content := text[1 : len(text)-1]
	parts := strings.Split(content, "|")
	args := make([]PromptNode, len(parts))
	for i, part := range parts {
		args[i] = PromptNode{Type: "plain", Value: strings.TrimSpace(part)}
	}
	return PromptNode{
		Type: "alternate",
		Args: args,
	}, nil
}

// This function is no longer needed and can be removed

func (p *PromptParser) parseScheduledFull(text string) (PromptNode, error) {
	content := text[1 : len(text)-1]
	parts := strings.SplitN(content, ":", 3)
	if len(parts) != 3 {
		return PromptNode{}, fmt.Errorf("invalid scheduled full format: %s", text)
	}
	from := PromptNode{Type: "plain", Value: strings.TrimSpace(parts[0])}
	to := PromptNode{Type: "plain", Value: strings.TrimSpace(parts[1])}
	number, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return PromptNode{}, fmt.Errorf("invalid number in scheduled full: %s", parts[2])
	}
	return PromptNode{
		Type:  "scheduled_full",
		Value: number,
		Args:  []PromptNode{from, to},
	}, nil
}

func (p *PromptParser) parseScheduledFrom(text string) (PromptNode, error) {
	content := text[1 : len(text)-1]
	parts := strings.SplitN(content, "::", 2)
	if len(parts) != 2 {
		return PromptNode{}, fmt.Errorf("invalid scheduled from format: %s", text)
	}
	from := PromptNode{Type: "plain", Value: strings.TrimSpace(parts[0])}
	number, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return PromptNode{}, fmt.Errorf("invalid number in scheduled from: %s", parts[1])
	}
	return PromptNode{
		Type:  "scheduled_from",
		Value: number,
		Args:  []PromptNode{from},
	}, nil
}

func (p *PromptParser) parseScheduledTo(text string) (PromptNode, error) {
	content := text[1 : len(text)-1]
	parts := strings.SplitN(content, ":", 3)
	if len(parts) == 2 {
		to := PromptNode{Type: "plain", Value: strings.TrimSpace(parts[0])}
		number, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return PromptNode{}, fmt.Errorf("invalid number in scheduled to: %s", parts[1])
		}
		return PromptNode{
			Type:  "scheduled_to",
			Value: number,
			Args:  []PromptNode{to},
		}, nil
	} else if len(parts) == 3 && parts[0] == "" {
		to := PromptNode{Type: "plain", Value: strings.TrimSpace(parts[1])}
		number, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return PromptNode{}, fmt.Errorf("invalid number in scheduled to: %s", parts[2])
		}
		return PromptNode{
			Type:  "scheduled_to",
			Value: number,
			Args:  []PromptNode{to},
		}, nil
	}
	return PromptNode{}, fmt.Errorf("invalid scheduled to format: %s", text)
}

// This function is no longer needed and can be removed

