<template>
  <AppLayout>
    <div class="bg-white rounded-2xl border border-slate-100 p-8 text-center space-y-3 shadow-xs">
      <div
        class="w-12 h-12 rounded-2xl bg-brand-50 text-brand-600 flex items-center justify-center mx-auto"
      >
        <ShoppingBag class="w-6 h-6" />
      </div>
      <h2 class="text-xl font-bold text-slate-900">Transaksi & Antrian Kasir</h2>
      <p class="text-sm text-slate-500 max-w-md mx-auto">
        Modul kasir walk-in, timbangan kuantitas desimal, 5 tab status antrian, dan cetak 3 nota
        thermal akan diaktifkan penuh pada <strong>Milestone 3</strong>.
      </p>
      <RouterLink
        to="/"
        class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-slate-900 text-white text-xs font-bold hover:bg-slate-800 transition-all"
      >
        Kembali ke Beranda
      </RouterLink>
    </div>
  </AppLayout>
</template>

<script setup>
import { ShoppingBag } from 'lucide-vue-next'
import AppLayout from '../components/layout/AppLayout.vue'
</script>
