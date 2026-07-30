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
        <form.Field name="id">
          {({ value }) => (
            <input type="hidden" value={value} />
          )}
        </form.Field>
        <form.Field name="name">
          {({ value, onChange }) => (
            <div className="grid gap-1.5">
              <label className="text-sm font-medium">Name</label>
              <input
                type="text"
                value={value}
                onChange={(e) => onChange(e.target.value)}
                className="rounded border px-3 py-2"
                placeholder="Product name"
              />
            </div>
          )}
        </form.Field>
        <form.Field name="description">
          {({ value, onChange }) => (
            <div className="grid gap-1.5">
              <label className="text-sm font-medium">Description</label>
              <textarea
                value={value}
                onChange={(e) => onChange(e.target.value)}
                className="rounded border px-3 py-2"
                placeholder="Product description"
                rows={3}
              />
            </div>
          )}
        </form.Field>
        <form.Field name="price">
          {({ value, onChange }) => (
            <div className="grid gap-1.5">
              <label className="text-sm font-medium">Price</label>
              <input
                type="number"
                step="0.01"
                value={value}
                onChange={(e) => onChange(parseFloat(e.target.value) || 0)}
                className="rounded border px-3 py-2"
                placeholder="Price"
              />
            </div>
          )}
        </form.Field>
        <form.Field name="stock">
          {({ value, onChange }) => (
            <div className="grid gap-1.5">
              <label className="text-sm font-medium">Stock</label>
              <input
                type="number"
                value={value}
                onChange={(e) => onChange(parseInt(e.target.value) || 0)}
                className="rounded border px-3 py-2"
                placeholder="Stock quantity"
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
