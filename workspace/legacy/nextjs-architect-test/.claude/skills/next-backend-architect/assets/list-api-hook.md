```tsx
// useList[Entity]s.tsx
"use client"

import { useQuery } from "@tanstack/react-query"
import { list[Entity]s } from "../server/[entity].api"

const useList[Entity]s = () => {
  const { data, isLoading, error } = useQuery({
      queryKey: ["[entity]s", "list"],
      queryFn: () => list[Entity]s(),
    })

  return {
    [entity]s: data ?? [],
    isLoading,
    error,
  }
}

export default useList[Entity]s
```
