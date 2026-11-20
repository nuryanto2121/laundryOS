<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
  >
    <div
      class="bg-white rounded-2xl max-w-2xl w-full p-4 sm:p-6 shadow-2xl border border-slate-100 flex flex-col max-h-[90vh]"
    >
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-100 shrink-0">
        <div class="flex items-center gap-2.5">
          <div
            class="w-10 h-10 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center shadow-xs"
          >
            <Wallet class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-black text-lg text-slate-900">Rincian Omset & Penjualan Hari Ini</h3>
            <p class="text-xs text-slate-500">{{ todayFormatted }} • Arus kas & metode bayar</p>
          </div>
        </div>
        <button
          class="p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-all cursor-pointer"
          @click="close"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Scrollable Body -->
      <div class="flex-1 overflow-y-auto py-3 space-y-4 pr-1 text-slate-800">
        <!-- Top Summary Cards -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
          <!-- Total Omset Bruto -->
          <div class="p-3.5 rounded-xl bg-slate-50 border border-slate-200">
            <span class="text-[11px] font-bold text-slate-500 uppercase block"
              >Total Omset Bruto</span
            >
            <span class="text-lg font-black text-slate-900 block mt-1">
              Rp {{ formatRupiah(summary.totalOmset) }}
            </span>
            <span class="text-[10px] text-slate-400">Total nilai seluruh order masuk hari ini</span>
          </div>

          <!-- Kas Diterima (Lunas) -->
          <div class="p-3.5 rounded-xl bg-emerald-50/70 border border-emerald-200">
            <span class="text-[11px] font-bold text-emerald-800 uppercase block"
              >Kas Diterima (Lunas)</span
            >
            <span class="text-lg font-black text-emerald-700 block mt-1">
              Rp {{ formatRupiah(summary.totalPaid) }}
            </span>
            <span class="text-[10px] text-emerald-600">Uang masuk ke kasir hari ini</span>
          </div>

          <!-- Belum Lunas (Piutang) -->
          <div class="p-3.5 rounded-xl bg-amber-50/70 border border-amber-200">
            <span class="text-[11px] font-bold text-amber-800 uppercase block"
              >Belum Lunas (Piutang)</span
            >
            <span class="text-lg font-black text-amber-700 block mt-1">
              Rp {{ formatRupiah(summary.totalUnpaid) }}
            </span>
            <span class="text-[10px] text-amber-600">Bayar saat ambil cucian</span>
          </div>
        </div>

        <!-- Breakdown per Metode Pembayaran -->
        <div class="space-y-2">
          <h4 class="text-xs font-black text-slate-700 uppercase tracking-wider">
            Realisasi Kas per Metode Pembayaran
          </h4>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <!-- Tunai -->
            <div
              class="p-3 rounded-xl border border-slate-200 bg-white flex items-center justify-between"
            >
              <div class="flex items-center gap-2.5">
                <div
                  class="w-8 h-8 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center font-bold text-xs"
                >
                  💵
                </div>
                <div>
                  <span class="text-xs font-bold text-slate-800 block">Tunai (Cash)</span>
                  <span class="text-[10px] text-slate-400"
                    >{{ summary.methodCounts.TUNAI || 0 }} Transaksi</span
                  >
                </div>
              </div>
              <span class="font-extrabold text-xs text-slate-900">
                Rp {{ formatRupiah(summary.methods.TUNAI) }}
              </span>
            </div>

            <!-- Transfer BCA -->
            <div
              class="p-3 rounded-xl border border-slate-200 bg-white flex items-center justify-between"
            >
              <div class="flex items-center gap-2.5">
                <div
                  class="w-8 h-8 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center font-bold text-xs"
                >
                  🏦
                </div>
                <div>
                  <span class="text-xs font-bold text-slate-800 block">Transfer BCA</span>
                  <span class="text-[10px] text-slate-400"
                    >{{ summary.methodCounts.BCA || 0 }} Transaksi</span
                  >
                </div>
              </div>
              <span class="font-extrabold text-xs text-slate-900">
                Rp {{ formatRupiah(summary.methods.BCA) }}
              </span>
            </div>

            <!-- QRIS AstraPay -->
            <div
              class="p-3 rounded-xl border border-slate-200 bg-white flex items-center justify-between"
            >
              <div class="flex items-center gap-2.5">
                <div
                  class="w-8 h-8 rounded-lg bg-purple-50 text-purple-600 flex items-center justify-center font-bold text-xs"
                >
                  📱
                </div>
                <div>
                  <span class="text-xs font-bold text-slate-800 block">QRIS AstraPay</span>
                  <span class="text-[10px] text-slate-400"
                    >{{ summary.methodCounts.QRIS || 0 }} Transaksi</span
                  >
                </div>
              </div>
              <span class="font-extrabold text-xs text-slate-900">
                Rp {{ formatRupiah(summary.methods.QRIS) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Daftar Transaksi Hari Ini -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <h4 class="text-xs font-black text-slate-700 uppercase tracking-wider">
              Daftar Transaksi Hari Ini ({{ todayOrders.length }})
            </h4>
            <span class="text-[11px] text-slate-400">Klik transaksi untuk melihat detail</span>
          </div>

          <div
            v-if="loading"
            class="py-8 text-center text-xs text-slate-400 flex items-center justify-center gap-2"
          >
            <Loader2 class="w-4 h-4 animate-spin text-brand-600" />
            <span>Memuat rincian transaksi...</span>
          </div>

          <div
            v-else-if="todayOrders.length === 0"
            class="p-6 text-center text-xs text-slate-400 border border-dashed border-slate-200 rounded-xl"
          >
            Belum ada transaksi tercatat hari ini.
          </div>

          <div
            v-else
            class="divide-y divide-slate-100 rounded-xl border border-slate-200 bg-white overflow-hidden text-xs"
          >
            <div
              v-for="o in todayOrders"
              :key="o.id"
              class="p-3 hover:bg-slate-50 transition-colors flex items-center justify-between cursor-pointer group"
              @click="openOrder(o.id)"
            >
              <div class="flex items-center gap-2.5">
                <div
                  class="w-8 h-8 rounded-lg flex items-center justify-center font-bold text-[10px]"
                  :class="
                    o.payment_status === 'PAID'
                      ? 'bg-emerald-50 text-emerald-700'
                      : 'bg-amber-50 text-amber-700'
                  "
                >
                  {{ o.payment_status === 'PAID' ? 'LUNAS' : 'TEMPO' }}
                </div>
                <div>
                  <div class="flex items-center gap-1.5">
                    <span
                      class="font-extrabold text-slate-900 font-mono group-hover:text-brand-600 transition-colors"
                    >
                      {{ formatDisplayOrder(o.order_number) }}
                    </span>
                    <span class="text-[10px] text-slate-400">• {{ formatTime(o.created_at) }}</span>
                  </div>
                  <p class="text-[11px] text-slate-600">
                    <span class="font-bold text-slate-800">{{ o.customer_name }}</span>
                    <span v-if="o.primary_service_name" class="text-slate-400 ml-1">
                      ({{ o.primary_service_name }})
                    </span>
                  </p>
                </div>
              </div>

              <div class="text-right">
                <span class="font-black text-xs text-slate-900 block">
                  Rp {{ formatRupiah(o.total_price) }}
                </span>
                <span
                  class="text-[10px] font-semibold"
                  :class="o.payment_status === 'PAID' ? 'text-emerald-600' : 'text-amber-600'"
                >
                  {{
                    o.payment_status === 'PAID' ? o.payment_method_name || 'Lunas' : 'Belum Lunas'
                  }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Footer -->
      <div class="pt-3 border-t border-slate-100 flex items-center justify-between gap-2 shrink-0">
        <button
          type="button"
          class="px-4 py-2.5 rounded-xl border border-slate-200 text-slate-700 font-bold text-xs hover:bg-slate-50 cursor-pointer"
          @click="close"
        >
          Tutup
        </button>

        <button
          type="button"
          class="px-4 py-2.5 rounded-xl bg-slate-900 hover:bg-slate-800 text-white font-bold text-xs flex items-center gap-1.5 transition-all cursor-pointer"
          @click="goToReports"
        >
          <span>Buka Laporan Keuangan Lengkap</span>
          <ArrowRight class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Wallet, X, ArrowRight, Loader2 } from 'lucide-vue-next'
import { api, formatRupiah } from '../../api/client'

const props = defineProps({
  modelValue: Boolean,
  metrics: {
    type: Object,
    default: () => ({}),
  },
})

const emit = defineEmits(['update:modelValue'])
const router = useRouter()

const loading = ref(false)
const todayOrders = ref([])

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

const formatTime = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

const summary = computed(() => {
  let totalOmset = 0
  let totalPaid = 0
  let totalUnpaid = 0

  const methods = {
    TUNAI: 0,
    BCA: 0,
    QRIS: 0,
  }

  const methodCounts = {
    TUNAI: 0,
    BCA: 0,
    QRIS: 0,
  }

  todayOrders.value.forEach((o) => {
    const price = Number(o.total_price) || 0
    const paid = Number(o.paid_amount) || 0

    totalOmset += price
    totalPaid += paid
    if (o.payment_status !== 'PAID') {
      totalUnpaid += Math.max(0, price - paid)
    }

    if (paid > 0) {
      const m = (o.payment_method_name || '').toUpperCase()
      if (m.includes('TUNAI') || m.includes('CASH')) {
        methods.TUNAI += paid
        methodCounts.TUNAI++
      } else if (m.includes('BCA') || m.includes('TRANSFER')) {
        methods.BCA += paid
        methodCounts.BCA++
      } else if (m.includes('QRIS') || m.includes('ASTRAPAY')) {
        methods.QRIS += paid
        methodCounts.QRIS++
      } else {
        methods.TUNAI += paid
        methodCounts.TUNAI++
      }
    }
  })

  return {
    totalOmset: totalOmset || props.metrics?.omset_hari_ini || 0,
    totalPaid,
    totalUnpaid,
    methods,
    methodCounts,
  }
})

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

const loadTodayOrders = async () => {
  loading.value = true
  try {
    const data = await api.get('/pos/orders')
    const all = data || []
    todayOrders.value = all.filter((o) => isToday(o.created_at) || isToday(o.paid_at))
  } catch (err) {
    console.warn('Failed loading today orders:', err)
  } finally {
    loading.value = false
  }
}

const openOrder = (orderId) => {
  close()
  router.push({ path: '/orders', query: { id: orderId } })
}

const goToReports = () => {
  close()
  router.push('/reports')
}

const close = () => {
  emit('update:modelValue', false)
}

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      loadTodayOrders()
    }
  },
  { immediate: true }
)
</script>
