<template>
  <header
    class="sticky top-0 z-30 bg-white/90 backdrop-blur-md border-b border-slate-100 shadow-xs"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
      <!-- Left: Brand Logo & Title -->
      <div class="flex items-center gap-3">
        <div
          class="w-10 h-10 rounded-xl bg-gradient-to-tr from-brand-600 to-sky-400 flex items-center justify-center text-white shadow-md shadow-brand-500/20"
        >
          <Sparkles class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <span class="font-extrabold text-lg tracking-tight text-slate-900">FAN Laundry</span>
            <span
              class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-brand-50 text-brand-600 border border-brand-200 uppercase tracking-wider"
            >
              POS Station
            </span>
          </div>
          <p class="text-xs text-slate-500 hidden sm:block">
            Kalibaru Barat I, Cilincing — Jakarta Utara
          </p>
        </div>
      </div>

      <!-- Right: Outlet Switcher, Live Status & Role Profile -->
      <div class="flex items-center gap-1.5 sm:gap-3 shrink-0">
        <!-- Outlet Switcher Trigger -->
        <button
          class="flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200/80 transition-all text-xs font-semibold text-slate-700 max-w-[110px] sm:max-w-none cursor-pointer"
          @click="$emit('open-outlet-modal')"
        >
          <Store class="w-3.5 h-3.5 text-brand-600 shrink-0" />
          <span class="truncate">{{ activeOutletName }}</span>
          <ChevronDown class="w-3 h-3 text-slate-400 shrink-0" />
        </button>

        <!-- Live Shift & Attendance Trigger Pill -->
        <button
          v-if="activeAttendance"
          class="flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-xl bg-emerald-50 hover:bg-emerald-100 border border-emerald-200 text-emerald-800 text-xs font-bold transition-all cursor-pointer shadow-xs"
          title="Klik untuk Tutup Shift & Serah Terima Kasir"
          @click="openCheckOutModal"
        >
          <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          <span class="truncate max-w-[110px] sm:max-w-none">{{ activeAttendance.shift_name }} (Aktif)</span>
          <LogOut class="w-3 h-3 text-emerald-600 shrink-0 ml-0.5" />
        </button>
        <button
          v-else
          class="flex items-center gap-1.5 px-2.5 sm:px-3 py-1.5 rounded-xl bg-amber-50 hover:bg-amber-100 border border-amber-200 text-amber-900 text-xs font-bold transition-all cursor-pointer shadow-xs animate-bounce"
          title="Klik untuk Check-In Foto & Mulai Shift"
          @click="openCheckInModal"
        >
          <Camera class="w-3.5 h-3.5 text-amber-600 shrink-0" />
          <span class="truncate max-w-[100px] sm:max-w-none">Buka Shift (Absen)</span>
        </button>

        <!-- Leave Request (Ajukan Libur) Trigger -->
        <button
          class="hidden sm:flex items-center gap-1 px-2.5 py-1.5 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 text-slate-700 text-xs font-semibold transition-all cursor-pointer"
          title="Ajukan Jatah Libur Karyawan (Maks 2x/Bulan)"
          @click="showLeaveModal = true"
        >
          <CalendarOff class="w-3.5 h-3.5 text-indigo-600 shrink-0" />
          <span class="hidden lg:inline">Ajukan Libur</span>
        </button>

        <!-- Live Sync Pill -->
        <div
          class="hidden md:flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-slate-50 text-slate-600 border border-slate-200 text-xs font-medium"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
          <span>Online</span>
        </div>

        <!-- Mobile App Prototype Link -->
        <router-link
          to="/mobile"
          class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl bg-purple-50 hover:bg-purple-100 text-purple-700 border border-purple-200 text-xs font-bold transition-all"
          title="Buka Desain Prototype Mobile Native LMS"
        >
          <Smartphone class="w-3.5 h-3.5 text-purple-600 shrink-0" />
          <span class="hidden lg:inline">Mobile LMS</span>
        </router-link>

        <!-- Role Profile Pill & Switcher Trigger -->
        <button
          class="flex items-center gap-2 pl-2 border-l border-slate-200 hover:opacity-85 transition-opacity cursor-pointer text-left"
          @click="openRoleModal"
        >
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold shadow-xs text-white transition-all"
            :class="
              settings.role === 'SUPER_ADMIN'
                ? 'bg-gradient-to-tr from-amber-500 to-rose-600 ring-2 ring-amber-300'
                : settings.role === 'OWNER'
                  ? 'bg-gradient-to-tr from-purple-600 to-indigo-600 ring-2 ring-purple-300'
                  : 'bg-slate-800'
            "
          >
            {{ settings.role === 'SUPER_ADMIN' ? '⚡' : settings.role === 'OWNER' ? 'OW' : 'K1' }}
          </div>
          <div class="hidden md:block">
            <div class="flex items-center gap-1">
              <p class="text-xs font-bold text-slate-800 leading-none">
                {{
                  auth.currentUser.value?.full_name ||
                  (settings.role === 'SUPER_ADMIN'
                    ? 'Super Admin'
                    : settings.role === 'OWNER'
                      ? 'Owner (Admin)'
                      : 'Kasir Utama')
                }}
              </p>
              <Crown v-if="settings.role === 'SUPER_ADMIN'" class="w-3.5 h-3.5 text-amber-500 inline" />
              <Shield v-else-if="settings.role === 'OWNER'" class="w-3 h-3 text-purple-600 inline" />
            </div>
            <p class="text-[10px] text-slate-400 leading-none mt-0.5">
              {{
                settings.role === 'SUPER_ADMIN'
                  ? 'Super Admin (All Tenant)'
                  : settings.role === 'OWNER'
                    ? 'Akses Penuh Owner'
                    : 'Karyawan / Kasir'
              }}
            </p>
          </div>
        </button>

        <!-- Logout Button -->
        <button
          type="button"
          class="flex items-center gap-1 px-2.5 py-1.5 rounded-xl bg-slate-50 hover:bg-rose-50 hover:text-rose-600 text-slate-500 border border-slate-200 text-xs font-bold transition-all cursor-pointer"
          title="Keluar dari Akun (Logout)"
          @click="handleLogout"
        >
          <LogOut class="w-3.5 h-3.5" />
          <span class="hidden xl:inline">Keluar</span>
        </button>
      </div>
    </div>

    <!-- Role Switcher & PIN Auth Modal (Teleported to body to prevent backdrop-blur clipping) -->
    <Teleport to="body">
      <div
        v-if="showRoleModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
        @click.self="closeRoleModal"
      >
        <div class="bg-white rounded-3xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4">
          <div class="flex items-center justify-between pb-3 border-b border-slate-100">
            <div class="flex items-center gap-2">
              <ShieldCheck class="w-5 h-5 text-purple-600" />
              <h3 class="font-black text-sm text-slate-900">Ganti Hak Akses (RBAC)</h3>
            </div>
            <button
              type="button"
              class="p-1.5 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 cursor-pointer transition-colors"
              @click="closeRoleModal"
            >
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="space-y-3.5">
            <!-- Peran Saat Ini -->
            <div class="p-3 rounded-2xl bg-slate-50 border border-slate-200/80 flex items-center justify-between">
              <div>
                <span class="text-[10px] text-slate-400 uppercase font-black block">Peran Aktif Saat Ini</span>
                <span class="font-extrabold text-xs text-slate-900 flex items-center gap-1.5 mt-0.5">
                  <span v-if="settings.role === 'SUPER_ADMIN'">⚡ Super Admin (Semua Outlet)</span>
                  <span v-else-if="settings.role === 'OWNER'">👑 Owner (Owner Toko)</span>
                  <span v-else>👤 Kasir (Mode Kasir)</span>
                </span>
              </div>
              <span
                class="text-[10px] font-black px-2 py-0.5 rounded-full"
                :class="
                  settings.role === 'SUPER_ADMIN'
                    ? 'bg-amber-100 text-amber-800'
                    : settings.role === 'OWNER'
                      ? 'bg-purple-100 text-purple-800'
                      : 'bg-slate-200 text-slate-700'
                "
              >
                {{ settings.role }}
              </span>
            </div>

            <!-- Role Selector Targets -->
            <div class="grid grid-cols-3 gap-1.5 p-1 bg-slate-100 rounded-xl text-[11px] font-bold text-center">
              <button
                type="button"
                class="py-1.5 rounded-lg transition-all cursor-pointer"
                :class="targetRole === 'CASHIER' ? 'bg-white text-slate-900 shadow-xs' : 'text-slate-500 hover:text-slate-900'"
                @click="targetRole = 'CASHIER'; pinError = ''"
              >
                👤 Kasir
              </button>
              <button
                type="button"
                class="py-1.5 rounded-lg transition-all cursor-pointer"
                :class="targetRole === 'OWNER' ? 'bg-white text-purple-700 shadow-xs' : 'text-slate-500 hover:text-slate-900'"
                @click="targetRole = 'OWNER'; pinError = ''"
              >
                👑 Owner
              </button>
              <button
                type="button"
                class="py-1.5 rounded-lg transition-all cursor-pointer"
                :class="targetRole === 'SUPER_ADMIN' ? 'bg-white text-amber-700 shadow-xs' : 'text-slate-500 hover:text-slate-900'"
                @click="targetRole = 'SUPER_ADMIN'; pinError = ''"
              >
                ⚡ Super
              </button>
            </div>

            <!-- Mode Kasir Selection -->
            <div v-if="targetRole === 'CASHIER'" class="space-y-3 pt-1">
              <p class="text-xs text-slate-600 leading-relaxed">
                Beralih ke <strong>Mode Kasir</strong> untuk membatasi akses omset, laporan keuangan, dan menu master data saat operasional kasir.
              </p>
              <div class="flex items-center gap-2">
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-xl border border-slate-200 text-slate-600 font-bold text-xs hover:bg-slate-50 cursor-pointer"
                  @click="closeRoleModal"
                >
                  Batal
                </button>
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-xl bg-slate-900 text-white font-bold text-xs hover:bg-slate-800 transition-all cursor-pointer"
                  @click="handleSwitchToCashier"
                >
                  Aktifkan Kasir
                </button>
              </div>
            </div>

            <!-- Mode Owner Selection -->
            <div v-else-if="targetRole === 'OWNER'" class="space-y-3 pt-1">
              <p class="text-xs text-slate-600 leading-relaxed">
                Masukkan <strong>PIN Owner</strong> untuk akses penuh laporan keuangan, omset, dan pengaturan outlet.
              </p>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">PIN Owner (Default: 1234)</label>
                <input
                  v-model="inputPin"
                  type="password"
                  maxlength="8"
                  placeholder="Masukkan PIN Owner..."
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-center tracking-widest text-base font-bold focus:ring-2 focus:ring-purple-500 focus:outline-hidden"
                  @keyup.enter="handleSwitchToTargetRole"
                />
                <p v-if="pinError" class="text-xs font-semibold text-rose-600 mt-1">
                  {{ pinError }}
                </p>
              </div>
              <div class="flex items-center gap-2">
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-xl border border-slate-200 text-slate-600 font-bold text-xs hover:bg-slate-50 cursor-pointer"
                  @click="closeRoleModal"
                >
                  Batal
                </button>
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-xl bg-purple-600 text-white font-bold text-xs hover:bg-purple-700 transition-all cursor-pointer flex items-center justify-center gap-1.5 shadow-xs"
                  @click="handleSwitchToTargetRole"
                >
                  <Lock class="w-3.5 h-3.5" />
                  <span>Buka Owner</span>
                </button>
              </div>
            </div>

            <!-- Mode Super Admin Selection -->
            <div v-else class="space-y-3 pt-1">
              <p class="text-xs text-slate-600 leading-relaxed">
                Masukkan <strong>PIN Super Admin</strong> untuk mengelola seluruh data outlet, seluruh owner, dan konfigurasi SaaS sistem.
              </p>
              <div>
                <label class="block text-xs font-bold text-slate-700 mb-1">PIN Super Admin (Default: 9999)</label>
                <input
                  v-model="inputPin"
                  type="password"
                  maxlength="8"
                  placeholder="Masukkan PIN Super Admin..."
                  class="w-full px-3 py-2 rounded-xl border border-slate-200 text-center tracking-widest text-base font-bold focus:ring-2 focus:ring-amber-500 focus:outline-hidden"
                  @keyup.enter="handleSwitchToTargetRole"
                />
                <p v-if="pinError" class="text-xs font-semibold text-rose-600 mt-1">
                  {{ pinError }}
                </p>
              </div>
              <div class="flex items-center gap-2">
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-xl border border-slate-200 text-slate-600 font-bold text-xs hover:bg-slate-50 cursor-pointer"
                  @click="closeRoleModal"
                >
                  Batal
                </button>
                <button
                  type="button"
                  class="flex-1 py-2.5 rounded-xl bg-gradient-to-r from-amber-600 to-rose-600 text-white font-bold text-xs hover:opacity-90 transition-all cursor-pointer flex items-center justify-center gap-1.5 shadow-xs"
                  @click="handleSwitchToTargetRole"
                >
                  <Crown class="w-3.5 h-3.5" />
                  <span>Buka Super Admin</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Attendance Modal (Check-In & Check-Out with Geofencing & Photo) -->
    <AttendanceModal
      v-model="showAttendanceModal"
      :mode="attendanceModalMode"
      :active-attendance="activeAttendance"
      :employee-name="currentEmployeeName"
      @attendance-updated="handleAttendanceUpdated"
    />

    <!-- Leave Request Modal -->
    <LeaveRequestModal
      v-model="showLeaveModal"
      :default-employee="currentEmployeeName"
    />
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  Sparkles,
  Store,
  ChevronDown,
  Shield,
  ShieldCheck,
  Crown,
  User,
  Lock,
  X,
  Smartphone,
  Camera,
  LogOut,
  CalendarOff,
} from 'lucide-vue-next'
import { useAppSettings } from '../../composables/useAppSettings'
import { useToast } from '../../composables/useToast'
import { useAttendance } from '../../composables/useAttendance'
import { useAuth } from '../../composables/useAuth'
import { useRouter } from 'vue-router'
import AttendanceModal from '../attendance/AttendanceModal.vue'
import LeaveRequestModal from '../attendance/LeaveRequestModal.vue'

defineProps({
  activeOutletName: {
    type: String,
    default: 'FAN Laundry Kalibaru',
  },
})

defineEmits(['open-outlet-modal'])

const router = useRouter()
const auth = useAuth()
const { settings, switchRole } = useAppSettings()
const { success, error } = useToast()

const handleLogout = () => {
  if (confirm('Keluar dari LaundryOS?')) {
    auth.logout()
    router.push('/login')
  }
}
const {
  activeAttendance,
  showAttendanceModal,
  showLeaveModal,
  fetchActiveAttendance,
} = useAttendance()

const showRoleModal = ref(false)
const targetRole = ref('OWNER')
const inputPin = ref('')
const pinError = ref('')

const attendanceModalMode = ref('CHECK_IN')
const currentEmployeeName = ref('Nuryanto')

const openCheckInModal = () => {
  attendanceModalMode.value = 'CHECK_IN'
  showAttendanceModal.value = true
}

const openCheckOutModal = () => {
  attendanceModalMode.value = 'CHECK_OUT'
  showAttendanceModal.value = true
}

const handleAttendanceUpdated = async (attendance) => {
  await fetchActiveAttendance(currentEmployeeName.value)
}

const openRoleModal = () => {
  inputPin.value = ''
  pinError.value = ''
  targetRole.value = settings.role === 'CASHIER' ? 'OWNER' : 'CASHIER'
  showRoleModal.value = true
}

const closeRoleModal = () => {
  showRoleModal.value = false
  inputPin.value = ''
  pinError.value = ''
}

const handleSwitchToCashier = () => {
  switchRole('CASHIER')
  success('Beralih ke Mode Kasir. Omset & Laporan Keuangan kini disembunyikan.')
  closeRoleModal()
}

const handleSwitchToTargetRole = () => {
  if (targetRole.value === 'CASHIER') {
    handleSwitchToCashier()
    return
  }
  const res = switchRole(targetRole.value, inputPin.value)
  if (res.success) {
    success(
      targetRole.value === 'SUPER_ADMIN'
        ? 'Akses Super Admin aktif! Anda dapat mengelola seluruh data owner dan cabang.'
        : 'Akses Owner berhasil diaktifkan!'
    )
    closeRoleModal()
  } else {
    pinError.value = res.message || 'PIN salah'
    error('Autentikasi gagal: ' + (res.message || 'PIN salah'))
  }
}

const handleSwitchToOwner = () => {
  targetRole.value = 'OWNER'
  handleSwitchToTargetRole()
}

onMounted(async () => {
  await fetchActiveAttendance(currentEmployeeName.value)
})
</script>
