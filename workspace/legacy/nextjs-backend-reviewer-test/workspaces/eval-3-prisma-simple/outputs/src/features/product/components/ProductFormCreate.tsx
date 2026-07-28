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
            <div>
              <label htmlFor={field.name} className="block text-sm font-medium text-gray-700">
                Name
              </label>
              <input
                id={field.name}
                name={field.name}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                className="mt-1 block w-full rounded border px-3 py-1.5 text-sm"
              />
              {field.state.meta.errors && (
                <p className="mt-1 text-xs text-red-600">{field.state.meta.errors.join(", ")}</p>
              )}
            </div>
          )}
        </form.Field>

        <form.Field name="description">
          {(field) => (
            <div>
              <label htmlFor={field.name} className="block text-sm font-medium text-gray-700">
                Description
              </label>
              <input
                id={field.name}
                name={field.name}
                value={field.state.value ?? ""}
                onChange={(e) => field.handleChange(e.target.value || null)}
                className="mt-1 block w-full rounded border px-3 py-1.5 text-sm"
              />
              {field.state.meta.errors && (
                <p className="mt-1 text-xs text-red-600">{field.state.meta.errors.join(", ")}</p>
              )}
            </div>
          )}
        </form.Field>

        <form.Field name="price">
          {(field) => (
            <div>
              <label htmlFor={field.name} className="block text-sm font-medium text-gray-700">
                Price
              </label>
              <input
                id={field.name}
                name={field.name}
                type="number"
                step="0.01"
                value={field.state.value}
                onChange={(e) => field.handleChange(Number(e.target.value))}
                className="mt-1 block w-full rounded border px-3 py-1.5 text-sm"
              />
              {field.state.meta.errors && (
                <p className="mt-1 text-xs text-red-600">{field.state.meta.errors.join(", ")}</p>
              )}
            </div>
          )}
        </form.Field>

        <form.Field name="stock">
          {(field) => (
            <div>
              <label htmlFor={field.name} className="block text-sm font-medium text-gray-700">
                Stock
              </label>
              <input
                id={field.name}
                name={field.name}
                type="number"
                step="1"
                value={field.state.value}
                onChange={(e) => field.handleChange(Number(e.target.value))}
                className="mt-1 block w-full rounded border px-3 py-1.5 text-sm"
              />
              {field.state.meta.errors && (
                <p className="mt-1 text-xs text-red-600">{field.state.meta.errors.join(", ")}</p>
              )}
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
