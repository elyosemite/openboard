# API Design Standards — OpenBoard

## 1. Core Principle

API endpoints must be intentionally designed, consistently organized, and aligned with the OpenBoard domain model.

AI agents and developers MUST NOT create endpoints arbitrarily, introduce inconsistent naming conventions, duplicate existing resources, or expose internal implementation details through the API.

Every endpoint must represent a meaningful business capability or resource operation.

The API is a product contract, not merely a collection of HTTP handlers.

Before creating, modifying, or removing an endpoint, the developer MUST inspect the existing API structure, domain model, resource relationships, and established conventions.

## 2. API Architecture

OpenBoard must expose a RESTful HTTP API with resource-oriented URLs.

Use the following base structure:

```text
/api/v1/{resource}
/api/v1/{resource}/{resource_id}
/api/v1/{resource}/{resource_id}/{sub-resource}
/api/v1/{resource}/{resource_id}/{sub-resource}/{sub-resource_id}
```

Requirements:

- Use plural nouns for resource names.
- Use lowercase kebab-case for multiword resource names.
- Use nouns instead of verbs in resource URLs.
- Use HTTP methods to express the intended operation.
- Use consistent resource identifiers.
- Avoid unnecessary nesting.
- Avoid exposing database table names or internal implementation details.
- Keep URL structures predictable and discoverable.

Examples:

```http
GET    /api/v1/organizations
POST   /api/v1/organizations
GET    /api/v1/organizations/{organization_id}
PATCH  /api/v1/organizations/{organization_id}
DELETE /api/v1/organizations/{organization_id}
```

## 3. Domain-Oriented Resource Organization

API resources must follow the OpenBoard domain architecture.

The primary domain hierarchy is:

```text
Organization
    └── Workspace
          └── Team
          └── Project
                ├── Work Items
                ├── Eras
                ├── Modules
                ├── Views
                ├── Pages
                └── Analytics
```

The exact ownership and association rules must follow the approved domain model. Do not infer relationships solely from this illustrative hierarchy.

Organize endpoints by domain:

### Organizations

```http
GET    /api/v1/organizations
POST   /api/v1/organizations
GET    /api/v1/organizations/{organization_id}
PATCH  /api/v1/organizations/{organization_id}
DELETE /api/v1/organizations/{organization_id}
```

### Workspaces

```http
GET    /api/v1/workspaces
POST   /api/v1/workspaces
GET    /api/v1/workspaces/{workspace_id}
PATCH  /api/v1/workspaces/{workspace_id}
DELETE /api/v1/workspaces/{workspace_id}
```

### Teams

```http
GET    /api/v1/teams
POST   /api/v1/teams
GET    /api/v1/teams/{team_id}
PATCH  /api/v1/teams/{team_id}
DELETE /api/v1/teams/{team_id}
```

### Projects

```http
GET    /api/v1/projects
POST   /api/v1/projects
GET    /api/v1/projects/{project_id}
PATCH  /api/v1/projects/{project_id}
DELETE /api/v1/projects/{project_id}
```

These examples illustrate resource conventions, not an instruction to implement every endpoint immediately.

## 4. HTTP Method Semantics

Use HTTP methods consistently:

| Method | Purpose |
|---|---|
| GET | Retrieve resources |
| POST | Create resources or execute explicitly modeled commands |
| PUT | Replace a resource representation |
| PATCH | Partially update a resource |
| DELETE | Remove a resource |

Requirements:

- GET requests must not cause business-state changes.
- POST must not be used as a generic substitute for every operation.
- PUT must represent replacement semantics.
- PATCH must represent partial updates.
- DELETE must follow the resource's documented deletion policy.
- Use appropriate HTTP status codes.
- Do not return HTTP 200 for every operation.

Common response codes:

```text
200 OK
201 Created
202 Accepted
204 No Content
400 Bad Request
401 Unauthorized
403 Forbidden
404 Not Found
409 Conflict
422 Unprocessable Entity
429 Too Many Requests
500 Internal Server Error
```

## 5. Resource Relationships and Nesting

Nested endpoints should only be introduced when the relationship is meaningful to the API consumer.

For example:

```http
GET /api/v1/projects/{project_id}/work-items
```

This is appropriate because work items are queried within a project context.

However, do not automatically nest every resource under its parent.

Prefer:

```http
GET /api/v1/work-items/{work_item_id}
```

instead of requiring every client to retrieve a work item through its complete parent hierarchy.

Rules:

- Use nesting for contextual collection queries and relationship-specific operations.
- Preserve canonical top-level resource endpoints.
- Avoid deeply nested URLs.
- Do not duplicate the same operation under multiple URL structures without a clear justification.
- Validate that nested resources actually belong to their specified parent.

## 6. Query, Filtering, Sorting, and Pagination

Collection endpoints must support consistent query conventions.

Example:

```http
GET /api/v1/projects?page=1&page_size=20
GET /api/v1/projects?sort_by=created_at&sort_order=desc
GET /api/v1/projects?status=active
GET /api/v1/work-items?project_id={id}&status=open
```

Standard query parameters:

| Parameter | Purpose |
|---|---|
| page | Page number, where supported |
| page_size | Number of returned records |
| sort_by | Field used for sorting |
| sort_order | asc or desc |
| search | Text search |
| filter | Additional filtering, when appropriate |

Requirements:

- Define pagination behavior consistently.
- Apply reasonable maximum page sizes.
- Validate all filter and sorting fields.
- Never interpolate user-provided sorting or filtering expressions directly into SQL.
- Define a consistent response structure for paginated collections.
- Use cursor pagination when the access pattern and dataset justify it.

## 7. Actions and Business Operations

Not every business operation is naturally represented by CRUD.

When an operation represents a meaningful domain command, model it explicitly.

Examples:

```http
POST /api/v1/projects/{project_id}/archive
POST /api/v1/projects/{project_id}/restore
POST /api/v1/work-items/{work_item_id}/assign
POST /api/v1/work-items/{work_item_id}/transitions
POST /api/v1/eras/{era_id}/complete
```

These command-style endpoints are acceptable when they represent explicit domain operations that cannot be expressed cleanly as ordinary resource updates.

However:

- Do not create action endpoints for ordinary field updates.
- Prefer PATCH when changing resource attributes.
- Do not create generic endpoints such as `/execute`, `/process`, `/manage`, or `/update-status`.
- Business commands must have explicit validation, authorization, and documented side effects.
- Model complex state transitions through domain rules rather than arbitrary status updates.

## 8. API Contracts and Documentation

OpenBoard must maintain an OpenAPI specification as the authoritative API contract.

Requirements:

- Every endpoint must be documented.
- Every request and response must have a defined schema.
- Every endpoint must document its possible HTTP status codes.
- Authentication and authorization requirements must be explicit.
- Error response formats must be consistent.
- Request validation rules must be documented.
- Breaking changes must be handled through an explicit versioning strategy.

Do not implement undocumented endpoints.

Do not modify an existing API contract without evaluating its compatibility implications.

## 9. Multi-Tenancy and Authorization

All resource access must respect the OpenBoard organizational boundaries.

Requirements:

- Enforce tenant isolation on the backend.
- Never trust organization, workspace, team, or project identifiers supplied by the client without authorization checks.
- Validate resource ownership and relationships.
- Apply authorization to individual resource operations, not only to collection endpoints.
- Prevent insecure direct object reference (IDOR) vulnerabilities.
- Avoid leaking the existence or contents of resources belonging to unauthorized tenants.

Authorization must be enforced server-side regardless of frontend visibility.

## 10. API Error Standards

All API errors must follow one consistent response format.

Prefer RFC 9457 Problem Details for HTTP APIs.

Example:

```json
{
  "type": "https://openboard.dev/problems/validation-error",
  "title": "Validation failed",
  "status": 422,
  "detail": "The supplied project data is invalid.",
  "instance": "/api/v1/projects",
  "errors": [
    {
      "field": "name",
      "code": "required",
      "message": "Project name is required."
    }
  ]
}
```

Never expose stack traces, SQL errors, internal filesystem paths, secrets, or sensitive implementation details.

## 11. Mandatory Endpoint Design Process

Before implementing any API endpoint, the AI agent MUST:

1. Inspect the existing OpenAPI specification and route definitions.
2. Identify the domain responsible for the operation.
3. Identify whether the resource or operation already exists.
4. Determine the canonical resource URL.
5. Evaluate whether the operation is CRUD, a relationship operation, or a domain command.
6. Define request and response contracts.
7. Define authorization and tenant-isolation requirements.
8. Define validation and error behavior.
9. Evaluate backward compatibility.
10. Update the OpenAPI specification.
11. Implement automated tests for the endpoint.
12. Verify that naming, HTTP semantics, and response structures follow existing conventions.

If the correct resource boundary or endpoint design is unclear, the agent MUST stop and document the ambiguity rather than inventing an API structure.

## 12. Prohibited Practices

The following practices are prohibited:

- Creating endpoints without inspecting existing resources.
- Using inconsistent naming conventions.
- Using verbs indiscriminately in URLs.
- Creating duplicate endpoints for the same business operation.
- Exposing database implementation details.
- Returning different response structures for equivalent operations.
- Using arbitrary HTTP methods.
- Creating deeply nested resource URLs without justification.
- Allowing frontend-only authorization.
- Implementing undocumented endpoints.
- Changing established API conventions without an architectural decision.

## 13. API Architecture Governance

All significant API design decisions must be documented.

Maintain the following documentation:

```text
docs/
└── api/
    ├── api-design-standards.md
    ├── resource-catalog.md
    ├── endpoint-catalog.md
    ├── error-handling.md
    ├── pagination.md
    └── openapi.yaml
```

Use Architecture Decision Records (ADRs) for significant decisions involving resource boundaries, API versioning, authentication, authorization, and compatibility.

**Final rule: An AI agent must extend the existing API architecture, not invent a new one for each feature. Consistency, domain integrity, predictability, and long-term maintainability take precedence over implementation convenience.**
