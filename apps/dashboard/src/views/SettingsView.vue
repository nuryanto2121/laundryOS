<template>
  <AppLayout>
    <div class="bg-white rounded-2xl border border-slate-100 p-8 text-center space-y-3 shadow-xs">
      <div
        class="w-12 h-12 rounded-2xl bg-brand-50 text-brand-600 flex items-center justify-center mx-auto"
      >
        <Settings class="w-6 h-6" />
      </div>
      <h2 class="text-xl font-bold text-slate-900">Pengaturan & Master Data</h2>
      <p class="text-sm text-slate-500 max-w-md mx-auto">
        Manajemen hak akses pegawai (RBAC privasi omset), custom printer thermal 58mm/80mm, master
        parfum, dan diskon loyalitas akan diaktifkan pada <strong>Milestone 6</strong>.
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
import { Settings } from 'lucide-vue-next'
import AppLayout from '../components/layout/AppLayout.vue'
</script>
