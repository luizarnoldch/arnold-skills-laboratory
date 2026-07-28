---
name: prd-writer
description: "Write structured PRDs following the LLM-PRD v1 format. Provides the full template, EARS acceptance criteria syntax, section writing rules, clarification question bank, file naming conventions, and Implementation Triggers format. Use when any agent needs to produce a PRD document from a feature description."
---

# PRD Writer Skill

Write a Product Requirements Document (PRD) following **LLM-PRD v1** — structured for LLM parsing and direct use by implementation agents and engineers.

---

## LLM-PRD v1 Template

````markdown
---
prd_id: prd-[feature-name]-v1
feature: [Feature Display Name]
status: draft
created: [YYYY-MM-DD]
authors: [agent-name]
stack: [tech stack summary]
---

## Overview

One paragraph. What is this feature and why does it exist?

## Problem Statement

What problem does this solve? Who experiences it? What is the cost of not solving it?

## Goals

- [ ] Goal 1
- [ ] Goal 2
- [ ] Goal 3

## Non-Goals

- Not in scope: X
- Not in scope: Y

## User Stories

| ID | As a... | I want to... | So that... |
|---|---|---|---|
| US-01 | [persona] | [action] | [benefit] |
| US-02 | [persona] | [action] | [benefit] |

## Acceptance Criteria

Each criterion maps to one or more User Stories. All criteria follow EARS syntax.

| ID | Story | Criterion |
|---|---|---|
| AC-01 | US-01 | When [trigger], the system shall [response]. |
| AC-02 | US-01 | When [trigger], the system shall [response]. |

## Data Model

### Entity: [EntityName]

| Field | Type | Required | Notes |
|---|---|---|---|
| id | [id type] | yes | Primary key |
| [field] | [type] | [yes/no] | [note] |

> Repeat block for each entity.

## UI/UX

### Pages

| Route | Component | Description |
|---|---|---|
| `/[entity]s` | `[Entity]ListView` | List all records |
| `/[entity]s/new` | `[Entity]CreateForm` | Create form |
| `/[entity]s/[id]` | `[Entity]DetailView` | Detail / edit |

### Component Flow

Describe the component tree or user interaction flow in plain text or a simple diagram.

## Technical Context

- **Stack**: [framework and runtime]
- **Transport**: [API transport layer]
- **Database ORM**: [ORM or query builder]
- **Auth**: [if relevant, otherwise "n/a"]
- **External services**: [if any, otherwise "none"]

## Implementation Triggers

> Machine-parseable section. One YAML block per entity. Do not add prose here.
> The calling agent fills in actual agent/skill names and tech values for the target project.

```yaml
entity: [EntityName]
backend:
  agent: [backend-scaffolding-agent]
  layers: [all|schema|server|hooks]
  transport: [transport-value]
  database: [orm-value]
frontend:
  agent: [frontend-scaffolding-agent]
  flag: [scaffold-flag]
```

> Repeat for each entity.

## Out of Scope

- [item]
- [item]

## Open Questions

| # | Question | Owner | Status |
|---|---|---|---|
| 1 | [question] | [role] | open |
````

---

## EARS Acceptance Criteria

Every acceptance criterion must follow one of these EARS patterns:

| EARS keyword | Pattern | Use when |
|---|---|---|
| When | `When [event], the system shall [response].` | Event-driven behaviour |
| While | `While [state], the system shall [behaviour].` | State-dependent behaviour |
| If / Then | `If [condition], then the system shall [response].` | Optional or conditional path |
| Where | `Where [feature is included], the system shall [behaviour].` | Feature-dependent behaviour |

**Default pattern:** `When <trigger>, the system shall <response>.`

**Examples:**

- `When a user submits the form with valid data, the system shall create the record and redirect to the list view.`
- `When the session expires, the system shall redirect the user to the login page.`
- `If an optional email field is provided, the system shall send a confirmation email.`
- `While a background job is running, the system shall display a progress indicator.`

---

## Section Writing Rules

### Overview
- One paragraph only.
- Answer: what it is + why it exists.
- No implementation details.

### Problem Statement
- Describe the current pain or gap.
- Identify who is affected.
- Quantify impact if possible.

### Goals
- 3–6 measurable goals.
- Start each with an action verb: "Enable", "Reduce", "Allow", "Ensure".
- Each goal should be traceable to at least one User Story.

### Non-Goals
- Explicitly list what this feature will NOT do.
- Include deferred items that may appear related.
- Minimum 2 items.

### User Stories
- Minimum 2, maximum 10.
- Keep persona, action, and benefit distinct.
- IDs: US-01, US-02, …

### Acceptance Criteria
- Every User Story must have at least one criterion.
- Use EARS syntax (see above).
- IDs: AC-01, AC-02, …

### Data Model
- One section per entity.
- Include all fields needed to satisfy the User Stories.
- Note constraints (unique, nullable, FK references) in the Notes column.

### UI/UX
- List every page this feature requires.
- Component Flow: describe hierarchy — do not design layout details.

### Technical Context
- Fill all fields; write "n/a" or "none" for empty ones.
- Do not invent technology choices — only document what the calling agent detects or the user specifies.

### Implementation Triggers
- One YAML block per entity.
- The calling agent fills in `agent`, `layers`, `transport`, `database`, and `flag` values — no generic placeholders in the final PRD.
- Do NOT add any prose, comments, or explanation inside the YAML blocks.

### Out of Scope
- Minimum 2 items.
- If genuinely nothing is out of scope, add "Future enhancements" and "Performance optimisation".

### Open Questions
- Log any decision not resolved during clarification.
- Assign an owner (PM / eng / design) and set status to "open".

---

## Clarification Question Bank

Ask only what is genuinely missing. Never ask for information you can infer from context.

| Missing info | Question to ask |
|---|---|
| Target persona | "Who is the primary user — admin, end user, or both?" |
| Core action | "What is the single most important thing a user should be able to do?" |
| Entity shape | "What data needs to be stored? Any fields beyond the obvious?" |
| Transport preference | "Should the API use the project's default transport or something different?" |
| Auth scope | "Is this feature behind authentication, or public?" |
| Relational scope | "Does [Entity] belong to a user, org, or team — or is it standalone?" |
| Success metric | "How will you know this feature is working correctly in production?" |

---

## File Naming

| Artifact | Pattern |
|---|---|
| PRD (first version) | `specs/prd-[feature-name]-v1.md` |
| Revised PRD | `specs/prd-[feature-name]-v2.md` |
| Feature name (file/URL) | kebab-case — `user-invitations` |
| Feature name (display) | Title Case — `User Invitations` |
| Entity name (code) | PascalCase singular — `UserInvitation` |
| Entity name (URL segment) | lowercase plural — `/user-invitations` |

---

## Error Handling

| Situation | Action |
|---|---|
| `specs/` directory does not exist | Calling agent must run `mkdir -p specs` before writing |
| Entity name given in plural form | Convert to PascalCase singular in Data Model and Implementation Triggers |
| Open question blocks a PRD section | Write the section with a `[TBD: <question>]` placeholder; log in Open Questions table |
| User requests revisions after writing | Overwrite the same file; bump to v2 only if the user explicitly says "new version" |
