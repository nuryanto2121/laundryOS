import { reactive, computed } from 'vue'
import { api } from '../api/client'

const STORAGE_KEY = 'laundryos_active_outlet_id'

const state = reactive({
  outlets: [],
  activeOutletId: localStorage.getItem(STORAGE_KEY) || '',
  loading: false,
})

export function useActiveOutlet() {
  const activeOutlet = computed(() => {
    if (!state.outlets.length) return null
    return (
      state.outlets.find((o) => o.id === state.activeOutletId) ||
      state.outlets.find((o) => o.is_active) ||
      state.outlets[0]
    )
  })

  const activeOutletName = computed(() => {
    return activeOutlet.value?.name || 'FAN Laundry Kalibaru'
  })

  const loadOutlets = async () => {
    state.loading = true
    try {
      const data = await api.get('/master/outlets')
      if (Array.isArray(data)) {
        state.outlets = data
        // If stored ID is invalid or empty, default to first active outlet
        if (!state.activeOutletId || !data.some((o) => o.id === state.activeOutletId)) {
          const first = data.find((o) => o.is_active) || data[0]
          if (first) {
            setActiveOutlet(first.id)
          }
        }
      }
    } catch (err) {
      console.error('Failed to load outlets in useActiveOutlet:', err)
    } finally {
      state.loading = false
    }
  }

  const setActiveOutlet = (id) => {
    state.activeOutletId = id
    try {
      localStorage.setItem(STORAGE_KEY, id)
    } catch (e) {
      console.error('Failed to save active outlet to localStorage:', e)
    }
  }

  return {
    state,
    outlets: computed(() => state.outlets),
    activeOutlet,
    activeOutletId: computed(() => state.activeOutletId),
    activeOutletName,
    loadOutlets,
    setActiveOutlet,
  }
}
