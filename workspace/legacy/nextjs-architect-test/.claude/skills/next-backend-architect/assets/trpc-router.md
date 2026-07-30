```ts
// [entity].router.ts
import { createTRPCRouter, protectedProcedure, baseProcedure } from "@/trpc/init"
import { TRPCError } from "@trpc/server"
import * as [entity]Service from "./[entity].service"
import { create[Entity]Schema, update[Entity]Schema, delete[Entity]Schema } from "../schemas/[entity].schema"

export const [entity]Router = createTRPCRouter({
  list: baseProcedure.query(async () => {
    try {
      return await [entity]Service.list()
    } catch (err) {
      throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to list [entity]s" })
    }
  }),

  get: baseProcedure
    .input(delete[Entity]Schema)
    .query(async ({ input }) => {
      try {
        return await [entity]Service.get(input.id)
      } catch (err) {
        if (err instanceof Error && err.message.includes("not found")) {
          throw new TRPCError({ code: "NOT_FOUND", message: "[Entity] not found" })
        }
        throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to get [entity]" })
      }
    }),

  create: protectedProcedure
    .input(create[Entity]Schema)
    .mutation(async ({ input }) => {
      try {
        return await [entity]Service.create(input)
      } catch (err) {
        if (err instanceof Error && err.message.includes("already exists")) {
          throw new TRPCError({ code: "CONFLICT", message: err.message })
        }
        throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to create [entity]" })
      }
    }),

  update: protectedProcedure
    .input(update[Entity]Schema)
    .mutation(async ({ input }) => {
      try {
        return await [entity]Service.update(input.id, input)
      } catch (err) {
        if (err instanceof Error && err.message.includes("not found")) {
          throw new TRPCError({ code: "NOT_FOUND", message: "[Entity] not found" })
        }
        throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to update [entity]" })
      }
    }),

  delete: protectedProcedure
    .input(delete[Entity]Schema)
    .mutation(async ({ input }) => {
      try {
        await [entity]Service.remove(input.id)
        return { success: true }
      } catch (err) {
        if (err instanceof Error && err.message.includes("not found")) {
          throw new TRPCError({ code: "NOT_FOUND", message: "[Entity] not found" })
        }
        throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Failed to delete [entity]" })
      }
    }),
})
```
