# Security Architecture - Visual Diagrams

This document provides visual diagrams to help understand the security architecture.

---

## Request Flow Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          CLIENT REQUEST                                 │
│                                                                         │
│  Browser (Cookie)              API Client (Header)                      │
│  Cookie: auth_token=...        Authorization: Bearer ...                │
└────────────────────────────────┬────────────────────────────────────────┘
                                 │
                                 ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                       AUTHENTICATION LAYER                              │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │  RequireAuth() / RequireAuthFromCookie() / RequireAuthFromHeader│  │
│  │                                                                   │  │
│  │  1. Extract token from request (cookie or header)                │  │
│  │  2. Validate JWT signature                                       │  │
│  │  3. Check expiration, issuer, audience                           │  │
│  │  4. Parse claims → User object                                   │  │
│  │  5. Store User in context                                        │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  Success: User in Context  │  Failure: 401 Unauthorized                │
└────────────────────────────┼────────────────────────────────────────────┘
                             │
                             ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                       AUTHORIZATION LAYER                               │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │  RequirePermission() / RequirePermissionsMatchAll/Any()          │  │
│  │                                                                   │  │
│  │  1. Get User from context                                        │  │
│  │  2. Get User.Permissions array                                   │  │
│  │  3. Check against required permissions                           │  │
│  │     - MatchAll: User has ALL required permissions?               │  │
│  │     - MatchAny: User has ANY required permission?                │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  Success: Continue          │  Failure: 403 Forbidden                  │
└────────────────────────────┼────────────────────────────────────────────┘
                             │
                             ↓
┌─────────────────────────────────────────────────────────────────────────┐
│                       HANDLER / RESOLVER                                │
│                                                                         │
│  Business logic executes - User is authenticated and authorized         │
│                                                                         │
│  Can access: auth.GetCurrentUser(ctx)                                  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## JWT Token Flow

```
┌──────────────────────┐
│  External Auth       │
│  Service             │
│  (e.g., Okta, Auth0, │
│   Custom Service)    │
└──────────┬───────────┘
           │
           │ 1. User logs in
           │
           ↓
    ┌──────────────┐
    │   Issues     │
    │   JWT Token  │
    └──────┬───────┘
           │
           │ 2. JWT returned to client
           │
           ↓
┌──────────────────────┐
│  Client              │
│  Stores Token:       │
│  - Cookie (web)      │
│  - Storage (mobile)  │
└──────┬───────────────┘
       │
       │ 3. Include in requests
       │
       ↓
┌──────────────────────────────────┐
│  RxIntake App                    │
│                                  │
│  auth.ValidateToken()            │
│  - Verify signature              │
│  - Check expiration              │
│  - Parse claims                  │
│  - Extract user + permissions    │
│                                  │
│  → User in Context               │
└──────────────────────────────────┘
```

---

## JWT Token Structure

```
┌─────────────────────────────────────────────────────────────┐
│                      JWT TOKEN                              │
├─────────────────────────────────────────────────────────────┤
│  Header (Algorithm & Token Type)                            │
│  {                                                           │
│    "alg": "HS256",                                          │
│    "typ": "JWT"                                             │
│  }                                                           │
├─────────────────────────────────────────────────────────────┤
│  Payload (Claims)                                            │
│  {                                                           │
│    "user_id": "12345",              ← User identification   │
│    "email": "doctor@example.com",   ← User info             │
│    "name": "Dr. Smith",             ← Display name          │
│    "permissions": [                 ← Authorization         │
│      "patient:read",                                        │
│      "patient:write",                                       │
│      "prescription:read",                                   │
│      "doctor:role"                                          │
│    ],                                                        │
│    "iss": "rxintake",               ← Issuer                │
│    "aud": "rxintake",               ← Audience              │
│    "exp": 1728691200,               ← Expiration            │
│    "iat": 1728604800                ← Issued at             │
│  }                                                           │
├─────────────────────────────────────────────────────────────┤
│  Signature                                                   │
│  HMACSHA256(                                                │
│    base64UrlEncode(header) + "." +                          │
│    base64UrlEncode(payload),                                │
│    secret                                                    │
│  )                                                           │
└─────────────────────────────────────────────────────────────┘
```

---

## Permission Checking Logic

### MatchAny (OR Logic)

```
Required: ["patient:read", "admin:all"]

User Permissions: ["patient:read", "prescription:read"]
                        ↓
                  patient:read ✓
                        ↓
                   MATCH → Allow


User Permissions: ["prescription:read", "prescription:write"]
                        ↓
                  No match with required
                        ↓
                   NO MATCH → 403 Forbidden
```

### MatchAll (AND Logic)

```
Required: ["patient:read", "patient:export"]

User Permissions: ["patient:read", "patient:export", "admin:all"]
                        ↓
                  patient:read ✓
                  patient:export ✓
                        ↓
                   ALL MATCH → Allow


User Permissions: ["patient:read", "prescription:read"]
                        ↓
                  patient:read ✓
                  patient:export ✗
                        ↓
                   NOT ALL MATCH → 403 Forbidden
```

---

## Component Interaction

```
┌─────────────────────────────────────────────────────────────────────┐
│                          APPLICATION                                │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌─────────────────────┐  │
│  │  HTTP Routes │    │   GraphQL    │    │  Business Logic     │  │
│  │              │    │   Schema     │    │  (Services)         │  │
│  │  @middleware │───▶│  @directives │───▶│                     │  │
│  │              │    │              │    │  auth.GetUser(ctx)  │  │
│  └──────────────┘    └──────────────┘    └─────────────────────┘  │
│         │                    │                                     │
│         │                    │                                     │
│         ↓                    ↓                                     │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │          internal/platform/auth/                             │ │
│  │                                                                │ │
│  │  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌────────────┐│ │
│  │  │   JWT    │  │  Context │  │ Permissions│  │  Middleware││ │
│  │  │ Handler  │  │  Helpers │  │  Checker   │  │            ││ │
│  │  └──────────┘  └──────────┘  └────────────┘  └────────────┘│ │
│  │       ↑                                                       │ │
│  │       │                                                       │ │
│  │  ┌────────────────────┐                                      │ │
│  │  │     Config         │                                      │ │
│  │  │  (JWT secret, etc) │                                      │ │
│  │  └────────────────────┘                                      │ │
│  └──────────────────────────────────────────────────────────────┘ │
│         ↑                                                          │
│         │                                                          │
│  ┌──────────────────────────────────────────────────────────────┐ │
│  │          domain/*/security/                                  │ │
│  │                                                                │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐                   │ │
│  │  │ Patient  │  │Prescription│ │Dashboard │                   │ │
│  │  │Permissions│ │Permissions │ │Permissions│                  │ │
│  │  └──────────┘  └──────────┘  └──────────┘                   │ │
│  │                                                                │ │
│  │  Constants + Permission Sets                                  │ │
│  └──────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Route Protection Patterns

### Pattern 1: Public + Protected Routes

```
Router
  │
  ├─ / (home)                     [No Auth]
  ├─ /login                       [No Auth]
  │
  └─ /patients                    [Auth Required]
      │
      ├─ GET /                    [patient:read OR admin:all]
      ├─ POST /                   [patient:write OR admin:all]
      ├─ GET /:id                 [patient:read OR admin:all]
      ├─ PUT /:id                 [patient:write OR admin:all]
      ├─ DELETE /:id              [patient:write AND patient:delete]
      └─ GET /export              [patient:read AND patient:export]
```

### Pattern 2: Different Auth Methods

```
Router
  │
  ├─ /patients (UI)                [Cookie Auth]
  │   └─ Uses RequireAuthFromCookie()
  │
  └─ /api/v1/patients (API)        [Header Auth]
      └─ Uses RequireAuthFromHeader()
```

---

## GraphQL Directive Flow

```
Client sends GraphQL query
         ↓
┌────────────────────────┐
│  GraphQL Server        │
│  receives request      │
└────────┬───────────────┘
         │
         ↓
┌────────────────────────┐
│  RequireAuth()         │
│  middleware            │
│  - Sets user in context│
└────────┬───────────────┘
         │
         ↓
┌────────────────────────────────────────┐
│  GraphQL Resolver Execution            │
│                                        │
│  For each field:                       │
│    ↓                                   │
│  ┌──────────────────────────────────┐ │
│  │  @auth directive?                │ │
│  │  → Check user in context         │ │
│  │  → 401 if not authenticated      │ │
│  └──────────────────────────────────┘ │
│    ↓                                   │
│  ┌──────────────────────────────────┐ │
│  │  @permissionAny or @permissionAll│ │
│  │  → Check user permissions        │ │
│  │  → 403 if insufficient           │ │
│  └──────────────────────────────────┘ │
│    ↓                                   │
│  ┌──────────────────────────────────┐ │
│  │  Field Resolver                  │ │
│  │  → Execute business logic        │ │
│  └──────────────────────────────────┘ │
└────────────────────────────────────────┘
         ↓
Return response to client
```

---

## Error Response Flow

```
Request
   ↓
Auth Middleware
   │
   ├─ No token ────────────→ 401 Unauthorized
   │                         {
   │                           "error": "unauthorized",
   │                           "message": "Authentication required"
   │                         }
   │
   ├─ Invalid token ───────→ 401 Unauthorized
   │                         {
   │                           "error": "unauthorized",
   │                           "message": "Invalid or expired token"
   │                         }
   │
   └─ Valid token
       ↓
   Permission Middleware
       │
       ├─ Missing permission ──→ 403 Forbidden
       │                         {
       │                           "error": "forbidden",
       │                           "message": "Insufficient permissions",
       │                           "required_permissions": [...]
       │                         }
       │
       └─ Has permission
           ↓
       Handler
           ↓
       200 OK
```

---

## Permission Hierarchy Example

```
System Permissions
├─ admin:all ───────────────────────┐
│                                   │ (Super admin - has all permissions)
│                                   │
Domain Permissions                  │
│                                   │
├─ Patient Domain                   │
│  ├─ patient:read ◄─────────────┬──┘
│  ├─ patient:write ◄────────────┤
│  ├─ patient:delete ◄───────────┤
│  └─ patient:export ◄───────────┤
│                                 │
├─ Prescription Domain             │
│  ├─ prescription:read ◄─────────┤
│  ├─ prescription:write ◄────────┤
│  ├─ prescription:approve ◄──────┤
│  └─ prescription:dispense ◄─────┤
│                                 │
└─ Dashboard Domain                │
   ├─ dashboard:view ◄─────────────┤
   ├─ dashboard:analytics ◄────────┤
   └─ dashboard:reports ◄──────────┘

Role Permissions (Optional)
├─ doctor:role
│  → Can have: patient:*, prescription:write
│
├─ pharmacist:role
│  → Can have: prescription:read, prescription:dispense
│
└─ nurse:role
   → Can have: patient:read, prescription:read
```

---

## Deployment Considerations

```
┌─────────────────────────────────────────────────────────────┐
│                    Production Setup                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐        ┌──────────────────────────┐  │
│  │ External Auth   │        │   RxIntake App           │  │
│  │ Service         │        │                          │  │
│  │ (Okta, Auth0)   │        │  Environment Variables:  │  │
│  │                 │        │  - RX_AUTH_JWT_SECRET    │  │
│  │ Issues JWT ────────────▶ │  - RX_AUTH_JWT_ISSUER    │  │
│  │ with permissions│        │                          │  │
│  └─────────────────┘        │  Config:                 │  │
│                             │  - cookie.secure: true   │  │
│                             │  - HTTPS only            │  │
│                             │  - Short expiration      │  │
│                             └──────────────────────────┘  │
│                                                             │
│  Load Balancer (HTTPS)                                     │
│  ├─ Terminates SSL                                         │
│  ├─ Multiple App Instances (stateless)                     │
│  └─ Health checks                                          │
│                                                             │
│  Monitoring                                                 │
│  ├─ Auth failures (401/403)                                │
│  ├─ Token validation errors                                │
│  └─ Permission denied events                               │
└─────────────────────────────────────────────────────────────┘
```

---

These diagrams provide a visual understanding of how the security system works. For detailed implementation, see [SECURITY_ARCHITECTURE.md](./SECURITY_ARCHITECTURE.md).

