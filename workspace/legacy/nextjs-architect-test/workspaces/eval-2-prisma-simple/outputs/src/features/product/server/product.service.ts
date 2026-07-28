// product.service.ts
import * as productRepository from "./product.repository"
import type { CreateProductInput, UpdateProductInput } from "../schemas/product.schema"
import type { Product } from "../schemas/product.schema"

export const list = async (): Promise<Product[]> => {
  return productRepository.findAll()
}

export const get = async (id: string): Promise<Product> => {
  const product = await productRepository.findById(id)
  if (!product) {
    throw new Error("Product not found")
  }
  return product
}

export const create = async (data: CreateProductInput): Promise<Product> => {
  const existing = await productRepository.findByName(data.name)
  if (existing) {
    throw new Error("Product already exists")
  }
  return productRepository.create(data)
}

export const update = async (id: string, data: UpdateProductInput): Promise<Product> => {
  const product = await productRepository.findById(id)
  if (!product) {
    throw new Error("Product not found")
  }
  return productRepository.update(id, data)
}

export const remove = async (id: string): Promise<void> => {
  const product = await productRepository.findById(id)
  if (!product) {
    throw new Error("Product not found")
  }
  await productRepository.remove(id)
}