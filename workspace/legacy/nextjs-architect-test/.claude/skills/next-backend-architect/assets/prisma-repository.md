```ts
// [entity].repository.ts
import { db } from "@/lib/db"
import type { Create[Entity]Input, Update[Entity]Input } from "../schemas/[entity].schema"
import type { [Entity] } from "../schemas/[entity].schema"

export const findAll = async (): Promise<[Entity][]> => {
  return db.[entityTable].findMany({
    orderBy: { createdAt: "desc" },
  })
}

export const findById = async (id: string): Promise<[Entity] | null> => {
  return db.[entityTable].findUnique({
    where: { id },
  })
}

export const findByName = async (name: string): Promise<[Entity] | null> => {
  return db.[entityTable].findFirst({
    where: { name },
  })
}

export const create = async (data: Create[Entity]Input): Promise<[Entity]> => {
  return db.[entityTable].create({
    data,
  })
}

export const update = async (id: string, data: Update[Entity]Input): Promise<[Entity]> => {
  return db.[entityTable].update({
    where: { id },
    data,
  })
}

export const remove = async (id: string): Promise<void> => {
  await db.[entityTable].delete({
    where: { id },
  })
}
```
