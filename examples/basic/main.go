// Package main demonstrates basic usage of the snyk-api library.
//
// This example shows how to:
// - Create a client
// - List organizations
// - List projects in an organization
//
// Prerequisites:
// - Set SNYK_TOKEN environment variable
//
// Run:
//
//	go run examples/basic/main.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sam1el/snyk-api/pkg/apiclients/orgs"
	"github.com/sam1el/snyk-api/pkg/apiclients/projects"
	"github.com/sam1el/snyk-api/pkg/client"
)

func main() {
	ctx := context.Background()

	// Create base client (uses SNYK_TOKEN from environment)
	baseClient, err := client.New(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() {
		if err := baseClient.Close(); err != nil {
			log.Printf("Failed to close client: %v", err)
		}
	}()

	// List organizations
	fmt.Println("📋 Listing organizations...")
	orgsClient := orgs.NewOrgsClient(baseClient)
	limit := 5
	orgsList, err := orgsClient.ListOrganizations(ctx, &orgs.ListOrganizationsParams{
		Limit: &limit,
	})
	if err != nil {
		log.Fatalf("Failed to list organizations: %v", err)
	}

	fmt.Printf("Found %d organizations:\n\n", len(orgsList.Data))
	for i, org := range orgsList.Data {
		fmt.Printf("%d. %s\n", i+1, org.Attributes.Name)
		fmt.Printf("   ID: %s\n", org.Id)
		fmt.Printf("   Slug: %s\n", org.Attributes.Slug)
		if org.Attributes.Created != nil {
			fmt.Printf("   Created: %s\n", org.Attributes.Created.Format("2006-01-02"))
		}
		fmt.Println()
	}

	// List projects in first organization
	if len(orgsList.Data) > 0 {
		orgID := orgsList.Data[0].Id.String()
		fmt.Printf("📦 Listing projects in '%s'...\n\n", orgsList.Data[0].Attributes.Name)

		projectsClient := projects.NewProjectsClient(baseClient)
		projectLimit := 5
		projectsList, err := projectsClient.ListProjects(ctx, orgID, &projects.ListProjectsParams{
			Limit: &projectLimit,
		})
		if err != nil {
			log.Fatalf("Failed to list projects: %v", err)
		}

		if len(projectsList.Data) == 0 {
			fmt.Println("No projects found in this organization.")
		} else {
			fmt.Printf("Found %d projects:\n\n", len(projectsList.Data))
			for i, project := range projectsList.Data {
				fmt.Printf("%d. %s\n", i+1, project.Attributes.Name)
				fmt.Printf("   ID: %s\n", project.Id)
				fmt.Printf("   Type: %s\n", project.Attributes.Type)
				fmt.Printf("   Origin: %s\n", project.Attributes.Origin)
				fmt.Println()
			}
		}
	}

	fmt.Println("✅ Example completed successfully!")
}

