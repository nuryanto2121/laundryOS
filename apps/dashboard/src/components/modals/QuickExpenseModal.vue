<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs animate-in fade-in"
  >
    <div class="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl border border-slate-100">
      <!-- Header -->
      <div class="flex items-center justify-between pb-4 border-b border-slate-100">
        <div class="flex items-center gap-2.5">
          <div
            class="w-9 h-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center"
          >
            <Receipt class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-bold text-base text-slate-900">Catat Kas Kecil (Pengeluaran)</h3>
            <p class="text-xs text-slate-500">Biaya operasional harian FAN Laundry</p>
          </div>
        </div>
        <button class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600" @click="close">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Form -->
      <form class="py-4 space-y-4" @submit.prevent="submitExpense">
        <!-- Category -->
        <div>
          <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1.5"
            >Kategori Pengeluaran</label
          >
          <select
            v-model="form.category"
            required
            class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-slate-900 text-sm focus:outline-hidden focus:ring-2 focus:ring-brand-500 focus:bg-white transition-all font-medium"
          >
            <option value="BAHAN_BAKU">Bahan Baku (Deterjen, Pewangi, Plastik)</option>
            <option value="LISTRIK_AIR">Listrik & Air</option>
            <option value="GAS">Gas Mesin Pengering</option>
            <option value="GAJI">Gaji Pegawai / Staf</option>
            <option value="BONUS">Bonus / Uang Makan</option>
            <option value="SEWA">Sewa Tempat / Ruko</option>
            <option value="MAINTENANCE">Perawatan / Servis Mesin</option>
            <option value="LAIN_LAIN">Biaya Lain-Lain</option>
          </select>
        </div>

        <!-- Amount -->
        <div>
          <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1.5"
            >Nominal (Rp)</label
          >
          <div class="relative">
            <span class="absolute left-3.5 top-2.5 text-sm font-bold text-slate-400">Rp</span>
            <input
              v-model.number="form.amount"
              type="number"
              min="100"
              step="500"
              required
              placeholder="Contoh: 50000"
              class="w-full pl-11 pr-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-slate-900 font-bold text-base focus:outline-hidden focus:ring-2 focus:ring-brand-500 focus:bg-white transition-all"
            />
          </div>
        </div>

        <!-- Description -->
        <div>
          <label class="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1.5"
            >Keterangan Pengeluaran</label
          >
          <textarea
            v-model="form.description"
            rows="2"
            required
            placeholder="Contoh: Beli sabun deterjen cair 5 Liter"
            class="w-full px-3.5 py-2.5 rounded-xl border border-slate-200 bg-slate-50 text-slate-900 text-sm focus:outline-hidden focus:ring-2 focus:ring-brand-500 focus:bg-white transition-all resize-none"
          ></textarea>
        </div>

        <!-- Error Alert -->
        <p
          v-if="errorMessage"
          class="text-xs text-rose-600 font-medium bg-rose-50 p-2.5 rounded-xl border border-rose-200"
        >
          {{ errorMessage }}
        </p>

        <!-- Actions -->
        <div class="flex items-center gap-2 pt-2">
          <button
            type="button"
            class="w-1/2 py-2.5 rounded-xl border border-slate-200 text-slate-700 font-semibold text-sm hover:bg-slate-100 transition-all"
            @click="close"
          >
            Batal
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="w-1/2 py-2.5 rounded-xl bg-slate-900 text-white font-semibold text-sm hover:bg-slate-800 disabled:opacity-50 transition-all flex items-center justify-center gap-1.5 shadow-sm"
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
import { reactive, ref } from 'vue'
import { Receipt, X } from 'lucide-vue-next'
import { api } from '../../api/client'

defineProps({
  modelValue: Boolean,
})

const emit = defineEmits(['update:modelValue', 'expense-added'])

const loading = ref(false)
const errorMessage = ref('')

const form = reactive({
  category: 'BAHAN_BAKU',
  amount: null,
  description: '',
  created_by: 'Kasir 1',
})

const close = () => {
  emit('update:modelValue', false)
  form.amount = null
  form.description = ''
  errorMessage.value = ''
}

const submitExpense = async () => {
  if (!form.amount || form.amount <= 0 || !form.description.trim()) {
    errorMessage.value = 'Silakan isi nominal dan keterangan dengan benar.'
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
      expense_date: new Date().toISOString().split('T')[0],
    })
    emit('expense-added')
    close()
  } catch (err) {
    errorMessage.value = err.message || 'Gagal menyimpan pengeluaran.'
  } finally {
    loading.value = false
  }
}
</script>
