<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
  >
    <div
      class="bg-white rounded-2xl max-w-4xl w-full p-4 sm:p-6 shadow-2xl border border-slate-100 flex flex-col max-h-[96vh] sm:max-h-[92vh]"
    >
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-100 shrink-0">
        <div class="flex items-center gap-2.5">
          <div
            class="w-10 h-10 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center shadow-xs"
          >
            <PlusCircle class="w-6 h-6" />
          </div>
          <div>
            <h3 class="font-black text-lg text-slate-900">Transaksi Kasir Walk-in</h3>
            <p class="text-xs text-slate-500">Pencatatan cucian masuk, timbangan & pembayaran</p>
          </div>
        </div>
        <button
          class="p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-all"
          @click="close"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Main Form (Scrollable) -->
      <div class="flex-1 overflow-y-auto py-3 space-y-5 pr-1 text-slate-800">
        <!-- STEP 1: Pelanggan (Search / Add New) -->
        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <User class="w-4 h-4 text-brand-600" />
              <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
                1. Data Pelanggan
              </h4>
            </div>
            <button
              type="button"
              class="text-xs font-bold text-brand-600 hover:underline"
              @click="isNewCustomer = !isNewCustomer"
            >
              {{ isNewCustomer ? 'Pilih Pelanggan Terdaftar' : '+ Pelanggan Baru' }}
            </button>
          </div>

          <!-- Existing Customer Selected Pill -->
          <div
            v-if="!isNewCustomer && selectedCustomer"
            class="p-3 rounded-xl bg-white border border-brand-200 flex items-center justify-between shadow-xs"
          >
            <div>
              <p class="font-bold text-sm text-slate-900">{{ selectedCustomer.full_name }}</p>
              <p class="text-xs text-slate-500">
                {{ selectedCustomer.phone }} • {{ selectedCustomer.address || 'Kalibaru' }}
              </p>
            </div>
            <button
              class="text-xs font-bold text-rose-600 hover:underline"
              @click="selectedCustomer = null"
            >
              Ganti
            </button>
          </div>

          <!-- Existing Customer Autocomplete Input -->
          <div v-else-if="!isNewCustomer" class="relative">
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-3" />
            <input
              v-model="customerSearchQuery"
              type="text"
              placeholder="Ketik nama atau nomor WhatsApp pelanggan..."
              class="w-full pl-9 pr-4 py-2.5 rounded-xl border border-slate-200 bg-white text-xs font-medium focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              @input="onCustomerSearch"
            />
            <!-- Search Results Dropdown -->
            <div
              v-if="customerResults.length > 0"
              class="absolute top-full left-0 right-0 mt-1 bg-white rounded-xl border border-slate-200 shadow-lg z-20 max-h-40 overflow-y-auto divide-y divide-slate-100"
            >
              <div
                v-for="c in customerResults"
                :key="c.id"
                class="p-2.5 hover:bg-brand-50 cursor-pointer text-xs flex justify-between items-center"
                @click="selectCustomer(c)"
              >
                <div>
                  <span class="font-bold text-slate-900">{{ c.full_name }}</span>
                  <span class="text-slate-400 ml-1.5 font-normal">({{ c.phone }})</span>
                </div>
                <span class="text-[10px] text-slate-400 truncate max-w-[150px]">{{
                  c.address
                }}</span>
              </div>
            </div>
          </div>

          <!-- Inline Form New Customer -->
          <div v-else class="grid grid-cols-1 sm:grid-cols-3 gap-2.5 pt-1">
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">Nama Lengkap *</label>
              <input
                v-model="newCustomer.full_name"
                type="text"
                required
                placeholder="Nama pelanggan"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Nomor WhatsApp *</label
              >
              <input
                v-model="newCustomer.phone"
                type="tel"
                required
                placeholder="0852... / 628..."
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Alamat (Opsional)</label
              >
              <input
                v-model="newCustomer.address"
                type="text"
                placeholder="Rt/Rw, patokan rumah"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
          </div>
        </div>

        <!-- STEP 2: Pilih Layanan & Keranjang Multi-Layanan -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <ShoppingBag class="w-4 h-4 text-brand-600" />
              <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
                2. Pilih Layanan Laundry
              </h4>
            </div>
            <span class="text-[11px] text-slate-500"
              >{{ cartItems.length }} Layanan di keranjang</span
            >
          </div>

          <!-- Category Filter Pills -->
          <div class="flex gap-1.5 overflow-x-auto pb-1 text-xs">
            <button
              v-for="cat in serviceCategories"
              :key="cat.id"
              type="button"
              class="px-3 py-1.5 rounded-xl font-bold whitespace-nowrap transition-all"
              :class="
                activeCategory === cat.id
                  ? 'bg-slate-900 text-white shadow-xs'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              "
              @click="activeCategory = cat.id"
            >
              {{ cat.label }}
            </button>
          </div>

          <!-- Service Grid -->
          <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-2">
            <button
              v-for="srv in filteredServices"
              :key="srv.id"
              type="button"
              class="p-2.5 rounded-xl border border-slate-200 bg-white hover:border-brand-400 hover:bg-brand-50/40 text-left transition-all group flex flex-col justify-between"
              @click="openQuantityModal(srv)"
            >
              <div>
                <span
                  class="text-[10px] font-bold px-1.5 py-0.2 rounded-md bg-slate-100 text-slate-600"
                >
                  {{ srv.unit }}
                </span>
                <p
                  class="font-bold text-xs text-slate-900 mt-1 line-clamp-2 leading-tight group-hover:text-brand-600"
                >
                  {{ srv.service_name }}
                </p>
              </div>
              <div
                class="mt-2 pt-1 border-t border-slate-100 flex justify-between items-center text-xs"
              >
                <span class="font-extrabold text-brand-700">Rp {{ formatRupiah(srv.price) }}</span>
                <Plus class="w-3.5 h-3.5 text-slate-400 group-hover:text-brand-600" />
              </div>
            </button>
          </div>

          <!-- Active Cart Table -->
          <div
            v-if="cartItems.length > 0"
            class="mt-3 rounded-xl border border-slate-200 overflow-hidden shadow-xs"
          >
            <div
              class="bg-slate-100 px-3.5 py-2 flex justify-between text-xs font-bold text-slate-700"
            >
              <span>Keranjang Cucian</span>
              <span>Subtotal</span>
            </div>
            <div class="divide-y divide-slate-100 bg-white text-xs">
              <div
                v-for="(item, idx) in cartItems"
                :key="idx"
                class="px-3.5 py-2.5 flex items-center justify-between"
              >
                <div class="flex items-center gap-2">
                  <button
                    class="p-1 text-slate-300 hover:text-rose-600 transition-colors"
                    title="Hapus"
                    @click="removeItem(idx)"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                  <div>
                    <span class="font-bold text-slate-900">{{ item.service_name }}</span>
                    <p class="text-[11px] text-slate-500">
                      {{ item.quantity }} {{ item.unit }} x Rp
                      {{ formatRupiah(item.price_per_unit) }}
                      <span v-if="item.notes" class="text-amber-700 ml-1">({{ item.notes }})</span>
                    </p>
                  </div>
                </div>
                <span class="font-black text-slate-900">Rp {{ formatRupiah(item.subtotal) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- STEP 3: Konfigurasi Order (Parfum, Diskon, Rak, SLA) -->
        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-3">
          <div class="flex items-center gap-2">
            <Sparkles class="w-4 h-4 text-brand-600" />
            <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
              3. Pilihan Tambahan & Lokasi Rak
            </h4>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-4 gap-2.5">
            <!-- Parfum -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">Aroma Parfum</label>
              <select
                v-model="orderConfig.perfume_name"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-medium"
              >
                <option value="Akasia">Akasia (Segar Favorit)</option>
                <option value="Sakura">Sakura (Manis Lembut)</option>
                <option value="Ocean Fresh">Ocean Fresh (Semilir Laut)</option>
                <option value="Snappy">Snappy (Elegan Segar)</option>
                <option value="Downy Black">Downy Black (Mewah Modern)</option>
              </select>
            </div>

            <!-- Diskon / Promo -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">Diskon / Promo</label>
              <select
                v-model="orderConfig.discount_type"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-medium"
                @change="applyDiscount"
              >
                <option value="NONE">Tanpa Diskon</option>
                <option value="PROMO_10">Promo Opening (10%)</option>
                <option value="LOYALTY">Loyalitas (Gratis 1x Cuci)</option>
                <option value="CUSTOM">Custom Nominal Rp</option>
              </select>
            </div>

            <!-- Custom Diskon input if selected -->
            <div v-if="orderConfig.discount_type === 'CUSTOM'">
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Potongan Diskon (Rp)</label
              >
              <input
                v-model.number="orderConfig.custom_discount"
                type="number"
                placeholder="5000"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
                @input="applyDiscount"
              />
            </div>

            <!-- Nomor Rak Simpan -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Nomor Rak Simpan</label
              >
              <input
                v-model="orderConfig.rack_number"
                type="text"
                placeholder="Contoh: A-01"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-bold"
              />
            </div>

            <!-- Estimasi Selesai (SLA) -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Estimasi Selesai (SLA)</label
              >
              <input
                v-model="orderConfig.estimated_completion_at"
                type="datetime-local"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-medium"
              />
            </div>
          </div>

          <!-- Catatan Order -->
          <div>
            <label class="block text-[11px] font-bold text-slate-600 mb-1"
              >Catatan Khusus (Opsional)</label
            >
            <input
              v-model="orderConfig.notes"
              type="text"
              placeholder="Contoh: Ada noda oli di lengan kiri, baju luntur dipisah..."
              class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
            />
          </div>
        </div>

        <!-- STEP 4: Opsi Pembayaran Awal -->
        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <CreditCard class="w-4 h-4 text-brand-600" />
              <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
                4. Pembayaran Awal
              </h4>
            </div>
            <!-- Toggle Bayar Sekarang vs Bayar Nanti -->
            <div
              class="flex items-center gap-1.5 p-1 bg-white rounded-xl border border-slate-200 text-xs"
            >
              <button
                type="button"
                class="px-2.5 py-1 rounded-lg font-bold transition-all"
                :class="
                  paymentOption === 'SEKARANG'
                    ? 'bg-brand-600 text-white'
                    : 'text-slate-500 hover:text-slate-800'
                "
                @click="paymentOption = 'SEKARANG'"
              >
                Bayar Sekarang
              </button>
              <button
                type="button"
                class="px-2.5 py-1 rounded-lg font-bold transition-all"
                :class="
                  paymentOption === 'NANTI'
                    ? 'bg-slate-900 text-white'
                    : 'text-slate-500 hover:text-slate-800'
                "
                @click="paymentOption = 'NANTI'"
              >
                Bayar Saat Ambil
              </button>
            </div>
          </div>

          <div
            v-if="paymentOption === 'SEKARANG'"
            class="grid grid-cols-1 sm:grid-cols-3 gap-2.5 pt-1"
          >
            <!-- Pilihan Metode -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">Metode Bayar</label>
              <select
                v-model="paymentMethod"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-bold"
              >
                <option value="TUNAI">Tunai / Cash</option>
                <option value="BCA">Transfer BCA (4140500068)</option>
                <option value="QRIS">QRIS AstraPay</option>
              </select>
            </div>

            <!-- Nominal Uang Diterima -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Nominal Diterima (Rp)</label
              >
              <input
                v-model.number="paidAmount"
                type="number"
                min="0"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs font-bold focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>

            <!-- Kembalian Info -->
            <div
              v-if="paymentMethod === 'TUNAI' && paidAmount > grandTotal"
              class="flex flex-col justify-end"
            >
              <div
                class="p-2 rounded-xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-800"
              >
                <span class="text-[10px] text-emerald-600 block">Kembalian:</span>
                <span class="font-extrabold text-sm text-emerald-700"
                  >Rp {{ formatRupiah(paidAmount - grandTotal) }}</span
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quantity Decimal Modal Sub-Dialog -->
      <div
        v-if="showQtyModal"
        class="fixed inset-0 z-60 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs"
      >
        <div
          class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4"
        >
          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h4 class="font-bold text-sm text-slate-900">{{ currentService?.service_name }}</h4>
              <p class="text-xs text-brand-600 font-extrabold">
                Rp {{ formatRupiah(currentService?.price) }} / {{ currentService?.unit }}
              </p>
            </div>
            <button class="p-1 text-slate-400 hover:text-slate-600" @click="showQtyModal = false">
              <X class="w-4 h-4" />
            </button>
          </div>

          <!-- Quantity Decimal Input -->
          <div>
            <label class="block text-xs font-bold text-slate-600 mb-1"
              >Kuantitas / Timbangan ({{ currentService?.unit }})</label
            >
            <div class="flex items-center gap-2">
              <button
                type="button"
                class="w-10 h-10 rounded-xl bg-slate-100 font-bold text-slate-700 hover:bg-slate-200"
                @click="qtyInput = Math.max(0.1, Number((qtyInput - 0.5).toFixed(2)))"
              >
                -
              </button>
              <input
                v-model.number="qtyInput"
                type="number"
                step="0.05"
                min="0.1"
                class="flex-1 py-2 text-center text-lg font-black rounded-xl border border-slate-200 focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
              <button
                type="button"
                class="w-10 h-10 rounded-xl bg-slate-100 font-bold text-slate-700 hover:bg-slate-200"
                @click="qtyInput = Number((qtyInput + 0.5).toFixed(2))"
              >
                +
              </button>
            </div>
          </div>

          <!-- Quick Presets -->
          <div class="flex gap-1.5 justify-center">
            <button
              v-for="p in [1, 2, 3, 5, 8]"
              :key="p"
              type="button"
              class="px-2.5 py-1 rounded-lg bg-slate-50 border border-slate-200 hover:bg-brand-50 text-xs font-bold text-slate-700"
              @click="qtyInput = p"
            >
              {{ p }} {{ currentService?.unit }}
            </button>
          </div>

          <!-- Subtotal Calculation -->
          <div class="p-2.5 rounded-xl bg-slate-50 flex justify-between items-center text-xs">
            <span class="text-slate-500">Subtotal:</span>
            <span class="font-extrabold text-sm text-slate-900"
              >Rp {{ formatRupiah(qtyInput * (currentService?.price || 0)) }}</span
            >
          </div>

          <button
            type="button"
            class="w-full py-2.5 rounded-xl bg-brand-600 hover:bg-brand-700 text-white font-bold text-xs shadow-sm transition-all"
            @click="addCurrentServiceToCart"
          >
            Tambahkan ke Keranjang
          </button>
        </div>
      </div>

      <!-- Footer Bar (Total & Submit) -->
      <div
        class="pt-3 border-t border-slate-100 flex flex-col sm:flex-row items-center justify-between gap-3 shrink-0"
      >
        <div class="flex items-center gap-4 w-full sm:w-auto justify-between sm:justify-start">
          <div>
            <span class="text-[10px] text-slate-400 block font-semibold">TOTAL TAGIHAN</span>
            <span class="font-black text-xl text-slate-900">Rp {{ formatRupiah(grandTotal) }}</span>
          </div>
          <div v-if="discountAmount > 0" class="text-xs text-emerald-600 font-bold">
            Hemat: Rp {{ formatRupiah(discountAmount) }}
          </div>
        </div>

        <div class="flex items-center gap-2 w-full sm:w-auto">
          <button
            type="button"
            class="w-1/3 sm:w-auto px-4 py-2.5 rounded-xl border border-slate-200 text-slate-700 font-bold text-xs hover:bg-slate-50"
            @click="close"
          >
            Batal
          </button>
          <button
            type="button"
            :disabled="creating || cartItems.length === 0"
            class="w-2/3 sm:w-auto px-6 py-2.5 rounded-xl bg-gradient-to-r from-brand-600 to-sky-500 hover:from-brand-700 hover:to-sky-600 text-white font-extrabold text-xs shadow-md shadow-brand-500/20 flex items-center justify-center gap-1.5 disabled:opacity-50 transition-all"
            @click="submitOrder"
          >
            <span v-if="creating">Menyimpan...</span>
            <span v-else>BUAT ORDERAN KASIR</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import {
  PlusCircle,
  User,
  ShoppingBag,
  Plus,
  Trash2,
  Sparkles,
  CreditCard,
  Search,
  X,
} from 'lucide-vue-next'
import { api, formatRupiah } from '../../api/client'
import { useToast } from '../../composables/useToast'

defineProps({
  modelValue: Boolean,
})

const emit = defineEmits(['update:modelValue', 'order-created'])
const toast = useToast()

const creating = ref(false)
const isNewCustomer = ref(false)
const customerSearchQuery = ref('')
const customerResults = ref([])
const selectedCustomer = ref(null)

const newCustomer = reactive({
  full_name: '',
  phone: '',
  address: '',
})

const services = ref([])
const activeCategory = ref('ALL')
const cartItems = ref([])

const showQtyModal = ref(false)
const currentService = ref(null)
const qtyInput = ref(1.0)

const orderConfig = reactive({
  perfume_name: 'Akasia',
  discount_type: 'NONE',
  custom_discount: 0,
  rack_number: 'A-01',
  estimated_completion_at: '',
  notes: '',
})

const paymentOption = ref('SEKARANG')
const paymentMethod = ref('TUNAI')
const paidAmount = ref(0)

const serviceCategories = [
  { id: 'ALL', label: 'Semua' },
  { id: 'TABUNG', label: 'Self Service' },
  { id: 'KILOAN', label: 'Cuci Setrika / Kiloan' },
  { id: 'SATUAN', label: 'Satuan & Bed Cover' },
  { id: 'METERAN', label: 'Karpet Permadani' },
]

const filteredServices = computed(() => {
  if (activeCategory.value === 'ALL') return services.value
  return services.value.filter((s) => s.category === activeCategory.value)
})

const subtotalCart = computed(() => {
  return cartItems.value.reduce((acc, it) => acc + (it.subtotal || 0), 0)
})

const discountAmount = computed(() => {
  if (orderConfig.discount_type === 'PROMO_10') {
    return Math.round(subtotalCart.value * 0.1)
  }
  if (orderConfig.discount_type === 'LOYALTY') {
    return Math.min(subtotalCart.value, 20000)
  }
  if (orderConfig.discount_type === 'CUSTOM') {
    return orderConfig.custom_discount || 0
  }
  return 0
})

const grandTotal = computed(() => {
  return Math.max(0, subtotalCart.value - discountAmount.value)
})

const close = () => {
  emit('update:modelValue', false)
  cartItems.value = []
}

const onCustomerSearch = async () => {
  if (!customerSearchQuery.value || customerSearchQuery.value.length < 2) {
    customerResults.value = []
    return
  }
  try {
    const data = await api.get('/customers', { search: customerSearchQuery.value })
    customerResults.value = data || []
  } catch {
    customerResults.value = []
  }
}

const selectCustomer = (c) => {
  selectedCustomer.value = c
  customerResults.value = []
  customerSearchQuery.value = ''
}

const openQuantityModal = (srv) => {
  currentService.value = srv
  qtyInput.value = srv.category === 'KILOAN' ? 4.0 : 1.0
  showQtyModal.value = true
}

const addCurrentServiceToCart = () => {
  if (!currentService.value || qtyInput.value <= 0) return

  const itemSubtotal = qtyInput.value * currentService.value.price
  cartItems.value.push({
    service_id: currentService.value.id,
    service_name: currentService.value.service_name,
    category: currentService.value.category,
    quantity: qtyInput.value,
    unit: currentService.value.unit,
    price_per_unit: currentService.value.price,
    subtotal: itemSubtotal,
    notes: '',
  })

  // Auto calculate SLA if not manually set
  recalculateSLA()
  showQtyModal.value = false
}

const removeItem = (idx) => {
  cartItems.value.splice(idx, 1)
  recalculateSLA()
}

const recalculateSLA = () => {
  const now = new Date()
  let hasSelfService = cartItems.value.some((it) => it.category === 'TABUNG')

  if (hasSelfService) {
    // 3 hours SLA
    now.setHours(now.getHours() + 3)
  } else {
    // default 3 days
    now.setDate(now.getDate() + 3)
  }

  // Format for datetime-local: YYYY-MM-DDTHH:MM
  const pad = (n) => String(n).padStart(2, '0')
  orderConfig.estimated_completion_at = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())}T${pad(now.getHours())}:${pad(now.getMinutes())}`

  // Set paid amount default if bayar sekarang
  if (paymentOption.value === 'SEKARANG') {
    paidAmount.value = grandTotal.value
  }
}

const applyDiscount = () => {
  if (paymentOption.value === 'SEKARANG') {
    paidAmount.value = grandTotal.value
  }
}

const loadMasterServices = async () => {
  try {
    const data = await api.get('/master/services')
    services.value = data || []
  } catch {
    services.value = []
  }
}

const submitOrder = async () => {
  // Validate customer
  let custID = ''
  let custName = ''
  let custPhone = ''
  let custAddress = ''

  if (isNewCustomer.value) {
    if (!newCustomer.full_name || !newCustomer.phone) {
      toast.warning('Nama pelanggan dan nomor WhatsApp wajib diisi!')
      return
    }
    custName = newCustomer.full_name
    custPhone = newCustomer.phone
    custAddress = newCustomer.address
  } else if (selectedCustomer.value) {
    custID = selectedCustomer.value.id
    custName = selectedCustomer.value.full_name
    custPhone = selectedCustomer.value.phone
    custAddress = selectedCustomer.value.address
  } else {
    toast.warning('Silakan pilih atau tambahkan pelanggan terlebih dahulu!')
    return
  }

  if (cartItems.value.length === 0) {
    toast.warning('Keranjang cucian masih kosong! Silakan pilih layanan.')
    return
  }

  creating.value = true
  try {
    const totalWeight = cartItems.value
      .filter((it) => it.unit === 'Kg')
      .reduce((acc, it) => acc + it.quantity, 0)

    let estDate = null
    if (orderConfig.estimated_completion_at) {
      estDate = new Date(orderConfig.estimated_completion_at).toISOString()
    }

    const payload = {
      customer_id: custID,
      customer_name: custName,
      customer_phone: custPhone,
      customer_address: custAddress,
      items: cartItems.value,
      perfume_name: orderConfig.perfume_name,
      discount_amount: discountAmount.value,
      total_weight: totalWeight,
      total_price: grandTotal.value,
      paid_amount: paymentOption.value === 'SEKARANG' ? paidAmount.value : 0,
      payment_method_name: paymentOption.value === 'SEKARANG' ? paymentMethod.value : '',
      rack_number: orderConfig.rack_number,
      estimated_completion_at: estDate,
      notes: orderConfig.notes,
      pickup_required: false,
      created_by: 'Kasir Utama',
    }

    const res = await api.post('/pos/orders', payload)
    toast.success('Transaksi kasir walk-in berhasil dibuat!')
    emit('order-created', res)
    close()
  } catch (err) {
    toast.error(err.message || 'Gagal membuat order.')
  } finally {
    creating.value = false
  }
}

onMounted(() => {
  loadMasterServices()
  recalculateSLA()
})
</script>
