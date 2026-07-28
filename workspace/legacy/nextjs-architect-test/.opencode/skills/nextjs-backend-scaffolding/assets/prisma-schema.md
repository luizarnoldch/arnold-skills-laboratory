```ts
import { z } from "zod"

// Define base schema manually or use generated Prisma types
export const [entity]Schema = z.object({
  id: z.uuid(),
  ...,
  createdAt: z.date(),
})

export const create[Entity]Schema = [entity]Schema.omit({ id: true, createdAt: true })
export const update[Entity]Schema = [entity]Schema.partial().required({ id: true })
export const delete[Entity]Schema = z.object({ id: z.uuid() })

export type [Entity] = z.infer<typeof [entity]Schema>
export type Create[Entity]Input = z.infer<typeof create[Entity]Schema>
export type Update[Entity]Input = z.infer<typeof update[Entity]Schema>
export type Delete[Entity]Input = z.infer<typeof delete[Entity]Schema>
```
