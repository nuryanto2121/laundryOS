import { reactive, watch } from 'vue'

const STORAGE_KEY = 'laundryos_app_settings'

const defaultSettings = {
  role: 'OWNER', // 'SUPER_ADMIN' | 'OWNER' | 'CASHIER'
  superAdminPin: '9999',
  ownerPin: '1234',
  hideFinancialsFromCashier: true,
  printer: {
    paperWidth: '58mm', // '58mm' | '80mm'
    storeName: 'FAN LAUNDRY',
    addressLine1: 'Jl. Kalibaru Barat I Rt 007 Rw 008',
    addressLine2: 'Kalibaru, Cilincing — Jakarta Utara 14110',
    phone: '0852-8597-7834',
    footerNote:
      'Terima kasih atas kepercayaan Anda!\nCucian tidak diambil > 30 hari di luar tanggung jawab kasir.',
    showQrCode: false,
    autoPrintOnComplete: false,
  },
  machine: {
    tubeCapacityKg: 8.0,
    toleranceKg: 0.5,
  },
  workflow: {
    allowSkipProcess: true,
  },
}

// Load from localStorage if available
function loadSavedSettings() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      return {
        ...defaultSettings,
        ...parsed,
        printer: { ...defaultSettings.printer, ...(parsed.printer || {}) },
        machine: { ...defaultSettings.machine, ...(parsed.machine || {}) },
        workflow: { ...defaultSettings.workflow, ...(parsed.workflow || {}) },
      }
    }
  } catch (e) {
    console.error('Failed to load settings from localStorage:', e)
  }
  return JSON.parse(JSON.stringify(defaultSettings))
}

const state = reactive(loadSavedSettings())

// Auto persist whenever state changes
watch(
  state,
  (newVal) => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(newVal))
    } catch (e) {
      console.error('Failed to save settings to localStorage:', e)
    }
  },
  { deep: true }
)

export function useAppSettings() {
  const switchRole = (newRole, pin = '') => {
    if (newRole === 'SUPER_ADMIN') {
      if (pin === (state.superAdminPin || '9999')) {
        state.role = 'SUPER_ADMIN'
        return { success: true }
      }
      return { success: false, message: 'PIN Super Admin tidak valid' }
    }
    if (newRole === 'OWNER') {
      if (pin === state.ownerPin || pin === (state.superAdminPin || '9999')) {
        state.role = 'OWNER'
        return { success: true }
      }
      return { success: false, message: 'PIN Owner tidak valid' }
    }
    state.role = 'CASHIER'
    return { success: true }
  }

  const changePin = (currentPin, newPin) => {
    if (currentPin !== state.ownerPin) {
      return { success: false, message: 'PIN saat ini salah' }
    }
    if (!newPin || newPin.length < 4) {
      return { success: false, message: 'PIN baru minimal 4 digit' }
    }
    state.ownerPin = newPin
    return { success: true, message: 'PIN Owner berhasil diperbarui' }
  }

  const resetToDefaults = () => {
    Object.assign(state, JSON.parse(JSON.stringify(defaultSettings)))
  }

  return {
    settings: state,
    switchRole,
    changePin,
    resetToDefaults,
  }
}
