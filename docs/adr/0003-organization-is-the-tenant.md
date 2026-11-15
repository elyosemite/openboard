# 0003. The Organization is the tenant

- Status: Accepted
- Date: 2026-11-15

The Organization, not the Workspace, is the tenant: the outermost boundary of data, administration and access. Every domain record carries its Organization, and no data ever crosses Organizations. A Workspace only restricts visibility inside its Organization, and Organization Owners can always enter every Workspace. We picked this over per-Workspace tenancy because companies need one administrative boundary over all their areas, and because the same model then serves a self-hosted Instance with one Organization and a Cloud Instance with thousands. Changing the tenant later would mean migrating every table and every authorization rule.
