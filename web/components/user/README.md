# User Info Component

## Overview

Reusable templ component for displaying authenticated user information throughout the application.

---

## Components

### 1. `UserInfo` - Full User Profile Card

Displays complete user information with avatar, name, email, and optional dev mode badge.

**Signature:**
```go
templ UserInfo(ctx context.Context, params UserInfoParams)
```

**Parameters:**
```go
type UserInfoParams struct {
    ShowEmail    bool  // Show email address
    ShowDevBadge bool  // Show dev mode badge if enabled
    Compact      bool  // Compact layout (reserved for future use)
}
```

**Usage:**
```go
// In your templ file
import usercomponents "pharmacy-modernization-project-model/web/components/user"

@usercomponents.UserInfo(ctx, usercomponents.UserInfoParams{
    ShowEmail:    true,
    ShowDevBadge: true,
    Compact:      false,
})
```

**Renders:**
```html
<div class="rounded-lg bg-base-200 p-3">
  <div class="flex items-center gap-3">
    <div class="avatar placeholder">
      <div class="bg-primary text-primary-content w-10 rounded-full">
        <span class="text-sm">D</span>
      </div>
    </div>
    <div class="flex-1 overflow-hidden">
      <p class="truncate text-sm font-semibold">Dr. Dev</p>
      <p class="truncate text-xs text-base-content/70">doctor@dev.local</p>
    </div>
  </div>
  <div class="mt-2 text-center">
    <span class="badge badge-warning badge-xs">Dev Mode</span>
  </div>
</div>
```

---

### 2. `UserInfoCompact` - Inline User Display

Displays user info in a compact, inline format suitable for headers or toolbars.

**Signature:**
```go
templ UserInfoCompact(ctx context.Context)
```

**Usage:**
```go
@usercomponents.UserInfoCompact(ctx)
```

**Renders:**
```html
<div class="flex items-center gap-2">
  <div class="avatar placeholder">
    <div class="bg-primary text-primary-content w-8 rounded-full">
      <span class="text-xs">D</span>
    </div>
  </div>
  <span class="text-sm font-medium">Dr. Dev</span>
  <span class="badge badge-warning badge-xs">Dev</span>
</div>
```

---

### 3. `UserAvatar` - Avatar Only

Displays just the user's avatar with customizable size.

**Signature:**
```go
templ UserAvatar(ctx context.Context, size string)
```

**Usage:**
```go
// Small avatar
@usercomponents.UserAvatar(ctx, "w-8")

// Medium avatar
@usercomponents.UserAvatar(ctx, "w-10")

// Large avatar
@usercomponents.UserAvatar(ctx, "w-16")
```

**Renders:**
```html
<div class="avatar placeholder">
  <div class="bg-primary text-primary-content rounded-full w-10">
    <span class="text-sm">D</span>
  </div>
</div>
```

---

### 4. `UserName` - Name Only

Displays just the user's name.

**Signature:**
```go
templ UserName(ctx context.Context)
```

**Usage:**
```go
<p>Welcome, @usercomponents.UserName(ctx)!</p>
```

**Renders:**
```html
<span>Dr. Dev</span>
```

---

### 5. `UserEmail` - Email Only

Displays just the user's email address.

**Signature:**
```go
templ UserEmail(ctx context.Context)
```

**Usage:**
```go
<p>Email: @usercomponents.UserEmail(ctx)</p>
```

**Renders:**
```html
<span>doctor@dev.local</span>
```

---

## Usage Examples

### Example 1: Sidebar (Current Implementation)

```go
// web/components/layouts/sidebar.templ
import usercomponents "pharmacy-modernization-project-model/web/components/user"

templ Sidebar(ctx context.Context) {
    <aside class="...">
        <div class="mb-4">
            <h2>Daisy Pharmacy</h2>
            <div class="mt-4">
                @usercomponents.UserInfo(ctx, usercomponents.UserInfoParams{
                    ShowEmail:    true,
                    ShowDevBadge: true,
                    Compact:      false,
                })
            </div>
        </div>
        <!-- Rest of sidebar -->
    </aside>
}
```

### Example 2: Header/Navbar

```go
templ Header(ctx context.Context) {
    <header class="navbar bg-base-100">
        <div class="navbar-start">
            <h1>Pharmacy Modernization</h1>
        </div>
        <div class="navbar-end">
            @usercomponents.UserInfoCompact(ctx)
        </div>
    </header>
}
```

### Example 3: Profile Page

```go
templ ProfilePage(ctx context.Context) {
    <div class="container">
        <h1>User Profile</h1>
        <div class="card">
            @usercomponents.UserInfo(ctx, usercomponents.UserInfoParams{
                ShowEmail:    true,
                ShowDevBadge: false,  // Don't show dev badge on profile page
                Compact:      false,
            })
            <!-- More profile details -->
        </div>
    </div>
}
```

### Example 4: Welcome Message

```go
templ WelcomeBanner(ctx context.Context) {
    <div class="hero">
        <div class="hero-content">
            <div>
                <h1>Welcome back, @usercomponents.UserName(ctx)!</h1>
                <p>Your email: @usercomponents.UserEmail(ctx)</p>
            </div>
        </div>
    </div>
}
```

### Example 5: User Dropdown Menu

```go
templ UserDropdown(ctx context.Context) {
    <div class="dropdown dropdown-end">
        <label tabindex="0" class="btn btn-ghost">
            @usercomponents.UserAvatar(ctx, "w-8")
        </label>
        <ul tabindex="0" class="dropdown-content menu">
            <li><a>Profile</a></li>
            <li><a>Settings</a></li>
            <li><a>Logout</a></li>
        </ul>
    </div>
}
```

---

## Customization

### Custom Avatar Colors

You can customize the component by modifying `user_info.component.templ`:

```go
// Example: Different colors for different roles
if hasRole(user, "admin") {
    <div class="bg-error text-error-content w-10 rounded-full">
} else if hasRole(user, "doctor") {
    <div class="bg-primary text-primary-content w-10 rounded-full">
} else if hasRole(user, "pharmacist") {
    <div class="bg-success text-success-content w-10 rounded-full">
} else {
    <div class="bg-neutral text-neutral-content w-10 rounded-full">
}
```

### Show Full Initials

Modify `getInitials()` function to show first and last name initials:

```go
func getInitials(name string) string {
    if name == "" {
        return "?"
    }
    
    parts := strings.Fields(name)
    if len(parts) == 0 {
        return "?"
    }
    
    if len(parts) == 1 {
        // Single word - first character
        return string([]rune(parts[0])[0])
    }
    
    // Multiple words - first char of first and last word
    first := []rune(parts[0])[0]
    last := []rune(parts[len(parts)-1])[0]
    return string(first) + string(last)
}
```

---

## Testing

### View Component in Different Contexts

```bash
# Start app
go run cmd/server/main.go

# Navigate to pages - sidebar shows user info
http://localhost:8080/
http://localhost:8080/patients
http://localhost:8080/prescriptions
```

### Test with Different Mock Users

```bash
# In browser with ModHeader extension:
# Set X-Mock-User: doctor → See "Dr. Dev"
# Set X-Mock-User: nurse → See "Dev Nurse"
# Set X-Mock-User: admin → See "Dev Admin"
```

---

## Benefits

- ✅ **Reusable** - Use anywhere in your application
- ✅ **Flexible** - 5 variants for different use cases
- ✅ **Automatic** - Gets user from context
- ✅ **Dev-aware** - Shows dev mode indicator
- ✅ **Graceful** - Handles missing user
- ✅ **Type-safe** - Structured parameters
- ✅ **DaisyUI styled** - Matches your theme

---

## File Location

```
web/components/user/
└── user_info.component.templ    ← Component definitions
    └── user_info.component_templ.go    ← Generated Go code (auto-generated)
```

---

## Summary

You now have a complete, reusable user component system with:

1. **UserInfo** - Full profile card
2. **UserInfoCompact** - Inline display
3. **UserAvatar** - Avatar only
4. **UserName** - Name only
5. **UserEmail** - Email only

Use them anywhere in your templ templates by importing:
```go
import usercomponents "pharmacy-modernization-project-model/web/components/user"
```

🎉 Clean, reusable, and production-ready!

