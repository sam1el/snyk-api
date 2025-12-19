// Package main demonstrates advanced usage of the snyk-api library.
//
// This example shows how to:
// - Use custom client configuration
// - Handle pagination
// - Filter results
// - Use specific API versions
//
// Prerequisites:
// - Set SNYK_TOKEN environment variable
//
// Run:
//
//	go run examples/advanced/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sam1el/snyk-api/pkg/apiclients/projects"
	"github.com/sam1el/snyk-api/pkg/client"
)

func main() {
	ctx := context.Background()

	// Create client with custom configuration
	fmt.Println("🔧 Creating client with custom configuration...")
	baseClient, err := client.New(ctx,
		// Use specific API version
		client.WithVersion(client.DefaultAPIVersion),
		// Custom rate limiting (20 requests per second)
		client.WithRateLimit(20, time.Second),
		// Custom retry policy
		client.WithRetryPolicy(10, 200*time.Millisecond, 10*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() {
		if err := baseClient.Close(); err != nil {
			log.Printf("Failed to close client: %v", err)
		}
	}()

	// Example: List all projects with pagination
	orgID := "your-org-id-here" // Replace with actual org ID
	if orgID == "your-org-id-here" {
		fmt.Println("⚠️  Please set a valid organization ID in the code")
		fmt.Println("    Run: snyk-api orgs list")
		return
	}

	fmt.Printf("\n📦 Listing all projects in organization %s...\n\n", orgID)

	projectsClient := projects.NewProjectsClient(baseClient)
	allProjects := []projects.Project{}
	var cursor *string
	pageNum := 1

	for {
		fmt.Printf("Fetching page %d...\n", pageNum)

		limit := 10
		params := &projects.ListProjectsParams{
			Limit:         &limit,
			StartingAfter: cursor,
			// Filter by origin (optional)
			// Origin: ptr("github"),
			// Filter by type (optional)
			// Type: ptr("npm"),
		}

		projectsList, err := projectsClient.ListProjects(ctx, orgID, params)
		if err != nil {
			log.Fatalf("Failed to list projects: %v", err)
		}

		allProjects = append(allProjects, projectsList.Data...)
		fmt.Printf("  Found %d projects on this page\n", len(projectsList.Data))

		// Check if there's a next page
		if projectsList.Links == nil || projectsList.Links.Next == nil {
			break
		}

		// Extract cursor from next link (simplified)
		// In production, parse the URL properly
		cursor = ptr("next-page-cursor") // Placeholder
		pageNum++

		// Safety limit for example
		if pageNum > 5 {
			fmt.Println("  (Limiting to 5 pages for example)")
			break
		}
	}

	fmt.Printf("\n✅ Total projects fetched: %d\n\n", len(allProjects))

	// Group projects by type
	byType := make(map[string]int)
	for _, project := range allProjects {
		byType[project.Attributes.Type]++
	}

	fmt.Println("Projects by type:")
	for projType, count := range byType {
		fmt.Printf("  %s: %d\n", projType, count)
	}

	// Group projects by origin
	byOrigin := make(map[string]int)
	for _, project := range allProjects {
		byOrigin[project.Attributes.Origin]++
	}

	fmt.Println("\nProjects by origin:")
	for origin, count := range byOrigin {
		fmt.Printf("  %s: %d\n", origin, count)
	}

	fmt.Println("\n✅ Advanced example completed successfully!")
}

// ptr is a helper function to create string pointers
func ptr(s string) *string {
	return &s
}

