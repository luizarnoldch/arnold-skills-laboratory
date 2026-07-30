```ts
// useDelete[Entity].tsx
"use client"

import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { delete[Entity] as delete[Entity]Api } from "../server/[entity].api"

type UseDelete[Entity]Props = {
  [entity]Id: string
  onSuccess?: () => void
  onError?: (error: unknown) => void
}

const useDelete[Entity] = ({ [entity]Id, onSuccess, onError }: UseDelete[Entity]Props) => {
  const queryClient = useQueryClient()

  const mutation = useMutation({
    mutationFn: () => delete[Entity]Api([entity]Id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["[entity]s", "list"] })
      await queryClient.invalidateQueries({ queryKey: ["[entity]s", "detail", [entity]Id] })
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