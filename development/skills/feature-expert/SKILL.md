---
name: feature-expert
description: >
  Feature expert manager. Manages the active and archived feature lifecycle in FEATURES.yml and spect/archive/FEATURES.yml. Use this skill when the user asks to initialize feature tracking, add or create a feature, update a feature, delete or remove a feature from the tracker, archive a completed feature, list or query tracked features, set sprint metadata (current_sprint, sprint_start, sprint_end), manage feature relationships (blocked_by, blocking, related_with), or link an existing feature to an already-authored PRD by reference only (including creating a feature and linking it to a PRD in one request). Does NOT apply to creating, authoring, editing, or reading PRD documents — creating a new PRD is never a trigger even when the request mentions a feature — nor to creating, editing, or reading ADR files, coding tasks, git operations including branch creation, CI/CD, package management, or general project architecture questions.
---
<skill_content name='feature-expert'>

# Feature Expert

This skill manages the active and archived feature lifecycle across `FEATURES.yml` and `spect/archive/FEATURES.yml`. It is the single source of truth for feature tracking, feature–PRD linkages, and sprint alignment.

## Responsibilities & Scope

* **In Scope:**
  * Initializing root `FEATURES.yml` and `spect/archive/FEATURES.yml` when they do not exist.
  * Adding / creating features in `FEATURES.yml`.
  * Updating feature metadata (`id`, `name`, `path_to_prds`, `short_description`, `blocked_by`, `blocking`, `related_with`, `sprint`).
  * Deleting / removing a feature entry from the active tracker (tracker only — never delete PRD/ADR files).
  * Linking a feature to an existing PRD path under `spect/features/prd_XXXXX/index.md` (on create or later).
  * Archiving completed features into `spect/archive/FEATURES.yml`.
  * Managing sprint metadata (`current_sprint`, `sprint_start`, `sprint_end`) and listing/querying tracked features.
  * Validating that `path_to_prds` entries point at existing PRD index files.

* **Out of Scope:**
  * Creating, modifying, or deleting PRD files in `spect/features/`.
  * Creating, modifying, or deleting ADR files in `spect/adrs/`.

## Directory Conventions

* **Active Features:** `./FEATURES.yml`
* **Archived Features:** `./spect/archive/FEATURES.yml`
* **PRD Locations:** `./spect/features/prd_XXXXX/index.md`
* **ADR Locations:** `./spect/adrs/adr_XXXXX.md` (or `.json` / `.yml`)

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

### Archive (`spect/archive/FEATURES.yml`)

Prefer a flat `archived_features` list (matches `scripts/main.sh init`). Nested `archived_sprints` is also acceptable if already present.

```yaml
archived_sprints: []
archived_features:
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

## Workflows

Use the bundled helper when applicable: `bash scripts/main.sh …` (skill directory). For YAML edits, keep valid structure and copy results to the requested outputs directory.

### 1. Initialize (no tracker yet)

```bash
bash scripts/main.sh init
```

Creates `FEATURES.yml` (`features: []` + sprint metadata) and `spect/archive/FEATURES.yml` if missing.

### 2. Create a feature

Edit `FEATURES.yml`: append under `features` with a new `id`, `name`, `short_description`, empty lists for relationships/`path_to_prds` unless linking, and `sprint` = `sprint_metadata.current_sprint` unless specified.

### 3. Create a feature and link a PRD

Same as create, but set `path_to_prds` to the existing PRD index path. Before linking:

```bash
bash scripts/main.sh validate-prd spect/features/prd_XXXXX/index.md
```

Do **not** create or edit the PRD file.

### 4. Link an existing feature to a PRD

Find the feature by `id` or `name`, append the PRD path to `path_to_prds` after `validate-prd`. Do not invent PRDs.

### 5. Update a feature

Change only the requested fields; keep `id` stable unless the user asks to rename the id.

### 6. Delete a feature

Remove the feature object from the active `features` list only. Do not delete PRD/ADR files or archive entries unless the user asked to archive.

### 7. Archive a completed feature

1. Copy the feature object from active `FEATURES.yml` into `spect/archive/FEATURES.yml` under `archived_features` (add `archived_at` ISO timestamp when practical).
2. Remove it from the active `features` list.
3. Leave other active features unchanged.

## Helper script

`scripts/main.sh` supports:

* `init` — scaffold active + archive trackers
* `validate-prd <path>` — warn/fail if the PRD index path is missing

Add / update / delete / archive / link are done by editing the YAML files as above.

<skill_resources>
  <file>scripts/main.sh</file>
</skill_resources>
</skill_content>
