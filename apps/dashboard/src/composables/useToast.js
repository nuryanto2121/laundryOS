import { ref } from 'vue'

const toasts = ref([])
let toastCounter = 0

/**
 * Modern Toast Notification Composable for FAN Laundry Dashboard.
 */
export function useToast() {
  const addToast = (type, message, title = '', duration = 4000) => {
    const id = ++toastCounter
    const toast = {
      id,
      type,
      title: title || getDefaultTitle(type),
      message,
      duration,
    }

    toasts.value.push(toast)

    if (duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, duration)
    }

    return id
  }

  const removeToast = (id) => {
    const idx = toasts.value.findIndex((t) => t.id === id)
    if (idx !== -1) {
      toasts.value.splice(idx, 1)
    }
  }

  const getDefaultTitle = (type) => {
    switch (type) {
      case 'success':
        return 'Berhasil'
      case 'error':
        return 'Terjadi Kesalahan'
      case 'warning':
        return 'Perhatian'
      case 'info':
        return 'Informasi'
      default:
        return 'Notifikasi'
    }
  }

  return {
    toasts,
    removeToast,
    success: (message, title, duration) => addToast('success', message, title, duration),
    error: (message, title, duration) => addToast('error', message, title, duration),
    warning: (message, title, duration) => addToast('warning', message, title, duration),
    info: (message, title, duration) => addToast('info', message, title, duration),
  }
}
