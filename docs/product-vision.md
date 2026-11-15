# OpenBoard: product vision

**English** | [Português (Brasil)](product-vision.pt-BR.md)

## 1. Overview

**OpenBoard** will be an open source project management and work tracking platform, inspired by established tools such as Azure DevOps Boards, Jira and Plane.

Its purpose is to offer a modern, flexible and extensible alternative so that companies, engineering teams and organizations can plan, organize, track and analyze their work in a single environment.

The platform is designed for two deployment models: **Self-hosted and Cloud**. The self-hosted model lets organizations install and run their own OpenBoard instance on internal servers or on the cloud infrastructure of their choice, keeping control over their data, configuration and access policies.

OpenBoard should serve everyone from small teams to organizations with complex structures, multiple departments, teams and projects, offering features for collaboration, planning, execution, documentation and activity analysis.

Its architecture is organized into functional subdomains with well-defined responsibilities, so the platform can evolve modularly without compromising the consistency of the core domain.

## 2. Organizational structure

The organizational structure is the main foundation of OpenBoard. It defines how users, teams and projects relate to each other within the platform.

The four core concepts are Organizations, Workspaces, Teams and Projects.

### 2.1. Organizations

Organizations represent the organizations or companies that use OpenBoard. Each organization has its own administrative context and can manage its people and their access, settings and different workspaces. An Organization centralizes administration and establishes the organizational boundaries for access and collaboration.

### 2.2. Workspaces

Workspaces are independent work areas within an Organization, used to separate different areas, departments, business units or work contexts. Each Workspace can have its own teams, projects, settings and organizational structures. The goal is to let a single organization keep different work environments without mixing activities, information and responsibilities.

### 2.3. Teams

Teams represent groups of people working together within a Workspace. Each team can have people, leads, settings and associated projects. Teams are essential for organizing responsibilities, distributing work and tracking execution.

### 2.4. Projects

Projects are the central units of work management and execution. Each project can have its own Work Items, Eras, Modules, Views, Pages and Analytics indicators. Projects can be associated with one or more teams, enabling different collaboration models according to each organization's needs.

## 3. Functional subdomains

The platform is organized into subdomains that represent different system capabilities. Some form the initial core of the product; others will be built progressively.

### 3.1. Work Items

Work Items are the fundamental unit of work in OpenBoard.

They represent activities that need to be planned, executed and tracked throughout a project's lifecycle.

Each item can contain a title, description, assignee, priority, status, due date, effort estimate, labels, attachments and comments.

A rich text editor lets users document requirements, record technical decisions, share images and attach files.

Work Items can also be related to each other through dependencies, blockers, related items, subtasks and work hierarchies.

Work Item types may include Epics, Features, Tasks, Bugs and Subtasks, customizable to each organization's needs.

### 3.2. Eras

Eras are execution periods or cycles used to organize a team's work.

Conceptually, they are similar to Scrum Sprints and to Cycles in tools such as Plane.

Each Era has a start date, an end date, goals and a set of associated Work Items.

Teams can plan their activities, track progress and review what was completed during a given period.

Tracking includes indicators such as overall progress, completed and pending items, effort estimates, burndown charts and burnup charts.

### 3.3. Modules

Modules split larger projects into smaller, more manageable units.

Each module can represent a feature, subsystem, software component, department or any other logical division defined by the organization.

Modules group related Work Items, making it possible to track specific initiatives without losing sight of the overall project.

### 3.4. Views

Views offer different ways to visualize and organize Work Items.

Views may include:

- List View: tabular organization of items.
- Board View: Kanban board organized by status.
- Calendar View: tracking of deadlines and dates.
- Timeline View: time-based view of activities and dependencies.

Users can create custom filters using properties such as assignee, priority, status, module, Era, labels and dates.

Views can be saved, shared with the team and reused.

### 3.5. Pages

Pages are documentation spaces integrated with project management.

They let teams document ideas, requirements, architectural decisions, meeting notes, technical specifications and other project knowledge.

Each page has a rich text editor with support for images, hyperlinks, attachments, advanced formatting and references to existing Work Items.

An idea recorded on a page can be turned into a Work Item, bringing documentation and execution closer together.

The architecture should also allow future integration with artificial intelligence features to help create, organize and summarize documents.

### 3.6. Analytics

Analytics provides a consolidated view of project data, so teams and managers can follow how work evolves through indicators and visualizations.

Planned indicators include:

- Distribution of Work Items by status, priority and assignee.
- Team velocity.
- Era progress.
- Burndown and burnup.
- Lead time and cycle time.
- Identification of blocked work.
- Delivery trends over time.

The goal is to make bottlenecks easier to spot, help understand team capacity and support data-driven decisions.

## 4. Supporting subdomains and functional infrastructure

Beyond the subdomains directly tied to organizing and executing work, OpenBoard will have cross-cutting capabilities that keep the platform running.

These subdomains will be built progressively.

### 4.1. Workflow Engine

Defines and runs Work Item workflows.

It will allow configuring states, transitions, rules and constraints, so each organization can adapt its operational processes.

### 4.2. Automations

An automation system based on events and conditions.

It will run actions automatically when certain things happen, such as status changes, assignments, due dates and Work Item creation.

### 4.3. Integrations

The layer responsible for integrating with external tools and services.

It may offer integrations with GitHub, GitLab, Bitbucket, Slack, Microsoft Teams and other platforms teams use.

### 4.4. RBAC

Role-Based Access Control.

It manages permissions and controls which resources and operations are available to each user within an Organization, Workspace, Team or Project.

### 4.5. Audit Logs

An audit system that records relevant events and changes within the platform.

It will make it possible to trace operations, identify who made changes and browse activity history.

### 4.6. Notifications

A notification system that informs users about relevant events.

It will include notifications about assignments, comments, mentions, status changes, deadlines and other project activity.

## 5. Project principles

OpenBoard development is guided by the following principles:

- **Open Source:** open source code, enabling auditing, collaboration and community-driven evolution.
- **Self-hosted First:** designed to run on your own infrastructure, with no mandatory dependency on external services.
- **Cloud Ready:** architecture prepared to offer a Cloud offering.
- **Data Ownership:** organizations own and control their data.
- **Modular Architecture:** subdomains with well-defined responsibilities and boundaries.
- **Extensibility:** architecture prepared for integrations, plugins, automations and community-built features.
- **Multi-tenant:** many organizations on a single installation, each with its data isolated.
- **Customizable Workflows:** workflows, properties and processes can be adapted to each organization's needs.
- **Developer Experience:** a particularly polished experience for software engineering teams, without limiting use to that audience.

## 6. Long-term goal

OpenBoard aims to become a complete work management platform, able to support different organizational models and execution methodologies.

Rather than copying features from commercial tools, the project seeks to offer an open, independent and extensible alternative in which organizations can build their own work processes without giving up data ownership, operational flexibility or technological freedom.

## 7. Technical specifications and architecture

OpenBoard should be built following sound software engineering practices, prioritizing maintainability, security, scalability, testability and architectural simplicity.

Technical decisions must consider both the self-hosted and the Cloud models, avoiding unnecessary dependencies on specific providers and ensuring the application can run on different infrastructure environments.

### 7.1. Technology stack

The main technologies are:

- **Backend:** Go (Golang).
- **Frontend:** React.
- **Frontend language:** TypeScript.
- **Frontend build and dependencies:** Node.js, with a consistent package manager.
- **Containerization:** Docker and Docker Compose.
- **CI/CD:** GitHub Actions.
- **Version control:** Git and GitHub.
- **Primary database:** PostgreSQL.

Technology choices should prioritize stability, ecosystem maturity, ease of maintenance and the availability of development and observability tooling.

### 7.2. Backend architecture

The backend is written in Go, with an architecture organized by domains and responsibilities.

The project should adopt Domain-Driven Design (DDD) principles where they add value to the domain model, keeping a clear separation between business rules, application, infrastructure and communication interfaces.

Code organization should favor modularity, low coupling and high cohesion, avoiding both a disorganized monolith and premature fragmentation into microservices.

The initial architecture follows a **Modular Monolith** approach, in which subdomains have well-defined boundaries and can evolve independently within a single application.

The main modules include:

- Organizations
- Workspaces
- Teams
- Projects
- Work Items
- Eras
- Modules
- Views
- Pages
- Analytics
- Workflow Engine
- Automations
- Integrations
- RBAC
- Audit Logs
- Notifications

Each module keeps its own responsibilities, business rules and communication interfaces.

The backend exposes a well-defined HTTP API, preferably RESTful, with consistent contracts, versioning and OpenAPI documentation.

Additional requirements:

- Centralized error handling.
- Input and business rule validation.
- Structured logging.
- Context propagation and support for cancelling operations.
- Proper concurrency management.
- Configuration through environment variables.
- Versioned database migrations.
- Unit, integration and architecture tests.
- Graceful shutdown.
- Health checks and readiness checks.

### 7.3. Frontend architecture

The frontend is written in React with TypeScript.

The application should follow a modular architecture organized by features and domains, avoiding an overly centralized structure or one based only on generic components.

The interface should be responsive, accessible and consistent, prioritizing an intuitive experience for teams that work with project management every day.

The structure should include:

- Reusable components.
- A centralized design system.
- Consistent state management.
- An API communication layer.
- Standardized error handling.
- Authentication and authorization controls in the UI.
- Routes organized by functional context.
- Validated forms.
- A data caching and synchronization strategy.
- Unit and component tests.
- Integration and end-to-end tests.

The frontend should keep a clear separation between presentation, application state, interaction rules and backend communication.

### 7.4. Database and persistence

PostgreSQL is the platform's primary relational database.

The data model should reflect domain boundaries and preserve the integrity of relationships between Organizations, Workspaces, Teams, Projects and other entities.

Requirements:

- Versioned, reproducible migrations.
- Referential integrity.
- Indexes planned around query patterns.
- Transactions for operations that require consistency.
- Pagination strategies for large collections.
- Auditing of relevant changes.
- Backup and restore strategies.

Database access goes through a well-defined persistence layer, so business rules do not depend directly on PostgreSQL-specific details.

### 7.5. Containerization and infrastructure

OpenBoard is distributed as Docker containers.

The repository has dedicated Dockerfiles for the backend and the frontend, prioritizing small images, reproducible builds and runtime security.

Dockerfiles should follow good practices such as:

- Multi-stage builds.
- Official, trusted base images.
- Running as a non-privileged user.
- Separating build and runtime dependencies.
- Reducing the attack surface.
- Configuration through environment variables.
- Health checks where applicable.

The project provides a Docker Compose setup to make local execution and self-hosted installation easier.

The initial infrastructure should run the platform's essential components without requiring mandatory external services.

### 7.6. CI/CD and development automation

Continuous integration and delivery are implemented with GitHub Actions.

The pipeline automates quality checks, tests and artifact generation, ensuring code changes do not compromise project stability.

The pipeline should cover:

- Formatting and lint checks.
- Static code analysis.
- Backend compilation.
- Frontend build.
- Unit tests.
- Integration tests.
- Dependency and vulnerability scanning.
- Docker image builds.
- Image publishing to the GitHub Container Registry (GHCR).
- Versioned artifacts.

The flow should account for different stages, such as pull requests, integration into the main branch, release creation and version publishing.

Releases should be versioned consistently, preferably following Semantic Versioning.

### 7.7. Quality and testing

Quality is an integral part of the development process.

The project has a testing strategy at several levels:

- **Unit Tests:** isolated validation of business rules and components.
- **Integration Tests:** validation of communication between modules and infrastructure.
- **End-to-End Tests:** validation of the platform's main user flows.
- **Architecture Tests:** verification of dependencies and architectural boundaries.

Code must pass automated checks before being merged into the main branch.

There should also be contribution guidelines, code review standards and technical documentation to guide contributors.

### 7.8. Security

Security is considered from the start of the architecture.

The platform should implement proper authentication, authorization and data isolation.

Initial requirements include:

- Secure credential storage.
- Protection against unauthorized access.
- Input validation and sanitization.
- Protection against common vulnerabilities, including the relevant OWASP Top 10 categories.
- Secure secrets management.
- Principle of least privilege.
- Access control across Organizations, Workspaces, Teams and Projects.
- Logging of sensitive operations.
- Updating and monitoring vulnerable dependencies.

Authorization must be validated on the backend and can never rely solely on controls implemented in the UI.

### 7.9. Observability and operations

OpenBoard should provide observability mechanisms to understand the application's behavior in development and production.

This includes:

- Structured logs.
- Application and infrastructure metrics.
- Distributed tracing when needed.
- Health checks.
- Error and failure monitoring.
- Identification of slow operations.
- Request correlation.

The architecture should allow future integration with OpenTelemetry and tools such as Prometheus and Grafana.

### 7.10. Documentation and maintenance

The repository keeps up-to-date technical documentation, including:

- Installation and local run instructions.
- Self-hosted configuration guide.
- Overall application architecture.
- Module documentation and boundaries.
- API contracts and documentation.
- Testing strategy.
- Contribution process.
- Build, release and upgrade procedures.

Relevant architectural decisions are documented as Architecture Decision Records (ADRs), making it possible to understand the reasoning behind technical choices and how they evolve over time.

### 7.11. Architectural principles

OpenBoard's technical development follows these principles:

1. **Simplicity before complexity:** do not introduce abstractions, services or technologies without a concrete need.
2. **Modularity:** keep domains isolated with well-defined responsibilities.
3. **Low coupling:** avoid unnecessary dependencies between modules.
4. **High cohesion:** keep related rules and features within their respective domains.
5. **Testability:** design components that can be tested in isolation.
6. **Secure by default:** treat security as an architectural requirement, not a later feature.
7. **Portability:** avoid mandatory dependencies on specific infrastructure providers.
8. **Incremental evolution:** let the platform grow without constant architectural rewrites.
