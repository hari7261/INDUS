# INDUS v1.5.2 - GUI Release Hardening

**Release Date:** April 22, 2026
**Tag:** v1.5.2

---

## Highlights

- GUI-first Windows release flow centered on `indus.exe`
- Fixed `ind dev bench --command "ind sys stats"` parsing
- Restored legacy HTTP command compatibility:
  `ind http get ...`
  `ind net http get ...`
  `ind net http post ... --data ...`
- Added the new animated tricolor startup banner
- Aligned build, installer, smoke-test, and packaging scripts with the current app layout

## Release Artifacts

- `dist/indus.exe`
- `dist/indus-setup-v1.5.2-windows-amd64.exe`
  This requires Inno Setup locally to generate.

## Verification

- `go test ./...`
- `build.bat`
- `scripts/release_smoke.ps1`

The smoke suite exercises the documented native command surface, including the restored HTTP aliases.

## Notes

- The current documented native INDUS command catalog is 60 commands.
- Historical `v1.5.0` release notes are preserved under `docs/version/v1.5.0.html`.
- `v1.5.2` is the current release page for the docs site.

## Changelog

- Release pipeline now builds from `cmd/indus-terminal` only.
- Installer now targets `indus.exe` as the primary application.
- Versioned installer name is `indus-setup-v1.5.2-windows-amd64.exe`.
- Release smoke testing defaults to `dist/indus.exe`.
- Docs updated to point at `v1.5.2` as the current release.
