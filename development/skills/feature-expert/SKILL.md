---
name: feature-expert
description: >
  Feature expert manager. Manages the active and archived feature lifecycle in FEATURES.yml and spec/archive/FEATURES.yml. Triggers when the user asks to initialize, add, update, list, query, link, or archive features; view or set up tracking of features (e.g., "what features do we have tracked?"); set or change sprint metadata (current_sprint, sprint_start, sprint_end); manage feature relationships (blocked_by, blocking, related_with); or link an existing feature to an already-authored PRD by reference only. Does NOT apply to creating, authoring, or reading PRD documents — creating a new PRD is never a trigger even when the request mentions a feature — nor to creating, editing, or reading ADR files, coding tasks, git operations including branch creation, CI/CD, package management, or general project architecture questions.
---
<skill_content name='feature-expert'>

# Feature Expert

This skill manages the active and archived feature lifecycle across `FEATURES.yml` and `spect/archive/FEATURES.yml`. It acts as the single source of truth for feature tracking, feature-PRD linkages, and sprint alignment.

## Responsibilities & Scope
* **In Scope:**
  * Initializing root `FEATURES.yml` and `spect/archive/FEATURES.yml` structures.
  * Adding, updating, and archiving features in `FEATURES.yml`.
  * Tracking feature metadata (`id`, `name`, `path_to_prds`, `short_description`, `blocked_by`, `blocking`, `related_with`, `sprint`).
  * Managing current sprint metadata (`current_sprint`, `sprint_start`, `sprint_end`) and archiving finished sprints.
  * Validating paths pointing to `spect/features/prd_XXXXX/index.md` and referenced ADRs in `spect/adrs/adr_XXXXX`.
* **Out of Scope:**
  * Creating, modifying, or deleting PRD files in `spect/features/`.
  * Creating, modifying, or deleting ADR files in `spect/adrs/`.

## Directory Conventions
* **Active Features:** `./FEATURES.yml`
* **Archived Features:** `./spect/archive/FEATURES.yml`
* **PRD Locations:** `./spect/features/prd_XXXXX/index.md`
* **ADR Locations:** `./spect/adrs/adr_XXXXX.md` (or `.json`/`.yml`)

## File Schemas

### Active Features (`FEATURES.yml`)
```yaml
sprint_metadata:
  current_sprint: 12
  sprint_start: "2026-07-20"
  sprint_end: "2026-08-03"

features:
  - id: "feat-auth-01"
    name: "OAuth2 Refresh Tokens"
    path_to_prds:
      - "spect/features/prd_00102/index.md"
    short_description: "Implement silent token rotation and multi-device revocation."
    blocked_by: []
    blocking: ["feat-dash-04"]
    related_with: ["feat-user-02"]
    sprint: 12
```

### Archive (spect/archive/FEATURES.yml)

```yaml
archived_sprints:
  - sprint: 11
    sprint_start: "2026-07-06"
    sprint_end: "2026-07-20"
    completed_features:
      - id: "feat-auth-00"
        name: "Basic Login Workflow"
        path_to_prds:
          - "spect/features/prd_00099/index.md"
        short_description: "Initial user authentication setup."
        blocked_by: []
        blocking: []
        related_with: []
        sprint: 11
        archived_at: "2026-07-20T10:00:00Z"
```

### Workflows & Commands
When performing operations, use the provided helper script located at ./scripts/feature_manager.sh.

1. Initialize Feature Tracking
- Run: ./scripts/feature_manager.sh init
- Generates blank FEATURES.yml and spect/archive/FEATURES.yml files if they do not exist.

2. Add / Update a Feature
- Modifies FEATURES.yml under the features: block.
- Verifies that any provided path_to_prds paths point to existing files under spect/features/prd_XXXXX/index.md.

3. Archive Completed Feature or Sprint
- Run: ./scripts/feature_manager.sh archive <feature_id> or ./scripts/feature_manager.sh complete-sprint
- Moves finished features from FEATURES.yml into spect/archive/FEATURES.yml.

## 2. Helper Script (`scripts/feature_manager.sh`)

Place this executable script in `./scripts/feature_manager.sh` to automate file creations and checks safely.

<skill_resources>
  <file>scripts/main</file>
</skill_resources>
</skill_content>
