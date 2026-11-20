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

        <!-- Live Sync Pill -->
        <div
          class="hidden sm:flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 text-xs font-medium"
        >
          <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
          <span>Online</span>
        </div>

        <!-- Role Profile Pill & Switcher Trigger -->
        <button
          class="flex items-center gap-2 pl-2 border-l border-slate-200 hover:opacity-85 transition-opacity cursor-pointer text-left"
          @click="openRoleModal"
        >
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold shadow-xs text-white transition-all"
            :class="
              settings.role === 'OWNER'
                ? 'bg-gradient-to-tr from-purple-600 to-indigo-600 ring-2 ring-purple-300'
                : 'bg-slate-800'
            "
          >
            {{ settings.role === 'OWNER' ? 'OW' : 'K1' }}
          </div>
          <div class="hidden md:block">
            <div class="flex items-center gap-1">
              <p class="text-xs font-bold text-slate-800 leading-none">
                {{ settings.role === 'OWNER' ? 'Owner (Admin)' : 'Kasir Utama' }}
              </p>
              <Shield v-if="settings.role === 'OWNER'" class="w-3 h-3 text-purple-600 inline" />
            </div>
            <p class="text-[10px] text-slate-400 leading-none mt-0.5">
              {{ settings.role === 'OWNER' ? 'Akses Penuh' : 'Mode Operasional' }}
            </p>
          </div>
        </button>
      </div>
    </div>

    <!-- Role Switcher & PIN Auth Modal -->
    <div
      v-if="showRoleModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
    >
      <div class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100">
        <div class="flex items-center justify-between pb-3 border-b border-slate-100">
          <div class="flex items-center gap-2">
            <ShieldCheck class="w-5 h-5 text-purple-600" />
            <h3 class="font-bold text-sm text-slate-900">Ganti Hak Akses (RBAC)</h3>
          </div>
          <button
            class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer"
            @click="closeRoleModal"
          >
            <X class="w-4 h-4" />
          </button>
        </div>

        <div class="py-4 space-y-3">
          <div
            class="p-3 rounded-xl bg-slate-50 border border-slate-200/80 flex items-center justify-between"
          >
            <div>
              <span class="text-[11px] text-slate-400 uppercase font-bold block"
                >Peran Aktif Saat Ini</span
              >
              <span class="font-black text-sm text-slate-900">
                {{ settings.role === 'OWNER' ? '👑 Owner (Akses Penuh)' : '👤 Kasir (Mode Kasir)' }}
              </span>
            </div>
            <span
              class="text-xs font-bold px-2 py-0.5 rounded-full"
              :class="
                settings.role === 'OWNER'
                  ? 'bg-purple-100 text-purple-800'
                  : 'bg-slate-200 text-slate-700'
              "
            >
              {{ settings.role }}
            </span>
          </div>

          <!-- Switch to Cashier Option -->
          <div v-if="settings.role === 'OWNER'" class="space-y-2">
            <p class="text-xs text-slate-600">
              Beralih ke <strong>Mode Kasir</strong> untuk membatasi tampilan omset dan laporan
              keuangan saat operasional toko.
            </p>
            <button
              type="button"
              class="w-full py-2.5 px-3 rounded-xl bg-slate-800 text-white font-bold text-xs hover:bg-slate-900 transition-all cursor-pointer flex items-center justify-center gap-2"
              @click="handleSwitchToCashier"
            >
              <User class="w-4 h-4" />
              <span>Beralih ke Mode Kasir</span>
            </button>
          </div>

          <!-- Switch to Owner Option (Requires PIN) -->
          <div v-else class="space-y-3">
            <p class="text-xs text-slate-600">
              Masukkan <strong>PIN Owner</strong> untuk membuka akses penuh ke data omset dan
              pembukuan kas.
            </p>
            <div>
              <label class="block text-xs font-bold text-slate-700 mb-1"
                >PIN Owner (Default: 1234)</label
              >
              <input
                v-model="inputPin"
                type="password"
                maxlength="8"
                placeholder="Masukkan 4 digit PIN..."
                class="w-full px-3 py-2 rounded-xl border border-slate-200 text-center tracking-widest text-base font-bold focus:ring-2 focus:ring-purple-500 focus:outline-hidden"
                @keyup.enter="handleSwitchToOwner"
              />
              <p v-if="pinError" class="text-xs font-semibold text-rose-600 mt-1">
                {{ pinError }}
              </p>
            </div>
            <button
              type="button"
              class="w-full py-2.5 px-3 rounded-xl bg-purple-600 text-white font-bold text-xs hover:bg-purple-700 transition-all cursor-pointer flex items-center justify-center gap-2 shadow-xs"
              @click="handleSwitchToOwner"
            >
              <Lock class="w-4 h-4" />
              <span>Buka Akses Owner</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { Sparkles, Store, ChevronDown, Shield, ShieldCheck, User, Lock, X } from 'lucide-vue-next'
import { useAppSettings } from '../../composables/useAppSettings'
import { useToast } from '../../composables/useToast'

defineProps({
  activeOutletName: {
    type: String,
    default: 'FAN Laundry Kalibaru',
  },
})

defineEmits(['open-outlet-modal'])

const { settings, switchRole } = useAppSettings()
const { success, error } = useToast()

const showRoleModal = ref(false)
const inputPin = ref('')
const pinError = ref('')

const openRoleModal = () => {
  inputPin.value = ''
  pinError.value = ''
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

const handleSwitchToOwner = () => {
  const res = switchRole('OWNER', inputPin.value)
  if (res.success) {
    success('Akses Owner berhasil diaktifkan!')
    closeRoleModal()
  } else {
    pinError.value = res.message || 'PIN salah'
    error('Autentikasi gagal: ' + (res.message || 'PIN salah'))
  }
}
</script>
