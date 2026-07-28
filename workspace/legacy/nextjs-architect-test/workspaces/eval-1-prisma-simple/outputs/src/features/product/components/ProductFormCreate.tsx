"use client"
import useCreateProduct from "../../hooks/useCreateProduct"

type ProductFormCreateProps = { onClose: () => void }

const ProductFormCreate = ({ onClose }: ProductFormCreateProps) => {
  const { form, isPending } = useCreateProduct({ onSuccess: onClose })

  return (
    <div className="rounded border bg-gray-50 p-4">
      <h3 className="mb-3 font-semibold">Create Product</h3>
      <form
        onSubmit={(e) => { e.preventDefault(); e.stopPropagation(); form.handleSubmit() }}
        className="grid gap-3"
      >
        <form.Field name="name">
          {(field) => (
            <div className="grid gap-1">
              <label className="text-sm font-medium">Name</label>
              <input
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                className="rounded border px-3 py-1.5 text-sm"
                placeholder="Product name"
              />
              {field.state.meta.errors && (
                <p className="text-xs text-red-600">{field.state.meta.errors.join(", ")}</p>
              )}
            </div>
          )}
        </form.Field>
        <form.Field name="description">
          {(field) => (
            <div className="grid gap-1">
              <label className="text-sm font-medium">Description</label>
              <input
                value={field.state.value ?? ""}
                onChange={(e) => field.handleChange(e.target.value)}
                className="rounded border px-3 py-1.5 text-sm"
                placeholder="Optional description"
              />
            </div>
          )}
        </form.Field>
        <form.Field name="price">
          {(field) => (
            <div className="grid gap-1">
              <label className="text-sm font-medium">Price</label>
              <input
                type="number"
                step="0.01"
                value={field.state.value}
                onChange={(e) => field.handleChange(parseFloat(e.target.value))}
                className="rounded border px-3 py-1.5 text-sm"
                placeholder="0.00"
              />
              {field.state.meta.errors && (
                <p className="text-xs text-red-600">{field.state.meta.errors.join(", ")}</p>
              )}
            </div>
          )}
        </form.Field>
        <form.Field name="stock">
          {(field) => (
            <div className="grid gap-1">
              <label className="text-sm font-medium">Stock</label>
              <input
                type="number"
                value={field.state.value}
                onChange={(e) => field.handleChange(parseInt(e.target.value, 10))}
                className="rounded border px-3 py-1.5 text-sm"
                placeholder="0"
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
            {isPending ? "Creating..." : "Create"}
          </button>
        </div>
      </form>
    </div>
  )
}

export default ProductFormCreate
