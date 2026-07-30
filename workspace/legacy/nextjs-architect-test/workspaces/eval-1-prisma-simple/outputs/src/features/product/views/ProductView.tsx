import { Suspense } from "react"
import { ErrorBoundary } from "react-error-boundary"
import ProductList from "../components/ProductList"

const ProductView = () => (
  <div className="grid min-h-dvh grid-rows-[auto_1fr] gap-4 p-4 md:p-6">
    <header className="border-b pb-4">
      <h1 className="text-2xl font-semibold">Products</h1>
    </header>
    <main>
      <ErrorBoundary fallback={<p className="text-red-600">Failed to load products. Refresh to try again.</p>}>
        <Suspense fallback={<p className="text-gray-500">Loading products...</p>}>
          <ProductList />
        </Suspense>
      </ErrorBoundary>
    </main>
  </div>
)

export default ProductView
