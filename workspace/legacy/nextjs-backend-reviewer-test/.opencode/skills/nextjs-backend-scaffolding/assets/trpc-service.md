```ts
// [entity].service.ts
import * as [entity]Repository from "./[entity].repository"
import type { Create[Entity]Input, Update[Entity]Input } from "../schemas/[entity].schema"
import type { [Entity] } from "../schemas/[entity].schema"

export const list = async (): Promise<[Entity][]> => {
  return [entity]Repository.findAll()
}

export const get = async (id: string): Promise<[Entity]> => {
  const [entity] = await [entity]Repository.findById(id)
  if (![entity]) {
    throw new Error("[Entity] not found")
  }
  return [entity]
}

export const create = async (data: Create[Entity]Input): Promise<[Entity]> => {
  const existing = await [entity]Repository.findByName(data.name)
  if (existing) {
    throw new Error("[Entity] already exists")
  }
  return [entity]Repository.create(data)
}

export const update = async (id: string, data: Update[Entity]Input): Promise<[Entity]> => {
  const [entity] = await [entity]Repository.findById(id)
  if (![entity]) {
    throw new Error("[Entity] not found")
  }
  return [entity]Repository.update(id, data)
}

export const remove = async (id: string): Promise<void> => {
  const [entity] = await [entity]Repository.findById(id)
  if (![entity]) {
    throw new Error("[Entity] not found")
  }
  await [entity]Repository.remove(id)
}
```
