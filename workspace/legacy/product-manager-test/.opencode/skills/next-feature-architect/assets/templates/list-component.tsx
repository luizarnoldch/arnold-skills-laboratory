"use client"
import { useState } from "react"
import useSuspenseList[Entity]s from "../../hooks/useSuspenseList[Entity]s"
import useDelete[Entity] from "../../hooks/useDelete[Entity]"
import type { [Entity] } from "../../schemas/[entity-kebab].schema"
import [Entity]FormCreate from "../[Entity]FormCreate"
import [Entity]FormUpdate from "../[Entity]FormUpdate"

const [Entity]List = () => {
  const { [entity]s } = useSuspenseList[Entity]s()
  const [showCreate, setShowCreate] = useState(false)
  const [editing, setEditing] = useState<[Entity] | null>(null)
  const { mutate: delete[Entity], isPending: isDeleting } = useDelete[Entity]()

  return (
    <div className="grid gap-4">
      <div className="flex items-center justify-between">
        <span className="text-sm text-gray-500">{[entity]s.length} [entity]s</span>
        <button
          onClick={() => { setShowCreate(true); setEditing(null) }}
          className="rounded bg-blue-600 px-3 py-1.5 text-sm text-white hover:bg-blue-700"
        >
          Create [Entity]
        </button>
      </div>

      {showCreate && <[Entity]FormCreate onClose={() => setShowCreate(false)} />}
      {editing && <[Entity]FormUpdate [entity]={editing} onClose={() => setEditing(null)} />}

      {[entity]s.length === 0 ? (
        <p className="py-8 text-center text-sm text-gray-500">No [entity]s found.</p>
      ) : (
        <table className="w-full border-collapse text-sm">
          <thead>
            <tr className="border-b bg-gray-50">
              {/* Add <th> for each schema field */}
              <th className="px-3 py-2 text-left font-medium text-gray-600">ID</th>
              <th className="px-3 py-2 text-right font-medium text-gray-600">Actions</th>
            </tr>
          </thead>
          <tbody>
            {[entity]s.map(([entity]) => (
              <tr key={[entity].id} className="border-b hover:bg-gray-50">
                {/* Add <td> for each schema field */}
                <td className="px-3 py-2 font-mono text-xs text-gray-500">{[entity].id}</td>
                <td className="px-3 py-2 text-right">
                  <div className="flex justify-end gap-3">
                    <button
                      onClick={() => { setEditing([entity]); setShowCreate(false) }}
                      className="text-xs text-blue-600 underline hover:no-underline"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => delete[Entity]({ [entity]Id: [entity].id })}
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

export default [Entity]List
