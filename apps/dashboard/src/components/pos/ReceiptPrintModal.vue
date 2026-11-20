<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs print:p-0 print:bg-white print:static"
  >
    <div
      class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl border border-slate-100 flex flex-col max-h-[90vh] print:max-w-none print:shadow-none print:border-none print:p-0"
    >
      <!-- Modal Header (Hidden on Print) -->
      <div
        class="flex items-center justify-between pb-4 border-b border-slate-100 print:hidden shrink-0"
      >
        <div class="flex items-center gap-2.5">
          <div
            class="w-9 h-9 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center"
          >
            <Printer class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-bold text-base text-slate-900">Cetak Nota Thermal</h3>
            <div class="flex items-center gap-1.5">
              <p class="text-xs text-slate-500">Pilih jenis nota untuk dicetak</p>
              <span
                class="text-[10px] font-bold px-1.5 py-0.2 rounded-md bg-slate-100 text-slate-600 border border-slate-200"
              >
                Format: {{ settings.printer?.paperWidth || '58mm' }}
              </span>
            </div>
          </div>
        </div>
        <button class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600" @click="close">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Receipt Type Tabs (Hidden on Print) -->
      <div class="flex items-center gap-1.5 p-1 bg-slate-100 rounded-xl my-4 print:hidden shrink-0">
        <button
          v-for="tab in receiptTabs"
          :key="tab.id"
          class="flex-1 py-1.5 text-xs font-bold rounded-lg transition-all text-center"
          :class="
            activeTab === tab.id
              ? 'bg-white text-slate-900 shadow-xs'
              : 'text-slate-500 hover:text-slate-800'
          "
          @click="activeTab = tab.id"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- Thermal Receipt Preview (This is what prints) -->
      <div
        class="flex-1 overflow-y-auto p-4 bg-slate-50 rounded-xl border border-slate-200/80 font-mono text-xs text-slate-900 print:p-0 print:border-none print:bg-white print:overflow-visible"
      >
        <div
          id="thermal-receipt"
          class="mx-auto bg-white p-4 shadow-xs print:shadow-none print:p-0 transition-all"
          :class="settings.printer?.paperWidth === '80mm' ? 'max-w-[420px]' : 'max-w-[320px]'"
        >
          <!-- Header Nota -->
          <div class="text-center space-y-1 pb-3 border-b border-dashed border-slate-400">
            <h2 class="font-extrabold text-base tracking-wider uppercase">
              {{ settings.printer?.storeName || 'FAN LAUNDRY' }}
            </h2>
            <p v-if="settings.printer?.addressLine1" class="text-[11px] leading-tight">
              {{ settings.printer.addressLine1 }}
            </p>
            <p v-if="settings.printer?.addressLine2" class="text-[11px] leading-tight">
              {{ settings.printer.addressLine2 }}
            </p>
            <p v-if="settings.printer?.phone" class="text-[11px] font-semibold">
              WA: {{ settings.printer.phone }}
            </p>
          </div>

          <!-- Type Badge -->
          <div
            class="py-2 text-center border-b border-dashed border-slate-400 font-bold uppercase tracking-wider text-[11px]"
          >
            *** {{ currentTabLabel }} ***
          </div>

          <!-- Order & Customer Meta -->
          <div class="py-2.5 space-y-1 border-b border-dashed border-slate-400 text-[11px]">
            <div class="flex justify-between">
              <span>No. Order:</span>
              <span class="font-bold font-mono">{{ formatDisplayOrder(order.order_number) }}</span>
            </div>
            <div class="flex justify-between">
              <span>Tanggal:</span>
              <span>{{ formatDateTime(order.created_at) }}</span>
            </div>
            <div class="flex justify-between">
              <span>Pelanggan:</span>
              <span class="font-bold">{{ order.customer_name }}</span>
            </div>
            <div class="flex justify-between">
              <span>No. Telp:</span>
              <span>{{ order.customer_phone }}</span>
            </div>
            <div v-if="order.perfume_name" class="flex justify-between">
              <span>Aroma Parfum:</span>
              <span class="font-semibold">{{ order.perfume_name }}</span>
            </div>
            <div v-if="order.rack_number" class="flex justify-between">
              <span>Lokasi Rak:</span>
              <span class="font-bold">Rak {{ order.rack_number }}</span>
            </div>
            <div class="flex justify-between">
              <span>Estimasi Selesai:</span>
              <span class="font-semibold">{{ formatDateTime(order.estimated_completion_at) }}</span>
            </div>
          </div>

          <!-- Items Table -->
          <div class="py-2.5 border-b border-dashed border-slate-400 text-[11px]">
            <div class="space-y-2">
              <div v-for="item in items" :key="item.id" class="space-y-0.5">
                <div class="flex justify-between font-bold">
                  <span>{{ item.service_name }}</span>
                  <span v-if="activeTab !== 'PRODUKSI'">Rp {{ formatRupiah(item.subtotal) }}</span>
                </div>
                <div class="flex justify-between text-slate-600 text-[10px]">
                  <span
                    >{{ item.quantity }} {{ item.unit }} x Rp
                    {{ formatRupiah(item.price_per_unit) }}</span
                  >
                  <span v-if="item.notes" class="italic">({{ item.notes }})</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Financial Calculation (Only for Konsumen & Pertinggal) -->
          <div
            v-if="activeTab !== 'PRODUKSI'"
            class="py-2.5 space-y-1 border-b border-dashed border-slate-400 text-[11px]"
          >
            <div class="flex justify-between">
              <span>Subtotal:</span>
              <span>Rp {{ formatRupiah(order.total_price + (order.discount_amount || 0)) }}</span>
            </div>
            <div v-if="order.discount_amount > 0" class="flex justify-between text-emerald-700">
              <span>Diskon:</span>
              <span>- Rp {{ formatRupiah(order.discount_amount) }}</span>
            </div>
            <div class="flex justify-between font-extrabold text-xs pt-1 border-t border-slate-200">
              <span>TOTAL TAGIHAN:</span>
              <span>Rp {{ formatRupiah(order.total_price) }}</span>
            </div>
            <div class="flex justify-between">
              <span>Dibayar ({{ order.payment_method_name || 'Tunai' }}):</span>
              <span>Rp {{ formatRupiah(order.paid_amount) }}</span>
            </div>
            <div
              class="flex justify-between font-bold text-[11px]"
              :class="order.payment_status === 'PAID' ? 'text-emerald-700' : 'text-rose-700'"
            >
              <span>STATUS:</span>
              <span>{{
                order.payment_status === 'PAID'
                  ? 'LUNAS'
                  : `SISA Rp ${formatRupiah(order.total_price - order.paid_amount)}`
              }}</span>
            </div>
          </div>

          <!-- Footer Information -->
          <div class="text-center space-y-1.5 pt-3 text-[10px] text-slate-500">
            <template v-if="activeTab === 'KONSUMEN'">
              <div class="whitespace-pre-line leading-relaxed text-slate-700">
                {{
                  settings.printer?.footerNote ||
                  'Terima kasih atas kepercayaan Anda!\nCucian tidak diambil > 30 hari di luar tanggung jawab kasir.'
                }}
              </div>
              <p class="font-bold pt-1">*** SIMPAN NOTA INI DENGAN BAIK ***</p>
            </template>
            <template v-else-if="activeTab === 'PERTINGGAL'">
              <p class="font-bold uppercase">ARSIP KASIR & PEMBUKUAN</p>
              <div class="flex justify-between pt-4 px-2">
                <span>( Kasir )</span>
                <span>( Konsumen )</span>
              </div>
            </template>
            <template v-else-if="activeTab === 'PRODUKSI'">
              <p class="font-extrabold text-xs uppercase bg-slate-100 py-1">
                TEMPEL PADA HANGER / KERANJANG
              </p>
              <p class="font-bold">
                Total Berat: {{ order.total_weight ? `${order.total_weight} Kg` : '-' }}
              </p>
              <p v-if="order.notes" class="text-amber-800 font-bold">Catatan: {{ order.notes }}</p>
            </template>
          </div>
        </div>
      </div>

      <!-- Action Buttons (Hidden on Print) -->
      <div class="flex items-center gap-3 pt-4 border-t border-slate-100 print:hidden shrink-0">
        <button
          class="w-1/3 py-2.5 rounded-xl border border-slate-200 text-slate-700 font-semibold text-xs hover:bg-slate-50 transition-all"
          @click="close"
        >
          Tutup
        </button>
        <button
          class="w-2/3 py-2.5 rounded-xl bg-brand-600 hover:bg-brand-700 text-white font-bold text-xs flex items-center justify-center gap-2 shadow-sm transition-all"
          @click="printReceipt"
        >
          <Printer class="w-4 h-4" />
          <span>Cetak {{ currentTabLabel }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Printer, X } from 'lucide-vue-next'
import { formatRupiah, formatDateTime } from '../../api/client'
import { useAppSettings } from '../../composables/useAppSettings'

const { settings } = useAppSettings()

defineProps({
  modelValue: Boolean,
  order: {
    type: Object,
    required: true,
  },
  items: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue'])

const activeTab = ref('KONSUMEN')

const receiptTabs = [
  { id: 'KONSUMEN', label: 'Nota Konsumen' },
  { id: 'PERTINGGAL', label: 'Nota Pertinggal' },
  { id: 'PRODUKSI', label: 'Nota Produksi' },
]

const currentTabLabel = computed(() => {
  const found = receiptTabs.find((t) => t.id === activeTab.value)
  return found ? found.label : 'Nota'
})

const formatDisplayOrder = (num) => {
  if (!num) return ''
  return num.startsWith('ORD') ? '#' + num.substring(3) : num
}

const close = () => {
  emit('update:modelValue', false)
}

const printReceipt = () => {
  window.print()
}
</script>

<style scoped>
@media print {
  body * {
    visibility: hidden;
  }
  #thermal-receipt,
  #thermal-receipt * {
    visibility: visible;
  }
  #thermal-receipt {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    max-width: 80mm;
    margin: 0;
    padding: 0;
  }
}
</style>
