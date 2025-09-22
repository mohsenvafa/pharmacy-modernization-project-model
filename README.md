# rxintake (Patient & Prescription Management - scaffold)

## Dev
- Install: Go 1.24, `templ`, Node.js 18+, `npm`
- Setup once: `make setup` (delegates to `Makefile.setup` to download the Tailwind standalone binary into `./bin` and install DaisyUI via npm)
- Run everything together: `make dev`
- Or run watchers separately: `make dev-watch` and `make tailwind-watch`
- Open: http://localhost:8080

## Notes
- Feature-based modules under `internal/domain/*` and UI under `internal/ui/*`.
- Viper YAML config in `internal/configs/` with env overrides (RX_*).
- Zap logging with request/correlation IDs.
- Tailwind source lives in `web/styles/input.css`; `make tailwind-watch` rebuilds `web/public/app.css` via the standalone Tailwind CLI with DaisyUI.
