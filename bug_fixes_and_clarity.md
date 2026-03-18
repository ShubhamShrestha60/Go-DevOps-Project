# bug_fixes_and_clarity.md

Addressing core application issues: Registration, Profile navigation, Project listing, and User onboarding.

## Proposed Changes

### [Backend] Routes & Handlers
#### [MODIFY] [router.go](file:///c:/Users/A%20C%20E%20R/Desktop/hello0/GoProject/Go-DevOps-Project/internal/router/router.go)
- Add `/profile` frontend route to the protected group.

#### [MODIFY] [dashboard_handler.go](file:///c:/Users/A%20C%20E%20R/Desktop/hello0/GoProject/Go-DevOps-Project/internal/handler/dashboard_handler.go)
- Implement `Profile` method to render the profile template.

#### [MODIFY] [project_handler.go](file:///c:/Users/A%20C%20E%20R/Desktop/hello0/GoProject/Go-DevOps-Project/internal/handler/project_handler.go)
- Fix `List` method to use the authenticated UserID from context instead of `uuid.Nil`.

### [Frontend] UI & UX
#### [MODIFY] [login.html](file:///c:/Users/A%20C%20E%20R/Desktop/hello0/GoProject/Go-DevOps-Project/web/templates/pages/login.html)
- Add a hidden registration form and JavaScript to toggle between Login and Register.
- Implement `/api/auth/register` call for new users.

#### [MODIFY] [dashboard.html](file:///c:/Users/A%20C%20E%20R/Desktop/hello0/GoProject/Go-DevOps-Project/web/templates/pages/dashboard.html)
- Add a "Quick Start Guide" section detailing project and task management capabilities.

## Verification Plan
1. **Registration**: Register a new user and login.
2. **Profile**: Navigate to the Profile page via the header and verify it loads (no 404).
3. **Projects**: Create a project and verify it appears in the list immediately.
4. **Dashboard**: Verify the new "Quick Start" section is visible and clear.
