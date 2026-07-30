```tsx
// HydrateProducts.tsx
import { HydrateClient, prefetch, trpc } from "@/trpc/server"
import { ReactNode } from "react"

type HydrateProductsProps = {
  children: ReactNode
}

const HydrateProducts = ({ children }: HydrateProductsProps) => {
  prefetch(trpc.product.list.queryOptions())

  return (
    <HydrateClient>
      {children}
    </HydrateClient>
  )
}

export default HydrateProducts
```