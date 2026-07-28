```tsx
// useList[Entity]s.tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useQuery } from "@tanstack/react-query"

const useList[Entity]s = () => {
  const trpc = useTRPC()
  const { data, isLoading, error } = useQuery(trpc.[entity].list.queryOptions())

  return {
    [entity]s: data ?? [],
    isLoading,
    error,
  }
}

export default useList[Entity]s
```