<template>
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
      @click.self="close"
    >
      <div
        class="bg-white rounded-3xl max-w-lg w-full p-5 sm:p-6 shadow-2xl border border-slate-100 space-y-4 max-h-[92vh] flex flex-col"
      >
      <!-- Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-100 shrink-0">
        <div class="flex items-center gap-2.5">
          <div class="w-10 h-10 rounded-2xl bg-indigo-50 text-indigo-600 flex items-center justify-center shadow-xs">
            <CalendarOff class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-black text-base text-slate-900">Pengajuan Libur Karyawan</h3>
            <p class="text-xs text-slate-500">Jatah libur 2x / bulan kalender</p>
          </div>
        </div>
        <button
          class="p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-colors cursor-pointer"
          @click="close"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Tab Switcher (Ajukan vs Riwayat) -->
      <div class="flex rounded-xl bg-slate-100 p-1 text-xs font-bold shrink-0">
        <button
          type="button"
          class="flex-1 py-1.5 rounded-lg transition-all text-center cursor-pointer"
          :class="activeTab === 'FORM' ? 'bg-white text-indigo-600 shadow-xs' : 'text-slate-600 hover:text-slate-900'"
          @click="activeTab = 'FORM'"
        >
          Form Pengajuan
        </button>
        <button
          type="button"
          class="flex-1 py-1.5 rounded-lg transition-all text-center cursor-pointer flex items-center justify-center gap-1.5"
          :class="activeTab === 'HISTORY' ? 'bg-white text-indigo-600 shadow-xs' : 'text-slate-600 hover:text-slate-900'"
          @click="loadHistoryTab"
        >
          <span>Riwayat Pengajuan</span>
          <span
            v-if="myLeaves.length > 0"
            class="px-1.5 py-0.2 rounded-full text-[10px] bg-indigo-100 text-indigo-700"
          >
            {{ myLeaves.length }}
          </span>
        </button>
      </div>

      <!-- Scrollable Body -->
      <div class="flex-1 overflow-y-auto space-y-4 pr-1">
        <!-- TAB 1: FORM -->
        <div v-if="activeTab === 'FORM'" class="space-y-4">
          <!-- Monthly Quota Card -->
          <div class="p-4 rounded-2xl bg-gradient-to-br from-indigo-50/90 to-blue-50/50 border border-indigo-100 space-y-2">
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-indigo-900 flex items-center gap-1.5">
                <ShieldAlert class="w-4 h-4 text-indigo-600" />
                Kuota Libur Bulan Ini ({{ currentMonthName }})
              </span>
              <span
                class="px-2.5 py-0.5 rounded-full text-xs font-black shadow-xs"
                :class="
                  quotaData.remaining_quota > 0
                    ? 'bg-emerald-500 text-white'
                    : 'bg-rose-500 text-white'
                "
              >
                {{ quotaData.remaining_quota }} / {{ quotaData.total_quota }} Hari Tersisa
              </span>
            </div>

            <!-- Progress Bar -->
            <div class="w-full bg-indigo-200/50 rounded-full h-2 overflow-hidden">
              <div
                class="bg-indigo-600 h-2 rounded-full transition-all duration-500"
                :style="{ width: `${(quotaData.used_quota / quotaData.total_quota) * 100}%` }"
              ></div>
            </div>

            <p class="text-[11px] text-indigo-700/80 leading-relaxed">
              Setiap karyawan FAN Laundry berhak mengajukan jatah libur maksimal <strong>2 hari per bulan</strong>. Pengajuan wajib disetujui oleh Owner terlebih dahulu.
            </p>
          </div>

          <!-- Form Fields -->
          <form @submit.prevent="handleSubmit" class="space-y-3.5">
            <!-- Employee Name -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Nama Karyawan</label>
              <input
                v-model="employeeNameInput"
                type="text"
                required
                placeholder="Nama Karyawan"
                class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm font-semibold text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
                @change="refreshQuota"
              />
            </div>

            <!-- Date Picker -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Tanggal Libur yang Diajukan</label>
              <input
                v-model="leaveDate"
                type="date"
                required
                :min="minDate"
                class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm font-semibold text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
              />
              <p class="text-[10px] text-slate-400 mt-1">Disarankan mengajukan minimal H-1 sebelum tanggal libur.</p>
            </div>

            <!-- Reason -->
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1">Alasan Pengajuan Libur</label>
              <textarea
                v-model="leaveReason"
                rows="3"
                required
                placeholder="Contoh: Acara keluarga di luar kota / Keperluan pribadi penting..."
                class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 text-sm font-normal text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 resize-none"
              ></textarea>
            </div>

            <!-- Notice if quota is empty -->
            <div
              v-if="quotaData.remaining_quota <= 0"
              class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-xs flex items-start gap-2"
            >
              <AlertCircle class="w-4 h-4 shrink-0 mt-0.5 text-rose-600" />
              <div>
                <strong>Kuota Libur Telah Habis!</strong> Anda sudah menggunakan 2 jatah libur untuk bulan ini. Pengajuan baru mungkin akan ditolak sistem.
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="pt-2 flex items-center justify-end gap-2.5">
              <button
                type="button"
                class="px-4 py-2.5 rounded-xl border border-slate-200 text-xs font-bold text-slate-600 hover:bg-slate-100 transition-colors cursor-pointer"
                @click="close"
              >
                Batal
              </button>
              <button
                type="submit"
                :disabled="submitting || quotaData.remaining_quota <= 0"
                class="px-5 py-2.5 rounded-xl bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold shadow-md shadow-indigo-600/20 disabled:opacity-50 disabled:cursor-not-allowed transition-all cursor-pointer flex items-center gap-2"
              >
                <Send class="w-3.5 h-3.5" />
                <span>{{ submitting ? 'Mengirim...' : 'Kirim Pengajuan' }}</span>
              </button>
            </div>
          </form>
        </div>

        <!-- TAB 2: RIWAYAT PENGAJUAN -->
        <div v-else class="space-y-3">
          <div v-if="loadingHistory" class="py-12 text-center text-slate-400 text-xs">
            Memuat riwayat pengajuan libur...
          </div>
          <div
            v-else-if="myLeaves.length === 0"
            class="py-12 text-center text-slate-400 space-y-2"
          >
            <Calendar class="w-8 h-8 mx-auto text-slate-300" />
            <p class="text-xs">Belum ada riwayat pengajuan libur untuk {{ employeeNameInput }}.</p>
          </div>
          <div v-else class="space-y-2.5">
            <div
              v-for="item in myLeaves"
              :key="item.id"
              class="p-3.5 rounded-2xl border border-slate-200/80 bg-slate-50/50 hover:bg-white hover:border-slate-300 transition-all space-y-2"
            >
              <div class="flex items-center justify-between">
                <span class="text-xs font-extrabold text-slate-800 flex items-center gap-1.5">
                  <Calendar class="w-3.5 h-3.5 text-indigo-600" />
                  {{ formatDate(item.date) }}
                </span>
                <span
                  class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                  :class="
                    item.status === 'APPROVED'
                      ? 'bg-emerald-100 text-emerald-800'
                      : item.status === 'REJECTED'
                        ? 'bg-rose-100 text-rose-800'
                        : 'bg-amber-100 text-amber-800 animate-pulse'
                  "
                >
                  {{ item.status === 'APPROVED' ? 'Disetujui ✅' : item.status === 'REJECTED' ? 'Ditolak ❌' : 'Menunggu Owner ⏳' }}
                </span>
              </div>

              <p class="text-xs text-slate-600 italic">"{{ item.reason }}"</p>

              <div
                v-if="item.reviewer_notes"
                class="pt-1.5 border-t border-slate-200/60 text-[11px] text-slate-500 flex items-center gap-1"
              >
                <span class="font-bold text-slate-700">Catatan Owner:</span>
                <span>{{ item.reviewer_notes }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import {
  CalendarOff,
  Calendar,
  X,
  ShieldAlert,
  Send,
  AlertCircle,
} from 'lucide-vue-next'
import { useAttendance } from '../../composables/useAttendance'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  defaultEmployee: {
    type: String,
    default: 'Nuryanto',
  },
})

const emit = defineEmits(['update:modelValue', 'request-submitted'])

const { getLeaveQuota, createLeaveRequest, fetchLeaveRequests } = useAttendance()

const activeTab = ref('FORM')
const employeeNameInput = ref(props.defaultEmployee)
const leaveDate = ref('')
const leaveReason = ref('')
const submitting = ref(false)
const loadingHistory = ref(false)
const myLeaves = ref([])

const quotaData = ref({
  total_quota: 2,
  used_quota: 0,
  remaining_quota: 2,
  can_request: true,
})

const minDate = computed(() => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  return tomorrow.toISOString().split('T')[0]
})

const currentMonthName = computed(() => {
  return new Intl.DateTimeFormat('id-ID', { month: 'long', year: 'numeric' }).format(new Date())
})

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }).format(d)
}

const refreshQuota = async () => {
  if (!employeeNameInput.value) return
  const q = await getLeaveQuota(employeeNameInput.value)
  if (q) {
    quotaData.value = q
  }
}

const loadHistoryTab = async () => {
  activeTab.value = 'HISTORY'
  loadingHistory.value = true
  try {
    const list = await fetchLeaveRequests({ employee: employeeNameInput.value })
    myLeaves.value = list
  } finally {
    loadingHistory.value = false
  }
}

const handleSubmit = async () => {
  if (!employeeNameInput.value || !leaveDate.value || !leaveReason.value) return

  submitting.value = true
  try {
    const res = await createLeaveRequest({
      employee_name: employeeNameInput.value,
      leave_date: leaveDate.value,
      date: leaveDate.value,
      reason: leaveReason.value,
    })
    if (res.success) {
      emit('request-submitted', res.data)
      leaveReason.value = ''
      leaveDate.value = ''
      await refreshQuota()
      await loadHistoryTab()
    }
  } finally {
    submitting.value = false
  }
}

const close = () => {
  emit('update:modelValue', false)
}

watch(
  () => props.modelValue,
  async (newVal) => {
    if (newVal) {
      activeTab.value = 'FORM'
      await refreshQuota()
    }
  }
)

onMounted(async () => {
  if (props.modelValue) {
    await refreshQuota()
  }
})
</script>
