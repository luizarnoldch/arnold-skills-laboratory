---
name: stamp-note
description: >
  Creates or updates STAMP.md with a short note and an ISO-8601 date when the user
  asks to stamp, mark, record, or write a note into STAMP.md (e.g. "stamp this",
  "update STAMP.md", "mark the file with today's note"). Does NOT apply to README
  edits, PRDs, ADRs, FEATURES.yml, source code, git, CI/CD, package management,
  or general documentation — only STAMP.md note stamping.
---

# Stamp Note

Write or refresh a single project stamp file: `STAMP.md`.

## Format

Always use exactly this structure:

```markdown
# Stamp
- date: <ISO-8601 datetime>
- note: <user message or summarized note>
```

## Rules

- Create `STAMP.md` if missing; overwrite the whole file when updating.
- Put the current date/time in ISO-8601 on the `date:` line.
- Put the user's note (or a faithful short paraphrase) on the `note:` line.
- Do not create or edit other files (README, code, FEATURES.yml, etc.).
- If the user asks for something outside stamping `STAMP.md`, refuse and say this skill only manages `STAMP.md`.
