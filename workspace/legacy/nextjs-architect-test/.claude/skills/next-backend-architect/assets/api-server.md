```ts
// [entity-kebab].api.ts
import { apiFetch } from "@/lib/api"
import type { Create[Entity]Input, Update[Entity]Input, [Entity] } from "../schemas/[entity-kebab].schema"

export async function list[Entity]s(): Promise<[Entity][]> {
  return apiFetch<[Entity][]>(`/[entity-kebab]s`, { method: "GET" })
}

export async function get[Entity](id: string): Promise<[Entity]> {
  return apiFetch<[Entity]>(`/[entity-kebab]s/${id}`, { method: "GET" })
}

export async function create[Entity](data: Create[Entity]Input): Promise<[Entity]> {
  return apiFetch<[Entity]>(`/[entity-kebab]s`, { method: "POST", body: JSON.stringify(data) })
}

export async function update[Entity](data: Update[Entity]Input): Promise<[Entity]> {
  return apiFetch<[Entity]>(`/[entity-kebab]s/${data.id}`, { method: "PUT", body: JSON.stringify(data) })
}

export async function delete[Entity](id: string): Promise<void> {
  return apiFetch<void>(`/[entity-kebab]s/${id}`, { method: "DELETE" })
}
```
