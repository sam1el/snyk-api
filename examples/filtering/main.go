// Package main demonstrates filtering and searching with the snyk-api library.
//
// This example shows how to:
// - Filter projects by type
// - Filter projects by origin
// - Search for specific projects
//
// Prerequisites:
// - Set SNYK_TOKEN environment variable
//
// Run:
//
//	go run examples/filtering/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/sam1el/snyk-api/pkg/apiclients/projects"
	"github.com/sam1el/snyk-api/pkg/client"
)

func main() {
	ctx := context.Background()

	baseClient, err := client.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() {
		if err := baseClient.Close(); err != nil {
			log.Printf("failed to close client: %v", err)
		}
	}()

	projectsClient := projects.NewProjectsClient(baseClient)
	orgID := "your-org-id-here" // Replace with actual org ID

	if orgID == "your-org-id-here" {
		fmt.Println("⚠️  Please set a valid organization ID in the code")
		fmt.Println("    Run: snyk-api orgs list")
		return
	}

	// Example 1: Filter projects by type (npm)
	fmt.Println("📦 Example 1: Finding all npm projects...")
	fmt.Println()
	npmType := "npm"
	limit := 20
	npmProjects, err := projectsClient.ListProjects(ctx, orgID, &projects.ListProjectsParams{
		Type:  &npmType,
		Limit: &limit,
	})
	if err != nil {
		log.Fatalf("Failed to list npm projects: %v", err)
	}

	fmt.Printf("Found %d npm projects:\n", len(npmProjects.Data))
	for i, project := range npmProjects.Data {
		fmt.Printf("%d. %s (Origin: %s)\n", i+1, project.Attributes.Name, project.Attributes.Origin)
	}

	// Example 2: Filter projects by origin (github)
	fmt.Println()
	fmt.Println("🐙 Example 2: Finding all GitHub projects...")
	fmt.Println()
	githubOrigin := "github"
	githubProjects, err := projectsClient.ListProjects(ctx, orgID, &projects.ListProjectsParams{
		Origin: &githubOrigin,
		Limit:  &limit,
	})
	if err != nil {
		log.Fatalf("Failed to list GitHub projects: %v", err)
	}

	fmt.Printf("Found %d GitHub projects:\n", len(githubProjects.Data))
	for i, project := range githubProjects.Data {
		fmt.Printf("%d. %s (Type: %s)\n", i+1, project.Attributes.Name, project.Attributes.Type)
	}

	// Example 3: Combined filter (GitHub + npm)
	fmt.Println()
	fmt.Println("🎯 Example 3: Finding GitHub npm projects...")
	fmt.Println()
	combinedProjects, err := projectsClient.ListProjects(ctx, orgID, &projects.ListProjectsParams{
		Origin: &githubOrigin,
		Type:   &npmType,
		Limit:  &limit,
	})
	if err != nil {
		log.Fatalf("Failed to list GitHub npm projects: %v", err)
	}

	fmt.Printf("Found %d GitHub npm projects:\n", len(combinedProjects.Data))
	for i, project := range combinedProjects.Data {
		fmt.Printf("%d. %s\n", i+1, project.Attributes.Name)
	}

	// Example 4: Client-side filtering (search for projects by name)
	fmt.Println()
	fmt.Println("🔍 Example 4: Searching for projects containing 'api'...")
	fmt.Println()
	allProjects, err := projectsClient.ListProjects(ctx, orgID, &projects.ListProjectsParams{
		Limit: &limit,
	})
	if err != nil {
		log.Fatalf("Failed to list projects: %v", err)
	}

	searchTerm := "api"
	matchingProjects := []projects.Project{}
	for _, project := range allProjects.Data {
		if strings.Contains(strings.ToLower(project.Attributes.Name), searchTerm) {
			matchingProjects = append(matchingProjects, project)
		}
	}

	fmt.Printf("Found %d projects containing '%s':\n", len(matchingProjects), searchTerm)
	for i, project := range matchingProjects {
		fmt.Printf("%d. %s (Type: %s, Origin: %s)\n",
			i+1, project.Attributes.Name, project.Attributes.Type, project.Attributes.Origin)
	}

	fmt.Println("\n✅ Filtering examples completed successfully!")
}
