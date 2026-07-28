# Detail Page: page.tsx → DetailView → Suspense → Detail

Full chain for a detail page (`src/app/[entity]s/[id]/page.tsx`).

## Transport: tRPC

### page.tsx

```tsx

import [Entity]DetailView from "@/features/[entity]/views/[Entity]DetailView"

type [Entity]DetailPageProps = {
  params: Promise<{ id: string }>
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>
}

const [Entity]DetailPage = async ({ params, searchParams }: [Entity]DetailPageProps) => {
  const { id } = await params
  return (
    <[Entity]DetailView [entity]Id={id} />
  )
}

export default [Entity]DetailPage
```

### [Entity]DetailView.tsx

```tsx
import { Suspense } from "react"
import { ErrorBoundary } from "react-error-boundary"

import [Entity]Detail from "../components/[Entity]Detail"
import [Entity]DetailViewLoader from "../components/loaders/[Entity]DetailViewLoader"
import [Entity]DetailViewError from "../components/error/[Entity]DetailViewError"
import Hydrate[Entity] from "@/features/[entity]/hooks/Hydrate[Entity]"

type [Entity]DetailViewProps = {
  [entity]Id: string
}

const [Entity]DetailView = ({ [entity]Id }: [Entity]DetailViewProps) => {
  return (
    <Hydrate[Entity] [entity]Id={[entity]Id}>
      <div className="container mx-auto p-6">
        <ErrorBoundary fallback={<[Entity]DetailViewError />}>
          <Suspense fallback={<[Entity]DetailViewLoader />}>
            <[Entity]Detail [entity]Id={[entity]Id} />
          </Suspense>
        </ErrorBoundary>
      </div>
    </Hydrate[Entity]>
  )
}

export default [Entity]DetailView
```

### [Entity]Detail/index.tsx (suspense)

```tsx
"use client"

import useSuspenseList[Entity]s from "../../hooks/useSuspenseList[Entity]s"
import Delete[Entity]Button from "../Delete[Entity]Button"

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
      <Delete[Entity]Button [entity]Id={item.id} />
    </div>
  )
}

export default [Entity]Detail
```

### Hydrate[Entity] by ID (server component)

```tsx
import { HydrateClient, prefetch, trpc } from "@/trpc/server"
import type { ReactNode } from "react"

type Hydrate[Entity]Props = {
  children: ReactNode
  [entity]Id: string
}

const Hydrate[Entity] = ({ children, [entity]Id }: Hydrate[Entity]Props) => {
  prefetch(trpc.[entity].getById.queryOptions({ id: [entity]Id }))

  return (
    <HydrateClient>
      {children}
    </HydrateClient>
  )
}

export default Hydrate[Entity]
```

## Transport: REST API

Same page/view/component structure. Only `Hydrate[Entity]` and `useSuspenseList[Entity]s` internals differ (REST API client vs tRPC).
```txt
## Anti-patterns

- **Do NOT** put business logic in `page.tsx`. Only await `params`, hydrate, and render the view.
- **Do NOT** skip passing `[entity]Id` from page through view to component. Each layer must thread the ID.
- **Do NOT** add complex UI — just render data fields.
```
