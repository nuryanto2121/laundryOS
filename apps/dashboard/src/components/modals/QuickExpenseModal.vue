<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-3 sm:p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
  >
    <div
      class="bg-white rounded-2xl max-w-lg w-full p-5 sm:p-6 shadow-2xl border border-slate-100 flex flex-col max-h-[92vh] overflow-y-auto"
    >
      <!-- Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-100 shrink-0">
        <div class="flex items-center gap-2.5">
          <div
            class="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shadow-xs"
          >
            <Receipt class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-black text-base sm:text-lg text-slate-900">
              Catat Kas Kecil (Petty Cash)
            </h3>
            <p class="text-xs text-slate-500">Pencatatan biaya operasional kasir & outlet</p>
          </div>
        </div>
        <button
          class="p-1.5 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-all cursor-pointer"
          @click="close"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Form -->
      <form class="py-4 space-y-4 text-slate-800" @submit.prevent="submitExpense">
        <!-- Tanggal & Dicatat Oleh -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1">
              Tanggal Pengeluaran *
            </label>
            <input
              v-model="form.expense_date"
              type="date"
              required
              class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white text-xs font-semibold text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-brand-500 transition-all"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1">
              Dicatat Oleh *
            </label>
            <input
              v-model="form.created_by"
              type="text"
              required
              placeholder="Kasir Utama"
              class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white text-xs font-semibold text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-brand-500 transition-all"
            />
          </div>
        </div>

        <!-- Category -->
        <div>
          <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1">
            Kategori Pengeluaran *
          </label>
          <select
            v-model="form.category"
            required
            class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white text-xs font-bold text-slate-800 focus:outline-hidden focus:ring-2 focus:ring-brand-500 transition-all cursor-pointer"
          >
            <option value="BAHAN_BAKU">🧴 Bahan Baku (Deterjen, Pewangi, Plastik, Softener)</option>
            <option value="GAS">🔥 Gas Mesin Pengering (Dryer Elpiji)</option>
            <option value="LISTRIK_AIR">⚡ Listrik & Air (Token PLN / PDAM)</option>
            <option value="GAJI">👥 Gaji Pegawai / Staf Operasional</option>
            <option value="BONUS">🍱 Uang Makan / Bonus Pegawai</option>
            <option value="MAINTENANCE">🔧 Perawatan & Servis Mesin Laundry</option>
            <option value="SEWA">🏠 Sewa Tempat / Ruko Laundry</option>
            <option value="LAIN_LAIN">📦 Biaya Lain-Lain Operasional</option>
          </select>
        </div>

        <!-- Amount -->
        <div>
          <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1">
            Nominal Pengeluaran (Rp) *
          </label>
          <div class="relative">
            <span class="absolute left-3.5 top-2.5 text-sm font-black text-slate-400">Rp</span>
            <input
              v-model.number="form.amount"
              type="number"
              min="500"
              step="500"
              required
              placeholder="Contoh: 35000"
              class="w-full pl-11 pr-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white text-slate-900 font-black text-base focus:outline-hidden focus:ring-2 focus:ring-brand-500 transition-all"
            />
          </div>

          <!-- Quick Preset Amounts -->
          <div class="flex items-center gap-1.5 flex-wrap mt-2">
            <button
              v-for="amt in [10000, 25000, 50000, 100000, 200000, 500000]"
              :key="amt"
              type="button"
              class="px-2.5 py-1 rounded-lg text-xs font-bold transition-all cursor-pointer"
              :class="
                form.amount === amt
                  ? 'bg-amber-600 text-white shadow-2xs'
                  : 'bg-slate-100 text-slate-700 hover:bg-slate-200'
              "
              @click="form.amount = amt"
            >
              {{ amt >= 1000000 ? amt / 1000000 + 'jt' : amt / 1000 + 'rb' }}
            </button>
          </div>
        </div>

        <!-- Description -->
        <div>
          <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1">
            Keterangan Keperluan *
          </label>
          <textarea
            v-model="form.description"
            rows="2"
            required
            placeholder="Contoh: Beli deterjen cair 5 liter dan plastik jinjing 3 pack"
            class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white text-slate-900 text-xs font-medium focus:outline-hidden focus:ring-2 focus:ring-brand-500 transition-all resize-none"
          ></textarea>
        </div>

        <!-- Error Alert -->
        <p
          v-if="errorMessage"
          class="text-xs text-rose-600 font-semibold bg-rose-50 p-2.5 rounded-xl border border-rose-200"
        >
          {{ errorMessage }}
        </p>

        <!-- Actions -->
        <div class="flex items-center gap-2 pt-2 border-t border-slate-100">
          <button
            type="button"
            class="w-1/3 py-2.5 rounded-xl border border-slate-200 text-slate-700 font-bold text-xs hover:bg-slate-50 transition-all cursor-pointer"
            @click="close"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="w-2/3 py-2.5 rounded-xl bg-gradient-to-r from-amber-600 to-amber-700 hover:from-amber-700 hover:to-amber-800 text-white font-extrabold text-xs shadow-md shadow-amber-600/20 disabled:opacity-50 transition-all flex items-center justify-center gap-1.5 cursor-pointer"
          >
            <span v-if="loading">Menyimpan...</span>
            <span v-else>Simpan Pengeluaran</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import { Receipt, X } from 'lucide-vue-next'
import { api } from '../../api/client'
import { useToast } from '../../composables/useToast'

const props = defineProps({
  modelValue: Boolean,
})

const emit = defineEmits(['update:modelValue', 'expense-added'])
const toast = useToast()

const loading = ref(false)
const errorMessage = ref('')

const getTodayDateStr = () => {
  const d = new Date()
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

const form = reactive({
  category: 'BAHAN_BAKU',
  amount: null,
  description: '',
  created_by: 'Kasir Utama',
  expense_date: getTodayDateStr(),
})

const close = () => {
  emit('update:modelValue', false)
  form.amount = null
  form.description = ''
  form.expense_date = getTodayDateStr()
  errorMessage.value = ''
}

const submitExpense = async () => {
  if (!form.amount || form.amount <= 0 || !form.description.trim()) {
    errorMessage.value = 'Silakan isi nominal (>0) dan keterangan dengan benar.'
    return
  }

  loading.value = true
  errorMessage.value = ''
  try {
    await api.post('/expenses', {
      category: form.category,
      amount: form.amount,
      description: form.description.trim(),
      created_by: form.created_by,
      expense_date: form.expense_date,
    })
    toast.success('Pengeluaran kas kecil berhasil dicatat!')
    emit('expense-added')
    close()
  } catch (err) {
    errorMessage.value = err.message || 'Gagal menyimpan pengeluaran.'
    toast.error(errorMessage.value)
  } finally {
    loading.value = false
  }
}

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      form.expense_date = getTodayDateStr()
      errorMessage.value = ''
    }
  }
)
</script>
