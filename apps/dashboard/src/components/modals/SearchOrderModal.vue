<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-start justify-center p-4 pt-16 sm:pt-24 bg-slate-900/50 backdrop-blur-xs animate-in fade-in"
  >
    <div
      class="bg-white rounded-2xl max-w-xl w-full p-5 shadow-2xl border border-slate-100 flex flex-col max-h-[80vh]"
    >
      <!-- Search Input Header -->
      <div class="relative flex items-center pb-3 border-b border-slate-100">
        <Search class="w-5 h-5 text-slate-400 absolute left-3" />
        <input
          v-model="query"
          type="text"
          autofocus
          placeholder="Cari order (#2611..., nama pelanggan, no WA, atau nomor rak)..."
          class="w-full pl-10 pr-10 py-2.5 rounded-xl bg-slate-50 border border-slate-200 text-sm font-medium focus:outline-hidden focus:ring-2 focus:ring-brand-500 focus:bg-white transition-all text-slate-900"
          @input="onSearchInput"
        />
        <button
          class="p-1 rounded-lg text-slate-400 hover:text-slate-600 absolute right-2.5"
          @click="close"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Results Body -->
      <div class="flex-1 overflow-y-auto py-3 space-y-2.5">
        <div
          v-if="loading"
          class="text-center py-8 text-slate-400 text-xs flex items-center justify-center gap-2"
        >
          <Loader2 class="w-4 h-4 animate-spin text-brand-600" />
          <span>Mencari transaksi...</span>
        </div>

        <div
          v-else-if="results.length === 0 && query.trim()"
          class="text-center py-8 text-slate-400 text-xs"
        >
          Tidak ditemukan transaksi yang cocok dengan "<span class="font-bold text-slate-700">{{
            query
          }}</span
          >".
        </div>

        <div v-else-if="results.length === 0" class="text-center py-8 text-slate-400 text-xs">
          Ketik minimal 2 karakter untuk mencari transaksi.
        </div>

        <!-- Order Result Card -->
        <div
          v-for="order in results"
          :key="order.id"
          class="p-3 rounded-xl border border-slate-100 hover:border-brand-300 hover:bg-brand-50/40 transition-all cursor-pointer flex items-center justify-between group"
          @click="selectOrder(order)"
        >
          <div class="space-y-1">
            <div class="flex items-center gap-2">
              <span class="font-extrabold text-sm text-brand-700 font-mono">{{
                formatDisplayOrder(order.order_number)
              }}</span>
              <span
                class="text-[10px] font-bold px-2 py-0.5 rounded-full uppercase"
                :class="getStatusClass(order.status)"
              >
                {{ order.status }}
              </span>
              <span
                v-if="order.rack_number"
                class="text-[10px] font-bold px-1.5 py-0.5 rounded-md bg-amber-50 text-amber-700 border border-amber-200"
              >
                Rak: {{ order.rack_number }}
              </span>
            </div>
            <p class="text-xs font-semibold text-slate-800">
              {{ order.customer_name }}
              <span class="text-slate-400 font-normal">({{ order.customer_phone }})</span>
            </p>
            <p class="text-[11px] text-slate-500">
              {{ order.primary_service_name || 'Layanan Laundry' }} •
              {{ formatDate(order.created_at) }}
            </p>
          </div>

          <div class="text-right flex flex-col items-end">
            <span class="font-extrabold text-sm text-slate-900"
              >Rp {{ formatRupiah(order.total_price) }}</span
            >
            <span
              class="text-[10px] font-bold px-1.5 py-0.5 rounded-full mt-1"
              :class="
                order.payment_status === 'PAID'
                  ? 'bg-emerald-50 text-emerald-700'
                  : 'bg-rose-50 text-rose-700'
              "
            >
              {{ order.payment_status === 'PAID' ? 'LUNAS' : 'BELUM LUNAS' }}
            </span>
          </div>
        </div>
      </div>

      <!-- Footer Info -->
      <div
        class="pt-3 border-t border-slate-100 flex items-center justify-between text-xs text-slate-400"
      >
        <span>Tekan Esc untuk menutup</span>
        <button class="text-xs font-semibold text-slate-600 hover:text-slate-900" @click="close">
          Tutup
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Search, X, Loader2 } from 'lucide-vue-next'
import { api, formatRupiah, formatDate } from '../../api/client'

defineProps({
  modelValue: Boolean,
})

const emit = defineEmits(['update:modelValue'])
const router = useRouter()

const query = ref('')
const loading = ref(false)
const results = ref([])
let searchTimeout = null

const close = () => {
  emit('update:modelValue', false)
  query.value = ''
  results.value = []
}

const onSearchInput = () => {
  clearTimeout(searchTimeout)
  if (!query.value.trim() || query.value.trim().length < 2) {
    results.value = []
    loading.value = false
    return
  }

  loading.value = true
  searchTimeout = setTimeout(async () => {
    try {
      const data = await api.get('/pos/orders', { search: query.value.trim() })
      results.value = data || []
    } catch {
      results.value = []
    } finally {
      loading.value = false
    }
  }, 300)
}

const selectOrder = (order) => {
  close()
  router.push({ path: '/orders', query: { id: order.id } })
}

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
    case 'READY':
    case 'SIAP AMBIL':
      return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
    case 'SELESAI':
      return 'bg-slate-100 text-slate-700 border border-slate-200'
    case 'BATAL':
      return 'bg-rose-50 text-rose-700 border border-rose-200'
    default:
      return 'bg-slate-50 text-slate-600'
  }
}
</script>
