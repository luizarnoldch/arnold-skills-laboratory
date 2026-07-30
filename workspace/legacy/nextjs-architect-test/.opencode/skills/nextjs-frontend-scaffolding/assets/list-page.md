# List Page: page.tsx → View → Suspense → Table

Full chain for a list page (`src/app/[entity]s/page.tsx`).

## Transport: tRPC

### page.tsx

```tsx
import [Entity]View from "@/features/[entity]/views/[Entity]View"

type [Entity]sPageProps = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const [Entity]sPage = async ({ searchParams }: [Entity]sPageProps) => {
  return (
    <[Entity]View />
  )
}

export default [Entity]sPage
```

### [Entity]View.tsx

```tsx
import { Suspense } from "react"
import { ErrorBoundary } from "react-error-boundary"

import [Entity]Table from "../components/[Entity]Table"
import Hydrate[Entity]s from "@/features/[entity]/hooks/Hydrate[Entity]s"
import [Entity]ViewLoader from "../components/loaders/[Entity]ViewLoader"
import [Entity]ViewError from "../components/error/[Entity]ViewError"

const [Entity]View = () => {
  return (
    <Hydrate[Entity]s>
      <div className="container mx-auto p-6">
        <ErrorBoundary fallback={<[Entity]ViewError />}>
          <Suspense fallback={<[Entity]ViewLoader />}>
            <[Entity]Table />
          </Suspense>
        </ErrorBoundary>
      </div>
    </Hydrate[Entity]s>
  )
}

export default [Entity]View
```

### [Entity]Table/index.tsx (suspense)

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

## Transport: REST API

### page.tsx (same as tRPC)

```tsx
import Hydrate[Entity]s from "@/features/[entity]/hooks/Hydrate[Entity]s"
import [Entity]View from "@/features/[entity]/views/[Entity]View"

type [Entity]sPageProps = {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const [Entity]sPage = async ({ searchParams }: [Entity]sPageProps) => {
  return (
    <Hydrate[Entity]s>
      <[Entity]View />
    </Hydrate[Entity]s>
  )
}

export default [Entity]sPage
```

View and table components are identical — only the hook internals differ (tRPC client vs apiFetch).
```txt
## Anti-patterns

- **Do NOT** put business logic in `page.tsx`. It is a thin pass-through.
- **Do NOT** skip `ErrorBoundary` when using `useSuspenseQuery`.
- **Do NOT** add complex UI (tables, cards, pagination) — just render the data.
```
