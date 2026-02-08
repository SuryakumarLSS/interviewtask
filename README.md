# Config-Driven Access Control System (RBAC)

This is a full-stack application demonstrating a Role-Based Access Control (RBAC) system with granular permission management (Table-level & Field-level).

## Technology Stack
- **Frontend:** React (Vite), Tailwind CSS, React Dictionary
- **Backend:** Node.js, Express.js (Chosen due to environment constraints - Golang was unavailable)
- **Database:** SQLite (Zero-config, embedded)
- **Authentication:** JWT (JSON Web Tokens)

## Features
1. **Authentication:**
   - Login system using JWT.
   - Protected routes.
   - User session management.

2. **Role & Permission Management (Admin Only):**
   - Create and manage Roles.
   - Create Users and assign Roles.
   - **Dynamic Permission Config:** Matrix UI to enable/disable access to specific resources (Employees, Projects, Orders) for each action (Read, Create, Update, Delete).
   - **Field-Level Security:** Define specific columns (attributes) allowed for Read/Write operations (e.g., `name, salary` vs `*`).

3. **Data Management:**
   - Generic API endpoints (`/api/data/:resource`) that enforce configured permissions dynamically.
   - Dashboard sidebar adapts to user's read permissions.
   - Actions (Edit/Delete/Create) are only shown if permitted.
   - Backend validation ensures no unauthorized fields are modified.

## Setup & Run

### Prerequisites
- Node.js installed.

### Installation
1. Install dependencies:
   ```bash
   npm install
   ```

### Running the App
Run both backend and frontend concurrently with a single command:
```bash
npm run dev
```

The app will open at `http://localhost:5173` (or 5174 if port is busy).
The backend server runs on `http://localhost:5000`.

## Login Credentials
**Admin User:**
- Username: `admin`
- Password: `admin123`

**Sample User:**
- You can create new users via the Admin Panel.

## Project Structure
- `server/` - Backend API & Database logic
  - `database.js` - SQLite schema & seeding
  - `middleware/rbac.js` - The core permission enforcement logic
  - `routes/` - API endpoints
- `src/` - React Frontend
  - `pages/Dashboard.tsx` - Dynamic resource viewer
  - `pages/Admin.tsx` - Role & Permission manager
  - `services/api.ts` - Axios API layer

## key Design Decisions
- **Config-Driven:** Permissions are stored in the database (`permissions` table) and loaded dynamically. The UI and Backend both query this configuration to determine access.
- **Generic Resource Handler:** Instead of writing separate controllers for each entity, a generic router (`server/routes/data.js`) handles CRUD for whitelisted resources, applying the RBAC middleware automatically. This makes adding new resources extremely easy (just add table & whitelist).
