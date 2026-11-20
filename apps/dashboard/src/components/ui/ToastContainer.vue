<template>
  <div
    class="fixed top-4 right-4 z-50 flex flex-col gap-2.5 max-w-sm w-full px-3 sm:px-0 pointer-events-none"
  >
    <TransitionGroup
      enter-active-class="transition duration-300 ease-out transform"
      enter-from-class="opacity-0 translate-y-[-12px] scale-95"
      enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition duration-200 ease-in transform"
      leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 translate-x-4 scale-95"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="pointer-events-auto flex items-start gap-3 p-4 rounded-2xl bg-white shadow-xl border border-slate-100/80 backdrop-blur-md transition-all group overflow-hidden relative"
        :class="getToastBorderClass(toast.type)"
      >
        <!-- Accent indicator line on left -->
        <div
          class="absolute left-0 top-0 bottom-0 w-1.5"
          :class="getAccentBarClass(toast.type)"
        ></div>

        <!-- Icon Avatar -->
        <div
          class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0 shadow-xs"
          :class="getIconAvatarClass(toast.type)"
        >
          <component :is="getToastIcon(toast.type)" class="w-5 h-5" />
        </div>

        <!-- Text Content -->
        <div class="flex-1 min-w-0 pr-2">
          <h4 class="text-xs font-bold text-slate-900 tracking-tight leading-tight">
            {{ toast.title }}
          </h4>
          <p class="text-xs text-slate-600 mt-0.5 leading-relaxed break-words">
            {{ toast.message }}
          </p>
        </div>

        <!-- Close Button -->
        <button
          class="p-1 rounded-lg text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors shrink-0 cursor-pointer"
          @click="removeToast(toast.id)"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup>
import { CheckCircle2, AlertCircle, AlertTriangle, Info, X } from 'lucide-vue-next'
import { useToast } from '../../composables/useToast'

const { toasts, removeToast } = useToast()

const getToastIcon = (type) => {
  switch (type) {
    case 'success':
      return CheckCircle2
    case 'error':
      return AlertCircle
    case 'warning':
      return AlertTriangle
    case 'info':
    default:
      return Info
  }
}

const getAccentBarClass = (type) => {
  switch (type) {
    case 'success':
      return 'bg-emerald-500'
    case 'error':
      return 'bg-rose-500'
    case 'warning':
      return 'bg-amber-500'
    case 'info':
    default:
      return 'bg-brand-500'
  }
}

const getIconAvatarClass = (type) => {
  switch (type) {
    case 'success':
      return 'bg-emerald-50 text-emerald-600'
    case 'error':
      return 'bg-rose-50 text-rose-600'
    case 'warning':
      return 'bg-amber-50 text-amber-600'
    case 'info':
    default:
      return 'bg-brand-50 text-brand-600'
  }
}

const getToastBorderClass = (type) => {
  switch (type) {
    case 'success':
      return 'hover:border-emerald-200'
    case 'error':
      return 'hover:border-rose-200'
    case 'warning':
      return 'hover:border-amber-200'
    case 'info':
    default:
      return 'hover:border-brand-200'
  }
}
</script>
