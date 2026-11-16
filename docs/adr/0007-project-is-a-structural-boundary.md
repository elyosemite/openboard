# 0007. Nothing structural crosses a Project

- Status: Accepted
- Date: 2026-11-15

Every structural relation between Work Items stays inside a single Project: a Parent and its Child, both ends of a Link (Blocks, Duplicates, Relates To), and every Work Item of an Era all belong to the same Project. Unlike Jira or Azure DevOps, there are no cross-project Epics, dependencies or sprints. We chose this because Projects can move between Workspaces, be Archived and enter Pending Deletion. Relations that span Projects would either break on those operations or leak Work Items from Restricted Workspaces to people who cannot enter them. Initiatives that span Projects are served by Views and Analytics, which can combine Work Items from several Projects without creating structural links. Allowing cross-project relations later would change the rules for moving, archiving and deleting Projects, so it needs a new ADR.
