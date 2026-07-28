# Mutation Hook Usage Patterns

Mutations expose `mutate`, `mutateAsync`, and `isPending`. Use `isPending` for button disabled state. Use `mutate()` for fire-and-forget (delete). Use `mutateAsync` when you need to await the result (create/update with form).

Never build complex form UI — only a button if a trigger is needed.

---

## Delete (button only)

```tsx
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

## Create (button + inline form)

```tsx
"use client"

import useCreate[Entity] from "../../hooks/useCreate[Entity]"

const Create[Entity]Form = () => {
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

export default Create[Entity]Form
```

## Update (button + form with initial data)

```tsx
"use client"

import useUpdate[Entity] from "../../hooks/useUpdate[Entity]"

type Update[Entity]FormProps = {
  [entity]Id: string
  initialName: string
}

const Update[Entity]Form = ({ [entity]Id, initialName }: Update[Entity]FormProps) => {
  const { form, isPending } = useUpdate[Entity]({
    id: [entity]Id,
    name: initialName,
  })

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
        {isPending ? "Saving..." : "Save"}
      </button>
    </form>
  )
}

export default Update[Entity]Form
```

---

## Rest API transport

Same component patterns. Only the hook internals differ (REST API client vs tRPC). Components are identical.

---

## Component Rules

- Always use `type` for Props (never `interface`).
- Always use arrow function `const X = () => {}`.
- Always `export default` the main component.
- Never add complex UI — single `<button>`, single `<input>`, no shadcn, no styling.
- For delete: only the button (no confirmation dialog unless explicit).
- For create/update: inline form with `<form.Field>` — no complex form layout.
