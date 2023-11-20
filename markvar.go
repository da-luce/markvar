package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func findMarkvarFile() (string, error) {
    cwd, err := os.Getwd()
    if err != nil {
        return "", err
    }

    filePath := filepath.Join(cwd, ".markvar")
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return "", fmt.Errorf(".markvar file not found in the current directory")
    }

    return filePath, nil
}

func processContent(content string, mappings map[string]string) (string, []string, []string, map[string]string) {
    usedMappings := make(map[string]bool)
    unmatchedPlaceholders := make([]string, 0)
    updatedContent := content

    for key := range mappings {
        placeholder := fmt.Sprintf("{{var:%s}}", key)
        if strings.Contains(updatedContent, placeholder) {
            updatedContent = strings.ReplaceAll(updatedContent, placeholder, mappings[key])
            usedMappings[key] = true
        } else {
            unmatchedPlaceholders = append(unmatchedPlaceholders, placeholder)
        }
    }

    // Removing unused mappings
    for key := range mappings {
        if !usedMappings[key] {
            delete(mappings, key)
        }
    }

    return updatedContent, unmatchedPlaceholders, make([]string, 0), mappings
}

func main() {
    markvarFilePath, err := findMarkvarFile()
    if err != nil {
        fmt.Println(err)
        return
    }

    jsonContent, err := ioutil.ReadFile(markvarFilePath)
    if err != nil {
        fmt.Println("Error reading .markvar file:", err)
        return
    }

    var mappings map[string]string
    if err := json.Unmarshal(jsonContent, &mappings); err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    mdFilePath := "example.md" // Replace with the actual markdown file path
    mdContent, err := ioutil.ReadFile(mdFilePath)
    if err != nil {
        fmt.Println("Error reading Markdown file:", err)
        return
    }

    updatedContent, unmatchedPlaceholders, _, updatedMappings := processContent(string(mdContent), mappings)

    if len(unmatchedPlaceholders) > 0 {
        fmt.Println("Warning: The following placeholders have no corresponding mapping and will be removed:", unmatchedPlaceholders)
    }

    // Write updated Markdown content
    if err := ioutil.WriteFile(mdFilePath+".updated", []byte(updatedContent), 0644); err != nil {
        fmt.Println("Error writing updated Markdown file:", err)
        return
    }

    // Update and write the .markvar file with removed unused mappings
    updatedJSON, err := json.MarshalIndent(updatedMappings, "", "    ")
    if err != nil {
        fmt.Println("Error marshalling updated JSON:", err)
        return
    }
    if err := ioutil.WriteFile(markvarFilePath, updatedJSON, 0644); err != nil {
        fmt.Println("Error writing updated .markvar file:", err)
        return
    }

    fmt.Println("Updated Markdown file and .markvar file successfully.")
}
