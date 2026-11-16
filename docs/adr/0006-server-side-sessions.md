# 0006. Server-side sessions with cookies, not JWT

- Status: Accepted
- Date: 2026-11-15

A signed-in User is identified by a server-side session stored in PostgreSQL and referenced by an `HttpOnly`, `Secure`, `SameSite=Lax` cookie. We chose this over JWTs because access must be revocable immediately: changing a password, receiving a new Temporary Password or signing out ends the User's sessions at once, which a self-contained token cannot do without a revocation list. The frontend and the API share one origin through the reverse proxy, so cookies need no CORS setup, and JavaScript never touches the credential, so XSS cannot steal it. The stateless validation JWTs offer buys nothing in a single modular monolith that already depends on PostgreSQL.
