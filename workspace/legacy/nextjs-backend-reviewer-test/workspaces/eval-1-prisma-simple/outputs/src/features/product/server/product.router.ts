// product.router.ts
import { createTRPCRouter, protectedProcedure, baseProcedure } from "@/trpc/init"
import { TRPCError } from "@trpc/server"
import * as productService from "./product.service"
import { createProductSchema, updateProductSchema, deleteProductSchema } from "../schemas/product.schema"

export const productRouter = createTRPCRouter({
  list: baseProcedure.query(async () => {
    try {
      return await productService.list()
    } catch (err) {
      throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to list products" })
    }
  }),

  get: baseProcedure
    .input(deleteProductSchema)
    .query(async ({ input }) => {
      try {
        return await productService.get(input.id)
      } catch (err) {
        if (err instanceof Error && err.message.includes("not found")) {
          throw new TRPCError({ code: "NOT_FOUND", message: "Product not found" })
        }
        throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to get product" })
      }
    }),

  create: protectedProcedure
    .input(createProductSchema)
    .mutation(async ({ input }) => {
      try {
        return await productService.create(input)
      } catch (err) {
        if (err instanceof Error && err.message.includes("already exists")) {
          throw new TRPCError({ code: "CONFLICT", message: err.message })
        }
        throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to create product" })
      }
    }),

  update: protectedProcedure
    .input(updateProductSchema)
    .mutation(async ({ input }) => {
      try {
        return await productService.update(input.id, input)
      } catch (err) {
        if (err instanceof Error && err.message.includes("not found")) {
          throw new TRPCError({ code: "NOT_FOUND", message: "Product not found" })
        }
        throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to update product" })
      }
    }),

  delete: protectedProcedure
    .input(deleteProductSchema)
    .mutation(async ({ input }) => {
      try {
        await productService.remove(input.id)
        return { success: true }
      } catch (err) {
        if (err instanceof Error && err.message.includes("not found")) {
          throw new TRPCError({ code: "NOT_FOUND", message: "Product not found" })
        }
        throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to delete product" })
      }
    }),
})
