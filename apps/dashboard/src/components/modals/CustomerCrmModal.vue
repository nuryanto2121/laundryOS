<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
    @click.self="close"
  >
    <div
      class="bg-white rounded-3xl max-w-xl w-full p-6 shadow-2xl border border-slate-100 space-y-5 text-slate-800 relative animate-in zoom-in-95 duration-200 max-h-[90vh] overflow-y-auto"
    >
      <!-- Close Button -->
      <button
        class="absolute right-5 top-5 p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-all cursor-pointer"
        @click="close"
      >
        <X class="w-5 h-5" />
      </button>

      <!-- Header: Customer Identity -->
      <div class="flex items-start gap-4 pr-10">
        <div
          class="w-14 h-14 rounded-2xl bg-gradient-to-br from-purple-500 to-indigo-600 text-white flex items-center justify-center font-black text-xl shadow-md shadow-purple-500/20 shrink-0"
        >
          {{ getInitials(customerData?.full_name || 'Pelanggan') }}
        </div>
        <div class="min-w-0">
          <div class="flex items-center gap-2">
            <h2 class="text-xl font-black text-slate-900 truncate">
              {{ customerData?.full_name || 'Profil Pelanggan' }}
            </h2>
            <span
              v-if="customerData?.customer_code"
              class="px-2 py-0.5 rounded-md bg-purple-50 text-purple-700 font-extrabold text-[10px] tracking-wider uppercase border border-purple-200/60"
            >
              {{ customerData.customer_code }}
            </span>
          </div>
          <p class="text-xs text-slate-500 mt-1 flex items-center gap-1.5">
            <Phone class="w-3.5 h-3.5 text-slate-400" />
            <span>{{ customerData?.phone || '-' }}</span>
          </p>
          <p
            v-if="customerData?.address"
            class="text-xs text-slate-500 mt-0.5 flex items-center gap-1.5 truncate"
          >
            <MapPin class="w-3.5 h-3.5 text-slate-400 shrink-0" />
            <span class="truncate">{{ customerData.address }}</span>
          </p>
        </div>
      </div>

      <!-- Mode Switcher: Tab Profil LTV vs Aktivitas Periode -->
      <div class="flex items-center gap-2 border-b border-slate-100 pb-2 text-xs font-bold">
        <button
          type="button"
          class="px-3 py-1.5 rounded-xl transition-all cursor-pointer"
          :class="
            modalTab === 'LIFETIME'
              ? 'bg-purple-100 text-purple-800'
              : 'text-slate-500 hover:text-slate-800'
          "
          @click="modalTab = 'LIFETIME'"
        >
          Profil & Seumur Hidup (LTV)
        </button>
        <button
          type="button"
          class="px-3 py-1.5 rounded-xl transition-all cursor-pointer flex items-center gap-1.5"
          :class="
            modalTab === 'PERIOD'
              ? 'bg-purple-100 text-purple-800'
              : 'text-slate-500 hover:text-slate-800'
          "
          @click="modalTab = 'PERIOD'"
        >
          <span>Aktivitas Periode</span>
          <span
            class="px-1.5 py-0.2 rounded-full text-[10px]"
            :class="
              modalTab === 'PERIOD'
                ? 'bg-purple-200 text-purple-900'
                : 'bg-slate-100 text-slate-600'
            "
          >
            {{ customerPeriodOrders.length }}
          </span>
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="py-12 flex flex-col items-center justify-center gap-3">
        <Loader2 class="w-8 h-8 text-purple-600 animate-spin" />
        <span class="text-xs text-slate-400 font-medium">Memuat statistik CRM pelanggan...</span>
      </div>

      <!-- 1. TAB: LIFETIME STATS -->
      <div v-else-if="modalTab === 'LIFETIME'" class="space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <!-- Total Spend -->
          <div
            class="p-4 rounded-2xl bg-slate-50 border border-slate-100 flex flex-col justify-between"
          >
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider"
              >Total Belanja (LTV)</span
            >
            <div class="mt-2">
              <span class="text-lg sm:text-xl font-black text-emerald-600 tracking-tight block">
                Rp {{ formatRupiah(profile?.lifetime_spend) }}
              </span>
              <span class="text-[10px] text-slate-400 block mt-0.5">Seumur hidup</span>
            </div>
          </div>

          <!-- Total Orders -->
          <div
            class="p-4 rounded-2xl bg-slate-50 border border-slate-100 flex flex-col justify-between"
          >
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider"
              >Frekuensi Order</span
            >
            <div class="mt-2">
              <span class="text-lg sm:text-xl font-black text-slate-900 tracking-tight block">
                {{ profile?.total_orders || 0 }}x
              </span>
              <span class="text-[10px] text-slate-400 block mt-0.5">Total transaksi</span>
            </div>
          </div>

          <!-- Total Weight -->
          <div
            class="p-4 rounded-2xl bg-slate-50 border border-slate-100 flex flex-col justify-between"
          >
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider"
              >Volume Cucian</span
            >
            <div class="mt-2">
              <span class="text-lg sm:text-xl font-black text-brand-600 tracking-tight block">
                {{ formatWeight(profile?.lifetime_weight) }} Kg
              </span>
              <span class="text-[10px] text-slate-400 block mt-0.5">Akumulasi kiloan</span>
            </div>
          </div>

          <!-- Total Unpaid -->
          <div
            class="p-4 rounded-2xl bg-slate-50 border border-slate-100 flex flex-col justify-between"
          >
            <span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider"
              >Sisa Piutang</span
            >
            <div class="mt-2">
              <span
                class="text-lg sm:text-xl font-black tracking-tight block"
                :class="Number(profile?.total_unpaid) > 0 ? 'text-rose-600' : 'text-slate-700'"
              >
                Rp {{ formatRupiah(profile?.total_unpaid) }}
              </span>
              <span
                class="text-[10px] block mt-0.5"
                :class="
                  Number(profile?.total_unpaid) > 0 ? 'text-rose-500 font-bold' : 'text-slate-400'
                "
              >
                {{ Number(profile?.total_unpaid) > 0 ? 'Perlu ditagih' : 'Lunas semua' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Member Since Info -->
        <div
          class="p-3.5 rounded-xl bg-purple-50/60 border border-purple-100 flex items-center justify-between text-xs"
        >
          <div class="flex items-center gap-2 text-purple-900 font-semibold">
            <Calendar class="w-4 h-4 text-purple-600" />
            <span>Terdaftar Sejak</span>
          </div>
          <span class="font-bold text-purple-700">
            {{ formatDate(customerData?.created_at || profile?.created_at) }}
          </span>
        </div>
      </div>

      <!-- 2. TAB: AKTIVITAS PERIODE TERPILIH -->
      <div v-else class="space-y-4">
        <!-- Date indicator banner -->
        <div
          class="p-3 bg-purple-50 rounded-xl border border-purple-100 flex items-center justify-between text-xs"
        >
          <span class="text-purple-900 font-bold"
            >Periode: {{ formatDate(startDate) }} &mdash; {{ formatDate(endDate) }}</span
          >
          <span class="text-purple-700 font-semibold"
            >{{ customerPeriodOrders.length }} transaksi ditemukan</span
          >
        </div>

        <!-- Mini Period KPI Summary -->
        <div class="grid grid-cols-3 gap-2.5 text-center">
          <div class="p-3 bg-slate-50 rounded-xl border border-slate-100">
            <span class="text-[10px] text-slate-400 font-bold uppercase block">Belanja</span>
            <span class="text-sm font-black text-emerald-600"
              >Rp {{ formatRupiah(periodSpend) }}</span
            >
          </div>
          <div class="p-3 bg-slate-50 rounded-xl border border-slate-100">
            <span class="text-[10px] text-slate-400 font-bold uppercase block">Berat</span>
            <span class="text-sm font-black text-brand-600">{{ periodWeight.toFixed(1) }} Kg</span>
          </div>
          <div class="p-3 bg-slate-50 rounded-xl border border-slate-100">
            <span class="text-[10px] text-slate-400 font-bold uppercase block">Layanan</span>
            <span class="text-sm font-black text-purple-600"
              >{{ periodServices.length }} Jenis</span
            >
          </div>
        </div>

        <!-- Services used badges in this period -->
        <div v-if="periodServices.length > 0" class="flex items-center gap-1.5 flex-wrap">
          <span class="text-xs text-slate-400 font-medium">Layanan Dipakai:</span>
          <span
            v-for="s in periodServices"
            :key="s"
            class="px-2 py-0.5 rounded-md bg-purple-100 text-purple-800 text-[10px] font-bold"
          >
            {{ formatCategoryName(s) }}
          </span>
        </div>

        <!-- Orders list in this period -->
        <div
          v-if="customerPeriodOrders.length === 0"
          class="py-6 text-center text-slate-400 text-xs"
        >
          Pelanggan tidak memiliki transaksi pada rentang tanggal ini.
        </div>
        <div v-else class="space-y-2 max-h-52 overflow-y-auto pr-1">
          <div
            v-for="ord in customerPeriodOrders"
            :key="ord.id"
            class="p-3 rounded-xl border border-slate-100 bg-slate-50/60 flex items-center justify-between text-xs"
          >
            <div>
              <div class="flex items-center gap-2">
                <span class="font-black text-slate-900">#{{ ord.order_number }}</span>
                <span
                  class="px-1.5 py-0.2 rounded text-[10px] font-black uppercase"
                  :class="
                    ord.payment_status === 'PAID'
                      ? 'bg-emerald-100 text-emerald-800'
                      : 'bg-rose-100 text-rose-800'
                  "
                >
                  {{ ord.payment_status }}
                </span>
              </div>
              <span class="text-[11px] text-slate-400 mt-0.5 block">{{
                formatDate(ord.created_at)
              }}</span>
            </div>
            <div class="text-right">
              <span class="font-black text-slate-900 block"
                >Rp {{ formatRupiah(ord.total_price) }}</span
              >
              <span class="text-[10px] text-slate-400">{{
                ord.total_weight ? `${ord.total_weight} Kg` : ''
              }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="pt-2 flex items-center gap-2">
        <a
          v-if="customerData?.phone"
          :href="whatsappUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="flex items-center justify-center gap-2 flex-1 py-3 px-4 rounded-2xl bg-emerald-600 hover:bg-emerald-700 active:scale-98 text-white font-black text-sm shadow-md shadow-emerald-600/20 transition-all cursor-pointer text-center"
        >
          <MessageCircle class="w-4 h-4" />
          <span>Sapa WhatsApp</span>
        </a>
        <a
          v-if="customerData?.phone"
          :href="`tel:${customerData.phone}`"
          class="p-3 rounded-2xl bg-slate-100 hover:bg-slate-200 text-slate-700 transition-all cursor-pointer"
          title="Telepon Langsung"
        >
          <Phone class="w-4 h-4" />
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { X, Phone, MapPin, Loader2, Calendar, MessageCircle } from 'lucide-vue-next'
import { api, formatRupiah, formatDate } from '../../api/client'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  customerId: {
    type: String,
    default: '',
  },
  customer: {
    type: Object,
    default: () => ({}),
  },
  startDate: {
    type: String,
    default: '',
  },
  endDate: {
    type: String,
    default: '',
  },
  customerPeriodOrders: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue'])

const modalTab = ref('LIFETIME') // 'LIFETIME' | 'PERIOD'
const loading = ref(false)
const profile = ref(null)

const customerData = computed(() => {
  return {
    ...props.customer,
    ...(profile.value || {}),
  }
})

const getInitials = (name) => {
  if (!name) return 'PL'
  return name
    .split(' ')
    .map((w) => w[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
}

const formatWeight = (val) => {
  if (!val) return '0'
  const num = Number(val)
  return isNaN(num) ? '0' : num.toFixed(1)
}

const formatCategoryName = (cat) => {
  switch (cat) {
    case 'KILOAN':
      return 'Cuci Kiloan'
    case 'SELF_SERVICE':
      return 'Self Service'
    case 'SATUAN':
      return 'Cuci Satuan'
    case 'SEPATU_TAS':
      return 'Sepatu & Tas'
    case 'KARPET':
      return 'Karpet'
    default:
      return cat || 'Layanan'
  }
}

const periodSpend = computed(() => {
  return props.customerPeriodOrders.reduce((sum, o) => {
    return o.status !== 'BATAL' ? sum + (Number(o.total_price) || 0) : sum
  }, 0)
})

const periodWeight = computed(() => {
  return props.customerPeriodOrders.reduce((sum, o) => {
    return o.status !== 'BATAL' ? sum + (Number(o.total_weight) || 0) : sum
  }, 0)
})

const periodServices = computed(() => {
  const set = new Set()
  props.customerPeriodOrders.forEach((o) => {
    if (Array.isArray(o.items)) {
      o.items.forEach((item) => {
        if (item.category) set.add(item.category)
      })
    }
  })
  return Array.from(set)
})

const loadProfile = async () => {
  if (!props.customerId) return
  loading.value = true
  try {
    const data = await api.get(`/reports/customers/${props.customerId}/crm`)
    profile.value = data
  } catch (err) {
    console.error('Failed to load customer CRM profile', err)
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  (open) => {
    if (open && props.customerId) {
      loadProfile()
      modalTab.value = 'LIFETIME'
    } else {
      profile.value = null
    }
  }
)

const whatsappUrl = computed(() => {
  const phone = customerData.value?.phone || ''
  const cleanPhone = phone.replace(/^0/, '62').replace(/[^0-9]/g, '')
  const name = customerData.value?.full_name || 'Kak'
  const text = encodeURIComponent(
    `Halo Kak ${name}! Terima kasih banyak telah mempercayakan pakaiannya di FAN Laundry Kalibaru 🙏 Senang bisa melayani Kakak. Hubungi kami kapan saja jika ada pakaian yang ingin dijemput atau dicuci lagi ya! ✨`
  )
  return `https://wa.me/${cleanPhone}?text=${text}`
})

const close = () => {
  emit('update:modelValue', false)
}
</script>
