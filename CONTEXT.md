# OpenBoard

Open source project management and work tracking platform. This glossary is the shared language for the domain; implementation details live elsewhere.

## Language

### Organizational structure

**Instance**:
One installation of OpenBoard, self-hosted or Cloud, hosting any number of Organizations.
_Avoid_: Server, deployment, site

**Organization**:
The tenant: the outermost boundary of data, administration and access, whose data never crosses into another Organization.
_Avoid_: Tenant (in UI and domain code), company, account

**Workspace**:
A division inside one Organization (an area, department or business unit) that can restrict visibility but stays under the Organization's administration.
_Avoid_: Space, area, department

**Default Workspace**:
The Open Workspace created together with its Organization; otherwise an ordinary Workspace with no special privileges.
_Avoid_: Initial workspace, main workspace, general

**Open Workspace**:
A Workspace that every User with an Organization Membership in its Organization can enter.

**Restricted Workspace**:
A Workspace that only Users with a Workspace Membership (and Organization Owners) can enter; the default for every new Workspace except the Default Workspace.
_Avoid_: Private workspace

**Team**:
A group of Users who work together inside one Workspace, used for responsibility and planning, never for access control.
_Avoid_: Group, squad, crew

**Project**:
The central unit of planned and executed work, owned by exactly one Workspace.
_Avoid_: Board, initiative

**Project Team**:
A Team associated with a Project; the association never transfers ownership.

### People and access

**User**:
A person's identity and credentials, global to the Instance; by itself it grants access to no Organization data.
_Avoid_: Account, login, profile

**Username**:
The unique identifier of a User within an Instance.
_Avoid_: Login, handle, user ID

**Deactivated**:
The state of a User who can no longer sign in, while remaining the author of everything they did; only an Administrator deactivates or reactivates a User.
_Avoid_: Deleted user, disabled, banned

**Temporary Password**:
A password set by someone else (or by the Instance) that the User must replace at their next sign-in before doing anything else.
_Avoid_: Initial password, default password, reset password

**Membership**:
A link that gives a User access at one level and carries the User's role there; always qualified as Organization Membership or Workspace Membership when the level matters.
_Avoid_: Member, participant, collaborator

**Organization Membership**:
The Membership that links a User to one Organization and grants access to it.

**Workspace Membership**:
The Membership that links a User to one Workspace, granting entry to it when restricted and a role inside it.

**Administrator**:
A User who manages the Instance itself (settings, Users, Organizations); the role grants no access to any Organization's data.
_Avoid_: Instance Administrator, superuser, root, sysadmin

**Organization Owner**:
The role, held through an Organization Membership, that manages one Organization; never called "administrator".
_Avoid_: Organization admin, org administrator

**Organization Contributor**:
The role, held through an Organization Membership, of everyone in an Organization who is not an Organization Owner; it does not manage the Organization.
_Avoid_: Member, organization user, regular user

**Workspace Owner**:
The role, held through a Workspace Membership, that manages one Workspace: its Workspace Memberships, Projects, Teams and visibility.
_Avoid_: Workspace admin, workspace manager

**Workspace Contributor**:
The role that works inside a Workspace's Projects (creating and editing Work Items) without managing the Workspace.
_Avoid_: Workspace member, editor

**Workspace Viewer**:
The read-only role in a Workspace: sees everything in it, changes nothing.
_Avoid_: Guest, reader, observer

**Team Assignment**:
The link that places a User in a Team; it grants no access.
_Avoid_: Team Membership, team member

### Work

**Work Item**:
A unit of work planned, executed and tracked inside one Project; every Work Item has exactly one Work Item Type.
_Avoid_: Ticket, issue, card, task (as a generic word)

**Work Item Type**:
The fixed kind of a Work Item: Epic, Feature, User Story, Task, Bug or Subtask.
_Avoid_: Item kind, category

**Epic**:
The largest Work Item Type: a broad initiative that groups Features.

**Feature**:
A product capability that groups User Stories.

**User Story**:
A piece of value deliverable to a user, small enough for one cycle, that groups Tasks and Bugs.
_Avoid_: Story, requirement, PBI

**Task**:
Work needed to deliver a User Story.

**Bug**:
A defect to be fixed.
_Avoid_: Defect, issue

**Subtask**:
A small step inside a Task or a Bug.
_Avoid_: Task Item, Bug Item, checklist item

**Parent**:
The Work Item directly above another in the hierarchy; its counterpart is the **Child**.
_Avoid_: Container, owner

**State**:
Where a Work Item stands in its life: Backlog, To Do, In Progress or Done.
_Avoid_: Status, stage, column

**Backlog**:
The State of work recorded but not committed to; it may never be done, and discarded work returns here.

**To Do**:
The State of work committed to soon (for example, chosen for the current cycle) but not started.

**In Progress**:
The State of work being done.

**Done**:
The State of work delivered; discarded work never goes here.
_Avoid_: Closed, resolved, cancelled

**Link**:
A non-hierarchical relation between two Work Items of the same Project, of a fixed kind: Blocks, Duplicates or Relates To.
_Avoid_: Relation, association, connection

**Blocks**:
The Link kind meaning one Work Item must be finished before another can proceed; read from the other side as "blocked by".
_Avoid_: Dependency, depends on

**Duplicates**:
The Link kind meaning a Work Item reports the same thing as another.

**Relates To**:
The symmetric Link kind with no rules attached.

### Lifecycle

**Archived**:
The state of a Workspace or Project that is hidden from default lists and read-only, and can be restored intact.
_Avoid_: Closed, inactive, soft-deleted

**Pending Deletion**:
The state of an Organization, Workspace or Project (never a Work Item) that someone deleted: invisible and inaccessible to everyone, but restorable until its grace period ends.
_Avoid_: Trash, recycle bin, soft-deleted

**Deleted**:
The final state reached when a Pending Deletion grace period ends; nothing can be restored after it. Deleting does not require being Archived first.
_Avoid_: Removed, purged

## Relationships

### Structure

- An **Instance** hosts zero or more **Organizations**
- A **Workspace** belongs to exactly one **Organization**; an **Organization** always has at least one **Workspace** that is not **Archived**, so its last one can be neither archived nor deleted
- Every **Organization** is created with its **Default Workspace**, and its creator receives an **Organization Membership** as **Organization Owner**
- A **Team** belongs to exactly one **Workspace**
- A **Project** belongs to exactly one **Workspace**
- A **Project** has zero or more **Project Teams**, and only **Teams** of the **Project**'s own **Workspace** can be associated with it
- Removing a **Team** removes its **Project Teams** but never the **Projects**
- A **Project** can move to another **Workspace** of the same **Organization**, taking its content and losing its **Project Teams**; it can never move to another **Organization**

### Access

- An **Organization** always has at least one **Organization Owner**: the last one can neither leave nor be removed, and a **User** who is the last **Organization Owner** somewhere cannot be **Deactivated**
- If an **Organization** has no reachable **Organization Owner**, an **Administrator** can appoint one, still without gaining access to its data
- Each **Instance** decides who can create **Organizations**: only **Administrators**, or any **User**
- An **Administrator** sees an **Organization**'s data only through an **Organization Membership**, like any other **User**
- A **User** can hold **Organization Memberships** in many **Organizations**; each links exactly one **User** to exactly one **Organization**
- Only an **Organization Owner** grants **Organization Memberships**, by adding a person's name and **Username**: an existing **User** is linked as is, and an unknown **Username** creates a new **User**
- A **Workspace Membership** exists only while the same **User** holds an **Organization Membership** in that **Workspace**'s **Organization**; removing the latter removes the former
- **Organization Owners** can enter every **Workspace** of their **Organization**, open or restricted
- Every **Organization Membership** holds exactly one role: **Organization Owner** or **Organization Contributor**
- Every **Workspace Membership** holds exactly one role: **Workspace Owner**, **Workspace Contributor** or **Workspace Viewer**
- **Organization Owners** act as **Workspace Owner** in every **Workspace** of their **Organization** without a **Workspace Membership**
- A **User** who enters an **Open Workspace** without a **Workspace Membership** acts as **Workspace Contributor**
- Only **Organization Owners** create **Workspaces**; the creator receives a **Workspace Membership** as **Workspace Owner**
- A new **Instance** with no **Users** creates one **Administrator**, either from the installer's configuration or as `admin` with the **Temporary Password** `admin`
- Adding an unknown **Username** gives the new **User** a **Temporary Password**, shown once to the **Organization Owner** who added them
- Only an **Administrator** can give an existing **User** a new **Temporary Password**, because a **User**'s password opens every **Organization** they belong to
- A **Deactivated** **User** keeps their **Memberships** and authorship, but the **Memberships** grant nothing and no new ones can be added until reactivation
- A **User** signs in with **Username** and password only; external identity providers are not part of the language yet
- A **Team Assignment** requires that the **User** can enter the **Team**'s **Workspace**; losing that access removes the **Team Assignment**

### Work

- A **Work Item** belongs to exactly one **Project**
- The hierarchy is fixed: an **Epic** is the **Parent** of **Features**, a **Feature** of **User Stories**, a **User Story** of **Tasks** and **Bugs**, and a **Task** or **Bug** of **Subtasks**
- A **Work Item** has at most one **Parent**; only a **Subtask** must have one (a **Task** or a **Bug**), every other **Work Item Type** can stand alone
- A **Work Item** can gain, lose or change its **Parent** later, as long as the types respect the hierarchy
- A **Parent** and its **Child** always belong to the same **Project**
- A **Link** always connects two **Work Items** of the same **Project**
- A **Work Item** is always in exactly one **State**, and can move between any two **States**
- Deleting a **Work Item** is immediate and final, with no **Pending Deletion**; its **Subtasks** go with it, its other **Children** lose their **Parent**, and its **Links** disappear

### Lifecycle

- Archiving a **Workspace** makes all of its **Projects** read-only; restoring it brings them back
- Deleting puts an **Organization**, **Workspace** or **Project** into **Pending Deletion**; each **Instance** sets the grace period before it becomes **Deleted**
- An **Organization Owner** can restore what is in **Pending Deletion**; an **Administrator** can restore a whole **Organization** without gaining access to its data

## Example dialogue

> **Dev:** "Ana consults for ACME and Globex. Does she need two **Users**?"
> **Domain expert:** "No: one **User** with an **Organization Membership** in each. What she sees in ACME depends only on her memberships in ACME."

> **Dev:** "Bruno left ACME but still has a **Workspace Membership** in Legal. Can he open it?"
> **Domain expert:** "That can't happen: removing his **Organization Membership** removed every **Workspace Membership** he had in ACME."

## Flagged ambiguities

- "status" and "state" were used interchangeably. Resolved: **State** is the term.
- Abandoned work had no State of its own. Resolved: it returns to **Backlog** or is deleted; **Done** always means delivered, so metrics stay honest.

- "dependency" and "blocker" were listed as two relations. Resolved: they are one, **Blocks**, read from opposite sides.

- "Task Item" and "Bug Item" were proposed as two types for the same small step, and both collide with "Work Item". Resolved: a single **Subtask** type, Child of a **Task** or a **Bug**.

- "removing a user" could mean leaving one Organization or leaving the Instance. Resolved: losing an **Organization Membership** is leaving an **Organization**; being **Deactivated** is leaving the **Instance**. Users are never deleted, so authorship survives.

- "isolated by Organization/Workspace" did not say which one is the tenant. Resolved: the **Organization** is the tenant; a **Workspace** only restricts visibility within it.
- "members" was used loosely for people in general. Resolved: a **User** is the global identity; access comes from **Memberships**. "Member" is not a domain term.
- Being "in a Team" sounded like a third access level. Resolved: a **Team** never controls access; joining one is a **Team Assignment**, not a Membership.
- "a Project is associated with one or more Teams" did not say who owns it. Resolved: the **Workspace** owns the **Project**; **Teams** are only associated with it, and there can be none.
- "administrator" could mean whoever runs the installation or whoever manages an Organization. Resolved: **Administrator** is only the Instance-level role; Organization- and Workspace-level management roles never use the word (exact role names may be refined with RBAC).
- "archive" and "delete" were both used for "make it go away". Resolved: **Archived** is for keeping and consulting (read-only, restorable); deleting starts a **Pending Deletion** grace period that ends in **Deleted**.
- "join an Organization" could mean self-service sign-up or being added. Resolved: a **User** never joins on their own; an **Organization Owner** adds them. Email invitations are a possible future addition, not part of the language yet.
- "Membership" alone is ambiguous between levels. Resolved: say **Organization Membership** or **Workspace Membership** whenever the level matters.
