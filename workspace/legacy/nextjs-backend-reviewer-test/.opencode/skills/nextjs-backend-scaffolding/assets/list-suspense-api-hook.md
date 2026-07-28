```tsx
// useSuspenseList[Entity]s.tsx
"use client"

import { useSuspenseQuery } from "@tanstack/react-query"
import { list[Entity]s } from "../server/[entity].api"

const useSuspenseList[Entity]s = () => {
  const { data } = useSuspenseQuery({
      queryKey: ["[entity]s", "list"],
      queryFn: () => list[Entity]s(),
    })

  return {
    [entity]s: data,
  }
}

export default useSuspenseList[Entity]s
```
