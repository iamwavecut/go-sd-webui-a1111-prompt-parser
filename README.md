# go-sd-webui-a1111-prompt-parser
Inspired by [StableCanvas/sd-webui-a1111-prompt-parser](https://github.com/StableCanvas/sd-webui-a1111-prompt-parser) written in TypeScript.

## Introduction
`go-sd-webui-a1111-prompt-parser` is a Stable Diffusion webUI (A1111) prompt parser package for Go. It parses Stable Diffusion model prompts into structured data for easy analysis and manipulation by developers.

## Features
- Parses A1111 format prompts, supporting the following syntax:
- - Plain text
- - Emphasis (parentheses)
- - Weight (brackets)
- - Lora models
- - Hypernetwork models
- - Negative prompts (square brackets)
- - Step Control (scheduling)
- Converts parsed results into JavaScript objects for easy manipulation and use
- Supports regenerating A1111 format prompts from JavaScript objects

## Installation
```bash
go get github.com/iamwavecut/go-sd-webui-a1111-prompt-parser
```

## Usage
```go
import (
	"github.com/iamwavecut/go-sd-webui-a1111-prompt-parser"
)

// ...

    // Parse a prompt
	parser := NewPromptParser()
	parsedNodes, err := parser.Parse("happy family, ([picnic]|beach day), [sunset:sunrise:0.5]")
    if err != nil {
        fmt.Println("Error parsing prompt:", err)
    }
    fmt.Println(parsedNodes)


    // Generate a prompt from the nodes slice
    prompt := GeneratePrompt(parsedNodes)
    fmt.Println(prompt)
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.