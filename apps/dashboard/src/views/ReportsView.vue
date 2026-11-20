<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- 1. Header & Navigation -->
      <div
        class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-4 sm:p-6 rounded-2xl border border-slate-100 shadow-xs"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shadow-xs"
          >
            <Receipt class="w-5 h-5" />
          </div>
          <div>
            <h1 class="text-xl sm:text-2xl font-black text-slate-900 tracking-tight">
              Buku Kas Kecil & Pengeluaran
            </h1>
            <p class="text-xs sm:text-sm text-slate-500 mt-0.5">
              Pencatatan biaya operasional harian outlet FAN Laundry Kalibaru
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2 self-start sm:self-auto">
          <button
            :disabled="loading"
            class="p-2.5 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 text-slate-600 hover:text-slate-900 transition-all disabled:opacity-50 cursor-pointer"
            title="Refresh Data"
            @click="loadExpensesData"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          </button>

          <button
            class="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-amber-600 to-amber-700 hover:from-amber-700 hover:to-amber-800 text-white font-bold text-xs sm:text-sm shadow-md shadow-amber-600/20 active:scale-98 transition-all cursor-pointer"
            @click="showCreateModal = true"
          >
            <PlusCircle class="w-4 h-4" />
            <span>+ Catat Pengeluaran</span>
          </button>
        </div>
      </div>

      <!-- 2. Navigation Tabs (Reports & Petty Cash) -->
      <div class="flex items-center gap-2 overflow-x-auto pb-1 sm:pb-0 scrollbar-none">
        <button
          type="button"
          class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap shrink-0 cursor-pointer"
          :class="
            activeTab === 'EXPENSES'
              ? 'bg-slate-900 text-white shadow-xs'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200/80'
          "
          @click="activeTab = 'EXPENSES'"
        >
          <Receipt class="w-4 h-4" />
          <span>Kas Kecil & Pengeluaran</span>
          <span
            class="px-2 py-0.5 rounded-full text-[10px] font-extrabold"
            :class="
              activeTab === 'EXPENSES' ? 'bg-white/20 text-white' : 'bg-slate-100 text-slate-700'
            "
          >
            {{ expenses.length }}
          </span>
        </button>

        <button
          type="button"
          class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap shrink-0 cursor-pointer"
          :class="
            activeTab === 'FINANCIAL'
              ? 'bg-slate-900 text-white shadow-xs'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200/80'
          "
          @click="activeTab = 'FINANCIAL'"
        >
          <BarChart3 class="w-4 h-4" />
          <span>Laporan Keuangan & Laba/Rugi (Milestone 5)</span>
        </button>

        <button
          type="button"
          class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap shrink-0 cursor-pointer"
          :class="
            activeTab === 'CRM'
              ? 'bg-slate-900 text-white shadow-xs'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200/80'
          "
          @click="activeTab = 'CRM'"
        >
          <Users class="w-4 h-4" />
          <span>CRM & Pelanggan Loyal (Milestone 5)</span>
        </button>
      </div>

      <!-- TAB 1: KAS KECIL & PENGELUARAN (Milestone 4) -->
      <div v-if="activeTab === 'EXPENSES'" class="space-y-6">
        <!-- 3. KPI Expense Summary Cards -->
        <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
          <!-- Card 1: Hari Ini -->
          <div
            class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
                >Pengeluaran Hari Ini</span
              >
              <div
                class="w-9 h-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center"
              >
                <Receipt class="w-5 h-5" />
              </div>
            </div>
            <div class="mt-3">
              <h3 class="text-lg sm:text-2xl font-black text-slate-900 tracking-tight">
                Rp {{ formatRupiah(summaryMetrics.todayTotal) }}
              </h3>
              <p class="text-[11px] text-slate-400 mt-1">
                {{ summaryMetrics.todayCount }} pengeluaran tercatat
              </p>
            </div>
          </div>

          <!-- Card 2: Periode Terpilih -->
          <div
            class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
                >Total Periode Ini</span
              >
              <div
                class="w-9 h-9 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center"
              >
                <TrendingDown class="w-5 h-5" />
              </div>
            </div>
            <div class="mt-3">
              <h3 class="text-lg sm:text-2xl font-black text-rose-600 tracking-tight">
                Rp {{ formatRupiah(summaryMetrics.rangeTotal) }}
              </h3>
              <p class="text-[11px] text-slate-400 mt-1">
                {{ filteredExpenses.length }} transaksi pengeluaran
              </p>
            </div>
          </div>

          <!-- Card 3: Kategori Terbesar -->
          <div
            class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
                >Kategori Terbesar</span
              >
              <div
                class="w-9 h-9 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center"
              >
                <PieChart class="w-5 h-5" />
              </div>
            </div>
            <div class="mt-3">
              <h3 class="text-base sm:text-lg font-black text-slate-900 truncate">
                {{ summaryMetrics.topCategoryName }}
              </h3>
              <p class="text-[11px] text-slate-500 font-semibold mt-1">
                Rp {{ formatRupiah(summaryMetrics.topCategoryAmount) }}
              </p>
            </div>
          </div>

          <!-- Card 4: Rata-Rata per Transaksi -->
          <div
            class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
                >Rata-Rata Kasbon/Biaya</span
              >
              <div
                class="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center"
              >
                <Wallet class="w-5 h-5" />
              </div>
            </div>
            <div class="mt-3">
              <h3 class="text-lg sm:text-2xl font-black text-slate-900 tracking-tight">
                Rp {{ formatRupiah(summaryMetrics.averageAmount) }}
              </h3>
              <p class="text-[11px] text-slate-400 mt-1">Per transaksi kas kecil</p>
            </div>
          </div>
        </div>

        <!-- 4. Category Visual Breakdown Cards -->
        <div class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs space-y-3">
          <div class="flex items-center justify-between">
            <h3
              class="text-xs font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
            >
              <PieChart class="w-4 h-4 text-brand-600" />
              <span>Komposisi Pengeluaran per Kategori</span>
            </h3>
            <span class="text-[11px] text-slate-400"
              >{{ categoryBreakdown.length }} Kategori Aktif</span
            >
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-4 gap-2.5">
            <div
              v-for="cat in categoryBreakdown"
              :key="cat.category"
              class="p-3 rounded-xl border border-slate-200/80 bg-slate-50/50 flex flex-col justify-between"
            >
              <div>
                <div class="flex items-center gap-1.5">
                  <span>{{ getCategoryIcon(cat.category) }}</span>
                  <span class="font-bold text-xs text-slate-800 truncate">
                    {{ getCategoryLabel(cat.category) }}
                  </span>
                </div>
                <span class="text-[10px] text-slate-400 block mt-0.5"
                  >{{ cat.count }} transaksi</span
                >
              </div>
              <div class="mt-2">
                <span class="text-xs sm:text-sm font-black text-slate-900 block">
                  Rp {{ formatRupiah(cat.total_amount) }}
                </span>
                <!-- Percentage Bar -->
                <div class="w-full bg-slate-200 h-1.5 rounded-full overflow-hidden mt-1.5">
                  <div
                    class="bg-amber-500 h-full rounded-full"
                    :style="{ width: `${cat.percent}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 5. Filter & Search Toolbar -->
        <div
          class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 bg-white p-3 sm:p-4 rounded-2xl border border-slate-100 shadow-xs"
        >
          <!-- Quick Date Range Pills -->
          <div class="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0 text-xs">
            <button
              v-for="p in periodPresets"
              :key="p.id"
              type="button"
              class="px-3 py-1.5 rounded-xl font-bold transition-all whitespace-nowrap cursor-pointer"
              :class="
                selectedPeriod === p.id
                  ? 'bg-slate-900 text-white shadow-xs'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              "
              @click="applyPeriod(p.id)"
            >
              {{ p.label }}
            </button>
          </div>

          <!-- Filters & Search -->
          <div class="flex items-center gap-2 flex-1 sm:justify-end">
            <!-- Category Filter -->
            <select
              v-model="selectedCategoryFilter"
              class="px-3 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-700 focus:outline-none focus:border-brand-500 cursor-pointer"
            >
              <option value="ALL">Semua Kategori</option>
              <option value="BAHAN_BAKU">🧴 Bahan Baku</option>
              <option value="GAS">🔥 Gas Dryer</option>
              <option value="LISTRIK_AIR">⚡ Listrik & Air</option>
              <option value="GAJI">👥 Gaji Pegawai</option>
              <option value="BONUS">🍱 Uang Makan</option>
              <option value="MAINTENANCE">🔧 Servis Mesin</option>
              <option value="SEWA">🏠 Sewa Tempat</option>
              <option value="LAIN_LAIN">📦 Lain-Lain</option>
            </select>

            <!-- Search Input -->
            <div class="relative min-w-[180px] max-w-xs flex-1">
              <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Cari keterangan..."
                class="w-full pl-9 pr-8 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-medium text-slate-800 placeholder-slate-400 focus:outline-none focus:border-brand-500 focus:bg-white transition-all"
              />
              <button
                v-if="searchQuery"
                class="absolute right-2.5 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                @click="searchQuery = ''"
              >
                <X class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>

        <!-- 6. Expense Table / List -->
        <div class="bg-white rounded-2xl border border-slate-100 shadow-xs overflow-hidden">
          <!-- Loading State -->
          <div v-if="loading" class="p-12 text-center text-slate-400 space-y-2">
            <RefreshCw class="w-6 h-6 animate-spin text-amber-600 mx-auto" />
            <p class="text-xs font-medium">Memuat data buku kas kecil...</p>
          </div>

          <!-- Empty State -->
          <div v-else-if="filteredExpenses.length === 0" class="p-12 text-center space-y-3">
            <div
              class="w-12 h-12 rounded-2xl bg-amber-50 text-amber-600 flex items-center justify-center mx-auto"
            >
              <Receipt class="w-6 h-6" />
            </div>
            <div>
              <h3 class="font-bold text-sm text-slate-800">Belum ada pengeluaran tercatat</h3>
              <p class="text-xs text-slate-400 max-w-sm mx-auto mt-1">
                Catat setiap biaya operasional harian kasir agar perhitungan laba bersih akurat.
              </p>
            </div>
            <button
              class="px-4 py-2 rounded-xl bg-amber-600 hover:bg-amber-700 text-white font-bold text-xs transition-all inline-flex items-center gap-1.5 cursor-pointer shadow-xs"
              @click="showCreateModal = true"
            >
              <PlusCircle class="w-4 h-4" />
              <span>+ Catat Pengeluaran Sekarang</span>
            </button>
          </div>

          <!-- Expenses Table -->
          <div v-else class="overflow-x-auto">
            <table class="w-full text-left text-xs text-slate-700">
              <thead
                class="bg-slate-50 border-b border-slate-100 text-slate-500 font-bold uppercase tracking-wider text-[10px]"
              >
                <tr>
                  <th class="py-3 px-4">Tanggal</th>
                  <th class="py-3 px-4">Kategori</th>
                  <th class="py-3 px-4">Keterangan Keperluan</th>
                  <th class="py-3 px-4">Dicatat Oleh</th>
                  <th class="py-3 px-4 text-right">Nominal (Rp)</th>
                  <th class="py-3 px-4 text-center">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="item in filteredExpenses"
                  :key="item.id"
                  class="hover:bg-slate-50/80 transition-colors group"
                >
                  <!-- Tanggal -->
                  <td class="py-3.5 px-4 font-semibold text-slate-900 whitespace-nowrap">
                    {{ formatDate(item.expense_date) }}
                  </td>

                  <!-- Kategori Badge -->
                  <td class="py-3.5 px-4 whitespace-nowrap">
                    <span
                      class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold"
                      :class="getCategoryBadgeClass(item.category)"
                    >
                      <span>{{ getCategoryIcon(item.category) }}</span>
                      <span>{{ getCategoryLabel(item.category) }}</span>
                    </span>
                  </td>

                  <!-- Keterangan -->
                  <td class="py-3.5 px-4">
                    <p class="font-bold text-slate-800">{{ item.description }}</p>
                  </td>

                  <!-- Staf -->
                  <td class="py-3.5 px-4 text-slate-500 whitespace-nowrap font-medium">
                    {{ item.created_by || 'Kasir' }}
                  </td>

                  <!-- Nominal -->
                  <td class="py-3.5 px-4 text-right whitespace-nowrap">
                    <span class="font-black text-sm text-rose-600">
                      - Rp {{ formatRupiah(item.amount) }}
                    </span>
                  </td>

                  <!-- Actions -->
                  <td class="py-3.5 px-4 text-center whitespace-nowrap">
                    <button
                      class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-all cursor-pointer"
                      title="Hapus Pengeluaran"
                      @click="promptDeleteExpense(item)"
                    >
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- TAB 2: FINANCIAL REPORT (Milestone 5 Placeholder Preview) -->
      <div
        v-else-if="activeTab === 'FINANCIAL'"
        class="bg-white rounded-2xl border border-slate-100 p-8 text-center space-y-4 shadow-xs"
      >
        <div
          class="w-12 h-12 rounded-2xl bg-brand-50 text-brand-600 flex items-center justify-center mx-auto"
        >
          <BarChart3 class="w-6 h-6" />
        </div>
        <h2 class="text-xl font-bold text-slate-900">Laporan Keuangan & Laba/Rugi</h2>
        <p class="text-sm text-slate-500 max-w-md mx-auto">
          Grafik tren omset harian, rekap laba bersih (Pendapatan Kas - Pengeluaran), dan laporan
          piutang macet akan diaktifkan pada <strong>Milestone 5</strong>.
        </p>
      </div>

      <!-- TAB 3: CRM (Milestone 5 Placeholder Preview) -->
      <div
        v-else-if="activeTab === 'CRM'"
        class="bg-white rounded-2xl border border-slate-100 p-8 text-center space-y-4 shadow-xs"
      >
        <div
          class="w-12 h-12 rounded-2xl bg-purple-50 text-purple-600 flex items-center justify-center mx-auto"
        >
          <Users class="w-6 h-6" />
        </div>
        <h2 class="text-xl font-bold text-slate-900">CRM & Peringkat Pelanggan Loyal</h2>
        <p class="text-sm text-slate-500 max-w-md mx-auto">
          Top 10 pelanggan paling sering cuci, leaderboard omset pelanggan, dan histori transaksi
          seumur hidup pelanggan akan aktif di <strong>Milestone 5</strong>.
        </p>
      </div>

      <!-- Modal Tambah Pengeluaran -->
      <QuickExpenseModal v-model="showCreateModal" @expense-added="loadExpensesData" />

      <!-- Sub-Modal Konfirmasi Hapus Pengeluaran -->
      <div
        v-if="showDeleteModal && expenseToDelete"
        class="fixed inset-0 z-60 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
      >
        <div
          class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4 text-slate-800"
        >
          <div class="flex items-center gap-2.5 text-rose-600">
            <div class="w-9 h-9 rounded-xl bg-rose-50 flex items-center justify-center">
              <Trash2 class="w-5 h-5" />
            </div>
            <div>
              <h4 class="font-black text-sm text-slate-900">Hapus Catatan Pengeluaran?</h4>
              <p class="text-xs text-slate-500">Tindakan ini tidak dapat dibatalkan</p>
            </div>
          </div>

          <div class="p-3 rounded-xl bg-slate-50 border border-slate-200 text-xs space-y-1">
            <p>
              <span class="text-slate-500">Keterangan:</span>
              <strong>{{ expenseToDelete.description }}</strong>
            </p>
            <p>
              <span class="text-slate-500">Nominal:</span>
              <strong class="text-rose-600">Rp {{ formatRupiah(expenseToDelete.amount) }}</strong>
            </p>
            <p>
              <span class="text-slate-500">Tanggal:</span>
              {{ formatDate(expenseToDelete.expense_date) }}
            </p>
          </div>

          <div class="flex items-center gap-2 pt-1">
            <button
              type="button"
              class="w-1/2 py-2.5 rounded-xl border border-slate-200 text-slate-700 font-bold text-xs hover:bg-slate-50 cursor-pointer"
              @click="showDeleteModal = false"
            >
              Batal
            </button>
            <button
              type="button"
              :disabled="deleting"
              class="w-1/2 py-2.5 rounded-xl bg-rose-600 hover:bg-rose-700 text-white font-bold text-xs transition-all disabled:opacity-50 cursor-pointer"
              @click="confirmDeleteExpense"
            >
              <span v-if="deleting">Menghapus...</span>
              <span v-else>Hapus</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import {
  Receipt,
  PlusCircle,
  RefreshCw,
  Search,
  X,
  Trash2,
  TrendingDown,
  PieChart,
  Wallet,
  BarChart3,
  Users,
} from 'lucide-vue-next'
import AppLayout from '../components/layout/AppLayout.vue'
import QuickExpenseModal from '../components/modals/QuickExpenseModal.vue'
import { api, formatRupiah, formatDate } from '../api/client'
import { useToast } from '../composables/useToast'

const route = useRoute()
const toast = useToast()

const activeTab = ref('EXPENSES')
const loading = ref(false)
const deleting = ref(false)
const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const expenseToDelete = ref(null)

const expenses = ref([])
const categoriesSummary = ref([])

const searchQuery = ref('')
const selectedCategoryFilter = ref('ALL')
const selectedPeriod = ref('THIS_MONTH')

const dateRange = reactive({
  start: '',
  end: '',
})

const periodPresets = [
  { id: 'TODAY', label: 'Hari Ini' },
  { id: 'LAST_7_DAYS', label: '7 Hari Terakhir' },
  { id: 'THIS_MONTH', label: 'Bulan Ini' },
  { id: 'ALL', label: 'Semua' },
]

const pad = (n) => String(n).padStart(2, '0')
const toDateStr = (d) => `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`

const applyPeriod = (presetId) => {
  selectedPeriod.value = presetId
  const now = new Date()

  if (presetId === 'TODAY') {
    dateRange.start = toDateStr(now)
    dateRange.end = toDateStr(now)
  } else if (presetId === 'LAST_7_DAYS') {
    const past = new Date()
    past.setDate(past.getDate() - 7)
    dateRange.start = toDateStr(past)
    dateRange.end = toDateStr(now)
  } else if (presetId === 'THIS_MONTH') {
    const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
    dateRange.start = toDateStr(firstDay)
    dateRange.end = toDateStr(now)
  } else {
    // ALL
    dateRange.start = '2026-01-01'
    dateRange.end = toDateStr(now)
  }

  loadExpensesData()
}

const loadExpensesData = async () => {
  loading.value = true
  try {
    const [expList, summary] = await Promise.all([
      api.get('/expenses', { start: dateRange.start, end: dateRange.end }),
      api.get('/expenses/summary', { start: dateRange.start, end: dateRange.end }),
    ])
    expenses.value = expList || []
    categoriesSummary.value = summary || []
  } catch (err) {
    console.warn('Failed loading expenses:', err)
  } finally {
    loading.value = false
  }
}

const filteredExpenses = computed(() => {
  return expenses.value.filter((item) => {
    // Category filter
    if (selectedCategoryFilter.value !== 'ALL' && item.category !== selectedCategoryFilter.value) {
      return false
    }

    // Search query
    if (searchQuery.value.trim()) {
      const q = searchQuery.value.toLowerCase().trim()
      const matchDesc = item.description?.toLowerCase().includes(q)
      const matchStaff = item.created_by?.toLowerCase().includes(q)
      if (!matchDesc && !matchStaff) return false
    }

    return true
  })
})

const summaryMetrics = computed(() => {
  const nowStr = toDateStr(new Date())

  const todayList = expenses.value.filter((e) => e.expense_date === nowStr)
  const todayTotal = todayList.reduce((acc, it) => acc + (Number(it.amount) || 0), 0)
  const todayCount = todayList.length

  const rangeTotal = filteredExpenses.value.reduce((acc, it) => acc + (Number(it.amount) || 0), 0)
  const averageAmount =
    filteredExpenses.value.length > 0 ? Math.round(rangeTotal / filteredExpenses.value.length) : 0

  // Find top category
  let topCategoryName = '-'
  let topCategoryAmount = 0
  if (categoriesSummary.value.length > 0) {
    const top = categoriesSummary.value[0]
    topCategoryName = getCategoryLabel(top.category)
    topCategoryAmount = Number(top.total_amount) || 0
  }

  return {
    todayTotal,
    todayCount,
    rangeTotal,
    averageAmount,
    topCategoryName,
    topCategoryAmount,
  }
})

const categoryBreakdown = computed(() => {
  const total = categoriesSummary.value.reduce((acc, c) => acc + (Number(c.total_amount) || 0), 0)

  return categoriesSummary.value.map((c) => {
    const amt = Number(c.total_amount) || 0
    const pct = total > 0 ? Math.round((amt / total) * 100) : 0
    return {
      ...c,
      percent: pct,
    }
  })
})

const getCategoryLabel = (cat) => {
  switch (cat) {
    case 'BAHAN_BAKU':
      return 'Bahan Baku'
    case 'GAS':
      return 'Gas Dryer'
    case 'LISTRIK_AIR':
      return 'Listrik & Air'
    case 'GAJI':
      return 'Gaji Pegawai'
    case 'BONUS':
      return 'Uang Makan'
    case 'MAINTENANCE':
      return 'Servis Mesin'
    case 'SEWA':
      return 'Sewa Tempat'
    case 'LAIN_LAIN':
      return 'Lain-Lain'
    default:
      return cat || 'Operasional'
  }
}

const getCategoryIcon = (cat) => {
  switch (cat) {
    case 'BAHAN_BAKU':
      return '🧴'
    case 'GAS':
      return '🔥'
    case 'LISTRIK_AIR':
      return '⚡'
    case 'GAJI':
      return '👥'
    case 'BONUS':
      return '🍱'
    case 'MAINTENANCE':
      return '🔧'
    case 'SEWA':
      return '🏠'
    default:
      return '📦'
  }
}

const getCategoryBadgeClass = (cat) => {
  switch (cat) {
    case 'BAHAN_BAKU':
      return 'bg-blue-50 text-blue-700 border border-blue-200'
    case 'GAS':
      return 'bg-orange-50 text-orange-700 border border-orange-200'
    case 'LISTRIK_AIR':
      return 'bg-yellow-50 text-yellow-700 border border-yellow-200'
    case 'GAJI':
      return 'bg-purple-50 text-purple-700 border border-purple-200'
    case 'BONUS':
      return 'bg-pink-50 text-pink-700 border border-pink-200'
    case 'MAINTENANCE':
      return 'bg-amber-50 text-amber-700 border border-amber-200'
    case 'SEWA':
      return 'bg-indigo-50 text-indigo-700 border border-indigo-200'
    default:
      return 'bg-slate-100 text-slate-700 border border-slate-200'
  }
}

const promptDeleteExpense = (item) => {
  expenseToDelete.value = item
  showDeleteModal.value = true
}

const confirmDeleteExpense = async () => {
  if (!expenseToDelete.value) return
  deleting.value = true
  try {
    await api.delete(`/expenses/${expenseToDelete.value.id}`)
    toast.success('Catatan pengeluaran berhasil dihapus.')
    showDeleteModal.value = false
    expenseToDelete.value = null
    loadExpensesData()
  } catch (err) {
    toast.error(err.message || 'Gagal menghapus pengeluaran.')
  } finally {
    deleting.value = false
  }
}

watch(
  () => route.query.tab,
  (newTab) => {
    if (newTab === 'expenses') activeTab.value = 'EXPENSES'
    else if (newTab === 'financial') activeTab.value = 'FINANCIAL'
    else if (newTab === 'crm') activeTab.value = 'CRM'
  },
  { immediate: true }
)

onMounted(() => {
  applyPeriod('THIS_MONTH')
})
</script>
