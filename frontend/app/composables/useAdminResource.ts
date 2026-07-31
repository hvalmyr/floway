import type { Ref } from 'vue'

/**
 * Generic CRUD data layer for the flat admin entities that follow the
 * List/Create/Update(PUT)/Delete shape (teachers, blog-posts, masterclasses,
 * courses, faq). Nested resources (course_blocks, lessons) and leads (which
 * only supports a status PATCH, not a full update) manage their own calls
 * via useApi() directly instead of this composable.
 */
export function useAdminResource<T extends { id: number }>(basePath: string) {
  const api = useApi()
  const items = ref<T[]>([]) as Ref<T[]>
  const loading = ref(false)
  const error = ref('')

  async function fetchAll() {
    loading.value = true
    error.value = ''
    try {
      items.value = (await api<T[]>(basePath)) ?? []
    } catch {
      error.value = 'Не удалось загрузить данные'
    } finally {
      loading.value = false
    }
  }

  async function create(payload: Partial<T>) {
    const created = await api<T>(basePath, { method: 'POST', body: payload })
    items.value = [...items.value, created]
    return created
  }

  async function update(id: number, payload: Partial<T>) {
    const updated = await api<T>(`${basePath}/${id}`, { method: 'PUT', body: payload })
    items.value = items.value.map((item) => (item.id === id ? updated : item))
    return updated
  }

  async function remove(id: number) {
    await api(`${basePath}/${id}`, { method: 'DELETE' })
    items.value = items.value.filter((item) => item.id !== id)
  }

  return { items, loading, error, fetchAll, create, update, remove }
}
