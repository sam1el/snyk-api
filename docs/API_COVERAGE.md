# Snyk API Coverage

This document details the API coverage provided by the snyk-api tool.

## Summary

| API | Operations | Coverage |
|-----|------------|----------|
| v1 API | 93 | Full |
| REST API (Manual) | 103 | Full |
| **Total** | **196** | **Full** |

## v1 API Coverage (93 Operations)

### Projects (22 operations)

- `ListAggregatedIssues` - List all aggregated issues for a project
- `GetIssuePaths` - Get issue paths for a project issue
- `ListHistory` - List all project snapshots
- `GetHistoryAggregatedIssues` - Get aggregated issues for a snapshot
- `GetHistoryIssuePaths` - Get issue paths for a historical issue
- `GetDepGraph` - Get project dependency graph
- `ListIgnores` - List all ignores for a project
- `CreateIgnore` - Create a new ignore
- `DeleteIgnore` - Delete an ignore
- `ListJiraIssues` - List all Jira issues for a project
- `CreateJiraIssue` - Create a new Jira issue
- `GetSettings` - Get project settings
- `UpdateSettings` - Update project settings
- `Move` - Move a project to a different org
- `ListTags` - List project tags
- `AddTags` - Add tags to a project
- `RemoveTags` - Remove tags from a project
- `GetAttributes` - Get project attributes
- `UpdateAttributes` - Update project attributes
- `Deactivate` - Deactivate a project
- `Activate` - Activate a project
- `Get` - Get project details

### Testing (17 operations)

- `TestMaven` - Test a Maven package
- `TestMavenPackage` - Test a specific Maven package version
- `TestNpm` - Test an npm package
- `TestNpmPackage` - Test a specific npm package version
- `TestGoPkg` - Test a Go package
- `TestGoPkgPackage` - Test a specific Go package version
- `TestPip` - Test a Python pip package
- `TestPipPackage` - Test a specific pip package version
- `TestComposer` - Test a Composer package
- `TestComposerPackage` - Test a specific Composer package version
- `TestRubyGems` - Test a RubyGems package
- `TestRubyGemsPackage` - Test a specific RubyGems package version
- `TestSBT` - Test an SBT project
- `TestSBTPackage` - Test a specific SBT package version
- `TestGradle` - Test a Gradle project
- `TestGradlePackage` - Test a specific Gradle package version
- `TestDepGraph` - Test a dependency graph

### Organizations (14 operations)

- `ListMembers` - List organization members
- `InviteUser` - Invite a user to the organization
- `ViewPendingInvites` - View pending user invites
- `RevokeInvite` - Revoke a pending invite
- `RemoveMember` - Remove a member from the organization
- `UpdateMemberRole` - Update a member's role
- `GetSettings` - Get organization settings
- `UpdateSettings` - Update organization settings
- `ListProvisioningDetails` - List provisioning details
- `SetProvisioningDetails` - Set provisioning details
- `DeleteProvisioningDetails` - Delete provisioning details
- `GetLicenses` - Get organization licenses
- `SendInvite` - Send an invite to a user
- `ListNotificationSettings` - List notification settings

### Integrations (11 operations)

- `ListIntegrations` - List all integrations
- `GetIntegration` - Get a specific integration
- `CreateIntegration` - Create a new integration
- `UpdateIntegration` - Update an integration
- `DeleteIntegration` - Delete an integration
- `GetCredentials` - Get integration credentials
- `UpdateCredentials` - Update integration credentials
- `GetBrokerToken` - Get the broker token
- `ProvisionToken` - Provision a new token
- `SwitchToken` - Switch to a different token
- `ImportProject` - Import a project

### Reporting (9 operations)

- `ListLatestIssues` - List latest issues
- `ListIssues` - List issues with filters
- `GetLatestIssueCounts` - Get latest issue counts
- `GetIssueCounts` - Get issue counts with filters
- `GetLatestProjectCounts` - Get latest project counts
- `GetProjectCounts` - Get project counts with filters
- `ListTestCounts` - List test counts
- `GetTestCounts` - Get test counts with filters
- `GetDependencyCounts` - Get dependency counts

### Webhooks (5 operations)

- `CreateWebhook` - Create a new webhook
- `ListWebhooks` - List all webhooks
- `GetWebhook` - Get a specific webhook
- `DeleteWebhook` - Delete a webhook
- `PingWebhook` - Ping a webhook to test it

### Groups (8 operations)

- `ListMembers` - List all group members
- `AddMember` - Add a member to the group
- `ListOrgs` - List all organizations in the group
- `GetSettings` - Get group settings
- `UpdateSettings` - Update group settings
- `ListTags` - List group tags
- `DeleteTag` - Delete a group tag
- `GetRole` - Get role details

### Users (6 operations)

- `GetMyDetails` - Get current user details
- `GetUserDetails` - Get details for a specific user
- `GetOrgNotificationSettings` - Get org notification settings
- `UpdateOrgNotificationSettings` - Update org notification settings
- `GetProjectNotificationSettings` - Get project notification settings
- `UpdateProjectNotificationSettings` - Update project notification settings

### Monitor (1 operation)

- `MonitorDepGraph` - Monitor a dependency graph

## REST API Coverage (103 Operations)

### Organizations (48 operations)

- List, Get, Update organizations
- Memberships: List, Get, Update, Delete
- Invites: List, Create, Delete
- Service Accounts: List, Get, Create, Update, Delete, Rotate Secret
- Policies: List, Get, Delete
- Collections: List, Get, Create, Delete
- Settings: IaC, SAST, OpenSource
- Audit Logs: Search
- Projects: List, Get, Delete, Get SBOM
- Targets: List, Get, Delete
- Issues: List, Get
- SBOM Testing: Create, Get Job, Get Results
- Container Images: List, Get
- Apps: List, List Creations, List Installs
- Cloud: List Environments, List Scans
- Export: Create, Get

### Groups (27 operations)

- List, Get groups
- Memberships: List, Get, Delete
- Org Memberships: List
- Organizations: List
- Service Accounts: List, Get, Delete
- Policies: List, Get, Delete
- Settings: IaC, PR Template
- SSO Connections: List, List Users, Delete User
- Audit Logs: Search
- Issues: List, Get
- Assets: Search, Get
- App Installs: List
- Export: Create, Get
- Users: Get

### Tenants (17 operations)

- List, Get tenants
- Memberships: List, Get, Delete
- Roles: List, Get
- Broker Deployments: List, Get, List by Install
- Broker Connections: List, Get, Delete
- Broker Credentials: List, Get, Delete
- Broker Integrations: List

### Self (11 operations)

- Get current user
- Access Requests: List
- Apps: List, Get, Delete
- App Installs: List, Get, Revoke
- App Sessions: List, Get, Revoke

## CLI Commands

### Top-level Commands (Minimal REST - Generated)

```sh
snyk-api orgs list|get
snyk-api projects list|get|delete
snyk-api targets list|get|delete
snyk-api issues list|get
```

### v1 API Commands

```sh
snyk-api v1 groups [operations]
snyk-api v1 integrations [operations]
snyk-api v1 orgs [operations]
snyk-api v1 projects [operations]
snyk-api v1 reporting [operations]
snyk-api v1 test [operations]
snyk-api v1 users [operations]
snyk-api v1 webhooks [operations]
```

### Full REST API Commands (Manual)

```sh
snyk-api rest orgs list|get|memberships|invites|service-accounts|policies|...
snyk-api rest groups list|get|memberships|orgs|policies|...
snyk-api rest tenants list|get|memberships|roles|broker-deployments|...
snyk-api rest self get|apps|access-requests
```

## Quality Metrics

- **Build**: ✅ Pass
- **Lint**: ✅ 0 issues
- **Snyk Code**: ✅ 0 security issues
- **Snyk SCA**: ✅ 0 fixable vulnerabilities (1 license notice on transitive dependency)

## Architecture

```tree
pkg/
├── apiclients/
│   ├── orgs/         # Generated REST client (minimal)
│   ├── projects/     # Generated REST client (minimal)
│   ├── targets/      # Generated REST client (minimal)
│   ├── issues/       # Generated REST client (minimal)
│   ├── rest/         # Manual REST clients (full coverage)
│   │   ├── client.go # Base REST client with versioning
│   │   ├── orgs/     # 48 operations
│   │   ├── groups/   # 27 operations
│   │   ├── tenants/  # 17 operations
│   │   └── self/     # 11 operations
│   └── v1/           # Manual v1 clients (full coverage)
│       ├── client.go # Base v1 client
│       ├── projects/ # 22 operations
│       ├── testing/  # 17 operations
│       ├── orgs/     # 14 operations
│       ├── integrations/ # 11 operations
│       ├── reporting/ # 9 operations
│       ├── webhooks/ # 5 operations
│       ├── groups/   # 8 operations
│       ├── users/    # 6 operations
│       └── monitor/  # 1 operation
└── client/           # Base client infrastructure
```
