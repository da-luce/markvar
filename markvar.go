package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
)

func main() {
    // Define command line flags
    markdownFilePtr := flag.String("md", "", "Path to the Markdown file")
    jsonFilePtr := flag.String("json", "", "Path to the JSON file")
    flag.Parse()

    // Check if both flags are provided
    if *markdownFilePtr == "" || *jsonFilePtr == "" {
        fmt.Println("Please specify both the Markdown file and the JSON file.")
        flag.Usage()
        return
    }

    // Read Markdown file
    mdContent, err := ioutil.ReadFile(*markdownFilePtr)
    if err != nil {
        fmt.Println("Error reading Markdown file:", err)
        return
    }

    // Read JSON file
    jsonContent, err := ioutil.ReadFile(*jsonFilePtr)
    if err != nil {
        fmt.Println("Error reading JSON file:", err)
        return
    }

    // Parse JSON into a map
    var mappings map[string]string
    if err := json.Unmarshal(jsonContent, &mappings); err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    // Process Markdown content
    updatedContent, unusedMappings, unmatchedTags, err := processMarkdown(string(mdContent), mappings)
    if err != nil {
        fmt.Println("Error processing Markdown:", err)
        return
    }

    // Warn about unused mappings
    if len(unusedMappings) > 0 {
        fmt.Println("Warning: Unused mappings in JSON file:", unusedMappings)
    }

    // Warn about unmatched tags
    if len(unmatchedTags) > 0 {
        fmt.Println("Warning: The following tags in the Markdown file have no corresponding mapping in the JSON file:", unmatchedTags)
    }

    // Write updated Markdown back to the original file
    if err := ioutil.WriteFile(*markdownFilePtr, []byte(updatedContent), 0644); err != nil {
        fmt.Println("Error writing updated Markdown file:", err)
    }
}

func processMarkdown(content string, mappings map[string]string) (string, []string, []string, error) {
    // Regular expression to match multi-line tags
    re := regexp.MustCompile(`<!--id:([a-zA-Z0-9]+)-->([\s\S]*?)<!---->`)
    usedMappings := make(map[string]bool)
    unmatchedTags := make([]string, 0)

    updatedContent := re.ReplaceAllStringFunc(content, func(match string) string {
        id := re.FindStringSubmatch(match)[1]
        val, exists := mappings[id]
        if exists {
            usedMappings[id] = true
            return fmt.Sprintf("<!--id:%s-->%s<!---->", id, val)
        } else {
            unmatchedTags = append(unmatchedTags, id)
            return match
        }
    })

    // Find unused mappings
    var unusedMappings []string
    for key := range mappings {
        if !usedMappings[key] {
            unusedMappings = append(unusedMappings, key)
        }
    }

    return updatedContent, unusedMappings, unmatchedTags, nil
}
