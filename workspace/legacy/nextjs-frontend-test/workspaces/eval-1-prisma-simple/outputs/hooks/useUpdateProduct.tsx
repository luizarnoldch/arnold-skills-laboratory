// useUpdateProduct.tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useForm } from "@tanstack/react-form"
import { toast } from "sonner"
import { updateProductSchema } from "../schemas/product.schema"
import type { UpdateProductInput, Product } from "../schemas/product.schema"

type UseUpdateProductProps = {
  product: Product
  onSuccess?: () => void
  onError?: (error: unknown) => void
}

const useUpdateProduct = ({ product, onSuccess, onError }: UseUpdateProductProps) => {
  const trpc = useTRPC()
  const queryClient = useQueryClient()

  const mutation = useMutation({
    ...trpc.product.update.mutationOptions(),
    onSuccess: async (data) => {
      await queryClient.invalidateQueries(trpc.product.list.queryOptions())
      await queryClient.invalidateQueries(trpc.product.get.queryOptions({ id: data.id }))
      toast.success("Product updated successfully")
      onSuccess?.()
    },
    onError: (error) => {
      toast.error("Failed to update product")
      onError?.(error)
    },
  })

  // product.id is included in UpdateProductInput so tRPC knows which record to update
  const form = useForm({
    defaultValues: {
      id: product.id,
      name: product.name,
      description: product.description,
      price: product.price,
      stock: product.stock,
    } as UpdateProductInput,
    validators: {
      onChange: updateProductSchema,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync(value)
    },
  })

  return {
    form,
    ...mutation,
  }
}

export default useUpdateProduct
