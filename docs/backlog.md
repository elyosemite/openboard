# Backlog

Suggestions recorded to do later. Each entry says what is wrong, why, and the proposed fix.

## Flaky integration tests

### "Database migrations" CI job still downloads the JDBC driver

- **Where:** the `database` job in `.github/workflows/ci.yml` builds `database/Dockerfile` with a plain `docker build`, without cache.
- **Symptom:** `RUN lpm add postgresql@42.7.8 --global` downloads the driver on every run. The same download already failed once on GitHub runners, in the Backend job.
- **Proposed fix:** build the image in that job as the Backend job does, with `docker/build-push-action`, `load: true` and `cache-from: type=gha,scope=database`, so it downloads only when `database/Dockerfile` changes.
