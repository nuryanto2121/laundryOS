<template>
  <div class="min-h-screen bg-slate-900 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Background Glow Effects -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-brand-500/20 rounded-full blur-3xl pointer-events-none"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-purple-500/20 rounded-full blur-3xl pointer-events-none"></div>

    <div class="sm:mx-auto sm:w-full sm:max-w-md z-10 px-4">
      <!-- App Brand Logo -->
      <div class="flex flex-col items-center mb-6">
        <div class="w-16 h-16 rounded-3xl bg-gradient-to-tr from-brand-600 via-indigo-600 to-purple-600 flex items-center justify-center text-white shadow-xl shadow-brand-500/30 mb-3 border border-white/10">
          <WashingMachine class="w-8 h-8 animate-pulse" />
        </div>
        <h2 class="text-2xl font-black text-white tracking-tight">LaundryOS</h2>
        <p class="text-xs text-slate-400 mt-1 text-center font-medium">
          Sistem Operasi POS, Manajemen Multi-Cabang & Layanan Laundry
        </p>
      </div>

      <!-- Main Login Card -->
      <div class="bg-white/95 backdrop-blur-md py-8 px-6 shadow-2xl rounded-3xl sm:px-10 border border-white/20">
        <div class="mb-6">
          <h3 class="text-lg font-black text-slate-900 tracking-tight">Masuk ke Akun</h3>
          <p class="text-xs text-slate-500 mt-0.5">
            Silakan masukkan username dan password Anda untuk memulai.
          </p>
        </div>

        <!-- Error Notification Alert -->
        <div
          v-if="errorMessage"
          class="mb-4 p-3.5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-700 text-xs font-bold flex items-center gap-2 animate-shake"
        >
          <AlertCircle class="w-4 h-4 shrink-0 text-rose-500" />
          <span>{{ errorMessage }}</span>
        </div>

        <!-- Form Login -->
        <form class="space-y-4" @submit.prevent="handleLogin">
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">
              Username atau No. Handphone
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <User class="w-4 h-4" />
              </div>
              <input
                v-model="form.username"
                type="text"
                required
                autocomplete="username"
                placeholder="superadmin / owner / kasir1"
                class="w-full pl-10 pr-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-xs font-bold text-slate-900 focus:outline-hidden focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all placeholder:font-normal placeholder:text-slate-400"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">
              Password
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400">
                <Lock class="w-4 h-4" />
              </div>
              <input
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                required
                autocomplete="current-password"
                placeholder="Masukkan password..."
                class="w-full pl-10 pr-10 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-xs font-bold text-slate-900 focus:outline-hidden focus:ring-2 focus:ring-brand-500/20 focus:border-brand-500 transition-all placeholder:font-normal placeholder:text-slate-400"
              />
              <button
                type="button"
                class="absolute inset-y-0 right-0 pr-3.5 flex items-center text-slate-400 hover:text-slate-600 transition-colors cursor-pointer"
                @click="showPassword = !showPassword"
              >
                <EyeOff v-if="showPassword" class="w-4 h-4" />
                <Eye v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full mt-2 py-3 px-4 rounded-xl font-black text-xs text-white bg-gradient-to-r from-brand-600 to-indigo-600 hover:from-brand-500 hover:to-indigo-500 shadow-md shadow-brand-500/20 transition-all cursor-pointer flex items-center justify-center gap-2 disabled:opacity-50"
          >
            <Loader2 v-if="loading" class="w-4 h-4 animate-spin" />
            <LogIn v-else class="w-4 h-4" />
            <span>{{ loading ? 'Memverifikasi...' : 'Masuk Sekarang' }}</span>
          </button>
        </form>

        <!-- Divider -->
        <div class="relative my-6">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-slate-200"></div>
          </div>
          <div class="relative flex justify-center text-xs uppercase">
            <span class="bg-white px-3 text-[10px] font-black text-slate-400 tracking-wider">
              Akses Cepat (Demo Akun)
            </span>
          </div>
        </div>

        <!-- Quick Demo Login Buttons -->
        <div class="grid grid-cols-3 gap-2">
          <button
            type="button"
            class="p-2.5 rounded-xl border border-amber-200 bg-amber-50/70 hover:bg-amber-100/80 transition-colors text-center cursor-pointer group"
            @click="fillQuickAccount('superadmin', 'admin123')"
          >
            <Crown class="w-4 h-4 text-amber-600 mx-auto mb-1 group-hover:scale-110 transition-transform" />
            <span class="text-[10px] font-black text-amber-900 block truncate">Super Admin</span>
            <span class="text-[9px] text-amber-700/80 block">All Tenants</span>
          </button>

          <button
            type="button"
            class="p-2.5 rounded-xl border border-purple-200 bg-purple-50/70 hover:bg-purple-100/80 transition-colors text-center cursor-pointer group"
            @click="fillQuickAccount('owner', 'owner123')"
          >
            <Shield class="w-4 h-4 text-purple-600 mx-auto mb-1 group-hover:scale-110 transition-transform" />
            <span class="text-[10px] font-black text-purple-900 block truncate">Owner</span>
            <span class="text-[9px] text-purple-700/80 block">Full Akses</span>
          </button>

          <button
            type="button"
            class="p-2.5 rounded-xl border border-brand-200 bg-brand-50/70 hover:bg-brand-100/80 transition-colors text-center cursor-pointer group"
            @click="fillQuickAccount('kasir1', 'kasir123')"
          >
            <User class="w-4 h-4 text-brand-600 mx-auto mb-1 group-hover:scale-110 transition-transform" />
            <span class="text-[10px] font-black text-brand-900 block truncate">Kasir POS</span>
            <span class="text-[9px] text-brand-700/80 block">Granular ACL</span>
          </button>
        </div>
      </div>

      <!-- Footer Info -->
      <p class="text-center text-[11px] text-slate-500 mt-6">
        &copy; 2026 LaundryOS Multi-Outlet System. Dilindungi hak cipta.
      </p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  WashingMachine,
  User,
  Lock,
  Eye,
  EyeOff,
  LogIn,
  AlertCircle,
  Crown,
  Shield,
  Loader2,
} from 'lucide-vue-next'
import { useAuth } from '../composables/useAuth'
import { useAppSettings } from '../composables/useAppSettings'

const router = useRouter()
const route = useRoute()
const auth = useAuth()
const { settings } = useAppSettings()

const form = reactive({
  username: '',
  password: '',
})

const showPassword = ref(false)
const loading = ref(false)
const errorMessage = ref('')

const fillQuickAccount = (username, password) => {
  form.username = username
  form.password = password
  errorMessage.value = ''
}

const handleLogin = async () => {
  if (!form.username || !form.password) {
    errorMessage.value = 'Silakan isi username dan password'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const res = await auth.login(form.username, form.password)
    if (res.success) {
      // Sync app settings role
      settings.role = res.user.role
      const redirectPath = route.query.redirect || '/'
      router.push(redirectPath)
    } else {
      errorMessage.value = res.message || 'Login gagal'
    }
  } catch (err) {
    errorMessage.value = 'Terjadi kesalahan sistem saat mencoba masuk'
  } finally {
    loading.value = false
  }
}
</script>
