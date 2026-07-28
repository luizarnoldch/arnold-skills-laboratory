## Setup: Add apiPrefetch

Add to `src/trpc/server.tsx`.

```ts

// rest of the code...

export function apiPrefetch(
  options: Parameters<ReturnType<typeof getQueryClient>["prefetchQuery"]>[0]
) {
  const queryClient = getQueryClient()
  void queryClient.prefetchQuery(options)
}
```


Shared `apiFetch` in `src/lib/api.ts`:

```ts
export async function apiFetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${process.env.API_BASE_URL}${endpoint}`
  const token = getSessionToken()
  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  })
  if (!response.ok) {
    const error = await response.json().catch(() => ({}))
    throw new Error(error.message || `API Error: ${response.statusText}`)
  }
  if (response.status === 204) return {} as T
  const text = await response.text()
  return text ? (JSON.parse(text) as T) : ({} as T)
}
