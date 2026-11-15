# 0004. The Administrator has no access to Organization data

- Status: Accepted
- Date: 2026-11-15

The Administrator manages the Instance itself (settings, Users, creating and restoring Organizations, appointing an Organization Owner when none is reachable), but that role never grants access to any Organization's data. To see an Organization's Projects, an Administrator needs an Organization Membership like any other User. This is deliberate, not an omission: unlike the usual "superadmin who sees everything", it backs the Data Ownership principle and is what makes a Cloud Instance trustworthy, since running the installation does not mean reading the customers' work. Do not add an Administrator bypass to authorization checks.
