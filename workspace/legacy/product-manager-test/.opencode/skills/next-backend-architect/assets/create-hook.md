```ts
// useCreate[Entity].tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useForm } from "@tanstack/react-form"
import { toast } from "sonner"
import { create[Entity]Schema } from "../schemas/[entity-kebab].schema"
import type { Create[Entity]Input } from "../schemas/[entity-kebab].schema"

type UseCreate[Entity]Props = {
  onSuccess?: () => void
  onError?: (error: unknown) => void
}

const useCreate[Entity] = ({ onSuccess, onError }: UseCreate[Entity]Props = {}) => {
  const trpc = useTRPC()
  const queryClient = useQueryClient()

  const mutation = useMutation({
    ...trpc.[entity].create.mutationOptions(),
    onSuccess: async () => {
      await queryClient.invalidateQueries(trpc.[entity].list.queryOptions())
      toast.success(`[Entity] created successfully`)
      onSuccess?.()
    },
    onError: (error) => {
      toast.error(`Failed to create [entity]`)
      onError?.(error)
    },
  })

  const form = useForm({
    defaultValues: {
      [requiredFields]: [defaultValue], // Use this for required fields like enums, dates, etc.
      [optionalFields]: [primitiveInitualValue], // Use this for optional fields
    } as Create[Entity]Input,
    validators: {
      onChange: create[Entity]Schema,
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

export default useCreate[Entity]
```
