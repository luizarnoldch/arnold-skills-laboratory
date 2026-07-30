import { z } from "zod"

export const productSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().nullable(),
  price: z.number(),
  stock: z.number().int().default(0),
  createdAt: z.date(),
  updatedAt: z.date(),
})

export const createProductSchema = z.object({
  name: z.string(),
  description: z.string().nullable().optional(),
  price: z.number(),
  stock: z.number().int().default(0),
})

export const updateProductSchema = createProductSchema.partial().extend({
  id: z.string(),
})

export const deleteProductSchema = z.object({ id: z.string() })

export type Product = z.infer<typeof productSchema>
export type CreateProductInput = z.infer<typeof createProductSchema>
export type UpdateProductInput = z.infer<typeof updateProductSchema>
export type DeleteProductInput = z.infer<typeof deleteProductSchema>
