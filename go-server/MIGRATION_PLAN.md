# Migration Plan: Domain-Based Structure + Advanced RBAC Database

## Phase 1: Database Schema Upgrade ✅
- Add `users.is_admin` column
- Create `resources` table
- Create `resource_fields` table  
- Create `role_resource_permissions` table
- Create `role_field_permissions` table
- Migrate existing `permissions` data to new tables
- Keep backward compatibility during migration

## Phase 2: Code Restructure 🔄
- Create `cmd/server/main.go`
- Create `internal/` directory structure:
  - `internal/config/` (database)
  - `internal/auth/` (handler, service, repository, model)
  - `internal/user/` (handler, service, repository, model)
  - `internal/role/` (handler, service, repository, model)
  - `internal/permission/` (handler, service, repository, model, evaluator)
  - `internal/resource/` (handler, service, repository)
  - `internal/middleware/`
  - `internal/router/`
- Move `pkg/utils/`

## Phase 3: Update Business Logic 🛠️
- Update permission evaluator to use new tables
- Add field-level filtering in resource handlers
- Update admin checks to use `is_admin` flag
- Maintain all existing features:
  - ✅ Superadmin privileges
  - ✅ Email invitation flow
  - ✅ JWT authentication
  - ✅ Dynamic role creation
  - ✅ RBAC middleware

## Phase 4: Frontend Updates (if needed) 🎨
- No breaking changes for frontend
- All APIs remain the same

## Estimated Time: 30-40 minutes
