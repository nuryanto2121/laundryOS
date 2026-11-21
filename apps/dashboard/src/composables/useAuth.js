import { reactive, computed } from 'vue'

const TOKEN_KEY = 'laundryos_auth_token'
const USER_KEY = 'laundryos_auth_user'

function loadSavedSession() {
  try {
    const rawUser = localStorage.getItem(USER_KEY)
    const token = localStorage.getItem(TOKEN_KEY)
    if (rawUser && token) {
      return {
        token,
        user: JSON.parse(rawUser),
      }
    }
  } catch (e) {
    console.error('Failed to load auth session:', e)
  }
  return {
    token: null,
    user: null,
  }
}

const state = reactive(loadSavedSession())

export function useAuth() {
  const isAuthenticated = computed(() => Boolean(state.token && state.user))
  const currentUser = computed(() => state.user)
  const userRole = computed(() => state.user?.role || 'KASIR')
  const isSuperAdmin = computed(() => userRole.value === 'SUPER_ADMIN')
  const isOwner = computed(() => userRole.value === 'OWNER' || isSuperAdmin.value)
  const isCashier = computed(() => userRole.value === 'KASIR')

  const can = (permissionKey) => {
    if (!state.user) return false
    // Super Admin and Owner have all permissions
    if (state.user.role === 'SUPER_ADMIN' || state.user.role === 'OWNER') {
      return true
    }
    if (state.user.permissions?.all === true) {
      return true
    }
    return Boolean(state.user.permissions?.[permissionKey])
  }

  const login = async (username, password) => {
    try {
      const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      })

      // If gateway proxies /auth/login directly without prefix fallback
      let data
      if (!res.ok && res.status === 404) {
        const fallbackRes = await fetch('/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, password }),
        })
        data = await fallbackRes.json()
        if (!fallbackRes.ok) {
          return { success: false, message: data.error || data.message || 'Login gagal' }
        }
      } else {
        data = await res.json()
        if (!res.ok) {
          return { success: false, message: data.error || data.message || 'Login gagal' }
        }
      }

      state.token = data.token
      state.user = data.user
      localStorage.setItem(TOKEN_KEY, data.token)
      localStorage.setItem(USER_KEY, JSON.stringify(data.user))

      return { success: true, user: data.user }
    } catch (err) {
      console.error('Login error:', err)
      return { success: false, message: 'Gagal menghubungi server. Periksa koneksi Anda.' }
    }
  }

  const logout = () => {
    state.token = null
    state.user = null
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  const setDemoUser = (role, outletId = 'a0000000-0000-0000-0000-000000000001') => {
    const demoProfiles = {
      SUPER_ADMIN: {
        id: 'f0000000-0000-0000-0000-000000000000',
        username: 'superadmin',
        full_name: 'Super Administrator',
        role: 'SUPER_ADMIN',
        permissions: { all: true },
      },
      OWNER: {
        id: 'f0000000-0000-0000-0000-000000000001',
        username: 'owner',
        full_name: 'Owner Utama (FAN Laundry)',
        role: 'OWNER',
        outlet_id: outletId,
        outlet_name: 'FAN Laundry Kalibaru (Pusat)',
        permissions: { all: true },
      },
      KASIR: {
        id: 'f0000000-0000-0000-0000-000000000002',
        username: 'kasir1',
        full_name: 'Kasir Utama Kalibaru',
        role: 'KASIR',
        outlet_id: outletId,
        outlet_name: 'FAN Laundry Kalibaru (Pusat)',
        permissions: {
          can_create_order: true,
          can_backdate_order: false,
          can_edit_order: false,
          can_cancel_order: false,
          can_create_expense: true,
          can_process_wash: true,
          can_process_dry: true,
          can_process_iron: true,
          can_process_pack: true,
          can_manage_services: false,
          can_manage_customers: true,
          can_manage_employees: false,
          can_view_financials: false,
          can_view_order_reports: false,
          can_view_financial_reports: false,
          can_view_performance_reports: false,
          can_view_customer_reports: false,
        },
      },
    }

    const selected = demoProfiles[role] || demoProfiles.KASIR
    const token = 'demo_token_' + selected.role
    state.token = token
    state.user = selected
    localStorage.setItem(TOKEN_KEY, token)
    localStorage.setItem(USER_KEY, JSON.stringify(selected))
    return selected
  }

  return {
    state,
    currentUser,
    userRole,
    isAuthenticated,
    isSuperAdmin,
    isOwner,
    isCashier,
    can,
    login,
    logout,
    setDemoUser,
  }
}
