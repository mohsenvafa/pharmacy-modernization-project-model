# rxintake (Patient & Prescription Management - scaffold)

## Dev
- Install: Go 1.24, `templ`
- Run: `make dev`
- Open: http://localhost:8080

## Notes
- Feature-based modules under `internal/domain/*` and UI under `internal/ui/*`.
- Viper YAML config in `configs/` with env overrides (RX_*).
- Zap logging with request/correlation IDs.
