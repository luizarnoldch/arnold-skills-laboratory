```ts
// product.repository.ts
import { db } from "@/lib/db"
import type { CreateProductInput, UpdateProductInput } from "../schemas/product.schema"
import type { Product } from "../schemas/product.schema"

export const findAll = async (): Promise<Product[]> => {
  return db.product.findMany({
    orderBy: { createdAt: "desc" },
  })
}

export const findById = async (id: string): Promise<Product | null> => {
  return db.product.findUnique({
    where: { id },
  })
}

export const findByName = async (name: string): Promise<Product | null> => {
  return db.product.findFirst({
    where: { name },
  })
}

export const create = async (data: CreateProductInput): Promise<Product> => {
  return db.product.create({
    data,
  })
}

export const update = async (id: string, data: UpdateProductInput): Promise<Product> => {
  return db.product.update({
    where: { id },
    data,
  })
}

export const remove = async (id: string): Promise<void> => {
  await db.product.delete({
    where: { id },
  })
}
```