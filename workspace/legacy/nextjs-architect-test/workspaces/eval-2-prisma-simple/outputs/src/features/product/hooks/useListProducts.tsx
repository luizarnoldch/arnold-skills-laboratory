// useListProducts.tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useQuery } from "@tanstack/react-query"

const useListProducts = () => {
  const trpc = useTRPC()
  const { data, isLoading, error } = useQuery(trpc.product.list.queryOptions())

  return {
    products: data ?? [],
    isLoading,
    error,
  }
}

export default useListProducts