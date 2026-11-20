<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs"
  >
    <div
      class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl border border-slate-100 space-y-4"
    >
      <!-- Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-100">
        <div class="flex items-center gap-2.5">
          <div
            class="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center"
          >
            <CreditCard class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-bold text-base text-slate-900">Pelunasan Pembayaran</h3>
            <p class="text-xs text-slate-500 font-mono">
              {{ formatDisplayOrder(order.order_number) }} • {{ order.customer_name }}
            </p>
          </div>
        </div>
        <button class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600" @click="close">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Tagihan Summary -->
      <div class="p-3.5 rounded-xl bg-slate-50 border border-slate-200/80 space-y-1.5 text-xs">
        <div class="flex justify-between text-slate-500">
          <span>Total Tagihan:</span>
          <span class="font-bold text-slate-900">Rp {{ formatRupiah(order.total_price) }}</span>
        </div>
        <div class="flex justify-between text-slate-500">
          <span>Sudah Dibayar Sebelumnya:</span>
          <span class="font-semibold text-slate-700">Rp {{ formatRupiah(order.paid_amount) }}</span>
        </div>
        <div
          class="flex justify-between text-rose-700 font-extrabold text-sm pt-1.5 border-t border-slate-200"
        >
          <span>SISA TAGIHAN:</span>
          <span>Rp {{ formatRupiah(sisaTagihan) }}</span>
        </div>
      </div>

      <!-- Payment Method Selection -->
      <div>
        <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1.5"
          >Metode Pembayaran</label
        >
        <div class="grid grid-cols-3 gap-2">
          <button
            v-for="method in methods"
            :key="method.id"
            type="button"
            class="p-2.5 rounded-xl border text-center transition-all flex flex-col items-center gap-1"
            :class="
              selectedMethod?.id === method.id
                ? 'border-brand-500 bg-brand-50/50 text-brand-700 font-bold'
                : 'border-slate-200 text-slate-600 hover:border-slate-300 font-medium'
            "
            @click="selectedMethod = method"
          >
            <component :is="getMethodIcon(method.category)" class="w-4 h-4" />
            <span class="text-xs">{{ method.name }}</span>
          </button>
        </div>
      </div>

      <!-- Transfer / QRIS Details Notice -->
      <div
        v-if="selectedMethod?.category === 'TRANSFER'"
        class="p-3 rounded-xl bg-blue-50/70 border border-blue-200 text-xs text-blue-900 space-y-0.5"
      >
        <p class="font-bold">BCA: 4140500068</p>
        <p class="text-[11px] text-blue-700">a.n. Siti Nurjanah (FAN Laundry)</p>
      </div>
      <div
        v-else-if="selectedMethod?.category === 'QRIS'"
        class="p-3 rounded-xl bg-emerald-50/70 border border-emerald-200 text-xs text-emerald-900 space-y-0.5"
      >
        <p class="font-bold">QRIS AstraPay / Semua E-Wallet</p>
        <p class="text-[11px] text-emerald-700">Scan QRIS statis di meja kasir</p>
      </div>

      <!-- Input Jumlah Bayar -->
      <div>
        <div class="flex justify-between items-center mb-1">
          <label class="text-xs font-bold text-slate-700 uppercase tracking-wider"
            >Nominal Diterima</label
          >
          <button
            type="button"
            class="text-[11px] font-bold text-brand-600 hover:text-brand-800"
            @click="payAmount = sisaTagihan"
          >
            Uang Pas (Rp {{ formatRupiah(sisaTagihan) }})
          </button>
        </div>
        <div class="relative">
          <span class="absolute left-3.5 top-2.5 text-sm font-bold text-slate-400">Rp</span>
          <input
            v-model.number="payAmount"
            type="number"
            min="0"
            step="1000"
            required
            class="w-full pl-11 pr-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-slate-900 font-bold text-base focus:outline-hidden focus:ring-2 focus:ring-brand-500 focus:bg-white transition-all"
          />
        </div>
      </div>

      <!-- Quick Preset Buttons -->
      <div class="flex gap-1.5 overflow-x-auto pb-1">
        <button
          v-for="preset in quickPresets"
          :key="preset"
          type="button"
          class="px-2.5 py-1 rounded-lg bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-semibold shrink-0 transition-all"
          @click="payAmount = preset"
        >
          Rp {{ formatRupiah(preset) }}
        </button>
      </div>

      <!-- Kembalian Calculation -->
      <div
        v-if="selectedMethod?.category === 'TUNAI' && payAmount > sisaTagihan"
        class="p-2.5 rounded-xl bg-emerald-50 border border-emerald-200 flex justify-between items-center text-xs"
      >
        <span class="font-semibold text-emerald-800">Kembalian:</span>
        <span class="font-extrabold text-emerald-700 text-sm"
          >Rp {{ formatRupiah(payAmount - sisaTagihan) }}</span
        >
      </div>

      <!-- Error -->
      <p
        v-if="errorMessage"
        class="text-xs text-rose-600 font-medium bg-rose-50 p-2.5 rounded-xl border border-rose-200"
      >
        {{ errorMessage }}
      </p>

      <!-- Action Buttons -->
      <div class="flex items-center gap-2 pt-2">
        <button
          type="button"
          class="w-1/3 py-2.5 rounded-xl border border-slate-200 text-slate-700 font-semibold text-xs hover:bg-slate-50 transition-all"
          @click="close"
        >
          Batal
        </button>
        <button
          type="button"
          :disabled="loading || payAmount <= 0"
          class="w-2/3 py-2.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs flex items-center justify-center gap-1.5 shadow-sm transition-all disabled:opacity-50"
          @click="submitPayment"
        >
          <span v-if="loading">Memproses...</span>
          <span v-else>Catat Pelunasan</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { CreditCard, Banknote, QrCode, X } from 'lucide-vue-next'
import { api, formatRupiah } from '../../api/client'

const props = defineProps({
  modelValue: Boolean,
  order: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue', 'payment-success'])

const loading = ref(false)
const errorMessage = ref('')
const payAmount = ref(0)
const methods = ref([])
const selectedMethod = ref(null)

const sisaTagihan = computed(() => {
  const total = Number(props.order.total_price) || 0
  const paid = Number(props.order.paid_amount) || 0
  return Math.max(0, total - paid)
})

const quickPresets = computed(() => {
  const presets = [sisaTagihan.value]
  if (sisaTagihan.value < 50000) presets.push(50000)
  if (sisaTagihan.value < 100000) presets.push(100000)
  presets.push(200000)
  return [...new Set(presets)]
})

const formatDisplayOrder = (num) => {
  if (!num) return ''
  return num.startsWith('ORD') ? '#' + num.substring(3) : num
}

const getMethodIcon = (category) => {
  switch (category) {
    case 'TUNAI':
      return Banknote
    case 'TRANSFER':
      return CreditCard
    case 'QRIS':
      return QrCode
    default:
      return CreditCard
  }
}

const loadPaymentMethods = async () => {
  try {
    const data = await api.get('/master/payment-methods')
    methods.value = data || []
    if (methods.value.length > 0) {
      selectedMethod.value = methods.value[0]
    }
  } catch {
    methods.value = [
      { id: '1', code: 'TUNAI', name: 'Tunai', category: 'TUNAI' },
      { id: '2', code: 'BCA', name: 'BCA', category: 'TRANSFER' },
      { id: '3', code: 'QRIS', name: 'QRIS', category: 'QRIS' },
    ]
    selectedMethod.value = methods.value[0]
  }
}

const close = () => {
  emit('update:modelValue', false)
  errorMessage.value = ''
}

const submitPayment = async () => {
  if (payAmount.value <= 0) {
    errorMessage.value = 'Nominal pembayaran harus lebih besar dari 0.'
    return
  }

  loading.value = true
  errorMessage.value = ''
  try {
    // New total paid is existing paid amount + this payment
    const newTotalPaid =
      (Number(props.order.paid_amount) || 0) + Math.min(payAmount.value, sisaTagihan.value)
    await api.patch(`/pos/orders/${props.order.id}/payment`, {
      paid_amount: newTotalPaid,
      payment_method_id: selectedMethod.value?.id || '',
      payment_method_name: selectedMethod.value?.name || 'Tunai',
      staff_name: 'Kasir Utama',
    })
    emit('payment-success')
    close()
  } catch (err) {
    errorMessage.value = err.message || 'Gagal memproses pembayaran.'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  payAmount.value = sisaTagihan.value
  loadPaymentMethods()
})
</script>
