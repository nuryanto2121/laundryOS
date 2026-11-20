<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- 1. Top Action & Navigation Header -->
      <div
        class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-4 sm:p-6 rounded-2xl border border-slate-100 shadow-xs"
      >
        <div>
          <div class="flex items-center gap-2.5">
            <div
              class="w-10 h-10 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center shadow-xs"
            >
              <ShoppingBag class="w-5 h-5" />
            </div>
            <div>
              <h1 class="text-xl sm:text-2xl font-black text-slate-900 tracking-tight">
                Transaksi & Antrian Kasir
              </h1>
              <p class="text-xs sm:text-sm text-slate-500 mt-0.5">
                Kelola 5 tahapan siklus cucian, timbangan kasir, pelunasan & cetak nota.
              </p>
            </div>
          </div>
        </div>

        <div class="flex items-center gap-2 self-start sm:self-auto">
          <button
            :disabled="loading"
            class="p-2.5 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 text-slate-600 hover:text-slate-900 transition-all disabled:opacity-50"
            title="Refresh Data Antrian"
            @click="loadOrders"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          </button>

          <button
            class="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-brand-600 to-sky-500 hover:from-brand-700 hover:to-sky-600 text-white font-bold text-xs sm:text-sm shadow-md shadow-brand-500/20 active:scale-98 transition-all"
            @click="showCreateModal = true"
          >
            <PlusCircle class="w-4 h-4" />
            <span>+ Transaksi Baru</span>
          </button>
        </div>
      </div>

      <!-- 2. Status Queue Tabs -->
      <div class="flex items-center gap-2 overflow-x-auto pb-1 sm:pb-0 scrollbar-none">
        <button
          v-for="tab in statusTabs"
          :key="tab.id"
          class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap shrink-0 cursor-pointer"
          :class="
            activeStatusTab === tab.id
              ? 'bg-slate-900 text-white shadow-xs'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200/80'
          "
          @click="activeStatusTab = tab.id"
        >
          <span>{{ tab.label }}</span>
          <span
            class="px-2 py-0.5 rounded-full text-[10px] font-extrabold"
            :class="
              activeStatusTab === tab.id ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-700'
            "
          >
            {{ tabCounts[tab.id] || 0 }}
          </span>
        </button>
      </div>

      <!-- 3. Filter & Search Toolbar -->
      <div
        class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 bg-white p-3 sm:p-4 rounded-2xl border border-slate-100 shadow-xs"
      >
        <!-- Search Input -->
        <div class="relative flex-1">
          <Search class="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Cari no order, nama pelanggan, no WA, atau nomor rak..."
            class="w-full pl-10 pr-9 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-medium text-slate-800 placeholder-slate-400 focus:outline-none focus:border-brand-500 focus:bg-white transition-all"
            @keyup.enter="loadOrders"
          />
          <button
            v-if="searchQuery"
            class="absolute right-3 top-1/2 -translate-y-1/2 p-0.5 text-slate-400 hover:text-slate-600"
            @click="clearSearch"
          >
            <X class="w-3.5 h-3.5" />
          </button>
        </div>

        <!-- Payment Filter Dropdown -->
        <div class="flex items-center gap-2">
          <select
            v-model="paymentFilter"
            class="px-3 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-700 focus:outline-none focus:border-brand-500"
          >
            <option value="ALL">Semua Pembayaran</option>
            <option value="PAID">Lunas</option>
            <option value="UNPAID">Belum Lunas</option>
          </select>
        </div>
      </div>

      <!-- Special Filter Alert Banner (from Home KPI Cards) -->
      <div
        v-if="specialFilter"
        class="p-3 sm:p-3.5 rounded-2xl flex items-center justify-between gap-3 shadow-xs transition-all"
        :class="
          specialFilter === 'OVERDUE'
            ? 'bg-rose-50 border border-rose-200 text-rose-800'
            : specialFilter === 'DUE_TODAY'
              ? 'bg-amber-50 border border-amber-200 text-amber-800'
              : 'bg-brand-50 border border-brand-200 text-brand-800'
        "
      >
        <div class="flex items-center gap-2.5 text-xs">
          <AlertCircle v-if="specialFilter === 'OVERDUE'" class="w-4 h-4 text-rose-600 shrink-0" />
          <Clock
            v-else-if="specialFilter === 'DUE_TODAY'"
            class="w-4 h-4 text-amber-600 shrink-0"
          />
          <ShoppingBag v-else class="w-4 h-4 text-brand-600 shrink-0" />
          <div>
            <span class="font-black">
              {{
                specialFilter === 'OVERDUE'
                  ? 'Filter Aktif: Cucian Terlambat (Melewati Batas SLA)'
                  : specialFilter === 'DUE_TODAY'
                    ? 'Filter Aktif: Order Harus Selesai Hari Ini'
                    : 'Filter Aktif: Order Masuk Hari Ini'
              }}
            </span>
            <span class="text-[11px] block opacity-80">
              {{
                specialFilter === 'OVERDUE'
                  ? 'Menampilkan cucian aktif yang butuh diprioritaskan pengerjaannya.'
                  : specialFilter === 'DUE_TODAY'
                    ? 'Menampilkan orderan dengan target selesai hari ini.'
                    : 'Menampilkan order yang masuk pada hari ini.'
              }}
            </span>
          </div>
        </div>

        <button
          type="button"
          class="px-2.5 py-1.5 rounded-xl bg-white text-xs font-bold shadow-2xs hover:bg-slate-50 transition-all cursor-pointer shrink-0"
          @click="clearSpecialFilter"
        >
          Hapus Filter (Tampilkan Semua)
        </button>
      </div>

      <!-- 4. Orders List / Table -->
      <div class="space-y-3">
        <!-- Loading State -->
        <div
          v-if="loading"
          class="bg-white rounded-2xl border border-slate-100 p-12 text-center text-slate-400 space-y-2 shadow-xs"
        >
          <RefreshCw class="w-6 h-6 animate-spin text-brand-600 mx-auto" />
          <p class="text-xs font-medium">Memuat data transaksi & antrian...</p>
        </div>

        <!-- Empty State -->
        <div
          v-else-if="filteredOrders.length === 0"
          class="bg-white rounded-2xl border border-slate-100 p-12 text-center space-y-3 shadow-xs"
        >
          <div
            class="w-12 h-12 rounded-2xl bg-slate-100 text-slate-400 flex items-center justify-center mx-auto"
          >
            <ShoppingBag class="w-6 h-6" />
          </div>
          <div>
            <h3 class="font-bold text-sm text-slate-800">Tidak ada pesanan ditemukan</h3>
            <p class="text-xs text-slate-400 max-w-sm mx-auto mt-1">
              {{
                searchQuery || paymentFilter !== 'ALL'
                  ? 'Coba ubah kata kunci pencarian atau filter pembayaran Anda.'
                  : 'Belum ada antrian cucian di status ini.'
              }}
            </p>
          </div>
          <button
            v-if="searchQuery || paymentFilter !== 'ALL'"
            class="px-3.5 py-1.5 rounded-xl bg-slate-100 hover:bg-slate-200 text-xs font-bold text-slate-700 transition-all inline-block cursor-pointer"
            @click="resetFilters"
          >
            Reset Filter
          </button>
        </div>

        <!-- Order Cards -->
        <div
          v-for="order in filteredOrders"
          v-else
          :key="order.id"
          class="bg-white rounded-2xl border border-slate-100 p-4 sm:p-5 shadow-xs hover:border-brand-200 transition-all flex flex-col sm:flex-row sm:items-center justify-between gap-4 group cursor-pointer"
          @click="openOrderDetail(order)"
        >
          <!-- Left: Order Identity & Customer Info -->
          <div class="flex items-start gap-3.5 flex-1 min-w-0">
            <!-- Status Avatar / Indicator -->
            <div
              class="w-11 h-11 rounded-xl flex items-center justify-center font-black text-xs shrink-0 shadow-xs"
              :class="getStatusAvatarClass(order.status)"
            >
              {{ getStatusInitials(order.status) }}
            </div>

            <!-- Details -->
            <div class="space-y-1 min-w-0 flex-1">
              <!-- Header Badges -->
              <div class="flex flex-wrap items-center gap-2">
                <span
                  class="font-black text-sm sm:text-base text-slate-900 font-mono tracking-tight group-hover:text-brand-600 transition-colors"
                >
                  {{ formatDisplayOrder(order.order_number) }}
                </span>

                <span
                  class="text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider"
                  :class="getStatusClass(order.status)"
                >
                  {{ order.status }}
                </span>

                <span
                  class="text-[10px] font-bold px-2 py-0.5 rounded-full uppercase tracking-wider"
                  :class="
                    order.payment_status === 'PAID'
                      ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                      : 'bg-rose-50 text-rose-700 border border-rose-200'
                  "
                >
                  {{ order.payment_status === 'PAID' ? 'LUNAS' : 'BELUM LUNAS' }}
                </span>

                <span
                  v-if="order.rack_number"
                  class="text-[10px] font-bold px-1.5 py-0.5 rounded-md bg-slate-100 text-slate-700 flex items-center gap-1"
                >
                  <Tag class="w-2.5 h-2.5 text-slate-400" />
                  <span>Rak: {{ order.rack_number }}</span>
                </span>
              </div>

              <!-- Customer Info -->
              <div class="flex items-center gap-2 text-xs">
                <span class="font-bold text-slate-800">{{ order.customer_name }}</span>
                <span class="text-slate-300">•</span>
                <a
                  :href="`https://wa.me/${order.customer_phone?.replace(/[^0-9]/g, '')}`"
                  target="_blank"
                  class="text-slate-500 hover:text-emerald-600 flex items-center gap-1 font-medium transition-colors"
                  @click.stop
                >
                  <span>{{ order.customer_phone }}</span>
                  <ExternalLink class="w-3 h-3 text-slate-400" />
                </a>
              </div>

              <!-- Services & Perfume Spec -->
              <div class="text-xs text-slate-500 flex flex-wrap items-center gap-2 pt-0.5">
                <span class="font-medium text-slate-700">
                  {{ order.primary_service_name || 'Layanan Laundry' }}
                </span>
                <span v-if="order.total_weight" class="text-slate-400 font-semibold">
                  ({{ order.total_weight }} Kg)
                </span>
                <span
                  v-if="order.perfume_name"
                  class="flex items-center gap-1 text-[11px] text-amber-700 bg-amber-50 px-1.5 py-0.5 rounded"
                >
                  <Sparkles class="w-2.5 h-2.5 text-amber-500" />
                  <span>Parfum: {{ order.perfume_name }}</span>
                </span>
                <span
                  v-if="order.estimated_completion_at"
                  class="text-[11px] text-slate-400 flex items-center gap-1"
                >
                  <Clock class="w-3 h-3" />
                  <span>Selesai: {{ formatDateTime(order.estimated_completion_at) }}</span>
                </span>
              </div>
            </div>
          </div>

          <!-- Middle: Price & Payment Details -->
          <div
            class="flex sm:flex-col items-baseline sm:items-end justify-between sm:justify-center border-t sm:border-t-0 border-slate-100 pt-2 sm:pt-0 shrink-0"
          >
            <div class="text-left sm:text-right">
              <span class="text-[10px] text-slate-400 uppercase tracking-wider block"
                >Total Tagihan</span
              >
              <span class="text-sm sm:text-base font-black text-slate-900 tracking-tight">
                Rp {{ formatRupiah(order.total_price) }}
              </span>
            </div>
            <div v-if="order.payment_status !== 'PAID'" class="text-right">
              <span class="text-[10px] font-semibold text-rose-600 block">
                Sisa Rp
                {{ formatRupiah(Number(order.total_price) - Number(order.paid_amount || 0)) }}
              </span>
            </div>
          </div>

          <!-- Right: Actions Toolbar -->
          <div
            class="flex items-center gap-2 shrink-0 border-t sm:border-t-0 border-slate-100 pt-2 sm:pt-0"
          >
            <!-- 1-Click Quick Progression Button -->
            <button
              v-if="canQuickProgress(order.status)"
              :disabled="actionInProgressId === order.id"
              class="px-3 py-1.5 rounded-xl font-bold text-xs flex items-center gap-1.5 transition-all shadow-xs disabled:opacity-50 cursor-pointer"
              :class="getQuickProgressButtonClass(order.status)"
              :title="getQuickProgressTitle(order.status)"
              @click.stop="quickProgress(order)"
            >
              <RefreshCw v-if="actionInProgressId === order.id" class="w-3.5 h-3.5 animate-spin" />
              <component :is="getQuickProgressIcon(order.status)" v-else class="w-3.5 h-3.5" />
              <span>{{ getQuickProgressLabel(order.status) }}</span>
            </button>

            <!-- Quick Print Thermal Button -->
            <button
              class="p-2 rounded-xl bg-slate-50 hover:bg-slate-100 text-slate-600 hover:text-slate-900 border border-slate-200 transition-all cursor-pointer"
              title="Cetak Nota Thermal"
              @click.stop="openPrintReceipt(order)"
            >
              <Printer class="w-4 h-4" />
            </button>

            <!-- Detail Drawer Button -->
            <button
              class="px-3 py-1.5 rounded-xl bg-slate-100 group-hover:bg-brand-50 group-hover:text-brand-600 text-slate-700 font-bold text-xs transition-all flex items-center gap-1 cursor-pointer"
              @click.stop="openOrderDetail(order)"
            >
              <span>Kelola</span>
              <ChevronRight class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>

        <!-- Server-side Pagination Bar -->
        <PaginationBar
          v-if="totalOrders > 0"
          v-model:current-page="currentPage"
          v-model:items-per-page="itemsPerPage"
          :total-items="totalOrders"
        />
      </div>

      <!-- Modals -->
      <CreateOrderModal v-model="showCreateModal" @order-created="handleOrderCreated" />

      <OrderDetailModal
        v-if="selectedOrder"
        v-model="showDetailModal"
        :order="selectedOrder"
        @order-updated="handleOrderUpdated"
      />

      <ReceiptPrintModal v-if="orderForPrint" v-model="showPrintModal" :order="orderForPrint" />
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ShoppingBag,
  PlusCircle,
  RefreshCw,
  Search,
  X,
  Tag,
  ExternalLink,
  Sparkles,
  Clock,
  Printer,
  ChevronRight,
  Play,
  CheckCircle2,
  PackageCheck,
  AlertCircle,
} from 'lucide-vue-next'
import AppLayout from '../components/layout/AppLayout.vue'
import CreateOrderModal from '../components/pos/CreateOrderModal.vue'
import OrderDetailModal from '../components/pos/OrderDetailModal.vue'
import ReceiptPrintModal from '../components/pos/ReceiptPrintModal.vue'
import PaginationBar from '../components/ui/PaginationBar.vue'
import { api, formatRupiah, formatDateTime } from '../api/client'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const loading = ref(false)
const orders = ref([])
const activeStatusTab = ref('ALL')
const searchQuery = ref('')
const paymentFilter = ref('ALL')
const specialFilter = ref('')
const actionInProgressId = ref(null)

const currentPage = ref(1)
const itemsPerPage = ref(10)
const totalOrders = ref(0)

const isToday = (dateStr) => {
  if (!dateStr) return false
  const d = new Date(dateStr)
  const now = new Date()
  return (
    d.getFullYear() === now.getFullYear() &&
    d.getMonth() === now.getMonth() &&
    d.getDate() === now.getDate()
  )
}

const isDueToday = (order) => {
  if (!order.estimated_completion_at) return false
  if (order.status === 'SELESAI' || order.status === 'BATAL') return false
  return isToday(order.estimated_completion_at)
}

const isOverdue = (order) => {
  if (!order.estimated_completion_at) return false
  if (order.status === 'SELESAI' || order.status === 'BATAL') return false
  return new Date(order.estimated_completion_at) < new Date()
}

const showCreateModal = ref(false)
const showDetailModal = ref(false)
const showPrintModal = ref(false)

const selectedOrder = ref(null)
const orderForPrint = ref(null)

const statusTabs = [
  { id: 'ALL', label: 'Semua Transaksi' },
  { id: 'ANTRIAN', label: 'Antrian Cucian' },
  { id: 'PROSES', label: 'Proses Cuci' },
  { id: 'SIAP AMBIL', label: 'Siap Diambil' },
  { id: 'SELESAI', label: 'Selesai' },
  { id: 'BATAL', label: 'Dibatalkan' },
]

const tabCounts = computed(() => {
  const counts = {
    ALL: totalOrders.value || orders.value.length,
    ANTRIAN: 0,
    PROSES: 0,
    'SIAP AMBIL': 0,
    SELESAI: 0,
    BATAL: 0,
  }

  orders.value.forEach((o) => {
    if (counts[o.status] !== undefined) {
      counts[o.status]++
    }
  })

  return counts
})

const filteredOrders = computed(() => {
  return orders.value.filter((order) => {
    // 1. Status Tab filter
    if (activeStatusTab.value !== 'ALL' && order.status !== activeStatusTab.value) {
      return false
    }

    // 2. Payment filter
    if (paymentFilter.value === 'PAID' && order.payment_status !== 'PAID') {
      return false
    }
    if (paymentFilter.value === 'UNPAID' && order.payment_status === 'PAID') {
      return false
    }

    // 3. Search text query
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase().trim()
      const matchOrderNum = order.order_number?.toLowerCase().includes(q)
      const matchCustomer = order.customer_name?.toLowerCase().includes(q)
      const matchPhone = order.customer_phone?.toLowerCase().includes(q)
      const matchRack = order.rack_number?.toLowerCase().includes(q)

      if (!matchOrderNum && !matchCustomer && !matchPhone && !matchRack) {
        return false
      }
    }

    // 4. Special Filter (from Home KPI Cards)
    if (specialFilter.value === 'TODAY' && !isToday(order.created_at)) {
      return false
    }
    if (specialFilter.value === 'DUE_TODAY' && !isDueToday(order)) {
      return false
    }
    if (specialFilter.value === 'OVERDUE' && !isOverdue(order)) {
      return false
    }

    return true
  })
})

const formatDisplayOrder = (num) => {
  if (!num) return ''
  return num.startsWith('ORD') ? '#' + num.substring(3) : num
}

const getStatusClass = (status) => {
  switch (status) {
    case 'ANTRIAN':
      return 'bg-blue-50 text-blue-700 border border-blue-200'
    case 'PROSES':
      return 'bg-amber-50 text-amber-700 border border-amber-200'
    case 'SIAP AMBIL':
      return 'bg-purple-50 text-purple-700 border border-purple-200'
    case 'SELESAI':
      return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
    case 'BATAL':
      return 'bg-rose-50 text-rose-700 border border-rose-200'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

const getStatusAvatarClass = (status) => {
  switch (status) {
    case 'ANTRIAN':
      return 'bg-blue-50 text-blue-700 border border-blue-200'
    case 'PROSES':
      return 'bg-amber-50 text-amber-700 border border-amber-200'
    case 'SIAP AMBIL':
      return 'bg-purple-50 text-purple-700 border border-purple-200'
    case 'SELESAI':
      return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
    case 'BATAL':
      return 'bg-rose-50 text-rose-700 border border-rose-200'
    default:
      return 'bg-slate-100 text-slate-600'
  }
}

const getStatusInitials = (status) => {
  switch (status) {
    case 'ANTRIAN':
      return 'ATR'
    case 'PROSES':
      return 'PRS'
    case 'SIAP AMBIL':
      return 'SAP'
    case 'SELESAI':
      return 'OK'
    case 'BATAL':
      return 'BTL'
    default:
      return 'ORD'
  }
}

const canQuickProgress = (status) => {
  return status === 'ANTRIAN' || status === 'PROSES' || status === 'SIAP AMBIL'
}

const getQuickProgressLabel = (status) => {
  switch (status) {
    case 'ANTRIAN':
      return 'Mulai Cuci'
    case 'PROSES':
      return 'Siap Diambil'
    case 'SIAP AMBIL':
      return 'Selesai / Ambil'
    default:
      return 'Progres'
  }
}

const getQuickProgressIcon = (status) => {
  switch (status) {
    case 'ANTRIAN':
      return Play
    case 'PROSES':
      return CheckCircle2
    case 'SIAP AMBIL':
      return PackageCheck
    default:
      return RefreshCw
  }
}

const getQuickProgressButtonClass = (status) => {
  switch (status) {
    case 'ANTRIAN':
      return 'bg-amber-500 hover:bg-amber-600 text-white shadow-amber-500/20'
    case 'PROSES':
      return 'bg-purple-600 hover:bg-purple-700 text-white shadow-purple-500/20'
    case 'SIAP AMBIL':
      return 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-emerald-500/20'
    default:
      return 'bg-slate-900 text-white'
  }
}

const getQuickProgressTitle = (status) => {
  switch (status) {
    case 'ANTRIAN':
      return 'Ubah status pesanan ke tahap PROSES pengerjaan'
    case 'PROSES':
      return 'Ubah status pesanan ke tahap SIAP AMBIL & otomatis kirim WA notifikasi'
    case 'SIAP AMBIL':
      return 'Tandai cucian sudah diserahkan / selesai'
    default:
      return ''
  }
}

const quickProgress = async (order) => {
  let targetStatus = ''
  if (order.status === 'ANTRIAN') targetStatus = 'PROSES'
  else if (order.status === 'PROSES') targetStatus = 'SIAP AMBIL'
  else if (order.status === 'SIAP AMBIL') {
    if (order.payment_status !== 'PAID') {
      openOrderDetail(order)
      toast.warning(
        'Pesanan belum lunas! Silakan selesaikan pembayaran terlebih dahulu sebelum diserahkan.'
      )
      return
    }
    targetStatus = 'SELESAI'
  }

  if (!targetStatus) return

  actionInProgressId.value = order.id
  try {
    await api.patch(`/pos/orders/${order.id}/status`, {
      status: targetStatus,
      changed_by: 'Kasir Utama',
      notes: `Progresi cepat kasir via dashboard`,
    })
    toast.success(`Status pesanan berhasil diubah menjadi ${targetStatus}!`)
    await loadOrders()
  } catch (err) {
    toast.error(err.message || 'Gagal mengubah status antrian.')
  } finally {
    actionInProgressId.value = null
  }
}

const openOrderDetail = (order) => {
  selectedOrder.value = order
  showDetailModal.value = true
}

const openPrintReceipt = (order) => {
  orderForPrint.value = order
  showPrintModal.value = true
}

const handleOrderCreated = async (newOrderData) => {
  await loadOrders()
  if (newOrderData?.order) {
    orderForPrint.value = newOrderData.order
    showPrintModal.value = true
  }
}

const handleOrderUpdated = async () => {
  await loadOrders()
  if (selectedOrder.value) {
    // Refresh selected order data
    try {
      const refreshed = await api.get(`/pos/orders/${selectedOrder.value.id}`)
      if (refreshed?.order) {
        selectedOrder.value = refreshed.order
      }
    } catch {
      // ignore
    }
  }
}

const clearSearch = () => {
  searchQuery.value = ''
}

const clearSpecialFilter = () => {
  specialFilter.value = ''
  router.replace({ path: '/orders' })
}

const resetFilters = () => {
  searchQuery.value = ''
  paymentFilter.value = 'ALL'
  activeStatusTab.value = 'ALL'
  specialFilter.value = ''
  router.replace({ path: '/orders' })
}

const loadOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      limit: itemsPerPage.value,
    }
    if (activeStatusTab.value !== 'ALL') {
      params.status = activeStatusTab.value
    }
    const q = searchQuery.value.trim()
    if (q.length >= 3) {
      params.search = q
    }

    const res = await api.get('/pos/orders', params)
    if (res && res.data && res.meta) {
      orders.value = res.data
      totalOrders.value = res.meta.total_items
    } else if (Array.isArray(res)) {
      orders.value = res
      totalOrders.value = res.length
    } else {
      orders.value = []
      totalOrders.value = 0
    }

    // Check if route query has an order id to open
    if (route.query.id) {
      const matched = orders.value.find((o) => o.id === route.query.id)
      if (matched) {
        openOrderDetail(matched)
      }
    }
  } catch (err) {
    console.warn('Failed loading orders:', err)
  } finally {
    loading.value = false
  }
}

watch([currentPage, itemsPerPage], () => {
  loadOrders()
})

watch(activeStatusTab, () => {
  currentPage.value = 1
  loadOrders()
})

// 3-character threshold filter watcher for order search
let orderSearchTimer = null
watch(searchQuery, (val) => {
  const trimmed = (val || '').trim()
  if (orderSearchTimer) clearTimeout(orderSearchTimer)

  if (trimmed.length === 0) {
    currentPage.value = 1
    loadOrders()
  } else if (trimmed.length >= 3) {
    orderSearchTimer = setTimeout(() => {
      currentPage.value = 1
      loadOrders()
    }, 350)
  }
  // If 1 or 2 characters: DO NOT hit API
})

watch(
  () => route.query.new,
  (newVal) => {
    if (newVal === 'true') {
      showCreateModal.value = true
      // Clear query without retriggering
      router.replace({ path: '/orders' })
    }
  },
  { immediate: true }
)

watch(
  () => route.query.filter,
  (newFilter) => {
    specialFilter.value = newFilter || ''
  },
  { immediate: true }
)

watch(
  () => route.query.id,
  (newId) => {
    if (newId && orders.value.length > 0) {
      const matched = orders.value.find((o) => o.id === newId)
      if (matched) {
        openOrderDetail(matched)
      }
    }
  }
)

onMounted(() => {
  loadOrders()
})
</script>
