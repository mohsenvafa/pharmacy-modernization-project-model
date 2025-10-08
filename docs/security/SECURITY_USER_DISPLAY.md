# User Display in Sidebar

## Overview

The sidebar now displays the current authenticated user's information, including their name, email, and a visual indicator when dev mode is active.

---

## 🎨 What's Displayed

### User Profile Card

The sidebar shows:
- **Avatar** - User's initial in a colored circle
- **Name** - Full name of the authenticated user
- **Email** - User's email address
- **Dev Mode Badge** - Yellow "Dev Mode" badge (only visible in development)

### Visual Design

```
┌─────────────────────────────┐
│  Daisy Pharmacy             │
├─────────────────────────────┤
│  ┌─────────────────────┐    │
│  │  ╔═╗  Dr. Dev       │    │
│  │  ║ D ║  doctor@dev. │    │
│  │  ╚═╝     local      │    │
│  │      [Dev Mode]     │    │
│  └─────────────────────┘    │
├─────────────────────────────┤
│  Dashboard                  │
│  Patients                   │
│  ...                        │
└─────────────────────────────┘
```

---

## 🔧 Implementation

### Files Modified

**Sidebar Template**: `web/components/layouts/sidebar.templ`
```go
templ Sidebar(ctx context.Context) {
    <aside class="flex h-screen w-64 flex-col border-r bg-base-100 p-4">
        <!-- User Profile Section -->
        <div class="mb-4">
            <h2 class="text-xl font-bold">Daisy Pharmacy</h2>
            if user, err := auth.GetCurrentUser(ctx); err == nil && user != nil {
                <div class="mt-4 rounded-lg bg-base-200 p-3">
                    <div class="flex items-center gap-3">
                        <div class="avatar placeholder">
                            <div class="bg-primary text-primary-content w-10 rounded-full">
                                <span class="text-sm">{ getInitials(user.Name) }</span>
                            </div>
                        </div>
                        <div class="flex-1 overflow-hidden">
                            <p class="truncate text-sm font-semibold">{ user.Name }</p>
                            <p class="truncate text-xs text-base-content/70">{ user.Email }</p>
                        </div>
                    </div>
                    if auth.IsDevModeEnabled() {
                        <div class="mt-2 text-center">
                            <span class="badge badge-warning badge-xs">Dev Mode</span>
                        </div>
                    }
                </div>
            }
        </div>
        <!-- Rest of sidebar -->
    </aside>
}
```

**Base Layout**: `web/components/layouts/base.templ`
- Updated to pass `ctx` to Sidebar component

---

## 👤 User Information Source

The user information comes from the **authenticated user in context**:

### In Development Mode
- Shows mock user information based on `X-Mock-User` header
- Default: `admin` user
- Badge indicates dev mode is active

### In Production Mode
- Shows real user information from JWT token
- No dev mode badge displayed

---

## 🎯 User Display by Mock User

### Admin User
```
┌───────────────────────┐
│  ╔═╗  Dev Admin       │
│  ║ D ║  admin@dev.    │
│  ╚═╝     local        │
│      [Dev Mode]       │
└───────────────────────┘
```

### Doctor User
```
┌───────────────────────┐
│  ╔═╗  Dr. Dev         │
│  ║ D ║  doctor@dev.   │
│  ╚═╝     local        │
│      [Dev Mode]       │
└───────────────────────┘
```

### Pharmacist User
```
┌───────────────────────┐
│  ╔═╗  Dev Pharmacist  │
│  ║ D ║  pharmacist@   │
│  ╚═╝     dev.local    │
│      [Dev Mode]       │
└───────────────────────┘
```

### Nurse User
```
┌───────────────────────┐
│  ╔═╗  Dev Nurse       │
│  ║ D ║  nurse@dev.    │
│  ╚═╝     local        │
│      [Dev Mode]       │
└───────────────────────┘
```

### Read-Only User
```
┌───────────────────────┐
│  ╔═╗  Dev Readonly    │
│  ║ D ║  readonly@dev. │
│  ╚═╝     local        │
│      [Dev Mode]       │
└───────────────────────┘
```

---

## 🧪 Testing

### Test Different Users

When running in dev mode, switch between users and see the sidebar update:

```bash
# Start the app
go run cmd/server/main.go

# In browser, open different pages with different users:
# - Navigate to http://localhost:8080/patients
# - Use browser extension to set X-Mock-User header
# - Or use curl to test:

curl -H "X-Mock-User: doctor" http://localhost:8080/patients
curl -H "X-Mock-User: pharmacist" http://localhost:8080/prescriptions
curl -H "X-Mock-User: nurse" http://localhost:8080/dashboard
```

### Browser User Switcher (Optional Enhancement)

For easier testing in the browser, you can add a user switcher. Add this to your main layout:

```html
<!-- Add to base.templ or sidebar for dev mode only -->
if auth.IsDevModeEnabled() {
    <div class="mb-4 rounded-lg border-2 border-warning bg-warning/10 p-2">
        <p class="text-xs font-bold text-warning">Dev Mode - Switch User</p>
        <select 
            class="select select-bordered select-xs mt-2 w-full"
            onchange="switchMockUser(this.value)"
            id="mock-user-select">
            <option value="admin">Admin</option>
            <option value="doctor">Doctor</option>
            <option value="pharmacist">Pharmacist</option>
            <option value="nurse">Nurse</option>
            <option value="readonly">Read Only</option>
        </select>
    </div>
    
    <script>
        // Store and apply mock user
        const currentUser = localStorage.getItem('mockUser') || 'admin';
        document.getElementById('mock-user-select').value = currentUser;
        
        function switchMockUser(user) {
            localStorage.setItem('mockUser', user);
            location.reload();
        }
        
        // Add to all HTMX requests
        document.body.addEventListener('htmx:configRequest', function(evt) {
            const mockUser = localStorage.getItem('mockUser') || 'admin';
            evt.detail.headers['X-Mock-User'] = mockUser;
        });
    </script>
}
```

---

## 🎨 Styling

The user profile card uses DaisyUI classes:

- **Container**: `rounded-lg bg-base-200 p-3`
- **Avatar**: `avatar placeholder` with `bg-primary text-primary-content`
- **Name**: `text-sm font-semibold`
- **Email**: `text-xs text-base-content/70`
- **Dev Badge**: `badge badge-warning badge-xs`

### Customize Appearance

You can customize the colors and styling:

```go
// Change avatar color based on role
if user, err := auth.GetCurrentUser(ctx); err == nil && user != nil {
    var avatarClass string
    if hasPermission(user.Permissions, "admin:all") {
        avatarClass = "bg-error"  // Red for admin
    } else if hasPermission(user.Permissions, "doctor:role") {
        avatarClass = "bg-primary"  // Blue for doctor
    } else if hasPermission(user.Permissions, "pharmacist:role") {
        avatarClass = "bg-success"  // Green for pharmacist
    } else {
        avatarClass = "bg-neutral"  // Gray for others
    }
    
    <div class={"avatar placeholder"}>
        <div class={avatarClass + " text-primary-content w-10 rounded-full"}>
            <span class="text-sm">{ getInitials(user.Name) }</span>
        </div>
    </div>
}
```

---

## 🔒 Security Considerations

### User Context Required

The sidebar requires:
1. **Authentication middleware** must run before rendering pages
2. **User must be in context** via `auth.SetUser()`
3. Routes must use `auth.RequireAuthFromCookie()` or similar

### Error Handling

If user is not in context:
- Sidebar gracefully handles missing user
- No error displayed (optional: could show login prompt)

### Production vs Development

**Development**:
- Shows mock user name/email
- Displays "Dev Mode" badge
- User info from mock users

**Production**:
- Shows real user from JWT token
- No dev mode badge
- User info from JWT claims

---

## 📱 Responsive Design

The sidebar is fixed width (w-64 = 16rem) and scrollable:

```html
<aside class="flex h-screen w-64 flex-col border-r bg-base-100 p-4">
  <!-- Fixed header with user -->
  <div class="mb-4">...</div>
  
  <!-- Scrollable navigation -->
  <nav class="flex-1 overflow-y-auto">...</nav>
  
  <!-- Fixed footer with theme -->
  <div class="mt-auto border-t border-base-300 pt-4">...</div>
</aside>
```

---

## 🎯 Benefits

- ✅ **Visual confirmation** - User sees who they're logged in as
- ✅ **Dev mode indicator** - Clear warning when in dev mode
- ✅ **Professional appearance** - Avatar + name + email
- ✅ **Automatic** - No extra code needed in pages
- ✅ **Context-aware** - Shows current authenticated user
- ✅ **Graceful** - Handles missing user without errors

---

## 🚀 What You See

### In Development (with dev mode enabled):

When you browse to any page:
```
┌────────────────────────────┐
│  Daisy Pharmacy            │
│  ┌──────────────────────┐  │
│  │  ╔═╗  Dr. Dev       │  │
│  │  ║ D ║  doctor@dev. │  │
│  │  ╚═╝     local      │  │
│  │      [Dev Mode]     │  │
│  └──────────────────────┘  │
│                            │
│  Dashboard                 │
│  Patients                  │
│  ...                       │
└────────────────────────────┘
```

### In Production (with real JWT):

```
┌────────────────────────────┐
│  Daisy Pharmacy            │
│  ┌──────────────────────┐  │
│  │  ╔═╗  Dr. Smith     │  │
│  │  ║ D ║  jsmith@     │  │
│  │  ╚═╝     hospital   │  │
│  └──────────────────────┘  │
│                            │
│  Dashboard                 │
│  Patients                  │
│  ...                       │
└────────────────────────────┘
```

---

## 📚 Related Documentation

- **[SECURITY_MOCK_USERS.md](./SECURITY_MOCK_USERS.md)** - Mock user reference
- **[SECURITY_DEV_MODE.md](./SECURITY_DEV_MODE.md)** - Dev mode guide
- **[SECURITY_ARCHITECTURE.md](./SECURITY_ARCHITECTURE.md)** - Overall security architecture

---

## 🎉 Summary

The sidebar now displays:
- ✅ Current user's name and email
- ✅ User's initial in a colored avatar
- ✅ Dev mode indicator (when active)
- ✅ Automatic - works with both mock users and real JWT
- ✅ Professional and clean design

**The user always knows who they're logged in as!** 👤✨

