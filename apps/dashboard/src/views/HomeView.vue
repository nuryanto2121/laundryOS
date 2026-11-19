<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- 1. Header Greeting & Date Bar -->
      <div
        class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-4 sm:p-6 rounded-2xl border border-slate-100 shadow-xs"
      >
        <div>
          <div class="flex items-center gap-2">
            <h1 class="text-xl sm:text-2xl font-black text-slate-900 tracking-tight">
              Halo, Kasir FAN Laundry! 👋
            </h1>
          </div>
          <p class="text-xs sm:text-sm text-slate-500 mt-1">
            Pantau arus transaksi kasir, operasional antrian mesin cuci, dan SLA harian.
          </p>
        </div>
        <div class="flex items-center gap-2 self-start sm:self-auto">
          <div
            class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-600"
          >
            <Calendar class="w-4 h-4 text-brand-600" />
            <span>{{ todayFormatted }}</span>
          </div>
          <button
            :disabled="loading"
            class="p-2 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 text-slate-600 hover:text-slate-900 transition-all disabled:opacity-50"
            title="Refresh Data"
            @click="loadHomeData"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          </button>
        </div>
      </div>

      <!-- 2. Real-time KPI Metrics Cards -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
        <!-- KPI 1: Omset Hari Ini -->
        <div
          class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between hover:border-brand-300 transition-all group"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
              >Omset Hari Ini</span
            >
            <div
              class="w-9 h-9 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center group-hover:scale-110 transition-transform"
            >
              <Wallet class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-3">
            <h3 class="text-lg sm:text-2xl font-black text-slate-900 tracking-tight">
              Rp {{ formatRupiah(metrics.omset_hari_ini) }}
            </h3>
            <p class="text-[11px] text-slate-400 mt-1 flex items-center gap-1">
              <TrendingUp class="w-3.5 h-3.5 text-emerald-500" />
              <span>Penjualan bruto tercatat</span>
            </p>
          </div>
        </div>

        <!-- KPI 2: Order Masuk -->
        <div
          class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between hover:border-emerald-300 transition-all group"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
              >Order Masuk</span
            >
            <div
              class="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center group-hover:scale-110 transition-transform"
            >
              <ShoppingBag class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-3">
            <h3 class="text-lg sm:text-2xl font-black text-slate-900 tracking-tight">
              {{ metrics.order_masuk_hari_ini || 0 }}
              <span class="text-xs sm:text-sm font-semibold text-slate-400">Order</span>
            </h3>
            <p class="text-[11px] text-slate-400 mt-1">Walk-in & Pickup bot</p>
          </div>
        </div>

        <!-- KPI 3: Harus Selesai -->
        <div
          class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between hover:border-amber-300 transition-all group"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
              >Harus Selesai</span
            >
            <div
              class="w-9 h-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center group-hover:scale-110 transition-transform"
            >
              <Clock class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-3">
            <h3 class="text-lg sm:text-2xl font-black text-slate-900 tracking-tight">
              {{ metrics.harus_selesai_hari_ini || 0 }}
              <span class="text-xs sm:text-sm font-semibold text-slate-400">Order</span>
            </h3>
            <p class="text-[11px] text-slate-400 mt-1">Jatuh tempo SLA hari ini</p>
          </div>
        </div>

        <!-- KPI 4: Terlambat -->
        <div
          class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between hover:border-rose-300 transition-all group"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider">Terlambat</span>
            <div
              class="w-9 h-9 rounded-xl flex items-center justify-center transition-transform group-hover:scale-110"
              :class="
                metrics.terlambat_hari_ini > 0
                  ? 'bg-rose-50 text-rose-600 animate-pulse'
                  : 'bg-slate-50 text-slate-400'
              "
            >
              <AlertCircle class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-3">
            <h3
              class="text-lg sm:text-2xl font-black tracking-tight"
              :class="metrics.terlambat_hari_ini > 0 ? 'text-rose-600' : 'text-slate-900'"
            >
              {{ metrics.terlambat_hari_ini || 0 }}
              <span class="text-xs sm:text-sm font-semibold text-slate-400">Order</span>
            </h3>
            <span
              v-if="metrics.terlambat_hari_ini === 0"
              class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 mt-1 inline-block"
            >
              Semua SLA Aman ✨
            </span>
            <span
              v-else
              class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-rose-50 text-rose-700 mt-1 inline-block"
            >
              Perlu Dikebut! ⚠️
            </span>
          </div>
        </div>
      </div>

      <!-- 3. Primary Quick Action Bar -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <!-- Action 1: Transaksi Baru (Hero Button) -->
        <RouterLink
          to="/orders?new=true"
          class="sm:col-span-1 p-4 rounded-2xl bg-gradient-to-r from-brand-600 to-sky-500 hover:from-brand-700 hover:to-sky-600 text-white shadow-lg shadow-brand-500/25 flex items-center justify-between transition-all transform active:scale-98 group"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-11 h-11 rounded-xl bg-white/20 backdrop-blur-xs flex items-center justify-center text-white"
            >
              <PlusCircle class="w-6 h-6" />
            </div>
            <div>
              <h4 class="font-bold text-base leading-tight">Transaksi Baru</h4>
              <p class="text-xs text-white/80 mt-0.5">Kasir Walk-in & Timbangan</p>
            </div>
          </div>
          <ArrowRight
            class="w-5 h-5 text-white/70 group-hover:translate-x-1 transition-transform"
          />
        </RouterLink>

        <!-- Action 2: Catat Kas Kecil (Pengeluaran) -->
        <button
          class="p-4 rounded-2xl bg-white hover:bg-slate-50/80 border border-slate-200 text-slate-800 shadow-xs flex items-center justify-between transition-all text-left group"
          @click="showExpenseModal = true"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-11 h-11 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center"
            >
              <Receipt class="w-5 h-5" />
            </div>
            <div>
              <h4 class="font-bold text-sm text-slate-900 leading-tight">Catat Pengeluaran</h4>
              <p class="text-xs text-slate-500 mt-0.5">Kas kecil operasional harian</p>
            </div>
          </div>
          <Plus class="w-5 h-5 text-slate-400 group-hover:text-slate-600 transition-colors" />
        </button>

        <!-- Action 3: Cari Transaksi -->
        <button
          class="p-4 rounded-2xl bg-white hover:bg-slate-50/80 border border-slate-200 text-slate-800 shadow-xs flex items-center justify-between transition-all text-left group"
          @click="showSearchModal = true"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-11 h-11 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center"
            >
              <Search class="w-5 h-5" />
            </div>
            <div>
              <h4 class="font-bold text-sm text-slate-900 leading-tight">Cari Transaksi</h4>
              <p class="text-xs text-slate-500 mt-0.5">Cari #Order, no WA, nama, rak</p>
            </div>
          </div>
          <Search class="w-5 h-5 text-slate-400 group-hover:text-slate-600 transition-colors" />
        </button>
      </div>

      <!-- 4. Live Queue Preview (Antrian Cucian Berjalan) -->
      <div class="bg-white rounded-2xl border border-slate-100 shadow-xs overflow-hidden">
        <!-- Table Header -->
        <div class="p-4 sm:p-5 border-b border-slate-100 flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <div class="w-2.5 h-2.5 rounded-full bg-brand-500 animate-ping"></div>
            <div>
              <h3 class="font-bold text-base text-slate-900">Antrian Cucian Berjalan</h3>
              <p class="text-xs text-slate-500">Order dalam tahap Antrian & Proses pengerjaan</p>
            </div>
          </div>
          <RouterLink
            to="/orders"
            class="text-xs font-bold text-brand-600 hover:text-brand-700 flex items-center gap-1 hover:underline"
          >
            <span>Semua Transaksi</span>
            <ChevronRight class="w-4 h-4" />
          </RouterLink>
        </div>

        <!-- Table / List Content -->
        <div
          v-if="loadingOrders"
          class="p-8 text-center text-xs text-slate-400 flex items-center justify-center gap-2"
        >
          <Loader2 class="w-4 h-4 animate-spin text-brand-600" />
          <span>Memuat antrian cucian...</span>
        </div>

        <div v-else-if="activeQueueOrders.length === 0" class="p-10 text-center text-slate-400">
          <CheckCircle2 class="w-10 h-10 text-emerald-400 mx-auto mb-2" />
          <p class="text-sm font-bold text-slate-700">Tidak ada antrian tertunda!</p>
          <p class="text-xs text-slate-400 mt-1">
            Semua cucian sudah selesai atau siap diambil oleh pelanggan.
          </p>
        </div>

        <div v-else class="divide-y divide-slate-100">
          <div
            v-for="order in activeQueueOrders"
            :key="order.id"
            class="p-4 sm:p-5 hover:bg-slate-50/80 transition-all cursor-pointer flex flex-col sm:flex-row sm:items-center justify-between gap-3 group"
            @click="navigateToOrder(order.id)"
          >
            <!-- Order Meta -->
            <div class="flex items-start sm:items-center gap-3.5">
              <div
                class="w-10 h-10 rounded-xl flex items-center justify-center font-bold text-xs shrink-0"
                :class="
                  order.status === 'ANTRIAN'
                    ? 'bg-blue-50 text-blue-700 border border-blue-200'
                    : 'bg-amber-50 text-amber-700 border border-amber-200'
                "
              >
                {{ order.status === 'ANTRIAN' ? 'ATR' : 'PRS' }}
              </div>
              <div class="space-y-0.5">
                <div class="flex items-center gap-2">
                  <span
                    class="font-extrabold text-sm text-slate-900 font-mono tracking-tight group-hover:text-brand-600 transition-colors"
                  >
                    {{ formatDisplayOrder(order.order_number) }}
                  </span>
                  <span
                    class="text-[10px] font-bold px-2 py-0.5 rounded-full uppercase"
                    :class="
                      order.status === 'ANTRIAN'
                        ? 'bg-blue-50 text-blue-700 border border-blue-200'
                        : 'bg-amber-50 text-amber-700 border border-amber-200'
                    "
                  >
                    {{ order.status }}
                  </span>
                  <span
                    v-if="order.rack_number"
                    class="text-[10px] font-bold px-1.5 py-0.5 rounded-md bg-slate-100 text-slate-700"
                  >
                    Rak: {{ order.rack_number }}
                  </span>
                </div>
                <p class="text-xs font-bold text-slate-800">
                  {{ order.customer_name }}
                  <span class="text-slate-400 font-normal">({{ order.customer_phone }})</span>
                </p>
                <p class="text-[11px] text-slate-500">
                  {{ order.primary_service_name || 'Layanan Laundry' }} •
                  {{ order.total_weight ? `${order.total_weight} Kg` : '' }}
                </p>
              </div>
            </div>

            <!-- Right Info & Action -->
            <div
              class="flex items-center justify-between sm:justify-end gap-4 shrink-0 pl-13 sm:pl-0"
            >
              <div class="text-left sm:text-right">
                <p class="text-xs font-bold text-slate-900">
                  Rp {{ formatRupiah(order.total_price) }}
                </p>
                <p
                  class="text-[10px] font-semibold mt-0.5"
                  :class="order.payment_status === 'PAID' ? 'text-emerald-600' : 'text-rose-600'"
                >
                  {{ order.payment_status === 'PAID' ? 'LUNAS' : 'BELUM LUNAS' }}
                </p>
              </div>
              <button
                class="px-3 py-1.5 rounded-xl bg-slate-100 group-hover:bg-brand-50 group-hover:text-brand-600 text-slate-600 text-xs font-bold transition-all flex items-center gap-1"
              >
                <span>Kelola</span>
                <ChevronRight class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Modals -->
      <QuickExpenseModal v-model="showExpenseModal" @expense-added="loadHomeData" />
      <SearchOrderModal v-model="showSearchModal" />
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Calendar,
  RefreshCw,
  Wallet,
  ShoppingBag,
  Clock,
  AlertCircle,
  TrendingUp,
  PlusCircle,
  ArrowRight,
  Receipt,
  Plus,
  Search,
  ChevronRight,
  Loader2,
  CheckCircle2,
} from 'lucide-vue-next'
import AppLayout from '../components/layout/AppLayout.vue'
import QuickExpenseModal from '../components/modals/QuickExpenseModal.vue'
import SearchOrderModal from '../components/modals/SearchOrderModal.vue'
import { api, formatRupiah } from '../api/client'

const router = useRouter()

const loading = ref(false)
const loadingOrders = ref(false)
const showExpenseModal = ref(false)
const showSearchModal = ref(false)

const metrics = ref({
  omset_hari_ini: 0,
  order_masuk_hari_ini: 0,
  harus_selesai_hari_ini: 0,
  terlambat_hari_ini: 0,
})

const activeQueueOrders = ref([])

const todayFormatted = computed(() => {
  const d = new Date()
  return d.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
})

const formatDisplayOrder = (num) => {
  if (!num) return ''
  return num.startsWith('ORD') ? '#' + num.substring(3) : num
}

const navigateToOrder = (id) => {
  router.push({ path: '/orders', query: { id } })
}

const loadHomeData = async () => {
  loading.value = true
  try {
    const data = await api.get('/reports/home-metrics')
    if (data) {
      metrics.value = data
    }
  } catch (err) {
    console.warn('Failed loading home metrics:', err)
  } finally {
    loading.value = false
  }
}

const loadActiveQueue = async () => {
  loadingOrders.value = true
  try {
    // Fetch active queue orders
    const antrian = await api.get('/pos/orders', { status: 'ANTRIAN' })
    const proses = await api.get('/pos/orders', { status: 'PROSES' })
    const combined = [...(antrian || []), ...(proses || [])]
    activeQueueOrders.value = combined.slice(0, 8)
  } catch (err) {
    console.warn('Failed loading active queue:', err)
  } finally {
    loadingOrders.value = false
  }
}

onMounted(() => {
  loadHomeData()
  loadActiveQueue()
})
</script>
