<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
  >
    <div
      class="bg-white rounded-2xl max-w-2xl w-full p-4 sm:p-6 shadow-2xl border border-slate-100 flex flex-col max-h-[92vh]"
    >
      <!-- Modal Header -->
      <div class="flex items-start justify-between pb-4 border-b border-slate-100 shrink-0">
        <div class="space-y-1">
          <div class="flex items-center gap-2">
            <span class="font-black text-lg sm:text-xl text-slate-900 font-mono tracking-tight">
              {{ formatDisplayOrder(order.order_number) }}
            </span>
            <span
              class="text-xs font-bold px-2.5 py-0.5 rounded-full uppercase tracking-wider"
              :class="getStatusClass(order.status)"
            >
              {{ order.status }}
            </span>
            <span
              class="text-xs font-bold px-2.5 py-0.5 rounded-full uppercase tracking-wider"
              :class="
                order.payment_status === 'PAID'
                  ? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                  : 'bg-rose-50 text-rose-700 border border-rose-200'
              "
            >
              {{ order.payment_status === 'PAID' ? 'LUNAS' : 'BELUM LUNAS' }}
            </span>
          </div>
          <div class="flex items-center gap-2 text-xs text-slate-500">
            <User class="w-3.5 h-3.5 text-slate-400" />
            <span class="font-bold text-slate-800">{{ order.customer_name }}</span>
            <span>•</span>
            <a
              :href="`https://wa.me/${order.customer_phone?.replace(/[^0-9]/g, '')}`"
              target="_blank"
              class="text-brand-600 hover:underline flex items-center gap-1 font-semibold"
            >
              <span>{{ order.customer_phone }}</span>
              <ExternalLink class="w-3 h-3" />
            </a>
          </div>
        </div>

        <button class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600" @click="close">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Action Bar (Progresi 1-Klik & Quick Tools) -->
      <div
        class="py-3 border-b border-slate-100 flex flex-wrap items-center justify-between gap-2 shrink-0 bg-slate-50/80 -mx-4 sm:-mx-6 px-4 sm:px-6"
      >
        <!-- Primary 1-Click Progression Button -->
        <div class="flex items-center gap-2">
          <!-- Antrian -> Proses -->
          <button
            v-if="order.status === 'ANTRIAN'"
            :disabled="actionLoading"
            class="px-4 py-2 rounded-xl bg-amber-500 hover:bg-amber-600 text-white font-bold text-xs flex items-center gap-1.5 shadow-sm transition-all"
            @click="progressStatus('PROSES')"
          >
            <Play class="w-4 h-4 fill-current" />
            <span>PROSES ORDERAN</span>
          </button>

          <!-- Proses -> Siap Ambil -->
          <button
            v-else-if="order.status === 'PROSES'"
            :disabled="actionLoading"
            class="px-4 py-2 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs flex items-center gap-1.5 shadow-sm transition-all"
            @click="progressStatus('READY')"
          >
            <CheckCircle class="w-4 h-4" />
            <span>ORDERAN SIAP DIAMBIL</span>
          </button>

          <!-- Siap Ambil -> Selesai -->
          <button
            v-else-if="order.status === 'READY' || order.status === 'SIAP AMBIL'"
            :disabled="actionLoading"
            class="px-4 py-2 rounded-xl bg-slate-900 hover:bg-slate-800 text-white font-bold text-xs flex items-center gap-1.5 shadow-sm transition-all"
            @click="completeOrder"
          >
            <CheckCheck class="w-4 h-4" />
            <span>TRANSAKSI TELAH SELESAI</span>
          </button>
        </div>

        <!-- Secondary Tools -->
        <div class="flex items-center gap-1.5">
          <!-- Cetak Nota -->
          <button
            class="p-2 rounded-xl bg-white hover:bg-slate-100 border border-slate-200 text-slate-700 text-xs font-bold flex items-center gap-1.5 transition-all"
            title="Cetak Nota Thermal"
            @click="showReceiptModal = true"
          >
            <Printer class="w-4 h-4 text-brand-600" />
            <span class="hidden sm:inline">Cetak Nota</span>
          </button>

          <!-- Foto Dokumentasi -->
          <button
            class="p-2 rounded-xl bg-white hover:bg-slate-100 border border-slate-200 text-slate-700 text-xs font-bold flex items-center gap-1.5 transition-all"
            title="Foto Dokumentasi"
            @click="showPhotoModal = true"
          >
            <Camera class="w-4 h-4 text-indigo-600" />
            <span class="hidden sm:inline">Foto</span>
          </button>

          <!-- Pelunasan Bayar (if unpaid) -->
          <button
            v-if="order.payment_status !== 'PAID' && order.status !== 'BATAL'"
            class="p-2 rounded-xl bg-emerald-50 hover:bg-emerald-100 border border-emerald-200 text-emerald-800 text-xs font-bold flex items-center gap-1.5 transition-all"
            title="Pelunasan Tagihan"
            @click="showPaymentModal = true"
          >
            <CreditCard class="w-4 h-4 text-emerald-600" />
            <span>Bayar Lunas</span>
          </button>

          <!-- Resend WhatsApp Invoice -->
          <button
            :disabled="resendingWA"
            class="p-2 rounded-xl bg-emerald-50 hover:bg-emerald-100 border border-emerald-200 text-emerald-700 text-xs font-bold flex items-center gap-1 transition-all"
            title="Kirim Ulang Nota ke WhatsApp"
            @click="resendWAInvoice"
          >
            <Send class="w-4 h-4 text-emerald-600" />
            <span class="hidden md:inline">Kirim WA</span>
          </button>

          <!-- Batalkan Order -->
          <button
            v-if="order.status !== 'SELESAI' && order.status !== 'BATAL'"
            class="p-2 rounded-xl bg-rose-50 hover:bg-rose-100 border border-rose-200 text-rose-700 text-xs font-bold flex items-center gap-1 transition-all"
            title="Batalkan Transaksi"
            @click="promptCancelOrder"
          >
            <Ban class="w-4 h-4" />
            <span class="hidden md:inline">Batal</span>
          </button>
        </div>
      </div>

      <!-- Modal Body (Scrollable) -->
      <div class="flex-1 overflow-y-auto py-4 space-y-4 pr-1">
        <!-- Summary Cards Grid -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-2.5 text-xs">
          <!-- Total Berat -->
          <div class="p-3 rounded-xl bg-slate-50 border border-slate-100">
            <span class="text-slate-400 font-medium">Total Berat:</span>
            <p class="font-bold text-slate-800 text-sm mt-0.5">
              {{ order.total_weight ? `${order.total_weight} Kg` : '-' }}
            </p>
          </div>

          <!-- Aroma Parfum -->
          <div class="p-3 rounded-xl bg-slate-50 border border-slate-100">
            <span class="text-slate-400 font-medium">Parfum:</span>
            <p class="font-bold text-slate-800 text-sm mt-0.5 truncate">
              {{ order.perfume_name || 'Standar' }}
            </p>
          </div>

          <!-- Lokasi Rak -->
          <div
            class="p-3 rounded-xl bg-slate-50 border border-slate-100 flex items-center justify-between"
          >
            <div>
              <span class="text-slate-400 font-medium">Nomor Rak:</span>
              <p class="font-bold text-slate-800 text-sm mt-0.5">{{ order.rack_number || '-' }}</p>
            </div>
            <button
              class="p-1 text-brand-600 hover:text-brand-800"
              title="Ganti Rak"
              @click="editRack"
            >
              <Edit2 class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Estimasi Selesai -->
          <div class="p-3 rounded-xl bg-slate-50 border border-slate-100">
            <span class="text-slate-400 font-medium">SLA Selesai:</span>
            <p class="font-bold text-slate-800 text-xs mt-0.5">
              {{ formatDateTime(order.estimated_completion_at) }}
            </p>
          </div>
        </div>

        <!-- Items Table -->
        <div class="rounded-xl border border-slate-200 overflow-hidden">
          <div
            class="bg-slate-100/70 px-4 py-2.5 flex justify-between items-center text-xs font-bold text-slate-700"
          >
            <span>Rincian Layanan & Produk ({{ items.length }})</span>
            <span>Subtotal</span>
          </div>
          <div class="divide-y divide-slate-100 text-xs">
            <div
              v-for="item in items"
              :key="item.id"
              class="px-4 py-2.5 flex justify-between items-center hover:bg-slate-50/50"
            >
              <div>
                <p class="font-bold text-slate-900">{{ item.service_name }}</p>
                <p class="text-[11px] text-slate-500">
                  {{ item.quantity }} {{ item.unit }} x Rp {{ formatRupiah(item.price_per_unit) }}
                  <span v-if="item.notes" class="text-amber-700 ml-1">({{ item.notes }})</span>
                </p>
              </div>
              <span class="font-extrabold text-slate-900"
                >Rp {{ formatRupiah(item.subtotal) }}</span
              >
            </div>
          </div>
        </div>

        <!-- Financial Summary -->
        <div class="p-3.5 rounded-xl bg-slate-50 border border-slate-200 space-y-1.5 text-xs">
          <div class="flex justify-between text-slate-600">
            <span>Total Tagihan:</span>
            <span class="font-bold text-slate-900">Rp {{ formatRupiah(order.total_price) }}</span>
          </div>
          <div v-if="order.discount_amount > 0" class="flex justify-between text-emerald-700">
            <span>Diskon:</span>
            <span>- Rp {{ formatRupiah(order.discount_amount) }}</span>
          </div>
          <div class="flex justify-between text-slate-600">
            <span>Sudah Dibayar ({{ order.payment_method_name || 'Tunai' }}):</span>
            <span class="font-semibold text-slate-800"
              >Rp {{ formatRupiah(order.paid_amount) }}</span
            >
          </div>
          <div
            class="flex justify-between pt-1.5 border-t border-slate-200 font-black text-sm"
            :class="order.payment_status === 'PAID' ? 'text-emerald-700' : 'text-rose-700'"
          >
            <span>{{ order.payment_status === 'PAID' ? 'STATUS LUNAS ✅' : 'SISA TAGIHAN:' }}</span>
            <span>{{
              order.payment_status === 'PAID'
                ? 'Rp 0'
                : `Rp ${formatRupiah(order.total_price - order.paid_amount)}`
            }}</span>
          </div>
        </div>

        <!-- Audit Trail Timeline -->
        <div class="space-y-2 pt-2">
          <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider">
            Riwayat Status & Petugas
          </h4>
          <div class="p-3 rounded-xl bg-white border border-slate-200 space-y-3">
            <div v-for="log in history" :key="log.id" class="flex items-start gap-2.5 text-xs">
              <div class="w-2 h-2 rounded-full bg-brand-500 mt-1.5 shrink-0"></div>
              <div class="flex-1">
                <div class="flex items-center justify-between">
                  <span class="font-bold text-slate-800">{{ log.new_status }}</span>
                  <span class="text-[10px] text-slate-400">{{
                    formatDateTime(log.created_at)
                  }}</span>
                </div>
                <p class="text-[11px] text-slate-500">
                  Oleh: <span class="font-semibold text-slate-700">{{ log.changed_by }}</span>
                  <span v-if="log.notes" class="text-slate-400"> — {{ log.notes }}</span>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer Buttons -->
      <div class="pt-3 border-t border-slate-100 flex items-center justify-end gap-2 shrink-0">
        <button
          class="px-5 py-2.5 rounded-xl bg-slate-900 hover:bg-slate-800 text-white font-bold text-xs transition-all"
          @click="close"
        >
          Tutup
        </button>
      </div>
    </div>

    <!-- Sub-Modals -->
    <ReceiptPrintModal v-model="showReceiptModal" :order="order" :items="items" />
    <PaymentModal
      v-model="showPaymentModal"
      :order="order"
      @payment-success="handlePaymentSuccess"
    />
    <PhotoDocModal
      v-model="showPhotoModal"
      :order="order"
      :photos="photos"
      @photo-added="refreshDetails"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import {
  User,
  ExternalLink,
  Play,
  CheckCircle,
  CheckCheck,
  Printer,
  Camera,
  CreditCard,
  Send,
  Ban,
  Edit2,
  X,
} from 'lucide-vue-next'
import ReceiptPrintModal from './ReceiptPrintModal.vue'
import PaymentModal from './PaymentModal.vue'
import PhotoDocModal from './PhotoDocModal.vue'
import { api, formatRupiah, formatDateTime } from '../../api/client'

const props = defineProps({
  modelValue: Boolean,
  order: {
    type: Object,
    required: true,
  },
  items: {
    type: Array,
    default: () => [],
  },
  photos: {
    type: Array,
    default: () => [],
  },
  history: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue', 'order-updated'])

const actionLoading = ref(false)
const resendingWA = ref(false)
const showReceiptModal = ref(false)
const showPaymentModal = ref(false)
const showPhotoModal = ref(false)

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

const close = () => {
  emit('update:modelValue', false)
}

const progressStatus = async (newStatus) => {
  if (newStatus === 'READY' && !props.order.rack_number) {
    const rack = prompt('Masukkan nomor rak simpan (misal: A-01, B-03):', 'A-01')
    if (rack) {
      await api.patch(`/pos/orders/${props.order.id}/rack`, { rack_number: rack })
    }
  }

  actionLoading.value = true
  try {
    await api.patch(`/pos/orders/${props.order.id}/status`, {
      status: newStatus,
      changed_by: 'Kasir Utama',
      notes: `Status diperbarui menjadi ${newStatus}`,
    })
    emit('order-updated')
  } catch (err) {
    alert(err.message || 'Gagal mengubah status order.')
  } finally {
    actionLoading.value = false
  }
}

const completeOrder = async () => {
  if (props.order.payment_status !== 'PAID') {
    if (confirm('Pesanan ini belum lunas! Apakah ingin melakukan pelunasan pembayaran sekarang?')) {
      showPaymentModal.value = true
      return
    }
  }
  await progressStatus('SELESAI')
}

const editRack = async () => {
  const current = props.order.rack_number || ''
  const newRack = prompt('Ubah nomor rak simpan:', current)
  if (newRack !== null) {
    try {
      await api.patch(`/pos/orders/${props.order.id}/rack`, { rack_number: newRack })
      emit('order-updated')
    } catch (err) {
      alert(err.message || 'Gagal mengubah nomor rak.')
    }
  }
}

const resendWAInvoice = async () => {
  resendingWA.value = true
  try {
    await api.post(`/pos/orders/${props.order.id}/resend-invoice`)
    alert('Invoice WhatsApp berhasil dikirim ulang ke pelanggan! ✅')
  } catch (err) {
    alert(err.message || 'Gagal mengirim ulang WhatsApp.')
  } finally {
    resendingWA.value = false
  }
}

const promptCancelOrder = async () => {
  const reason = prompt('Masukkan alasan pembatalan pesanan ini:')
  if (!reason) return

  actionLoading.value = true
  try {
    await api.patch(`/pos/orders/${props.order.id}/cancel`, {
      reason: reason,
      changed_by: 'Kasir Utama',
    })
    emit('order-updated')
  } catch (err) {
    alert(err.message || 'Gagal membatalkan pesanan.')
  } finally {
    actionLoading.value = false
  }
}

const handlePaymentSuccess = () => {
  emit('order-updated')
}

const refreshDetails = () => {
  emit('order-updated')
}
</script>
