# Single-Entity View with Prefetch Hydration

Standard single-entity view structure with server-side prefetch hydration.

```tsx
// src/features/[entity]/views/[Entity]View.tsx
import { Suspense } from "react"
import { ErrorBoundary } from "react-error-boundary"
import Hydrate[Entity]s from "@/features/[entity]/hooks/Hydrate[Entity]s"

import [Entity]Table from "../components/[Entity]Table"
import [Entity]ViewError from "../components/error/[Entity]ViewError"
import [Entity]ViewLoader from "../components/loaders/[Entity]ViewLoader"

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