# Issue tracker: local markdown

Issues and PRDs for this repository live as markdown files in `.scratch/`. The folder is ignored by Git, so its content exists only on the local machine.

## Conventions

- One feature per directory: `.scratch/<feature-slug>/`
- The PRD is `.scratch/<feature-slug>/PRD.md`
- Implementation issues are `.scratch/<feature-slug>/issues/<NN>-<slug>.md`, numbered from `01`
- Triage state is a `Status:` line near the top of each issue (the values are in `triage-labels.md`)
- Comments and conversation history are appended to the end of the file under a `## Comments` heading

## When a skill says "publish to the issue tracker"

Create a file under `.scratch/<feature-slug>/`, creating the directory if needed.

## When a skill says "fetch the relevant ticket"

Read the file at the given path. The user normally provides the path or the issue number.
