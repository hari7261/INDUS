# INDUS v1.5.5 - Workflow Automation Release

INDUS `v1.5.5` adds three major quality-of-life features for day-to-day usage and deployment.

## Highlights

- Adds `ind update` for release checks, downloads, and staged Windows self-update installs
- Adds `ind task` for saved multi-step INDUS workflows
- Adds `ind term profile` for persistent banner, prompt, compact mode, and animation settings
- Keeps the GUI-first Windows experience and UTF-8 console rendering improvements

## New Commands

- `ind update`
- `ind task create`
- `ind task list`
- `ind task show`
- `ind task run`
- `ind task remove`
- `ind term profile`

## Validation

- `go test ./...`
- Windows GUI build completed successfully
- Release smoke suite rerun against the packaged binary

## Notes

- Existing command compatibility from `v1.5.4` remains included
- `ind update install` stages the new binary and restarts INDUS on Windows
