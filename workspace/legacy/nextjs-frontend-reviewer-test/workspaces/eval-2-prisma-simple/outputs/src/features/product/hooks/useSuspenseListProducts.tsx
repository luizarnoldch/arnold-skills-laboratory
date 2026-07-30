// useSuspenseListProducts.tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useSuspenseQuery } from "@tanstack/react-query"

const useSuspenseListProducts = () => {
  const trpc = useTRPC()
  const { data } = useSuspenseQuery(trpc.product.list.queryOptions())

  return {
    products: data,
  }
}

export default useSuspenseListProducts
