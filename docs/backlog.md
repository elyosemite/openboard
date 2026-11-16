# Backlog

Suggestions recorded to do later. Each entry says what is wrong, why, and the proposed fix.

## Flaky integration tests

### `cmd/api` start-up tests depend on a 10-second deadline

- **Where:** `startAndStop` in `backend/cmd/api/main_test.go`, used by `TestRunStartsCreatesAdministratorOnceAndStops` and `TestRunCreatesAdministratorFromSeed`.
- **Symptom:** both fail with `server did not start:` and empty output. Seen locally on 2026-11-16, in 1 of 3 runs of `go test ./...`.
- **Cause:** the test gives `run` 10 seconds of wall-clock time to open the database, create the first Administrator (argon2id, expensive by design) and log `http server listening`. With every test package running in parallel, and Docker Desktop limited to about 4 GB, the machine takes longer. Nothing is broken; start-up is just slow under load.
- **Proposed fix:** keep the same checks, but stop depending on a short deadline:
  - fail at once, with the returned error, when `run` exits before logging `http server listening` (today the test waits for the deadline);
  - raise the deadline from 10 s to 60 s, kept only so a hung server cannot block the run forever.

### "Database migrations" CI job still downloads the JDBC driver

- **Where:** the `database` job in `.github/workflows/ci.yml` builds `database/Dockerfile` with a plain `docker build`, without cache.
- **Symptom:** `RUN lpm add postgresql@42.7.8 --global` downloads the driver on every run. The same download already failed once on GitHub runners, in the Backend job.
- **Proposed fix:** build the image in that job as the Backend job does, with `docker/build-push-action`, `load: true` and `cache-from: type=gha,scope=database`, so it downloads only when `database/Dockerfile` changes.
