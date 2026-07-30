```ts
// [entity].repository.ts
import { db } from "@/lib/db"
import { [entityTable] } from "@/lib/db/schema"
import { eq, desc } from "drizzle-orm"
import type { Create[Entity]Input, Update[Entity]Input } from "../schemas/[entity].schema"
import type { [Entity] } from "../schemas/[entity].schema"

export const findAll = async (): Promise<[Entity][]> => {
  return db.select().from([entityTable]).orderBy(desc([entityTable].createdAt))
}

export const findById = async (id: string): Promise<[Entity] | undefined> => {
  const result = await db.select().from([entityTable]).where(eq([entityTable].id, id))
  return result[0]
}

export const findByName = async (name: string): Promise<[Entity] | undefined> => {
  const result = await db.select().from([entityTable]).where(eq([entityTable].name, name))
  return result[0]
}

export const create = async (data: Create[Entity]Input): Promise<[Entity]> => {
  const result = await db.insert([entityTable]).values(data).returning()
  return result[0]
}

export const update = async (id: string, data: Update[Entity]Input): Promise<[Entity]> => {
  const result = await db.update([entityTable]).set(data).where(eq([entityTable].id, id)).returning()
  return result[0]
}

export const remove = async (id: string): Promise<void> => {
  await db.delete([entityTable]).where(eq([entityTable].id, id))
}
```
