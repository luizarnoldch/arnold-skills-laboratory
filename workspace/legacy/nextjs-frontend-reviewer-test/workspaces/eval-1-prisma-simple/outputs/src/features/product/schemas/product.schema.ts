import { z } from "zod"

// Define base schema manually or use generated Prisma types
export const productSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  description: z.string().nullable(),
  price: z.number(),
  stock: z.number().int(),
  createdAt: z.date(),
  updatedAt: z.date(),
})

export const createProductSchema = productSchema.omit({ id: true, createdAt: true, updatedAt: true })
export const updateProductSchema = productSchema.partial().required({ id: true })
export const deleteProductSchema = z.object({ id: z.uuid() })

export type Product = z.infer<typeof productSchema>
export type CreateProductInput = z.infer<typeof createProductSchema>
export type UpdateProductInput = z.infer<typeof updateProductSchema>
export type DeleteProductInput = z.infer<typeof deleteProductSchema>
