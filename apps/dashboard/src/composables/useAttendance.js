import { ref } from 'vue'
import { api } from '../api/client'
import { useToast } from './useToast'

const activeAttendance = ref(null)
const activeStaffList = ref([])
const shifts = ref([])
const leaveRequests = ref([])
const attendanceHistory = ref([])
const loadingAttendance = ref(false)
const showAttendanceModal = ref(false)
const showLeaveModal = ref(false)

export function useAttendance() {
  const toast = useToast()

  const fetchShifts = async () => {
    try {
      const data = await api.get('/master/shifts')
      shifts.value = Array.isArray(data) ? data : []
      return shifts.value
    } catch (err) {
      console.error('Failed to load shifts', err)
      return []
    }
  }

  const fetchActiveAttendance = async (employeeName = 'Nuryanto') => {
    loadingAttendance.value = true
    try {
      const data = await api.get('/attendances/today', { employee: employeeName })
      activeAttendance.value = data || null
      return activeAttendance.value
    } catch (err) {
      console.error('Failed to get active attendance', err)
      activeAttendance.value = null
      return null
    } finally {
      loadingAttendance.value = false
    }
  }

  const fetchActiveStaffList = async () => {
    try {
      const data = await api.get('/attendances/active')
      activeStaffList.value = Array.isArray(data) ? data : []
      return activeStaffList.value
    } catch (err) {
      console.error('Failed to get active staff list', err)
      return []
    }
  }

  const checkIn = async (payload) => {
    try {
      const res = await api.post('/attendances/check-in', payload)
      activeAttendance.value = res
      toast.success(`Check-in berhasil! Selamat bertugas di ${res.shift_name} ✨`)
      await fetchActiveStaffList()
      return { success: true, data: res }
    } catch (err) {
      const msg = err.response?.data?.error || err.message || 'Check-in gagal'
      toast.error(msg)
      return { success: false, message: msg }
    }
  }

  const checkOut = async (payload) => {
    try {
      const res = await api.post('/attendances/check-out', payload)
      activeAttendance.value = null
      toast.success('Check-out & Serah terima shift berhasil dicatat! Terima kasih atas kerjanya 🙏')
      await fetchActiveStaffList()
      return { success: true, data: res }
    } catch (err) {
      const msg = err.response?.data?.error || err.message || 'Check-out gagal'
      toast.error(msg)
      return { success: false, message: msg }
    }
  }

  const getLeaveQuota = async (employeeName = 'Nuryanto') => {
    try {
      return await api.get('/leaves/quota', { employee: employeeName })
    } catch (err) {
      console.error('Failed to get leave quota', err)
      return { total_quota: 2, used_quota: 0, remaining_quota: 2, can_request: true }
    }
  }

  const createLeaveRequest = async (payload) => {
    try {
      const res = await api.post('/leaves/request', payload)
      toast.success('Pengajuan libur berhasil dikirim! Menunggu konfirmasi dari Owner 🙏')
      await fetchLeaveRequests()
      return { success: true, data: res }
    } catch (err) {
      const msg = err.response?.data?.error || err.message || 'Gagal mengajukan libur'
      toast.error(msg)
      return { success: false, message: msg }
    }
  }

  const fetchLeaveRequests = async (params = {}) => {
    try {
      const data = await api.get('/leaves', params)
      leaveRequests.value = Array.isArray(data) ? data : []
      return leaveRequests.value
    } catch (err) {
      console.error('Failed to load leave requests', err)
      return []
    }
  }

  const reviewLeaveRequest = async (id, status, reviewerNotes = '') => {
    try {
      const res = await api.put(`/leaves/${id}/review`, {
        status,
        reviewer_notes: reviewerNotes,
      })
      toast.success(
        status === 'APPROVED'
          ? 'Pengajuan libur disetujui ✅'
          : 'Pengajuan libur telah ditolak ❌'
      )
      await fetchLeaveRequests()
      return { success: true, data: res }
    } catch (err) {
      const msg = err.response?.data?.error || err.message || 'Gagal mereview pengajuan libur'
      toast.error(msg)
      return { success: false, message: msg }
    }
  }

  const fetchAttendanceHistory = async (params = {}) => {
    try {
      const data = await api.get('/attendances/history', params)
      attendanceHistory.value = Array.isArray(data) ? data : []
      return attendanceHistory.value
    } catch (err) {
      console.error('Failed to load attendance history', err)
      return []
    }
  }

  const updateShift = async (id, payload) => {
    try {
      const res = await api.put(`/master/shifts/${id}`, payload)
      toast.success('Pengaturan jam shift berhasil diperbarui ✨')
      await fetchShifts()
      return { success: true, data: res }
    } catch (err) {
      const msg = err.response?.data?.error || err.message || 'Gagal update jam shift'
      toast.error(msg)
      return { success: false, message: msg }
    }
  }

  return {
    activeAttendance,
    activeStaffList,
    shifts,
    leaveRequests,
    attendanceHistory,
    loadingAttendance,
    showAttendanceModal,
    showLeaveModal,
    fetchShifts,
    updateShift,
    fetchActiveAttendance,
    fetchActiveStaffList,
    checkIn,
    checkOut,
    getLeaveQuota,
    createLeaveRequest,
    fetchLeaveRequests,
    reviewLeaveRequest,
    fetchAttendanceHistory,
  }
}
