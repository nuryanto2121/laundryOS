<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-xs"
  >
    <div
      class="bg-white rounded-2xl max-w-lg w-full p-6 shadow-2xl border border-slate-100 flex flex-col max-h-[85vh] space-y-4"
    >
      <!-- Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-100">
        <div class="flex items-center gap-2.5">
          <div
            class="w-9 h-9 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center"
          >
            <Camera class="w-5 h-5" />
          </div>
          <div>
            <h3 class="font-bold text-base text-slate-900">Dokumentasi Foto Pakaian</h3>
            <p class="text-xs text-slate-500 font-mono">
              {{ formatDisplayOrder(order.order_number) }} • {{ order.customer_name }}
            </p>
          </div>
        </div>
        <button class="p-1.5 rounded-lg text-slate-400 hover:text-slate-600" @click="close">
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Upload / Capture Area -->
      <div
        class="p-4 rounded-xl border-2 border-dashed border-slate-200 hover:border-brand-400 transition-all text-center space-y-2 bg-slate-50"
      >
        <input
          id="cameraInput"
          type="file"
          accept="image/*"
          capture="environment"
          class="hidden"
          @change="handleFileUpload"
        />
        <label for="cameraInput" class="cursor-pointer flex flex-col items-center">
          <div
            class="w-10 h-10 rounded-full bg-brand-50 text-brand-600 flex items-center justify-center mb-1"
          >
            <Camera class="w-5 h-5" />
          </div>
          <span class="text-xs font-bold text-brand-600 hover:underline"
            >Ambil Foto / Pilih Gambar</span
          >
          <span class="text-[10px] text-slate-400">Kamera HP atau file foto galeri</span>
        </label>
      </div>

      <!-- Preview Selected File before upload -->
      <div v-if="previewUrl" class="p-3 rounded-xl border border-slate-200 bg-white space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold text-slate-700">Preview Foto Baru:</span>
          <button class="text-xs text-rose-500 font-semibold hover:underline" @click="clearPreview">
            Batal
          </button>
        </div>
        <img
          :src="previewUrl"
          class="w-full h-40 object-cover rounded-lg border border-slate-100"
        />
        <div class="grid grid-cols-2 gap-2">
          <button
            type="button"
            class="py-1.5 rounded-lg text-xs font-bold border transition-all"
            :class="
              photoType === 'MASUK'
                ? 'bg-blue-50 text-blue-700 border-blue-400'
                : 'bg-slate-50 text-slate-600 border-slate-200'
            "
            @click="photoType = 'MASUK'"
          >
            Dokumentasi Masuk
          </button>
          <button
            type="button"
            class="py-1.5 rounded-lg text-xs font-bold border transition-all"
            :class="
              photoType === 'KELUAR'
                ? 'bg-emerald-50 text-emerald-700 border-emerald-400'
                : 'bg-slate-50 text-slate-600 border-slate-200'
            "
            @click="photoType = 'KELUAR'"
          >
            Dokumentasi Keluar
          </button>
        </div>
        <input
          v-model="photoNotes"
          type="text"
          placeholder="Catatan foto (misal: Noda minyak di kerah)..."
          class="w-full px-3 py-1.5 rounded-lg border border-slate-200 text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
        />
        <button
          :disabled="uploading"
          class="w-full py-2 rounded-xl bg-brand-600 hover:bg-brand-700 text-white font-bold text-xs flex items-center justify-center gap-1.5 disabled:opacity-50"
          @click="submitPhoto"
        >
          <span v-if="uploading">Mengunggah...</span>
          <span v-else>Simpan Foto</span>
        </button>
      </div>

      <!-- Existing Photos Gallery -->
      <div class="flex-1 overflow-y-auto space-y-2">
        <h4 class="text-xs font-bold text-slate-500 uppercase tracking-wider">
          Foto Tersimpan ({{ existingPhotos.length }})
        </h4>

        <div v-if="existingPhotos.length === 0" class="text-center py-6 text-slate-400 text-xs">
          Belum ada foto dokumentasi untuk order ini.
        </div>

        <div v-else class="grid grid-cols-2 sm:grid-cols-3 gap-2">
          <div
            v-for="photo in existingPhotos"
            :key="photo.id"
            class="relative rounded-xl overflow-hidden border border-slate-200 group"
          >
            <img
              :src="photo.photo_url"
              class="w-full h-28 object-cover group-hover:scale-105 transition-transform"
            />
            <div class="absolute top-1.5 left-1.5">
              <span
                class="text-[9px] font-extrabold px-1.5 py-0.5 rounded-md text-white uppercase shadow-xs"
                :class="photo.photo_type === 'MASUK' ? 'bg-blue-600' : 'bg-emerald-600'"
              >
                {{ photo.photo_type }}
              </span>
            </div>
            <p
              v-if="photo.notes"
              class="text-[10px] text-slate-600 truncate px-1 py-0.5 bg-white border-t border-slate-100"
            >
              {{ photo.notes }}
            </p>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="pt-3 border-t border-slate-100 text-right">
        <button
          class="px-4 py-2 rounded-xl bg-slate-900 text-white text-xs font-bold hover:bg-slate-800 transition-all"
          @click="close"
        >
          Selesai
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Camera, X } from 'lucide-vue-next'
import { api } from '../../api/client'

const props = defineProps({
  modelValue: Boolean,
  order: {
    type: Object,
    required: true,
  },
  photos: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:modelValue', 'photo-added'])

const previewUrl = ref('')
const photoType = ref('MASUK')
const photoNotes = ref('')
const uploading = ref(false)
const existingPhotos = ref([])

watch(
  () => props.photos,
  (val) => {
    existingPhotos.value = val || []
  },
  { immediate: true }
)

const formatDisplayOrder = (num) => {
  if (!num) return ''
  return num.startsWith('ORD') ? '#' + num.substring(3) : num
}

const handleFileUpload = (e) => {
  const file = e.target.files[0]
  if (!file) return

  const reader = new FileReader()
  reader.onload = (event) => {
    previewUrl.value = event.target.result
  }
  reader.readAsDataURL(file)
}

const clearPreview = () => {
  previewUrl.value = ''
  photoNotes.value = ''
}

const close = () => {
  emit('update:modelValue', false)
  clearPreview()
}

const submitPhoto = async () => {
  if (!previewUrl.value) return

  uploading.value = true
  try {
    const res = await api.post(`/pos/orders/${props.order.id}/photos`, {
      photo_type: photoType.value,
      photo_url: previewUrl.value,
      notes: photoNotes.value,
    })
    existingPhotos.value.push(res)
    emit('photo-added')
    clearPreview()
  } catch (err) {
    alert(err.message || 'Gagal menyimpan foto.')
  } finally {
    uploading.value = false
  }
}
</script>
