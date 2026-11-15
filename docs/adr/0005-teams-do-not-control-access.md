# 0005. Teams do not control access

- Status: Accepted
- Date: 2026-11-15

Access has exactly two levels: Organization Membership and Workspace Membership. A Team organizes work (responsibility, planning, capacity) and never controls access. Joining one is a Team Assignment, deliberately not called a Membership. Anyone who can enter a Workspace sees its Projects; to hide work, use a Restricted Workspace. We rejected Team-based access (the Jira and Azure DevOps approach) because three overlapping permission layers make "why can't this person see it?" hard to answer, and they turn a routine team change into a security change. Project-level restriction to its Project Teams was considered as an opt-in middle ground and left out until a concrete need appears.
