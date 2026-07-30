```tsx
// Hydrate[Entity].tsx
import { HydrateClient, prefetch, trpc } from "@/trpc/server"
import { ReactNode } from "react"

type Hydrate[Entity]sProps = {
  children: ReactNode
}

const Hydrate[Entity]s = ({ children }: Hydrate[Entity]sProps) => {
  prefetch(trpc.[entity].list.queryOptions())

  return (
    <HydrateClient>
      {children}
    </HydrateClient>
  )
}

export default Hydrate[Entity]s
```