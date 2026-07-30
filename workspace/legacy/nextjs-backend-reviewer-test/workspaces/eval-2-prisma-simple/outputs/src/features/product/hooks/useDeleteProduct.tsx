```ts
// useDeleteProduct.tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

type UseDeleteProductProps = {
  productId: string
  onSuccess?: () => void
  onError?: (error: unknown) => void
}

const useDeleteProduct = ({ productId, onSuccess, onError }: UseDeleteProductProps) => {
  const trpc = useTRPC()
  const queryClient = useQueryClient()

  const mutation = useMutation({
    ...trpc.product.delete.mutationOptions({ id: productId }),
    onSuccess: async () => {
      await queryClient.invalidateQueries(trpc.product.list.queryOptions())
      await queryClient.invalidateQueries(trpc.product.getById.queryOptions({ id: productId }))
      toast.success("Product deleted successfully")
      onSuccess?.()
    },
    onError: (error) => {
      toast.error("Failed to delete product")
      onError?.(error)
    },
  })

  return {
    ...mutation,
  }
}

export default useDeleteProduct
```