```tsx
// useUpdate[Entity].tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useForm } from "@tanstack/react-form"
import { toast } from "sonner"
import { update[Entity]Schema } from "../schemas/[entity].schema"
import type { Update[Entity]Input, [Entity] } from "../schemas/[entity].schema"

type UseUpdate[Entity]Props = {
  [entity]: [Entity]
  onSuccess?: () => void
  onError?: (error: unknown) => void
}

const useUpdate[Entity] = ({ [entity], onSuccess, onError }: UseUpdate[Entity]Props) => {
  const trpc = useTRPC()
  const queryClient = useQueryClient()

  const mutation = useMutation({
    ...trpc.[entity].update.mutationOptions(),
    onSuccess: async (data) => {
      await queryClient.invalidateQueries(trpc.[entity].list.queryOptions())
      await queryClient.invalidateQueries(trpc.[entity].getById.queryOptions({ id: data.id }))
      toast.success("[Entity] updated successfully")
      onSuccess?.()
    },
    onError: (error) => {
      toast.error("Failed to update [entity]")
      onError?.(error)
    },
  })

  // [entity].id is included in Update[Entity]Input so tRPC knows which record to update
  const form = useForm({
    defaultValues: {
      id: [entity].id,
      name: [entity].name,
      price: [entity].price,
      categoryId: [entity].categoryId,
    } as Update[Entity]Input,
    validators: {
      onChange: update[Entity]Schema,
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

export default useUpdate[Entity]
```