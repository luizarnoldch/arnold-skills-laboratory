import { Suspense } from "react"
import { ErrorBoundary } from "react-error-boundary"
import [Entity]List from "../components/[Entity]List"

const [Entity]View = () => (
  <div className="grid min-h-dvh grid-rows-[auto_1fr] gap-4 p-4 md:p-6">
    <header className="border-b pb-4">
      <h1 className="text-2xl font-semibold">[Entity]s</h1>
    </header>
    <main>
      <ErrorBoundary fallback={<p className="text-red-600">Failed to load [entity]s. Refresh to try again.</p>}>
        <Suspense fallback={<p className="text-gray-500">Loading [entity]s...</p>}>
          <[Entity]List />
        </Suspense>
      </ErrorBoundary>
    </main>
  </div>
)

export default [Entity]View
