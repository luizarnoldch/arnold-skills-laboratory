```ts
import { z } from "zod"
import { createInsertSchema, createSelectSchema, createUpdateSchema } from "drizzle-zod"
import { [entityTable] } from "@/db/schema"

export const [entity]Schema = createSelectSchema([entityTable])
export const create[Entity]Schema = createInsertSchema([entityTable])
export const update[Entity]Schema = createUpdateSchema([entityTable])
export const delete[Entity]Schema = z.object({ id: z.uuid() })
  
export type [Entity] = z.infer<typeof [entity]Schema>
export type Create[Entity]Input = z.infer<typeof create[Entity]Schema>
export type Update[Entity]Input = z.infer<typeof update[Entity]Schema>
export type Delete[Entity]Input = z.infer<typeof delete[Entity]Schema>
```