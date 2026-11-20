<template>
  <div
    v-if="totalItems > 0"
    class="flex flex-col sm:flex-row items-center justify-between gap-3 pt-3 border-t border-slate-100 text-xs text-slate-500"
  >
    <!-- Left: Display count info & Page size selector -->
    <div class="flex items-center gap-3">
      <span>
        Menampilkan
        <strong class="text-slate-800">{{ startIndex }}</strong>
        -
        <strong class="text-slate-800">{{ endIndex }}</strong>
        dari
        <strong class="text-slate-800">{{ totalItems }}</strong>
        data
      </span>

      <div class="flex items-center gap-1.5">
        <span class="text-slate-400">Baris:</span>
        <select
          :value="itemsPerPage"
          class="px-2 py-1 rounded-lg bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-700 focus:outline-none focus:border-brand-500 cursor-pointer"
          @change="changePageSize($event.target.value)"
        >
          <option v-for="size in pageSizeOptions" :key="size" :value="size">
            {{ size }}
          </option>
        </select>
      </div>
    </div>

    <!-- Right: Page buttons -->
    <div class="flex items-center gap-1">
      <!-- Previous -->
      <button
        :disabled="currentPage <= 1"
        class="px-2.5 py-1 rounded-lg border border-slate-200 bg-white text-slate-700 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-all font-semibold cursor-pointer"
        @click="goToPage(currentPage - 1)"
      >
        &larr; Prev
      </button>

      <!-- Page Numbers -->
      <button
        v-for="p in visiblePages"
        :key="p"
        class="w-7 h-7 rounded-lg text-xs font-black transition-all cursor-pointer flex items-center justify-center"
        :class="
          currentPage === p
            ? 'bg-slate-900 text-white shadow-xs'
            : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
        "
        @click="goToPage(p)"
      >
        {{ p }}
      </button>

      <!-- Next -->
      <button
        :disabled="currentPage >= totalPages"
        class="px-2.5 py-1 rounded-lg border border-slate-200 bg-white text-slate-700 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-all font-semibold cursor-pointer"
        @click="goToPage(currentPage + 1)"
      >
        Next &rarr;
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  totalItems: {
    type: Number,
    required: true,
  },
  currentPage: {
    type: Number,
    default: 1,
  },
  itemsPerPage: {
    type: Number,
    default: 10,
  },
  pageSizeOptions: {
    type: Array,
    default: () => [10, 25, 50],
  },
})

const emit = defineEmits(['update:currentPage', 'update:itemsPerPage'])

const totalPages = computed(() => {
  return Math.ceil(props.totalItems / props.itemsPerPage) || 1
})

const startIndex = computed(() => {
  if (props.totalItems === 0) return 0
  return (props.currentPage - 1) * props.itemsPerPage + 1
})

const endIndex = computed(() => {
  return Math.min(props.currentPage * props.itemsPerPage, props.totalItems)
})

const visiblePages = computed(() => {
  const pages = []
  const maxButtons = 5
  let start = Math.max(1, props.currentPage - Math.floor(maxButtons / 2))
  let end = Math.min(totalPages.value, start + maxButtons - 1)

  if (end - start + 1 < maxButtons) {
    start = Math.max(1, end - maxButtons + 1)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value && page !== props.currentPage) {
    emit('update:currentPage', page)
  }
}

const changePageSize = (newSize) => {
  emit('update:itemsPerPage', Number(newSize))
  emit('update:currentPage', 1)
}
</script>
