package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
    updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
    listCmd := flag.NewFlagSet("list", flag.ExitOnError)
    addCmd := flag.NewFlagSet("add", flag.ExitOnError)
    removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)

    updateMdFile := updateCmd.String("file", "", "Path to the Markdown file to update")
    addVarName := addCmd.String("name", "", "Variable name to add")
    addVarContent := addCmd.String("content", "", "Content of the variable")
    removeVarName := removeCmd.String("name", "", "Variable name to remove")

    if len(os.Args) < 2 {
        fmt.Println("expected 'update', 'list', 'add', or 'remove' subcommands")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "update":
        updateCmd.Parse(os.Args[2:])
        handleUpdate(*updateMdFile)
    case "list":
        listCmd.Parse(os.Args[2:])
        handleList()
    case "add":
        addCmd.Parse(os.Args[2:])
        handleAdd(*addVarName, *addVarContent)
    case "remove":
        removeCmd.Parse(os.Args[2:])
        handleRemove(*removeVarName)
    default:
        fmt.Println("expected 'update', 'list', 'add', or 'remove' subcommands")
        os.Exit(1)
    }
}

func readMappings() (map[string]string, error) {
    var mappings map[string]string

    data, err := ioutil.ReadFile(".markvar")
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(data, &mappings); err != nil {
        return nil, err
    }

    return mappings, nil
}

func writeMappings(mappings map[string]string) error {
    data, err := json.MarshalIndent(mappings, "", "  ")
    if err != nil {
        return err
    }

    return ioutil.WriteFile(".markvar", data, 0644)
}

func processContent(content string, mappings map[string]string) (string, []string, []string) {
    usedMappings := make(map[string]bool)
    unmatchedTags := make([]string, 0)

    updatedContent := content
    for key, val := range mappings {
        placeholder := fmt.Sprintf("{{var:%s}}", key)
        if strings.Contains(updatedContent, placeholder) {
            updatedContent = strings.ReplaceAll(updatedContent, placeholder, val)
            usedMappings[key] = true
        }
    }

    // Identifying unused and unmatched mappings
    var unusedMappings []string
    for key := range mappings {
        if !usedMappings[key] {
            unusedMappings = append(unusedMappings, key)
        }
    }

    for _, placeholder := range strings.Split(updatedContent, "{{var:") {
        if strings.Contains(placeholder, "}}") {
            varName := strings.Split(placeholder, "}}")[0]
            if _, exists := mappings[varName]; !exists {
                unmatchedTag := fmt.Sprintf("{{var:%s}}", varName)
                unmatchedTags = append(unmatchedTags, unmatchedTag)
                updatedContent = strings.ReplaceAll(updatedContent, unmatchedTag, "")
            }
        }
    }

    return updatedContent, unusedMappings, unmatchedTags
}


func handleUpdate(mdFile string) {
    if mdFile == "" {
        fmt.Println("Please specify the Markdown file path.")
        return
    }

    mdContent, err := ioutil.ReadFile(mdFile)
    if err != nil {
        fmt.Printf("Error reading Markdown file: %s\n", err)
        return
    }

    mappings, err := readMappings()
    if err != nil {
        fmt.Printf("Error reading mappings: %s\n", err)
        return
    }

    updatedContent, unusedMappings, unmatchedTags := processContent(string(mdContent), mappings)

    if len(unusedMappings) > 0 {
        fmt.Println("Warning: Unused mappings in JSON file:", unusedMappings)
    }

    if len(unmatchedTags) > 0 {
        fmt.Println("Warning: The following placeholders have no corresponding mapping and will be removed:", unmatchedTags)
    }

    outputFile := mdFile + ".updated"
    if err := ioutil.WriteFile(outputFile, []byte(updatedContent), 0644); err != nil {
        fmt.Printf("Error writing updated Markdown file: %s\n", err)
        return
    }

    fmt.Printf("File processed successfully. Updated file: %s\n", outputFile)
}


func handleRemove(varName string) {
    if varName == "" {
        fmt.Println("Please specify the variable name to remove.")
        return
    }

    mappings, err := readMappings()
    if err != nil {
        fmt.Printf("Error reading mappings: %s\n", err)
        return
    }

    if _, exists := mappings[varName]; exists {
        delete(mappings, varName)
        writeMappings(mappings)
    } else {
        fmt.Printf("Variable '%s' does not exist.\n", varName)
    }
}


func handleAdd(varName string, varContent string) {
    if varName == "" || varContent == "" {
        fmt.Println("Please specify both the variable name and content.")
        return
    }

    mappings, err := readMappings()
    if err != nil {
        fmt.Printf("Error reading mappings: %s\n", err)
        return
    }

    mappings[varName] = varContent
    writeMappings(mappings)
}


func handleList() {
    mappings, err := readMappings()
    if err != nil {
        fmt.Printf("Error reading mappings: %s\n", err)
        return
    }

    for key, value := range mappings {
        fmt.Printf("%s: %s\n", key, value)
    }
}
