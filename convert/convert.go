package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Remove the unused function convert()
var filePath string = "./schema.prisma"

// Open the file
func main() {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the contents of the file
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	var content = string(data)

	re := regexp.MustCompile(`(?m)^model\s+(\w+)\s+\{([^}]+)\}`)

	// Find all matches in the data
	matches := re.FindAllStringSubmatch(content, -1)

	// Initialize an empty slice to hold JSON objects
	var models []map[string]interface{}

	// Loop through the matches
	for _, match := range matches {
		// Create a map to hold the model data
		model := make(map[string]interface{})

		// Extract model name
		modelName := match[1]
		model["model"] = modelName

		// Extract properties
		properties := match[2]
		propRe := regexp.MustCompile(`(\w+):\s+"?([^"]+)"?,?`)
		propMatches := propRe.FindAllStringSubmatch(properties, -1)

		// Add properties to the model map
		for _, propMatch := range propMatches {
			key := propMatch[1]
			value := propMatch[2]
			model[key] = value
		}

		// Append the model to the models slice
		models = append(models, model)
	}

	// Convert models to JSON
	jsonData, err := json.MarshalIndent(models, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Print the JSON data
	fmt.Println(string(jsonData))
}
