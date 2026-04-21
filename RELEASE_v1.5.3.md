# INDUS v1.5.3 - Banner Refresh Release

**Release Date:** April 22, 2026
**Tag:** v1.5.3

---

## Highlights

- New tricolor startup banner styled closer to the latest INDUS visual direction
- GUI-first Windows release flow remains centered on `indus.exe`
- Existing command fixes from `v1.5.2` remain included:
  `ind dev bench --command "ind sys stats"`
  `ind http get ...`
  `ind net http get ...`
  `ind net http post ... --data ...`

## Release Artifacts

- `dist/indus.exe`
- `indus-v1.5.3-windows-amd64.exe`
- `indus-v1.5.3-windows-amd64.zip`
- `smoke-summary-v1.5.3.txt`
- `smoke-report-v1.5.3.json`

## Verification

- `go test ./...`
- GUI build to `dist/indus.exe`
- `scripts/release_smoke.ps1`

The smoke suite verifies the documented native command surface and the restored HTTP alias compatibility.

## Notes

- The current documented native INDUS command catalog is 60 commands.
- `v1.5.3` is now the current release page for the docs site.
- Installer generation still requires local Inno Setup to produce `indus-setup-v1.5.3-windows-amd64.exe`.
