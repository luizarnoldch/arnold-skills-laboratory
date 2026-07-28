```tsx
// useSuspenseList[Entity]s.tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useSuspenseQuery } from "@tanstack/react-query"

const useSuspenseList[Entity]s = () => {
  const trpc = useTRPC()
  const { data } = useSuspenseQuery(trpc.[entity].list.queryOptions())

  return {
    [entity]s: data,
  }
}

export default useSuspenseList[Entity]s
```