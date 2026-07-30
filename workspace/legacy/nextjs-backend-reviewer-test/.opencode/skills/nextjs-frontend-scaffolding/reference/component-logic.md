# Component Logic Integration

The entrypoint component (`components/[Entity]Form/index.tsx` or any trigger component) is where business logic lives. It imports hooks and passes event handlers down as props. There are two patterns:

---

## Pattern 1: Form submission (create/update)

Hook provides `form` + `mutateAsync`. The component wires TanStack Form's `<form.Field>` and calls `form.handleSubmit()`.

```tsx
// src/features/[entity]/components/[Entity]Form/index.tsx
"use client"

import useCreate[Entity] from "../../hooks/useCreate[Entity]"

const [Entity]Form = () => {
  const { form, isPending } = useCreate[Entity]()

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    form.handleSubmit()
  }

  return (
    <form onSubmit={handleSubmit}>
      <form.Field name="name">
        {(field) => (
          <input
            name={field.name}
            value={field.state.value}
            onChange={(e) => field.handleChange(e.target.value)}
          />
        )}
      </form.Field>
      <button type="submit" disabled={isPending}>
        {isPending ? "Creating..." : "Create"}
      </button>
    </form>
  )
}

export default [Entity]Form
```

---

## Pattern 2: Button-only mutation (delete)

Hook provides `mutate` + `isPending`. No form. Just a button with `disabled={isPending}`.

```tsx
// src/features/[entity]/components/Delete[Entity]Button/index.tsx
"use client"

import useDelete[Entity] from "../../hooks/useDelete[Entity]"

type Delete[Entity]ButtonProps = {
  [entity]Id: string
}

const Delete[Entity]Button = ({ [entity]Id }: Delete[Entity]ButtonProps) => {
  const { mutate, isPending } = useDelete[Entity]({ [entity]Id })

  return (
    <button onClick={() => mutate()} disabled={isPending}>
      {isPending ? "Deleting..." : "Delete"}
    </button>
  )
}

export default Delete[Entity]Button
```

---

## Rules

- The entrypoint component is the **only** place hooks are called and wired.
- Sub-components (if any) receive callbacks and data as props — never import hooks directly.
- Always use `type` for Props, arrow function, `export default`.
- Never build complex UI. One `<button>`, one `<input>` per field. No shadcn.
