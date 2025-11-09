## Prescription Info Microproduct

This document defines the standard layout, contract, and tooling for the `prescription-info` web component microproduct.

### Overview

- **Purpose**: Render prescription details within external host applications (React, Angular, etc.) using a thin web component layer backed by server-rendered HTML.
- **Inputs**:
  - `prescriptionId` – unique identifier string.
  - `env` – one of `local`, `dev`, `stg`, `prod`.
  - `auth_token` – opaque bearer token used by the backend fragment for authorization.
  - `base-path` *(optional)* – host/base URL only (no micro UI path). If omitted, the component resolves the base URL from `env` via the shared environment helper.

The component delegates business logic and rendering to a Go/templ fragment served by the main application. The TypeScript layer simply fetches and injects that HTML into the host page.

### Directory Layout

```
web_components/
└── prescription/
    └── prescription_info/
        ├── package.json
        ├── tsconfig.json
        ├── src/
        │   ├── index.ts
        │   └── prescription-info.element.ts
        └── dist/                       # build output (gitignored)

domain/
└── prescription/
    └── micro_ui/
        └── prescription_info/
            ├── prescription_info.component.templ
            ├── prescription_info.component.go
            └── routes.go
```

### Backend Contract

- **Endpoint**: `/micro-ui/prescriptions/:prescriptionId`
- **Query Parameters**:
  - `env`
  - `auth_token`
- **Response**: HTML fragment rendered from templ.
- **Responsibilities**:
  - Validate `prescriptionId`.
  - Switch downstream service base URL / config based on `env`.
  - Authenticate using the provided `auth_token`.
  - Return a rendered fragment or structured error HTML.

### Web Component Behavior

1. Reads attributes `prescription-id`, `env`, and `auth-token`.
2. Calls the backend endpoint with those inputs.
3. Injects the returned HTML into the component shadow root.
4. Emits lifecycle events: `prescription-info:loading`, `...:loaded`, `...:error`.

### Build & Distribution

- Bundled with esbuild via local `package.json` scripts (`build`, `watch`).
- Output published as ESM and UMD bundles (`dist/`).
- Consumers import from the published npm package (`@rx/micro-prescription-info`) and register the custom element.
- Shared utilities live under `web_components/shared/` so additional components can reuse them.
- To produce a distributable tarball: `npm run clean && npm run build && npm pack`.

### Consumer Usage

- Run `npm run pack` inside the microproduct directory. Tarballs are emitted to `web_components/temp_pack/`.
- Consumer apps can install locally with `npm install ../../web_components/temp_pack/rx-micro-prescription-info-0.1.0.tgz` (or via a published registry once available) and then `import '@rx/micro-prescription-info'`.
- Helper utilities (e.g., environment base path resolution) live under `web_components/shared/` so additional components can reuse them.

### Test Beds

- `test_bed/angular_app` on port `6200` and `test_bed/react_app` on port `7200` both import the distributable bundle and instantiate the custom element, passing the three required attributes.
- `test_bed/Makefile` offers convenience targets (`make angular`, `make react`, `make run`) to run the playground apps.

### Next Steps Checklist

- [x] Implement backend fragment in `domain/prescription/micro_ui/prescription_info`.
- [x] Build TypeScript custom element and registration entrypoints.
- [x] Wire build script to generate distributable bundle.
- [x] Configure Angular/React test beds with example usage.
- [ ] Document deployment process for packaging the component.

