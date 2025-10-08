# GraphQL Sample Queries

Copy and paste these queries into the GraphQL Playground at `http://localhost:8080/playground`

**Note:** GraphQL endpoint paths are defined in `internal/platform/paths/registry.go`

## 🚀 Quick Start Queries

### 1. Simple Patient Query

```graphql
query GetPatient {
  patient(id: "1") {
    id
    name
    phone
    state
  }
}
```

---

### 2. List All Patients

```graphql
query ListPatients {
  patients {
    id
    name
    phone
    state
    createdAt
  }
}
```

---

### 3. List Patients with Pagination

```graphql
query ListPatientsWithPagination {
  patients(limit: 10, offset: 0) {
    id
    name
    phone
  }
}
```

---

### 4. Search Patients

```graphql
query SearchPatients {
  patients(query: "John") {
    id
    name
    phone
  }
}
```

---

## 🔗 Nested Queries

### 5. Patient with Addresses

```graphql
query PatientWithAddresses {
  patient(id: "1") {
    id
    name
    phone
    addresses {
      id
      line1
      line2
      city
      state
      zip
    }
  }
}
```

---

### 6. Patient with Prescriptions

```graphql
query PatientWithPrescriptions {
  patient(id: "1") {
    id
    name
    phone
    prescriptions {
      id
      drug
      dose
      status
      createdAt
    }
  }
}
```

---

### 7. Patient with Everything (Nested)

```graphql
query PatientComplete {
  patient(id: "1") {
    id
    name
    phone
    state
    dob
    createdAt
    addresses {
      id
      line1
      city
      state
      zip
    }
    prescriptions {
      id
      drug
      dose
      status
      createdAt
    }
  }
}
```

---

## 💊 Prescription Queries

### 8. Get Single Prescription

```graphql
query GetPrescription {
  prescription(id: "1") {
    id
    drug
    dose
    status
    createdAt
  }
}
```

---

### 9. Prescription with Patient Info

```graphql
query PrescriptionWithPatient {
  prescription(id: "1") {
    id
    drug
    dose
    status
    patient {
      id
      name
      phone
    }
  }
}
```

---

### 10. List All Prescriptions

```graphql
query ListPrescriptions {
  prescriptions {
    id
    drug
    dose
    status
    createdAt
  }
}
```

---

### 11. Filter Prescriptions by Status

```graphql
query ActivePrescriptions {
  prescriptions(status: "Active") {
    id
    drug
    dose
    patient {
      name
    }
  }
}
```

---

### 12. Prescriptions with Pagination

```graphql
query PrescriptionsPaginated {
  prescriptions(status: "Active", limit: 20, offset: 0) {
    id
    drug
    dose
    status
  }
}
```

---

## 📊 Dashboard Queries

### 13. Dashboard Statistics

```graphql
query GetDashboardStats {
  dashboardStats {
    totalPatients
    activePrescriptions
  }
}
```

---

## 🎯 Complex Queries

### 14. Multiple Queries in One Request

```graphql
query GetMultipleData {
  # Get specific patient
  patient(id: "1") {
    name
    phone
  }
  
  # Get all patients
  allPatients: patients(limit: 5) {
    id
    name
  }
  
  # Get dashboard stats
  stats: dashboardStats {
    totalPatients
    activePrescriptions
  }
}
```

---

### 15. All Patients with Nested Data

```graphql
query AllPatientsWithDetails {
  patients(limit: 10) {
    id
    name
    phone
    state
    addresses {
      city
      state
    }
    prescriptions {
      drug
      status
    }
  }
}
```

---

### 16. Prescriptions with Full Patient Details

```graphql
query PrescriptionsWithFullPatient {
  prescriptions(status: "Active", limit: 10) {
    id
    drug
    dose
    status
    createdAt
    patient {
      id
      name
      phone
      dob
      state
      addresses {
        city
        state
      }
    }
  }
}
```

---

## 🔍 Using Query Variables

### 17. Patient by ID (with variables)

**Query:**
```graphql
query GetPatientById($patientId: ID!) {
  patient(id: $patientId) {
    id
    name
    phone
    prescriptions {
      drug
      dose
    }
  }
}
```

**Variables (in playground bottom panel):**
```json
{
  "patientId": "1"
}
```

---

### 18. Prescriptions with Filters (with variables)

**Query:**
```graphql
query GetPrescriptions($status: String, $limit: Int, $offset: Int) {
  prescriptions(status: $status, limit: $limit, offset: $offset) {
    id
    drug
    dose
    status
    patient {
      name
    }
  }
}
```

**Variables:**
```json
{
  "status": "Active",
  "limit": 10,
  "offset": 0
}
```

---

### 19. Search Patients (with variables)

**Query:**
```graphql
query SearchPatients($searchQuery: String, $limit: Int) {
  patients(query: $searchQuery, limit: $limit) {
    id
    name
    phone
    state
  }
}
```

**Variables:**
```json
{
  "searchQuery": "John",
  "limit": 5
}
```

---

## 🎨 Aliases (Custom Field Names)

### 20. Using Aliases

```graphql
query MultiplePatients {
  # Get patient 1
  john: patient(id: "1") {
    name
    phone
  }
  
  # Get patient 2
  jane: patient(id: "2") {
    name
    phone
  }
  
  # Get active prescriptions
  activeMeds: prescriptions(status: "Active") {
    drug
  }
  
  # Get all prescriptions
  allMeds: prescriptions {
    drug
  }
}
```

---

## 🧪 Introspection Queries

### 21. Schema Exploration

**List all types:**
```graphql
query GetAllTypes {
  __schema {
    types {
      name
      kind
      description
    }
  }
}
```

**Get Patient type details:**
```graphql
query GetPatientType {
  __type(name: "Patient") {
    name
    fields {
      name
      type {
        name
        kind
      }
    }
  }
}
```

**Get all queries:**
```graphql
query GetAllQueries {
  __schema {
    queryType {
      fields {
        name
        type {
          name
        }
      }
    }
  }
}
```

---

## 📝 Tips for Playground

### Auto-complete
- Press `Ctrl+Space` to see available fields
- Start typing and get suggestions

### Documentation
- Click "Docs" button on right side
- Explore schema interactively

### Prettify
- Click "Prettify" button to format query

### History
- Previous queries are saved
- Access via history panel

### Variables Panel
- Click "Query Variables" at bottom
- Add JSON variables for parameterized queries

---

## 🎯 Quick Test Sequence

**Copy-paste these in order to test everything:**

```graphql
# 1. Dashboard stats
query { dashboardStats { totalPatients activePrescriptions } }

# 2. List patients
query { patients(limit: 5) { id name phone } }

# 3. Get specific patient with nested data
query { 
  patient(id: "1") { 
    name 
    addresses { city } 
    prescriptions { drug status } 
  } 
}

# 4. List prescriptions
query { prescriptions(limit: 5) { id drug dose status } }

# 5. Prescription with patient
query { 
  prescription(id: "1") { 
    drug 
    patient { name phone } 
  } 
}
```

---

## 🚨 Common Issues

### "Cannot query field on type"
- Field doesn't exist in schema
- Check spelling
- Check schema.graphql for available fields

### "Cannot return null for non-nullable field"
- Query returned null for required field (marked with `!`)
- Data might not exist in database/memory

### No data returned
- Check if sample data exists in repositories
- Check server logs for errors

---

## 💡 Pro Tips

### Format Response
Queries return exactly what you ask for:
```graphql
# Ask for this:
query { patient(id: "1") { name } }

# Get this:
{
  "data": {
    "patient": {
      "name": "John Doe"
    }
  }
}
```

### Combine Multiple Queries
```graphql
query GetEverything {
  stats: dashboardStats { totalPatients }
  recentPatients: patients(limit: 3) { name }
  activeMeds: prescriptions(status: "Active", limit: 5) { drug }
}
```

### Use Fragments
```graphql
fragment PatientInfo on Patient {
  id
  name
  phone
}

query GetPatients {
  patient1: patient(id: "1") { ...PatientInfo }
  patient2: patient(id: "2") { ...PatientInfo }
}
```

---

## 🎉 Summary

**Start with these:**
1. `query { dashboardStats { totalPatients activePrescriptions } }`
2. `query { patients(limit: 5) { id name phone } }`
3. `query { patient(id: "1") { name addresses { city } } }`

**Then explore:**
- Add more fields
- Try nested queries
- Use variables
- Test pagination

**Happy querying! 🚀**

