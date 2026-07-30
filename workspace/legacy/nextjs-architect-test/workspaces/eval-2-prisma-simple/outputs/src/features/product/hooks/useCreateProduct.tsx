// useCreateProduct.tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useForm } from "@tanstack/react-form"
import { toast } from "sonner"
import { createProductSchema } from "../schemas/product.schema"
import type { CreateProductInput } from "../schemas/product.schema"

type UseCreateProductProps = {
  onSuccess?: () => void
  onError?: (error: unknown) => void
}

const useCreateProduct = ({ onSuccess, onError }: UseCreateProductProps = {}) => {
  const trpc = useTRPC()
  const queryClient = useQueryClient()

  const mutation = useMutation({
    ...trpc.product.create.mutationOptions(),
    onSuccess: async () => {
      await queryClient.invalidateQueries(trpc.product.list.queryOptions())
      toast.success(`Product created successfully`)
      onSuccess?.()
    },
    onError: (error) => {
      toast.error(`Failed to create product`)
      onError?.(error)
    },
  })

  const form = useForm({
    defaultValues: {
      name: "",
      description: "",
      price: 0,
      categoryId: undefined,
    } as CreateProductInput,
    validators: {
      onChange: createProductSchema,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync(value)
      form.reset()
    },
  })

  return {
    form,
    ...mutation,
  }
}

export default useCreateProduct