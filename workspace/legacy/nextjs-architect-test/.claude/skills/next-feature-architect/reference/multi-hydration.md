# Multi-Entity Hydration

When a view needs data from two or more entities, nest `Hydrate` components. Each entity gets its own `ErrorBoundary` + `Suspense` block.

```tsx
// src/features/[entity]/views/[Entity]View.tsx
import { Suspense } from "react"
import { ErrorBoundary } from "react-error-boundary"

import Hydrate[Entity]s from "@/features/[entity]/hooks/Hydrate[Entity]s"
import Hydrate[OtherEntity]s from "@/features/[other-entity]/hooks/Hydrate[OtherEntity]s"

import [Entity]Table from "../components/[Entity]Table"
import [Entity]ViewError from "../components/error/[Entity]ViewError"
import [Entity]ViewLoader from "../components/loaders/[Entity]ViewLoader"

import [OtherEntity]Table from "../components/[OtherEntity]Table"
import [OtherEntity]ViewError from "../components/error/[OtherEntity]ViewError"
import [OtherEntity]ViewLoader from "../components/loaders/[OtherEntity]ViewLoader"

const [Entity]View = () => {
  return (
    <Hydrate[Entity]s>
      <Hydrate[OtherEntity]s>
        <ErrorBoundary fallback={<[Entity]ViewError />}>
          <Suspense fallback={<[Entity]ViewLoader />}>
            <[Entity]Table />
          </Suspense>
        </ErrorBoundary>

        <ErrorBoundary fallback={<[OtherEntity]ViewError />}>
          <Suspense fallback={<[OtherEntity]ViewLoader />}>
            <[OtherEntity]Table />
          </Suspense>
        </ErrorBoundary>
      </Hydrate[OtherEntity]s>
    </Hydrate[Entity]s>
  )
}

export default [Entity]View
```

## Rules

- Each entity group is isolated — one entity failing does not crash the other.
- Hydrate nesting order does not matter for correctness, but nest related entities adjacently.
- Always place `ErrorBoundary` + `Suspense` inside the innermost `Hydrate` wrapping that entity.
