# Config-Driven Role-Based Access Control (RBAC) System

A production-grade full-stack application demonstrating an advanced Role-Based Access Control (RBAC) system with granular permission management (Table-level & Field-level).

## 🚀 Project Overview

This system provides a robust framework for managing user access to resources. It features a React-powered frontend and a high-performance Go backend, utilizing a domain-driven architecture. The core strength of the system lies in its config-driven nature, where permissions are not hardcoded but managed dynamically through an administrative dashboard.

### Key Features
- **Advanced RBAC:** Go beyond simple roles with resource-specific and field-specific permissions.
- **Config-Driven:** Add new resources and fields in the database without changing backend logic.
- **Superadmin Privileges:** Dedicated `is_admin` flag for bypassing standard checks when necessary.
- **Email Invitation System:** Secure user onboarding via email invitations.
- **Secure Authentication:** JWT-based authentication with protected routes.

---

## 🏗️ Technology Stack

- **Frontend:** React (Vite), Tailwind CSS, Framer Motion, Lucide React.
- **Backend:** Go (Gin Framework), domain-driven architecture.
- **Database:** PostgreSQL (with production-grade schema).
- **Authentication:** JWT (JSON Web Tokens).

---

## 🔐 Role & Permission Model

The system uses a hierarchical and granular permission model:

### 1. Table Level Permission (Global Resource Access)
Controls high-level CRUD (Create, Read, Update, Delete) access to entire resources (e.g., Employees, Projects, Orders).
- **Read:** Can view the list and details of a resource.
- **Create:** Can add new records.
- **Update:** Can modify existing records.
- **Delete:** Can remove records.

### 2. Field Level Permission
Provides fine-grained control over specific columns/attributes within a resource.
- **View:** Controls visibility of specific fields (e.g., hide 'Salary' from certain roles).
- **Edit:** Controls whether a field can be modified even if 'Update' is granted at the table level.

### 3. Superadmin Status
Users with the `is_admin` flag set to `TRUE` in the database have full access to all resources and the Admin Console.

---

## 🚦 Getting Started

### Prerequisites
- [Node.js](https://nodejs.org/) (v16+)
- [Go](https://golang.google.com/dl/) (v1.20+)
- [PostgreSQL](https://www.postgresql.org/)

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd project-new
   ```

2. **Frontend Setup:**
   ```bash
   npm install
   ```

3. **Backend Setup:**
   Navigate to the backend directory and set up your environment variables:
   ```bash
   cd go-server
   cp .env.example .env
   ```
   *Edit `.env` with your PostgreSQL credentials and JWT secret.*

### Running the Application

You can run both the frontend and backend concurrently from the root directory:

```bash
npm run dev
```

- **Frontend:** `http://localhost:5173`
- **Backend Server:** `http://localhost:5001`

---

## 🔑 Sample Login Credentials

| Role | Username | Password |
|------|----------|----------|
| **Superadmin** | `admin` | `admin123` |

*Note: You can create more users and assign different roles via the Admin Console.*

---

## 📂 Project Structure

- **`src/`**: React Frontend source code.
  - `pages/Admin.jsx`: Granular permission management UI.
  - `pages/Dashboard.jsx`: Dynamic resource viewer based on permissions.
- **`go-server/`**: Go Backend (Domain-Based).
  - `cmd/server/main.go`: Application entry point.
  - `internal/`: Domain modules (auth, user, role, resource).
  - `pkg/utils/`: Shared utilities (JWT, random generators).
- **`cmd/`**: Utility scripts (e.g., `verify_admin`, `debug_perms`).
