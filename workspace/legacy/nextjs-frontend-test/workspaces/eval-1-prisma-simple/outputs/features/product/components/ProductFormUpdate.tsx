"use client"
import useUpdateProduct from "../../hooks/useUpdateProduct"
import type { Product } from "../../schemas/product.schema"

type ProductFormUpdateProps = { product: Product; onClose: () => void }

const ProductFormUpdate = ({ product, onClose }: ProductFormUpdateProps) => {
  const { form, isPending } = useUpdateProduct({ product, onSuccess: onClose })

  return (
    <div className="rounded border bg-gray-50 p-4">
      <h3 className="mb-3 font-semibold">Edit Product</h3>
      <form
        onSubmit={(e) => { e.preventDefault(); e.stopPropagation(); form.handleSubmit() }}
        className="grid gap-3"
      >
        <form.Field name="name">
          {(field) => (
            <div className="grid gap-1">
              <label className="text-xs font-medium text-gray-600">Name</label>
              <input
                name={field.name}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                className="rounded border px-2 py-1.5 text-sm"
              />
            </div>
          )}
        </form.Field>
        <form.Field name="description">
          {(field) => (
            <div className="grid gap-1">
              <label className="text-xs font-medium text-gray-600">Description</label>
              <input
                name={field.name}
                value={field.state.value ?? ""}
                onChange={(e) => field.handleChange(e.target.value || null)}
                className="rounded border px-2 py-1.5 text-sm"
              />
            </div>
          )}
        </form.Field>
        <form.Field name="price">
          {(field) => (
            <div className="grid gap-1">
              <label className="text-xs font-medium text-gray-600">Price</label>
              <input
                name={field.name}
                type="number"
                step="0.01"
                value={field.state.value}
                onChange={(e) => field.handleChange(parseFloat(e.target.value))}
                className="rounded border px-2 py-1.5 text-sm"
              />
            </div>
          )}
        </form.Field>
        <form.Field name="stock">
          {(field) => (
            <div className="grid gap-1">
              <label className="text-xs font-medium text-gray-600">Stock</label>
              <input
                name={field.name}
                type="number"
                step="1"
                value={field.state.value}
                onChange={(e) => field.handleChange(parseInt(e.target.value, 10))}
                className="rounded border px-2 py-1.5 text-sm"
              />
            </div>
          )}
        </form.Field>
        <div className="flex justify-end gap-2 pt-1">
          <button
            type="button"
            onClick={onClose}
            disabled={isPending}
            className="rounded border px-3 py-1.5 text-sm hover:bg-gray-100 disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isPending}
            className="rounded bg-blue-600 px-3 py-1.5 text-sm text-white hover:bg-blue-700 disabled:opacity-50"
          >
            {isPending ? "Saving..." : "Save"}
          </button>
        </div>
      </form>
    </div>
  )
}

export default ProductFormUpdate
