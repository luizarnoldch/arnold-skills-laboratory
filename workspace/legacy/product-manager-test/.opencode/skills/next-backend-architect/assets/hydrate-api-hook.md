```tsx
// Hydrate[Entity]s.tsx
import { HydrateClient, apiPrefetch } from "@/trpc/server"
import { list[Entity]s } from "../server/[entity-kebab].api"
import { ReactNode } from "react"

type Hydrate[Entity]sProps = {
  children: ReactNode
}

const Hydrate[Entity]s = ({ children }: Hydrate[Entity]sProps) => {
  apiPrefetch({
    queryKey: ["[entity]s", "list"],
    queryFn: () => list[Entity]s(),
  })

  return (
    <HydrateClient>
      {children}
    </HydrateClient>
  )
}

export default Hydrate[Entity]s
```
