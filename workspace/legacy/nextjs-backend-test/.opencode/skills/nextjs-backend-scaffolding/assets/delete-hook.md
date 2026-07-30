```ts
// useDelete[Entity].tsx
"use client"

import { useTRPC } from "@/trpc/client"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"

type UseDelete[Entity]Props = {
  [entity]Id: string
  onSuccess?: () => void
  onError?: (error: unknown) => void
}

const useDelete[Entity] = ({ [entity]Id, onSuccess, onError }: UseDelete[Entity]Props) => {
  const trpc = useTRPC()
  const queryClient = useQueryClient()

  const mutation = useMutation({
    ...trpc.[entity].delete.mutationOptions({ id: [entity]Id }),
    onSuccess: async () => {
      await queryClient.invalidateQueries(trpc.[entity].list.queryOptions())
      await queryClient.invalidateQueries(trpc.[entity].getById.queryOptions({ id: [entity]Id }))
      toast.success("[Entity] deleted successfully")
      onSuccess?.()
    },
    onError: (error) => {
      toast.error("Failed to delete [entity]")
      onError?.(error)
    },
  })

  return {
    ...mutation,
  }
}

export default useDelete[Entity]
```