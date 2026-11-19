const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

/**
 * Universal fetch wrapper for LaundryOS API.
 */
async function request(endpoint, options = {}) {
  const url = new URL(endpoint.startsWith('http') ? endpoint : `${BASE_URL}${endpoint}`)
  if (options.params) {
    Object.entries(options.params).forEach(([key, val]) => {
      if (val !== undefined && val !== null && val !== '') {
        url.searchParams.append(key, val)
      }
    })
  }

  const defaultHeaders = {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  }

  const response = await fetch(url.toString(), {
    ...options,
    headers: {
      ...defaultHeaders,
      ...options.headers,
    },
  })

  if (!response.ok) {
    let errorMsg = `HTTP Error ${response.status}`
    try {
      const errorData = await response.json()
      errorMsg = errorData.error || errorData.message || errorMsg
    } catch {
      // ignore json parse error
    }
    throw new Error(errorMsg)
  }

  // If 204 No Content
  if (response.status === 204) {
    return null
  }

  return response.json()
}

export const api = {
  get: (endpoint, params) => request(endpoint, { method: 'GET', params }),
  post: (endpoint, body) => request(endpoint, { method: 'POST', body: JSON.stringify(body) }),
  patch: (endpoint, body) => request(endpoint, { method: 'PATCH', body: JSON.stringify(body) }),
  delete: (endpoint) => request(endpoint, { method: 'DELETE' }),
}

/**
 * Formats a number to Indonesian Rupiah currency string (e.g. 25000 -> "25.000").
 */
export function formatRupiah(val) {
  if (val === undefined || val === null || isNaN(val)) return '0'
  return Math.round(Number(val)).toLocaleString('id-ID')
}

/**
 * Formats an ISO date string to Indonesian human date (e.g. "19 Nov 2026").
 */
export function formatDate(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

/**
 * Formats an ISO date string to Indonesian date + time (e.g. "19 Nov 2026, 14:30").
 */
export function formatDateTime(dateStr) {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
