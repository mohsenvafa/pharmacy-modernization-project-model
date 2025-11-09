## Patient Info Microproduct

This document captures the standard layout and contract for the `patient-info` web component microproduct.

### Overview

- **Purpose**: Render patient demographic details for external consumer applications by delegating data fetching and rendering to the server-side micro UI fragment.
- **Inputs**:
  - `patient-id` – unique identifier string.
  - `env` – one of `local`, `dev`, `stg`, `prod`.
  - `auth-token` – opaque bearer token used by the backend fragment for authorization.
  - `base-path` *(optional)* – host/base URL only (no micro UI path). If omitted, the component resolves the base URL from `env` via the shared helper.

### Directory Layout

```
web_components/
└── patient/
    └── patient_info/
        ├── package.json
        ├── tsconfig.json
        ├── esbuild.config.mjs
        ├── src/
        │   ├── index.ts
        │   └── patient-info.element.ts
        └── dist/                      # build output (gitignored)

domain/
└── patient/
    └── micro_ui/
        └── patient_info/
            ├── patient_info.component.templ
            ├── patient_info.component.go
            └── patient_info.component_templ.go (generated)
```

### Backend Contract

- **Endpoint**: `/micro-ui/patients/:patientId`
- **Query Parameters**:
  - `auth_token`
- **Response**: HTML fragment rendered from templ with patient demographics and metadata.
- **Responsibilities**:
  - Validate `patientId`.
  - Authenticate using the provided `auth_token`.
  - Render a friendly error fragment if data retrieval fails.

### Web Component Behavior

1. Reads attributes `patient-id`, `env`, `auth-token`, optional `base-path`.
2. Resolves the host URL using shared helper (`web_components/shared/base-url.ts`).
3. Calls the backend endpoint with composed path and query parameters.
4. Injects the returned HTML into the shadow DOM.
5. Emits lifecycle events: `patient-info:loading`, `patient-info:loaded`, `patient-info:error`.

### Build & Distribution

- Bundled with esbuild via local `package.json` scripts (`build`, `watch`).
- Output published as ESM and CJS bundles (`dist/`).
- Run `npm run pack` in `web_components/patient/patient_info/` to generate a tarball under `web_components/temp_pack/`.
- Install locally with `npm install ../../web_components/temp_pack/rx-micro-patient-info-0.1.0.tgz` (or from a registry once published) and import with `import '@rx/micro-patient-info'`.

### Consumer Usage

```html
<patient-info
  patient-id="PAT-001"
  env="local"
  auth-token="demo-token"
></patient-info>
```

### Test Beds

- Angular and React sample apps now include patient ID controls and demonstrate usage of both `prescription-info` and `patient-info` components.

### Notes

- Shared helper functionality resides in `web_components/shared/base-url.ts` so additional components can reuse environment and URL normalization logic.
- Server-side micro UI fragments remain ignorant of environment concerns; the client component handles host resolution.

