# Query Hook Usage Patterns

## Default: `useSuspenseList[Entity]s` (preferred)

Use **always** unless the user explicitly asks for non-suspense. Parent must wrap with `<ErrorBoundary>` and `<Suspense>`.

```tsx
"use client"

import useSuspenseList[Entity]s from "../../hooks/useSuspenseList[Entity]s"

const [Entity]Table = () => {
  const { [entity]s } = useSuspenseList[Entity]s()

  return (
    <div>
      {[entity]s.map((item) => (
        <div key={item.id}>{item.name}</div>
      ))}
    </div>
  )
}

export default [Entity]Table
```

### Detail variant: filter by ID from list

```tsx
"use client"

import useSuspenseList[Entity]s from "../../hooks/useSuspenseList[Entity]s"

type [Entity]DetailProps = {
  [entity]Id: string
}

const [Entity]Detail = ({ [entity]Id }: [Entity]DetailProps) => {
  const { [entity]s } = useSuspenseList[Entity]s()
  const item = [entity]s.find((i) => i.id === [entity]Id)

  if (!item) {
    return <div>Not found</div>
  }

  return (
    <div>
      <h1>{item.name}</h1>
    </div>
  )
}

export default [Entity]Detail
```

---

## Alternate: `useList[Entity]s` (non-suspense)

Use **only** when explicitly requested. Handles `isLoading`/`error` inline with minimal markup.

```tsx
"use client"

import useList[Entity]s from "../../hooks/useList[Entity]s"

const [Entity]List = () => {
  const { [entity]s, isLoading, error } = useList[Entity]s()

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (error) {
    return <div>Failed to load</div>
  }

  return (
    <div>
      {[entity]s.map((item) => (
        <div key={item.id}>{item.name}</div>
      ))}
    </div>
  )
}

export default [Entity]List
```

---

## Default: `useSuspenseList[Entity]s` (REST API)

Same component pattern. Only the hook import changes.

```tsx
"use client"

import useSuspenseList[Entity]s from "../../hooks/useSuspenseList[Entity]s"

const [Entity]Table = () => {
  const { [entity]s } = useSuspenseList[Entity]s()

  return (
    <div>
      {[entity]s.map((item) => (
        <div key={item.id}>{item.name}</div>
      ))}
    </div>
  )
}

export default [Entity]Table
```

---

## Component Rules

- Always use `type` for Props (never `interface`).
- Always use arrow function `const X = () => {}`.
- Always `export default` the main component.
- Never add complex UI — one `<div>` per element, no tables, no cards, no shadcn.
- Loading state: single `<div>Loading...</div>`.
- Error state: single `<div>Failed to load</div>`.
