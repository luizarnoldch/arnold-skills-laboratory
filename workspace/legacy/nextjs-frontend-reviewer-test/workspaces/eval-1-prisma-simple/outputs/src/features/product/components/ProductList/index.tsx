"use client"
import { useState } from "react"
import useSuspenseListProducts from "../../hooks/useSuspenseListProducts"
import useDeleteProduct from "../../hooks/useDeleteProduct"
import type { Product } from "../../schemas/product.schema"
import ProductFormCreate from "../ProductFormCreate"
import ProductFormUpdate from "../ProductFormUpdate"

const ProductList = () => {
  const { products } = useSuspenseListProducts()
  const [showCreate, setShowCreate] = useState(false)
  const [editing, setEditing] = useState<Product | null>(null)
  const { mutate: deleteProduct, isPending: isDeleting } = useDeleteProduct()

  return (
    <div className="grid gap-4">
      <div className="flex items-center justify-between">
        <span className="text-sm text-gray-500">{products.length} products</span>
        <button
          onClick={() => { setShowCreate(true); setEditing(null) }}
          className="rounded bg-blue-600 px-3 py-1.5 text-sm text-white hover:bg-blue-700"
        >
          Create Product
        </button>
      </div>

      {showCreate && <ProductFormCreate onClose={() => setShowCreate(false)} />}
      {editing && <ProductFormUpdate product={editing} onClose={() => setEditing(null)} />}

      {products.length === 0 ? (
        <p className="py-8 text-center text-sm text-gray-500">No products found.</p>
      ) : (
        <table className="w-full border-collapse text-sm">
          <thead>
            <tr className="border-b bg-gray-50">
              <th className="px-3 py-2 text-left font-medium text-gray-600">ID</th>
              <th className="px-3 py-2 text-left font-medium text-gray-600">Name</th>
              <th className="px-3 py-2 text-left font-medium text-gray-600">Description</th>
              <th className="px-3 py-2 text-right font-medium text-gray-600">Price</th>
              <th className="px-3 py-2 text-right font-medium text-gray-600">Stock</th>
              <th className="px-3 py-2 text-right font-medium text-gray-600">Actions</th>
            </tr>
          </thead>
          <tbody>
            {products.map((product) => (
              <tr key={product.id} className="border-b hover:bg-gray-50">
                <td className="px-3 py-2 font-mono text-xs text-gray-500">{product.id}</td>
                <td className="px-3 py-2">{product.name}</td>
                <td className="px-3 py-2 text-gray-600 max-w-xs truncate">{product.description ?? "-"}</td>
                <td className="px-3 py-2 text-right">${product.price.toFixed(2)}</td>
                <td className="px-3 py-2 text-right">{product.stock}</td>
                <td className="px-3 py-2 text-right">
                  <div className="flex justify-end gap-3">
                    <button
                      onClick={() => { setEditing(product); setShowCreate(false) }}
                      className="text-xs text-blue-600 underline hover:no-underline"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => deleteProduct({ productId: product.id })}
                      disabled={isDeleting}
                      className="text-xs text-red-600 underline hover:no-underline disabled:opacity-50"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default ProductList
