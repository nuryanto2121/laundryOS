<template>
  <AppLayout>
    <div class="space-y-6">
      <!-- 1. Header & Navigation -->
      <div
        class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-4 sm:p-6 rounded-2xl border border-slate-100 shadow-xs"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center shadow-xs"
          >
            <BarChart3 class="w-5 h-5" />
          </div>
          <div>
            <h1 class="text-xl sm:text-2xl font-black text-slate-900 tracking-tight">
              Laporan Keuangan, Operasional & CRM
            </h1>
            <p class="text-xs sm:text-sm text-slate-500 mt-0.5">
              Pembukuan kas, rekap layanan, piutang tagihan, dan loyalitas pelanggan FAN Laundry
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2 self-start sm:self-auto flex-wrap">
          <!-- Refresh -->
          <button
            :disabled="loading"
            class="p-2.5 rounded-xl bg-slate-50 hover:bg-slate-100 border border-slate-200 text-slate-600 hover:text-slate-900 transition-all disabled:opacity-50 cursor-pointer"
            title="Refresh Data"
            @click="loadAllReports"
          >
            <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          </button>

          <!-- Export CSV -->
          <button
            class="flex items-center gap-1.5 px-3.5 py-2.5 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-700 font-bold text-xs sm:text-sm transition-all cursor-pointer"
            title="Download CSV Laporan Aktif"
            @click="exportCurrentTabCSV"
          >
            <Download class="w-4 h-4 text-slate-500" />
            <span class="hidden sm:inline">Export CSV</span>
          </button>

          <!-- Catat Pengeluaran Button -->
          <button
            class="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-amber-600 to-amber-700 hover:from-amber-700 hover:to-amber-800 text-white font-bold text-xs sm:text-sm shadow-md shadow-amber-600/20 active:scale-98 transition-all cursor-pointer"
            @click="showCreateExpenseModal = true"
          >
            <PlusCircle class="w-4 h-4" />
            <span>+ Catat Pengeluaran</span>
          </button>
        </div>
      </div>

      <!-- 2. Global Date Filter Toolbar (Default: Hari Ini + Manual Custom Range) -->
      <div class="bg-white p-3.5 sm:p-4 rounded-2xl border border-slate-100 shadow-xs space-y-3">
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
          <!-- Quick Date Range Pills -->
          <div
            class="flex items-center gap-1.5 overflow-x-auto pb-1 md:pb-0 text-xs scrollbar-none"
          >
            <span class="text-xs font-bold text-slate-400 mr-1 hidden sm:inline">Periode:</span>
            <button
              v-for="p in periodPresets"
              :key="p.id"
              type="button"
              class="px-3.5 py-2 rounded-xl font-bold transition-all whitespace-nowrap cursor-pointer text-xs"
              :class="
                selectedPeriod === p.id
                  ? 'bg-slate-900 text-white shadow-xs'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              "
              @click="applyPeriod(p.id)"
            >
              {{ p.label }}
            </button>
          </div>

          <!-- Date Display & Custom Toggle -->
          <div class="flex items-center gap-2 text-xs font-semibold text-slate-600 shrink-0">
            <Calendar class="w-4 h-4 text-brand-600 shrink-0" />
            <span>{{ formatDate(startDate) }}</span>
            <span>&mdash;</span>
            <span>{{ formatDate(endDate) }}</span>
          </div>
        </div>

        <!-- Custom Date Manual Inputs (Visible when 'CUSTOM' is selected) -->
        <div
          v-if="selectedPeriod === 'CUSTOM'"
          class="pt-3 border-t border-slate-100 flex flex-col sm:flex-row items-stretch sm:items-center gap-3 bg-slate-50/70 p-3 rounded-xl"
        >
          <span class="text-xs font-bold text-slate-700">Pilih Rentang Tanggal Manual:</span>
          <div class="flex items-center gap-2 flex-1">
            <input
              v-model="customStartDate"
              type="date"
              class="px-3 py-1.5 rounded-lg bg-white border border-slate-200 text-xs font-semibold text-slate-800 focus:outline-none focus:border-brand-500"
            />
            <span class="text-xs text-slate-400">s/d</span>
            <input
              v-model="customEndDate"
              type="date"
              class="px-3 py-1.5 rounded-lg bg-white border border-slate-200 text-xs font-semibold text-slate-800 focus:outline-none focus:border-brand-500"
            />
          </div>
          <button
            class="px-4 py-1.5 rounded-lg bg-brand-600 hover:bg-brand-700 text-white font-bold text-xs shadow-xs transition-all cursor-pointer text-center"
            @click="applyCustomDateRange"
          >
            Terapkan Filter
          </button>
        </div>
      </div>

      <!-- 3. TOP 4 KPI CARDS: Omset, Pendapatan Kas, Pengeluaran, Laba Bersih -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
        <!-- 1. Omset (Penjualan Bruto) -->
        <div
          class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between hover:border-blue-200 transition-all"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
              >Omset (Bruto)</span
            >
            <div
              class="w-9 h-9 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center"
            >
              <TrendingUp class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-3">
            <h3 class="text-lg sm:text-2xl font-black text-slate-900 tracking-tight">
              Rp {{ formatRupiah(financialReport.total_omset) }}
            </h3>
            <p class="text-[11px] text-slate-400 mt-1">Nilai total pesanan dibuat</p>
          </div>
        </div>

        <!-- 2. Pendapatan (Kas/Transfer Masuk Riil) -->
        <div
          class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between hover:border-emerald-200 transition-all"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
              >Pendapatan Kas</span
            >
            <div
              class="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center"
            >
              <Wallet class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-3">
            <h3 class="text-lg sm:text-2xl font-black text-emerald-600 tracking-tight">
              Rp {{ formatRupiah(financialReport.total_pendapatan) }}
            </h3>
            <p class="text-[11px] text-slate-400 mt-1">Uang riil diterima kasir</p>
          </div>
        </div>

        <!-- 3. Pengeluaran (Buku Kas Kecil) -->
        <div
          class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between hover:border-rose-200 transition-all"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
              >Pengeluaran</span
            >
            <div
              class="w-9 h-9 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center"
            >
              <Receipt class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-3">
            <h3 class="text-lg sm:text-2xl font-black text-rose-600 tracking-tight">
              - Rp {{ formatRupiah(financialReport.total_pengeluaran) }}
            </h3>
            <p class="text-[11px] text-slate-400 mt-1">Kas kecil operasional</p>
          </div>
        </div>

        <!-- 4. Laba Bersih Operasional -->
        <div
          class="bg-white p-4 sm:p-5 rounded-2xl border border-slate-100 shadow-xs flex flex-col justify-between hover:border-indigo-200 transition-all"
        >
          <div class="flex items-center justify-between">
            <span class="text-xs font-bold text-slate-500 uppercase tracking-wider"
              >Laba Bersih</span
            >
            <div
              class="w-9 h-9 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center"
            >
              <DollarSign class="w-5 h-5" />
            </div>
          </div>
          <div class="mt-3">
            <h3
              class="text-lg sm:text-2xl font-black tracking-tight"
              :class="financialReport.laba_bersih >= 0 ? 'text-indigo-600' : 'text-rose-600'"
            >
              Rp {{ formatRupiah(financialReport.laba_bersih) }}
            </h3>
            <p class="text-[11px] text-slate-400 mt-1">Pendapatan Kas - Biaya</p>
          </div>
        </div>
      </div>

      <!-- 4. TIGA KATEGORI LAPORAN UTAMA (Main Tabs) -->
      <div class="flex items-center gap-2 overflow-x-auto pb-1 sm:pb-0 scrollbar-none">
        <!-- Tab 1: Laporan Transaksi -->
        <button
          type="button"
          class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs sm:text-sm font-bold transition-all whitespace-nowrap shrink-0 cursor-pointer"
          :class="
            activeTab === 'TRANSACTIONS'
              ? 'bg-slate-900 text-white shadow-xs'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200/80'
          "
          @click="activeTab = 'TRANSACTIONS'"
        >
          <Package class="w-4 h-4" />
          <span>Laporan Transaksi</span>
          <span
            class="px-2 py-0.5 rounded-full text-[10px] font-extrabold"
            :class="
              activeTab === 'TRANSACTIONS'
                ? 'bg-white/20 text-white'
                : 'bg-slate-100 text-slate-700'
            "
          >
            {{ periodOrders.length }} Order
          </span>
        </button>

        <!-- Tab 2: Laporan Keuangan & Kas -->
        <button
          type="button"
          class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs sm:text-sm font-bold transition-all whitespace-nowrap shrink-0 cursor-pointer"
          :class="
            activeTab === 'FINANCIAL'
              ? 'bg-slate-900 text-white shadow-xs'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200/80'
          "
          @click="activeTab = 'FINANCIAL'"
        >
          <Wallet class="w-4 h-4" />
          <span>Laporan Keuangan & Kas</span>
          <span
            v-if="unpaidOrders.length > 0"
            class="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-rose-500 text-white"
            title="Order belum lunas"
          >
            {{ unpaidOrders.length }} Piutang
          </span>
        </button>

        <!-- Tab 3: Laporan Pelanggan (CRM) -->
        <button
          type="button"
          class="flex items-center gap-2 px-4 py-2.5 rounded-xl text-xs sm:text-sm font-bold transition-all whitespace-nowrap shrink-0 cursor-pointer"
          :class="
            activeTab === 'CRM'
              ? 'bg-slate-900 text-white shadow-xs'
              : 'bg-white text-slate-600 hover:bg-slate-50 border border-slate-200/80'
          "
          @click="activeTab = 'CRM'"
        >
          <Users class="w-4 h-4" />
          <span>Laporan Pelanggan (CRM)</span>
          <span
            class="px-2 py-0.5 rounded-full text-[10px] font-extrabold"
            :class="
              activeTab === 'CRM' ? 'bg-white/20 text-white' : 'bg-purple-100 text-purple-700'
            "
          >
            {{ allCustomers.length }} Pelanggan
          </span>
        </button>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 1: LAPORAN TRANSAKSI DENGAN 5 SUB-TAB                      -->
      <!-- ============================================================== -->
      <div v-if="activeTab === 'TRANSACTIONS'" class="space-y-4">
        <!-- Sub-Navigasi Transaksi -->
        <div
          class="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0 scrollbar-none bg-slate-100/80 p-1.5 rounded-2xl"
        >
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              txSubTab === 'ALL'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="txSubTab = 'ALL'"
          >
            Semua Data Transaksi ({{ periodOrders.length }})
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              txSubTab === 'SERVICES'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="txSubTab = 'SERVICES'"
          >
            Ringkasan Perlayanan
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              txSubTab === 'PICKED_UP'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="txSubTab = 'PICKED_UP'"
          >
            Diambil Customer ({{ completedOrders.length }})
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              txSubTab === 'READY_UNCOLLECTED'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="txSubTab = 'READY_UNCOLLECTED'"
          >
            Belum Diambil / di Rak ({{ readyOrders.length }})
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              txSubTab === 'CANCELLED'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="txSubTab = 'CANCELLED'"
          >
            Transaksi Batal ({{ cancelledOrders.length }})
          </button>
        </div>

        <!-- 1.a. Semua Data Transaksi -->
        <div
          v-if="txSubTab === 'ALL'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div
            class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3"
          >
            <div>
              <h3 class="text-sm font-bold text-slate-800 uppercase tracking-wider">
                Semua Transaksi Masuk
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Daftar seluruh order laundry pada periode terpilih
              </p>
            </div>
            <!-- Search -->
            <div class="relative min-w-[200px]">
              <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
              <input
                v-model="orderSearchQuery"
                type="text"
                placeholder="Cari #order / pelanggan..."
                class="w-full pl-9 pr-7 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-medium text-slate-800 focus:outline-none focus:border-brand-500 focus:bg-white"
                @input="pageOrders = 1"
              />
              <button
                v-if="orderSearchQuery"
                class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                @click="clearOrderSearch"
              >
                <X class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <div
            v-if="filteredPeriodOrders.length === 0"
            class="py-10 text-center text-slate-400 text-xs"
          >
            Tidak ada transaksi pada periode ini.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">No. Order</th>
                    <th class="py-3 px-4">Pelanggan</th>
                    <th class="py-3 px-4">Tanggal Masuk</th>
                    <th class="py-3 px-4">Status Pengerjaan</th>
                    <th class="py-3 px-4">Pembayaran</th>
                    <th class="py-3 px-4 text-right">Total (Rp)</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="order in paginatedPeriodOrders"
                    :key="order.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <td class="py-3.5 px-4 font-black text-slate-900 whitespace-nowrap">
                      #{{ order.order_number }}
                    </td>
                    <td class="py-3.5 px-4">
                      <p class="font-bold text-slate-800">
                        {{ order.customer?.full_name || 'Walk-in' }}
                      </p>
                      <p class="text-[11px] text-slate-400">{{ order.customer?.phone || '-' }}</p>
                    </td>
                    <td class="py-3.5 px-4 text-slate-500 whitespace-nowrap">
                      {{ formatDateTime(order.created_at) }}
                    </td>
                    <td class="py-3.5 px-4 whitespace-nowrap">
                      <span
                        class="px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-wider"
                        :class="getOrderStatusBadgeClass(order.status)"
                      >
                        {{ order.status }}
                      </span>
                    </td>
                    <td class="py-3.5 px-4 whitespace-nowrap">
                      <span
                        class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                        :class="
                          order.payment_status === 'PAID'
                            ? 'bg-emerald-100 text-emerald-800'
                            : 'bg-rose-100 text-rose-800'
                        "
                      >
                        {{ order.payment_status === 'PAID' ? 'LUNAS' : 'BELUM LUNAS' }}
                      </span>
                    </td>
                    <td
                      class="py-3.5 px-4 text-right font-black text-slate-900 text-sm whitespace-nowrap"
                    >
                      Rp {{ formatRupiah(order.total_price) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pageOrders"
              v-model:items-per-page="pageSizeOrders"
              :total-items="filteredPeriodOrders.length"
            />
          </div>
        </div>

        <!-- 1.b. Ringkasan Transaksi Perlayanan -->
        <div
          v-else-if="txSubTab === 'SERVICES'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex items-center justify-between">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <Package class="w-4 h-4 text-brand-600" />
                <span>Rekapitulasi Volume & Pendapatan Per Layanan</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Akumulasi kiloan, tabung self-service, dan satuan yang dikerjakan
              </p>
            </div>
          </div>

          <div
            v-if="serviceBreakdown.length === 0"
            class="py-10 text-center text-slate-400 text-xs"
          >
            Belum ada rekapitulasi layanan pada periode ini.
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full text-left border-collapse text-xs">
              <thead>
                <tr
                  class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                >
                  <th class="py-3 px-4">Kategori Layanan</th>
                  <th class="py-3 px-4">Satuan</th>
                  <th class="py-3 px-4 text-right">Volume Total</th>
                  <th class="py-3 px-4 text-right">Frekuensi Order</th>
                  <th class="py-3 px-4 text-right">Total Pendapatan (Rp)</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="item in serviceBreakdown"
                  :key="item.category + item.unit"
                  class="hover:bg-slate-50 transition-colors"
                >
                  <td class="py-3.5 px-4 font-bold text-slate-900 flex items-center gap-2">
                    <span class="text-base">{{ getServiceCategoryIcon(item.category) }}</span>
                    <span>{{ formatCategoryName(item.category) }}</span>
                  </td>
                  <td class="py-3.5 px-4 font-semibold text-slate-600">
                    <span
                      class="px-2 py-0.5 rounded-md bg-slate-100 font-bold uppercase text-[10px]"
                    >
                      {{ item.unit }}
                    </span>
                  </td>
                  <td class="py-3.5 px-4 text-right font-black text-slate-900 text-sm">
                    {{ Number(item.total_volume).toFixed(item.unit === 'KG' ? 2 : 0) }}
                    {{ item.unit }}
                  </td>
                  <td class="py-3.5 px-4 text-right font-semibold text-slate-600">
                    {{ item.total_order }} pesanan
                  </td>
                  <td class="py-3.5 px-4 text-right font-black text-emerald-600 text-sm">
                    Rp {{ formatRupiah(item.total_pendapatan) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- 1.c. Transaksi Diambil Customer -->
        <div
          v-else-if="txSubTab === 'PICKED_UP'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex items-center justify-between">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <CheckCircle2 class="w-4 h-4 text-emerald-600" />
                <span>Transaksi Selesai & Diambil Pelanggan</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Daftar cucian yang statusnya sudah SELESAI
              </p>
            </div>
            <span
              class="px-3 py-1 rounded-full text-xs font-bold bg-emerald-50 text-emerald-700 border border-emerald-200"
            >
              {{ completedOrders.length }} Order
            </span>
          </div>

          <div v-if="completedOrders.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Belum ada transaksi berstatus selesai pada rentang tanggal ini.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">No. Order</th>
                    <th class="py-3 px-4">Pelanggan</th>
                    <th class="py-3 px-4">Tanggal Order</th>
                    <th class="py-3 px-4">Status Pembayaran</th>
                    <th class="py-3 px-4 text-right">Total Transaksi</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="order in paginatedCompletedOrders"
                    :key="order.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <td class="py-3.5 px-4 font-black text-slate-900 whitespace-nowrap">
                      #{{ order.order_number }}
                    </td>
                    <td class="py-3.5 px-4">
                      <p class="font-bold text-slate-800">
                        {{ order.customer?.full_name || 'Pelanggan' }}
                      </p>
                      <p class="text-[11px] text-slate-400">{{ order.customer?.phone || '-' }}</p>
                    </td>
                    <td class="py-3.5 px-4 text-slate-500 whitespace-nowrap font-medium">
                      {{ formatDate(order.created_at) }}
                    </td>
                    <td class="py-3.5 px-4 whitespace-nowrap">
                      <span
                        class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                        :class="
                          order.payment_status === 'PAID'
                            ? 'bg-emerald-100 text-emerald-800'
                            : 'bg-amber-100 text-amber-800'
                        "
                      >
                        {{ order.payment_status }}
                      </span>
                    </td>
                    <td class="py-3.5 px-4 text-right font-black text-slate-900 whitespace-nowrap">
                      Rp {{ formatRupiah(order.total_price) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pagePickedUp"
              v-model:items-per-page="pageSizePickedUp"
              :total-items="completedOrders.length"
            />
          </div>
        </div>

        <!-- 1.d. Transaksi Belum Diambil (di Rak) -->
        <div
          v-else-if="txSubTab === 'READY_UNCOLLECTED'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <Clock class="w-4 h-4 text-amber-600" />
                <span>Transaksi Belum Diambil (Menunggu di Rak)</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Pakaian siap ambil di rak outlet yang belum diambil pelanggan
              </p>
            </div>
            <span
              class="px-2.5 py-1 rounded-full text-xs font-bold bg-amber-50 text-amber-700 border border-amber-200"
            >
              {{ readyOrders.length }} Order di Rak
            </span>
          </div>

          <div v-if="readyOrders.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Semua cucian yang siap sudah diambil pelanggan.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">No. Order</th>
                    <th class="py-3 px-4">Pelanggan</th>
                    <th class="py-3 px-4">Nomor Rak</th>
                    <th class="py-3 px-4">Status Bayar</th>
                    <th class="py-3 px-4 text-right">Total Tagihan</th>
                    <th class="py-3 px-4 text-center">Aksi Reminder</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="order in paginatedReadyOrders"
                    :key="order.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <td class="py-3.5 px-4 font-black text-slate-900 whitespace-nowrap">
                      #{{ order.order_number }}
                    </td>
                    <td class="py-3.5 px-4">
                      <p class="font-bold text-slate-800">
                        {{ order.customer?.full_name || 'Pelanggan' }}
                      </p>
                      <p class="text-[11px] text-slate-400">{{ order.customer?.phone }}</p>
                    </td>
                    <td class="py-3.5 px-4 whitespace-nowrap">
                      <span
                        class="px-2.5 py-1 rounded-lg text-xs font-black uppercase tracking-wider"
                        :class="
                          order.rack_number
                            ? 'bg-amber-100 text-amber-800 border border-amber-300'
                            : 'bg-slate-100 text-slate-500'
                        "
                      >
                        📍 {{ order.rack_number || 'Belum di Rak' }}
                      </span>
                    </td>
                    <td class="py-3.5 px-4 whitespace-nowrap">
                      <span
                        class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider"
                        :class="
                          order.payment_status === 'PAID'
                            ? 'bg-emerald-100 text-emerald-800'
                            : 'bg-rose-100 text-rose-800'
                        "
                      >
                        {{ order.payment_status === 'PAID' ? 'LUNAS' : 'BELUM LUNAS' }}
                      </span>
                    </td>
                    <td class="py-3.5 px-4 text-right font-black text-slate-900 whitespace-nowrap">
                      Rp {{ formatRupiah(order.total_price) }}
                    </td>
                    <td class="py-3.5 px-4 text-center whitespace-nowrap">
                      <a
                        :href="getReadyOrderWhatsappUrl(order)"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-emerald-600 hover:bg-emerald-700 text-white font-bold text-xs shadow-xs transition-all cursor-pointer"
                        title="Kirim Pesan WhatsApp Pengambilan"
                      >
                        <MessageCircle class="w-3.5 h-3.5" />
                        <span>Ingatkan WA</span>
                      </a>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pageReady"
              v-model:items-per-page="pageSizeReady"
              :total-items="readyOrders.length"
            />
          </div>
        </div>

        <!-- 1.e. Semua Transaksi Batal -->
        <div
          v-else-if="txSubTab === 'CANCELLED'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex items-center justify-between">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <AlertTriangle class="w-4 h-4 text-rose-600" />
                <span>Semua Transaksi Batal</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Daftar order yang dibatalkan oleh kasir atau pelanggan
              </p>
            </div>
            <span
              class="px-3 py-1 rounded-full text-xs font-bold bg-rose-50 text-rose-700 border border-rose-200"
            >
              {{ cancelledOrders.length }} Batal
            </span>
          </div>

          <div v-if="cancelledOrders.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Bagus sekali! Tidak ada transaksi yang dibatalkan pada periode ini.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">No. Order</th>
                    <th class="py-3 px-4">Pelanggan</th>
                    <th class="py-3 px-4">Waktu Batal</th>
                    <th class="py-3 px-4">Catatan</th>
                    <th class="py-3 px-4 text-right">Nominal Dibatalkan</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="order in paginatedCancelledOrders"
                    :key="order.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <td class="py-3.5 px-4 font-black text-rose-600 line-through whitespace-nowrap">
                      #{{ order.order_number }}
                    </td>
                    <td class="py-3.5 px-4">
                      <p class="font-bold text-slate-800">
                        {{ order.customer?.full_name || 'Pelanggan' }}
                      </p>
                      <p class="text-[11px] text-slate-400">{{ order.customer?.phone || '-' }}</p>
                    </td>
                    <td class="py-3.5 px-4 text-slate-500 whitespace-nowrap font-medium">
                      {{ formatDateTime(order.created_at) }}
                    </td>
                    <td class="py-3.5 px-4 text-slate-500 italic">
                      {{ order.notes || 'Dibatalkan kasir' }}
                    </td>
                    <td class="py-3.5 px-4 text-right font-black text-rose-600 whitespace-nowrap">
                      Rp {{ formatRupiah(order.total_price) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pageCancelled"
              v-model:items-per-page="pageSizeCancelled"
              :total-items="cancelledOrders.length"
            />
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 2: LAPORAN KEUANGAN & KAS DENGAN 5 SUB-TAB                 -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'FINANCIAL'" class="space-y-4">
        <!-- Sub-Navigasi Keuangan -->
        <div
          class="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0 scrollbar-none bg-slate-100/80 p-1.5 rounded-2xl"
        >
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              finSubTab === 'OMSET'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="finSubTab = 'OMSET'"
          >
            Laporan Omset (Bruto)
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              finSubTab === 'PENDAPATAN'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="finSubTab = 'PENDAPATAN'"
          >
            Laporan Pendapatan Kas
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              finSubTab === 'PENGELUARAN'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="finSubTab = 'PENGELUARAN'"
          >
            Laporan Pengeluaran ({{ filteredExpenses.length }})
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              finSubTab === 'PIUTANG'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="finSubTab = 'PIUTANG'"
          >
            Piutang ({{ unpaidOrders.length }})
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              finSubTab === 'LABA_RUGI'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="finSubTab = 'LABA_RUGI'"
          >
            Laporan Laba / Rugi
          </button>
        </div>

        <!-- 2.a. Laporan Omset -->
        <div
          v-if="finSubTab === 'OMSET'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex items-center justify-between">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <TrendingUp class="w-4 h-4 text-blue-600" />
                <span>Rincian Omset Penjualan Kotor</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Daftar transaksi aktif yang membentuk omset periode terpilih
              </p>
            </div>
            <div class="text-right">
              <span class="text-xs text-slate-400 block">Total Omset</span>
              <span class="text-lg font-black text-blue-600"
                >Rp {{ formatRupiah(financialReport.total_omset) }}</span
              >
            </div>
          </div>

          <div v-if="omsetOrders.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Belum ada transaksi omset pada rentang tanggal ini.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">No. Order</th>
                    <th class="py-3 px-4">Pelanggan</th>
                    <th class="py-3 px-4">Tanggal Masuk</th>
                    <th class="py-3 px-4">Status</th>
                    <th class="py-3 px-4 text-right">Nilai Omset (Rp)</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="order in paginatedOmsetOrders"
                    :key="order.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <td class="py-3.5 px-4 font-black text-slate-900 whitespace-nowrap">
                      #{{ order.order_number }}
                    </td>
                    <td class="py-3.5 px-4 font-bold text-slate-800">
                      {{ order.customer?.full_name || 'Walk-in' }}
                    </td>
                    <td class="py-3.5 px-4 text-slate-500 whitespace-nowrap">
                      {{ formatDate(order.created_at) }}
                    </td>
                    <td class="py-3.5 px-4">
                      <span
                        class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase"
                        :class="getOrderStatusBadgeClass(order.status)"
                      >
                        {{ order.status }}
                      </span>
                    </td>
                    <td
                      class="py-3.5 px-4 text-right font-black text-slate-900 text-sm whitespace-nowrap"
                    >
                      Rp {{ formatRupiah(order.total_price) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pageOmset"
              v-model:items-per-page="pageSizeOmset"
              :total-items="omsetOrders.length"
            />
          </div>
        </div>

        <!-- 2.b. Laporan Pendapatan -->
        <div
          v-else-if="finSubTab === 'PENDAPATAN'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex items-center justify-between">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <Wallet class="w-4 h-4 text-emerald-600" />
                <span>Penerimaan Kas Riil Masuk</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Uang kas & transfer yang sudah diterima kasir
              </p>
            </div>
            <div class="text-right">
              <span class="text-xs text-slate-400 block">Total Pendapatan</span>
              <span class="text-lg font-black text-emerald-600"
                >Rp {{ formatRupiah(financialReport.total_pendapatan) }}</span
              >
            </div>
          </div>

          <div v-if="paidOrders.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Belum ada penerimaan pembayaran pada periode ini.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">No. Order</th>
                    <th class="py-3 px-4">Pelanggan</th>
                    <th class="py-3 px-4">Metode Bayar</th>
                    <th class="py-3 px-4">Status</th>
                    <th class="py-3 px-4 text-right">Uang Diterima (Rp)</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="order in paginatedPaidOrders"
                    :key="order.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <td class="py-3.5 px-4 font-black text-slate-900 whitespace-nowrap">
                      #{{ order.order_number }}
                    </td>
                    <td class="py-3.5 px-4 font-bold text-slate-800">
                      {{ order.customer?.full_name || 'Walk-in' }}
                    </td>
                    <td class="py-3.5 px-4 font-bold text-slate-700 whitespace-nowrap">
                      {{ order.payment_method || 'CASH / TUNAI' }}
                    </td>
                    <td class="py-3.5 px-4">
                      <span
                        class="px-2 py-0.5 rounded-full text-[10px] font-black uppercase bg-emerald-100 text-emerald-800"
                      >
                        {{ order.payment_status }}
                      </span>
                    </td>
                    <td
                      class="py-3.5 px-4 text-right font-black text-emerald-600 text-sm whitespace-nowrap"
                    >
                      Rp {{ formatRupiah(order.paid_amount || order.total_price) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pagePaid"
              v-model:items-per-page="pageSizePaid"
              :total-items="paidOrders.length"
            />
          </div>
        </div>

        <!-- 2.c. Laporan Pengeluaran (Buku Kas Kecil) -->
        <div v-else-if="finSubTab === 'PENGELUARAN'" class="space-y-4">
          <!-- Komposisi Pengeluaran per Kategori -->
          <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4">
            <div class="flex items-center justify-between">
              <h3
                class="text-xs font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <PieChart class="w-4 h-4 text-amber-600" />
                <span>Komposisi Pengeluaran per Kategori Biaya</span>
              </h3>
              <span class="text-[11px] text-slate-400"
                >{{ categoryBreakdown.length }} Kategori Aktif</span
              >
            </div>

            <div
              v-if="categoryBreakdown.length === 0"
              class="py-6 text-center text-slate-400 text-xs"
            >
              Belum ada catatan pengeluaran operasional pada periode ini.
            </div>

            <div v-else class="grid grid-cols-2 sm:grid-cols-4 gap-2.5">
              <div
                v-for="cat in categoryBreakdown"
                :key="cat.category"
                class="p-3 rounded-xl border border-slate-200/80 bg-slate-50/50 flex flex-col justify-between"
              >
                <div>
                  <div class="flex items-center gap-1.5">
                    <span>{{ getCategoryIcon(cat.category) }}</span>
                    <span class="font-bold text-xs text-slate-800 truncate">
                      {{ getCategoryLabel(cat.category) }}
                    </span>
                  </div>
                  <span class="text-[10px] text-slate-400 block mt-0.5"
                    >{{ cat.count }} transaksi</span
                  >
                </div>
                <div class="mt-2">
                  <span class="text-xs sm:text-sm font-black text-slate-900 block">
                    Rp {{ formatRupiah(cat.total_amount) }}
                  </span>
                  <div class="w-full bg-slate-200 h-1.5 rounded-full overflow-hidden mt-1.5">
                    <div
                      class="bg-amber-500 h-full rounded-full"
                      :style="{ width: `${cat.percent}%` }"
                    ></div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Tabel Buku Kas Kecil -->
          <div
            class="bg-white rounded-2xl border border-slate-100 shadow-xs overflow-hidden space-y-4 p-5"
          >
            <div
              class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3"
            >
              <div>
                <h3
                  class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
                >
                  <Receipt class="w-4 h-4 text-amber-600" />
                  <span>Buku Kas Kecil Operasional Outlet</span>
                </h3>
                <p class="text-xs text-slate-400 mt-0.5">
                  Rincian belanja bahan baku, operasional mesin & kasbon
                </p>
              </div>

              <!-- Filters -->
              <div class="flex items-center gap-2 flex-wrap">
                <select
                  v-model="selectedCategoryFilter"
                  class="px-3 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-700 focus:outline-none focus:border-brand-500 cursor-pointer"
                  @change="pageExpenses = 1"
                >
                  <option value="ALL">Semua Kategori</option>
                  <option value="BAHAN_BAKU">🧴 Bahan Baku</option>
                  <option value="GAS">🔥 Gas Dryer</option>
                  <option value="LISTRIK_AIR">⚡ Listrik & Air</option>
                  <option value="GAJI">👥 Gaji Pegawai</option>
                  <option value="BONUS">🍱 Uang Makan</option>
                  <option value="MAINTENANCE">🔧 Servis Mesin</option>
                  <option value="SEWA">🏠 Sewa Tempat</option>
                  <option value="LAIN_LAIN">📦 Lain-Lain</option>
                </select>

                <div class="relative min-w-[160px]">
                  <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    v-model="expenseSearchQuery"
                    type="text"
                    placeholder="Cari keterangan..."
                    class="w-full pl-9 pr-7 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-medium text-slate-800 placeholder-slate-400 focus:outline-none focus:border-brand-500 focus:bg-white"
                    @input="pageExpenses = 1"
                  />
                  <button
                    v-if="expenseSearchQuery"
                    class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                    @click="clearExpenseSearch"
                  >
                    <X class="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            </div>

            <div
              v-if="filteredExpenses.length === 0"
              class="py-8 text-center text-slate-400 text-xs"
            >
              Tidak ada catatan pengeluaran yang cocok dengan filter.
            </div>

            <div v-else class="space-y-4">
              <div class="overflow-x-auto">
                <table class="w-full text-left border-collapse text-xs">
                  <thead>
                    <tr
                      class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                    >
                      <th class="py-3 px-4">Tanggal</th>
                      <th class="py-3 px-4">Kategori</th>
                      <th class="py-3 px-4">Keterangan Keperluan</th>
                      <th class="py-3 px-4">Dicatat Oleh</th>
                      <th class="py-3 px-4 text-right">Nominal (Rp)</th>
                      <th class="py-3 px-4 text-center">Aksi</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-100">
                    <tr
                      v-for="item in paginatedExpenses"
                      :key="item.id"
                      class="hover:bg-slate-50 transition-colors"
                    >
                      <td class="py-3.5 px-4 font-semibold text-slate-900 whitespace-nowrap">
                        {{ formatDate(item.expense_date) }}
                      </td>
                      <td class="py-3.5 px-4 whitespace-nowrap">
                        <span
                          class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold"
                          :class="getCategoryBadgeClass(item.category)"
                        >
                          <span>{{ getCategoryIcon(item.category) }}</span>
                          <span>{{ getCategoryLabel(item.category) }}</span>
                        </span>
                      </td>
                      <td class="py-3.5 px-4">
                        <p class="font-bold text-slate-800">{{ item.description }}</p>
                      </td>
                      <td class="py-3.5 px-4 text-slate-500 whitespace-nowrap font-medium">
                        {{ item.created_by || 'Kasir' }}
                      </td>
                      <td class="py-3.5 px-4 text-right whitespace-nowrap">
                        <span class="font-black text-sm text-rose-600">
                          - Rp {{ formatRupiah(item.amount) }}
                        </span>
                      </td>
                      <td class="py-3.5 px-4 text-center whitespace-nowrap">
                        <button
                          class="p-1.5 rounded-lg text-slate-400 hover:text-rose-600 hover:bg-rose-50 transition-all cursor-pointer"
                          title="Hapus Pengeluaran"
                          @click="promptDeleteExpense(item)"
                        >
                          <Trash2 class="w-4 h-4" />
                        </button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <!-- Pagination -->
              <PaginationBar
                v-model:current-page="pageExpenses"
                v-model:items-per-page="pageSizeExpenses"
                :total-items="filteredExpenses.length"
              />
            </div>
          </div>
        </div>

        <!-- 2.d. Laporan Piutang Tagihan -->
        <div
          v-else-if="finSubTab === 'PIUTANG'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <ShieldAlert class="w-4 h-4 text-rose-600" />
                <span>Laporan Piutang Belum Lunas</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Daftar transaksi pelanggan yang masih memiliki sisa pembayaran
              </p>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-xs font-bold text-rose-600">
                Total Piutang: Rp {{ formatRupiah(totalPiutangNominal) }}
              </span>
            </div>
          </div>

          <div v-if="unpaidOrders.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Luar biasa! Tidak ada piutang yang tertunggak. Semua order lunas.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">No. Order</th>
                    <th class="py-3 px-4">Nama Pelanggan</th>
                    <th class="py-3 px-4">Tanggal Order</th>
                    <th class="py-3 px-4 text-right">Total Tagihan</th>
                    <th class="py-3 px-4 text-right">Sudah Dibayar</th>
                    <th class="py-3 px-4 text-right">Sisa Piutang</th>
                    <th class="py-3 px-4 text-center">Tagih WA</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="item in paginatedUnpaidOrders"
                    :key="item.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <td class="py-3.5 px-4 font-black text-slate-900 whitespace-nowrap">
                      #{{ item.order_number }}
                    </td>
                    <td class="py-3.5 px-4 whitespace-nowrap">
                      <p class="font-bold text-slate-800">{{ item.customer_name }}</p>
                      <p class="text-[11px] text-slate-400">{{ item.customer_phone }}</p>
                    </td>
                    <td class="py-3.5 px-4 text-slate-500 whitespace-nowrap font-medium">
                      {{ formatDate(item.created_at) }}
                    </td>
                    <td
                      class="py-3.5 px-4 text-right font-semibold text-slate-700 whitespace-nowrap"
                    >
                      Rp {{ formatRupiah(item.total_price) }}
                    </td>
                    <td
                      class="py-3.5 px-4 text-right font-semibold text-emerald-600 whitespace-nowrap"
                    >
                      Rp {{ formatRupiah(item.paid_amount) }}
                    </td>
                    <td class="py-3.5 px-4 text-right font-black text-rose-600 whitespace-nowrap">
                      Rp {{ formatRupiah(item.sisa_tagihan) }}
                    </td>
                    <td class="py-3.5 px-4 text-center whitespace-nowrap">
                      <a
                        :href="getUnpaidWhatsappUrl(item)"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="inline-flex items-center gap-1 px-3 py-1.5 rounded-xl bg-rose-50 hover:bg-rose-100 text-rose-700 border border-rose-200 font-bold text-xs transition-all cursor-pointer"
                        title="Kirim Tagihan WhatsApp"
                      >
                        <MessageCircle class="w-3.5 h-3.5" />
                        <span>Tagih WA</span>
                      </a>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pageUnpaid"
              v-model:items-per-page="pageSizeUnpaid"
              :total-items="unpaidOrders.length"
            />
          </div>
        </div>

        <!-- 2.e. Laporan Laba/Rugi -->
        <div
          v-else-if="finSubTab === 'LABA_RUGI'"
          class="bg-white p-6 rounded-2xl border border-slate-100 shadow-xs space-y-6"
        >
          <div class="border-b border-slate-100 pb-4">
            <h3 class="text-base font-black text-slate-900 flex items-center gap-2">
              <DollarSign class="w-5 h-5 text-indigo-600" />
              <span>Laporan Laba / Rugi Kas Operasional</span>
            </h3>
            <p class="text-xs text-slate-400 mt-1">
              Periode {{ formatDate(startDate) }} s/d {{ formatDate(endDate) }} (FAN Laundry
              Kalibaru)
            </p>
          </div>

          <!-- Income Statement Grid -->
          <div
            class="max-w-xl mx-auto space-y-3 bg-slate-50/60 p-5 rounded-2xl border border-slate-100 text-sm"
          >
            <!-- Inflow -->
            <div class="flex items-center justify-between pb-2 border-b border-slate-200/60">
              <span class="font-bold text-slate-700">1. Total Pendapatan Kas Diterima (+)</span>
              <span class="font-black text-emerald-600"
                >Rp {{ formatRupiah(financialReport.total_pendapatan) }}</span
              >
            </div>
            <!-- Outflow -->
            <div class="flex items-center justify-between pb-2 border-b border-slate-200/60">
              <span class="font-bold text-slate-700">2. Total Pengeluaran Kas Kecil (-)</span>
              <span class="font-black text-rose-600"
                >- Rp {{ formatRupiah(financialReport.total_pengeluaran) }}</span
              >
            </div>
            <!-- Net Profit -->
            <div class="flex items-center justify-between pt-2">
              <span class="text-base font-black text-slate-900">Laba Bersih Operasional (=)</span>
              <span
                class="text-xl font-black"
                :class="financialReport.laba_bersih >= 0 ? 'text-indigo-600' : 'text-rose-600'"
              >
                Rp {{ formatRupiah(financialReport.laba_bersih) }}
              </span>
            </div>
            <!-- Profit Margin -->
            <div class="flex items-center justify-between text-xs text-slate-500 pt-1">
              <span>Margin Keuntungan Kas</span>
              <span class="font-bold text-slate-700"> {{ profitMarginPercent }}% </span>
            </div>
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 3: LAPORAN PELANGGAN (CRM) DENGAN 4 SUB-TAB                -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'CRM'" class="space-y-4">
        <!-- Sub-Navigasi CRM -->
        <div
          class="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0 scrollbar-none bg-slate-100/80 p-1.5 rounded-2xl"
        >
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              crmSubTab === 'SUMMARY'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="crmSubTab = 'SUMMARY'"
          >
            Ringkasan Pelanggan
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              crmSubTab === 'GROWTH'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="crmSubTab = 'GROWTH'"
          >
            Pertumbuhan Pelanggan Baru ({{ newCustomersInPeriod.length }})
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              crmSubTab === 'DETAIL'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="crmSubTab = 'DETAIL'"
          >
            Detail & Rekap Aktivitas ({{ filteredCustomerWithStats.length }})
          </button>
          <button
            type="button"
            class="px-3 py-1.5 rounded-xl text-xs font-bold transition-all whitespace-nowrap cursor-pointer"
            :class="
              crmSubTab === 'TOP_CUSTOMERS'
                ? 'bg-white text-slate-900 shadow-xs'
                : 'text-slate-600 hover:text-slate-900'
            "
            @click="crmSubTab = 'TOP_CUSTOMERS'"
          >
            Top Customer (Loyal)
          </button>
        </div>

        <!-- 3.a. Ringkasan Pelanggan -->
        <div v-if="crmSubTab === 'SUMMARY'" class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
          <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-1">
            <span class="text-xs font-bold text-slate-400 uppercase tracking-wider"
              >Total Pelanggan Terdaftar</span
            >
            <h4 class="text-2xl font-black text-slate-900">{{ allCustomers.length }} Orang</h4>
            <p class="text-[11px] text-slate-400">Database pelanggan aktif</p>
          </div>
          <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-1">
            <span class="text-xs font-bold text-slate-400 uppercase tracking-wider"
              >Pelanggan Baru Periode Ini</span
            >
            <h4 class="text-2xl font-black text-purple-600">
              {{ newCustomersInPeriod.length }} Baru
            </h4>
            <p class="text-[11px] text-slate-400">Bergabung di periode ini</p>
          </div>
          <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-1">
            <span class="text-xs font-bold text-slate-400 uppercase tracking-wider"
              >Pelanggan Bertransaksi</span
            >
            <h4 class="text-2xl font-black text-emerald-600">{{ activeCustomerCount }} Orang</h4>
            <p class="text-[11px] text-slate-400">Memiliki order aktif di periode</p>
          </div>
          <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-1">
            <span class="text-xs font-bold text-slate-400 uppercase tracking-wider"
              >Rata-rata Belanja / Order</span
            >
            <h4 class="text-2xl font-black text-brand-600">Rp {{ formatRupiah(avgOrderValue) }}</h4>
            <p class="text-[11px] text-slate-400">Nilai keranjang rata-rata</p>
          </div>
        </div>

        <!-- 3.b. Pertumbuhan Pelanggan Baru -->
        <div
          v-else-if="crmSubTab === 'GROWTH'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex items-center justify-between">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <UserPlus class="w-4 h-4 text-purple-600" />
                <span>Pelanggan Baru Mendaftar di Periode Ini</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Pertumbuhan akuisisi pelanggan baru FAN Laundry
              </p>
            </div>
            <span
              class="px-3 py-1 rounded-full text-xs font-bold bg-purple-50 text-purple-700 border border-purple-200"
            >
              {{ newCustomersInPeriod.length }} Pendaftar
            </span>
          </div>

          <div
            v-if="newCustomersInPeriod.length === 0"
            class="py-8 text-center text-slate-400 text-xs"
          >
            Tidak ada pelanggan baru yang mendaftar pada rentang tanggal ini.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">Kode Pelanggan</th>
                    <th class="py-3 px-4">Nama Lengkap</th>
                    <th class="py-3 px-4">No. WhatsApp</th>
                    <th class="py-3 px-4">Tanggal Daftar</th>
                    <th class="py-3 px-4 text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="c in paginatedGrowthCustomers"
                    :key="c.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <td class="py-3.5 px-4 font-black text-purple-700">
                      {{ c.customer_code || '-' }}
                    </td>
                    <td class="py-3.5 px-4 font-bold text-slate-900">
                      {{ c.full_name }}
                    </td>
                    <td class="py-3.5 px-4 text-slate-600">
                      {{ c.phone }}
                    </td>
                    <td class="py-3.5 px-4 text-slate-500 whitespace-nowrap">
                      {{ formatDateTime(c.created_at) }}
                    </td>
                    <td class="py-3.5 px-4 text-center whitespace-nowrap">
                      <button
                        class="px-3 py-1.5 rounded-xl bg-purple-50 hover:bg-purple-100 text-purple-700 font-bold text-xs border border-purple-200 cursor-pointer"
                        @click="openCustomerCRM(c)"
                      >
                        Lihat CRM
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pageGrowth"
              v-model:items-per-page="pageSizeGrowth"
              :total-items="newCustomersInPeriod.length"
            />
          </div>
        </div>

        <!-- 3.c. Detail & Rekap Aktivitas Pelanggan (dengan filter Layanan, Total Kg, Uang & Paging) -->
        <div
          v-else-if="crmSubTab === 'DETAIL'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div
            class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3"
          >
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <Users class="w-4 h-4 text-purple-600" />
                <span>Detail & Rekap Aktivitas Pelanggan (Periode Ini)</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Pantau total transaksi, akumulasi Kg, nominal belanja, dan layanan yang dipakai tiap
                customer
              </p>
            </div>

            <!-- Filters & Search Toolbar -->
            <div class="flex items-center gap-2 flex-wrap">
              <!-- Service Category Filter -->
              <select
                v-model="customerServiceFilter"
                class="px-3 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-700 focus:outline-none focus:border-brand-500 cursor-pointer"
                @change="pageCustomers = 1"
              >
                <option value="ALL">Semua Layanan</option>
                <option value="KILOAN">🧺 Cuci Kiloan</option>
                <option value="SELF_SERVICE">🪙 Self Service</option>
                <option value="SATUAN">👔 Cuci Satuan</option>
                <option value="SEPATU_TAS">👟 Sepatu & Tas</option>
                <option value="KARPET">🧶 Cuci Karpet</option>
              </select>

              <!-- Customer Search -->
              <div class="relative min-w-[180px]">
                <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
                <input
                  v-model="customerSearchQuery"
                  type="text"
                  placeholder="Cari nama / no. hp / kode..."
                  class="w-full pl-9 pr-7 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-medium text-slate-800 focus:outline-none focus:border-brand-500 focus:bg-white"
                  @input="pageCustomers = 1"
                />
                <button
                  v-if="customerSearchQuery"
                  class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                  @click="clearCustomerSearch"
                >
                  <X class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>

          <div
            v-if="filteredCustomerWithStats.length === 0"
            class="py-8 text-center text-slate-400 text-xs"
          >
            Tidak ada pelanggan yang cocok dengan kriteria filter.
          </div>

          <div v-else class="space-y-4">
            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse text-xs">
                <thead>
                  <tr
                    class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                  >
                    <th class="py-3 px-4">Pelanggan</th>
                    <th class="py-3 px-4 text-center">Order (Periode Ini)</th>
                    <th class="py-3 px-4 text-right">Berat (Periode Ini)</th>
                    <th class="py-3 px-4 text-right">Belanja (Periode Ini)</th>
                    <th class="py-3 px-4">Layanan yang Digunakan</th>
                    <th class="py-3 px-4 text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-100">
                  <tr
                    v-for="item in paginatedCustomers"
                    :key="item.customer.id"
                    class="hover:bg-slate-50 transition-colors"
                  >
                    <!-- Pelanggan info -->
                    <td class="py-3.5 px-4">
                      <div class="flex items-center gap-1.5">
                        <span class="font-black text-slate-900">{{ item.customer.full_name }}</span>
                        <span
                          v-if="item.customer.customer_code"
                          class="px-1.5 py-0.2 rounded text-[10px] font-black bg-purple-50 text-purple-700 border border-purple-200/60"
                        >
                          {{ item.customer.customer_code }}
                        </span>
                      </div>
                      <p class="text-[11px] text-slate-400 mt-0.5">
                        {{ item.customer.phone || '-' }}
                      </p>
                    </td>

                    <!-- Order count -->
                    <td
                      class="py-3.5 px-4 text-center font-black text-slate-800 text-xs whitespace-nowrap"
                    >
                      <span
                        v-if="item.periodOrderCount > 0"
                        class="px-2 py-0.5 rounded-full bg-slate-100 font-extrabold text-slate-900"
                      >
                        {{ item.periodOrderCount }}x order
                      </span>
                      <span v-else class="text-slate-300">-</span>
                    </td>

                    <!-- Total Kg -->
                    <td
                      class="py-3.5 px-4 text-right font-semibold text-slate-700 whitespace-nowrap"
                    >
                      <span v-if="item.periodWeight > 0" class="font-bold text-brand-600">
                        {{ item.periodWeight.toFixed(1) }} Kg
                      </span>
                      <span v-else class="text-slate-300">-</span>
                    </td>

                    <!-- Total Spend -->
                    <td
                      class="py-3.5 px-4 text-right font-black text-emerald-600 text-xs whitespace-nowrap"
                    >
                      <span v-if="item.periodSpend > 0">
                        Rp {{ formatRupiah(item.periodSpend) }}
                      </span>
                      <span v-else class="text-slate-300">-</span>
                    </td>

                    <!-- Layanan Digunakan -->
                    <td class="py-3.5 px-4">
                      <div
                        v-if="item.servicesUsed.length > 0"
                        class="flex items-center gap-1 flex-wrap"
                      >
                        <span
                          v-for="s in item.servicesUsed"
                          :key="s"
                          class="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-purple-50 text-purple-700 border border-purple-100 text-[10px] font-bold"
                        >
                          <span>{{ getServiceCategoryIcon(s) }}</span>
                          <span>{{ formatCategoryName(s) }}</span>
                        </span>
                      </div>
                      <span v-else class="text-slate-400 text-[11px] italic"
                        >Belum ada transaksi di periode ini</span
                      >
                    </td>

                    <!-- Aksi -->
                    <td class="py-3.5 px-4 text-center whitespace-nowrap">
                      <button
                        class="px-3 py-1.5 rounded-xl bg-purple-50 hover:bg-purple-100 text-purple-700 font-bold text-xs border border-purple-200 cursor-pointer transition-all"
                        @click="openCustomerCRM(item.customer, item.orders)"
                      >
                        Lihat CRM
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <PaginationBar
              v-model:current-page="pageCustomers"
              v-model:items-per-page="pageSizeCustomers"
              :total-items="filteredCustomerWithStats.length"
            />
          </div>
        </div>

        <!-- 3.d. Top Customer Leaderboard -->
        <div
          v-else-if="crmSubTab === 'TOP_CUSTOMERS'"
          class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4"
        >
          <div class="flex items-center justify-between">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <Award class="w-4 h-4 text-purple-600" />
                <span>Leaderboard 10 Pelanggan Paling Loyal</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Berdasarkan total belanja akumulatif & frekuensi cucian
              </p>
            </div>
          </div>

          <div v-if="topCustomers.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Belum ada data pelanggan yang tersedia.
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full text-left border-collapse text-xs">
              <thead>
                <tr
                  class="border-b border-slate-100 text-slate-400 uppercase tracking-wider font-extrabold text-[11px]"
                >
                  <th class="py-3 px-4 text-center w-12">Peringkat</th>
                  <th class="py-3 px-4">Nama Pelanggan</th>
                  <th class="py-3 px-4 text-center">Frekuensi Order</th>
                  <th class="py-3 px-4 text-right">Total Berat Cucian</th>
                  <th class="py-3 px-4 text-right">Total Belanja (LTV)</th>
                  <th class="py-3 px-4 text-center">Aksi CRM</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-100">
                <tr
                  v-for="(cust, idx) in topCustomers"
                  :key="cust.id"
                  class="hover:bg-slate-50 transition-colors"
                >
                  <td class="py-3.5 px-4 text-center font-black">
                    <span
                      v-if="idx === 0"
                      class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-amber-100 text-amber-700 text-sm shadow-xs"
                      >🥇</span
                    >
                    <span
                      v-else-if="idx === 1"
                      class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-slate-200 text-slate-700 text-sm shadow-xs"
                      >🥈</span
                    >
                    <span
                      v-else-if="idx === 2"
                      class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-amber-50 text-amber-800 text-sm shadow-xs"
                      >🥉</span
                    >
                    <span v-else class="text-slate-400 font-bold text-xs">#{{ idx + 1 }}</span>
                  </td>
                  <td class="py-3.5 px-4">
                    <div class="flex items-center gap-2">
                      <p class="font-black text-slate-900 text-sm">{{ cust.full_name }}</p>
                      <span
                        v-if="cust.customer_code"
                        class="px-1.5 py-0.5 rounded text-[10px] font-extrabold bg-purple-50 text-purple-700 border border-purple-200/60"
                      >
                        {{ cust.customer_code }}
                      </span>
                    </div>
                    <p class="text-[11px] text-slate-400 mt-0.5">{{ cust.phone }}</p>
                  </td>
                  <td class="py-3.5 px-4 text-center font-black text-slate-800 text-sm">
                    {{ cust.total_orders }}x
                  </td>
                  <td class="py-3.5 px-4 text-right font-semibold text-slate-600">
                    {{ Number(cust.total_kg || 0).toFixed(1) }} Kg
                  </td>
                  <td
                    class="py-3.5 px-4 text-right font-black text-emerald-600 text-sm whitespace-nowrap"
                  >
                    Rp {{ formatRupiah(cust.total_spend) }}
                  </td>
                  <td class="py-3.5 px-4 text-center whitespace-nowrap">
                    <div class="flex items-center justify-center gap-1.5">
                      <button
                        class="px-3 py-1.5 rounded-xl bg-purple-50 hover:bg-purple-100 text-purple-700 font-bold text-xs border border-purple-200 cursor-pointer"
                        @click="openCustomerCRM(cust)"
                      >
                        Lihat CRM
                      </button>
                      <a
                        :href="getCustomerGreetingWhatsappUrl(cust)"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="p-1.5 rounded-xl bg-emerald-50 hover:bg-emerald-100 text-emerald-700 border border-emerald-200 cursor-pointer"
                        title="Sapa WhatsApp"
                      >
                        <MessageCircle class="w-4 h-4" />
                      </a>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- Modal Tambah Pengeluaran (Milestone 4) -->
      <QuickExpenseModal v-model="showCreateExpenseModal" @expense-added="loadAllReports" />

      <!-- Sub-Modal Konfirmasi Hapus Pengeluaran -->
      <div
        v-if="showDeleteExpenseModal && expenseToDelete"
        class="fixed inset-0 z-60 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
      >
        <div
          class="bg-white rounded-2xl max-w-sm w-full p-5 shadow-2xl border border-slate-100 space-y-4 text-slate-800"
        >
          <div class="flex items-center gap-2.5 text-rose-600">
            <div class="w-9 h-9 rounded-xl bg-rose-50 flex items-center justify-center">
              <Trash2 class="w-5 h-5" />
            </div>
            <h3 class="font-black text-base">Hapus Catatan Pengeluaran?</h3>
          </div>
          <p class="text-xs text-slate-500 leading-relaxed">
            Pengeluaran sebesar
            <strong class="text-slate-900">Rp {{ formatRupiah(expenseToDelete.amount) }}</strong>
            untuk "{{ expenseToDelete.description }}" akan dihapus permanen dari buku kas.
          </p>
          <div class="flex items-center justify-end gap-2 pt-2">
            <button
              class="px-3 py-2 rounded-xl text-xs font-bold text-slate-600 hover:bg-slate-100 transition-all cursor-pointer"
              @click="showDeleteExpenseModal = false"
            >
              Batal
            </button>
            <button
              :disabled="deletingExpense"
              class="px-4 py-2 rounded-xl text-xs font-bold bg-rose-600 hover:bg-rose-700 text-white shadow-md shadow-rose-600/20 transition-all cursor-pointer disabled:opacity-50"
              @click="confirmDeleteExpense"
            >
              {{ deletingExpense ? 'Menghapus...' : 'Ya, Hapus' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Modal Detail CRM Pelanggan -->
      <CustomerCrmModal
        v-model="showCrmModal"
        :customer-id="selectedCrmCustomerId"
        :customer="selectedCrmCustomer"
        :start-date="startDate"
        :end-date="endDate"
        :customer-period-orders="selectedCrmCustomerOrders"
      />
    </div>
  </AppLayout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '../components/layout/AppLayout.vue'
import QuickExpenseModal from '../components/modals/QuickExpenseModal.vue'
import CustomerCrmModal from '../components/modals/CustomerCrmModal.vue'
import PaginationBar from '../components/ui/PaginationBar.vue'
import {
  BarChart3,
  Calendar,
  TrendingUp,
  Receipt,
  Wallet,
  DollarSign,
  Package,
  Users,
  Award,
  RefreshCw,
  PlusCircle,
  Download,
  Search,
  X,
  Trash2,
  PieChart,
  CheckCircle2,
  ShieldAlert,
  MessageCircle,
  Clock,
  AlertTriangle,
  UserPlus,
} from 'lucide-vue-next'
import { api, formatRupiah, formatDate, formatDateTime } from '../api/client'
import { useToast } from '../composables/useToast'

const route = useRoute()
const toast = useToast()

const activeTab = ref('TRANSACTIONS')
const loading = ref(false)

// Sub-Tab Navigation States
const txSubTab = ref('ALL') // 'ALL', 'SERVICES', 'PICKED_UP', 'READY_UNCOLLECTED', 'CANCELLED'
const finSubTab = ref('OMSET') // 'OMSET', 'PENDAPATAN', 'PENGELUARAN', 'PIUTANG', 'LABA_RUGI'
const crmSubTab = ref('SUMMARY') // 'SUMMARY', 'GROWTH', 'DETAIL', 'TOP_CUSTOMERS'

// Date Range Presets (Default: Hari Ini / TODAY)
const selectedPeriod = ref('TODAY')
const startDate = ref('')
const endDate = ref('')
const customStartDate = ref('')
const customEndDate = ref('')

const periodPresets = [
  { id: 'TODAY', label: 'Hari Ini' },
  { id: 'LAST_7_DAYS', label: '7 Hari Terakhir' },
  { id: 'LAST_30_DAYS', label: '30 Hari Terakhir' },
  { id: 'THIS_MONTH', label: 'Bulan Ini' },
  { id: 'CUSTOM', label: 'Kustom 📅' },
]

// Data States
const financialReport = ref({
  total_omset: 0,
  total_pendapatan: 0,
  total_pengeluaran: 0,
  laba_bersih: 0,
  daily_chart: [],
})

const serviceBreakdown = ref([])
const allOrders = ref([])
const expenses = ref([])
const unpaidOrders = ref([])
const allCustomers = ref([])
const topCustomers = ref([])

// Search States
const orderSearchQuery = ref('')
const expenseSearchQuery = ref('')
const customerSearchQuery = ref('')
const selectedCategoryFilter = ref('ALL')
const customerServiceFilter = ref('ALL')

// Pagination States
const pageOrders = ref(1)
const pageSizeOrders = ref(10)

const pagePickedUp = ref(1)
const pageSizePickedUp = ref(10)

const pageReady = ref(1)
const pageSizeReady = ref(10)

const pageCancelled = ref(1)
const pageSizeCancelled = ref(10)

const pageOmset = ref(1)
const pageSizeOmset = ref(10)

const pagePaid = ref(1)
const pageSizePaid = ref(10)

const pageExpenses = ref(1)
const pageSizeExpenses = ref(10)

const pageUnpaid = ref(1)
const pageSizeUnpaid = ref(10)

const pageCustomers = ref(1)
const pageSizeCustomers = ref(10)

const pageGrowth = ref(1)
const pageSizeGrowth = ref(10)

// Expense Modals
const showCreateExpenseModal = ref(false)
const showDeleteExpenseModal = ref(false)
const expenseToDelete = ref(null)
const deletingExpense = ref(false)

// CRM Modal
const showCrmModal = ref(false)
const selectedCrmCustomerId = ref('')
const selectedCrmCustomer = ref(null)
const selectedCrmCustomerOrders = ref([])

// Date helpers
const toYMD = (d) => {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const applyPeriod = (presetId) => {
  selectedPeriod.value = presetId
  const now = new Date()

  if (presetId === 'TODAY') {
    startDate.value = toYMD(now)
    endDate.value = toYMD(now)
  } else if (presetId === 'LAST_7_DAYS') {
    const past = new Date()
    past.setDate(past.getDate() - 6)
    startDate.value = toYMD(past)
    endDate.value = toYMD(now)
  } else if (presetId === 'LAST_30_DAYS') {
    const past = new Date()
    past.setDate(past.getDate() - 29)
    startDate.value = toYMD(past)
    endDate.value = toYMD(now)
  } else if (presetId === 'THIS_MONTH') {
    const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
    startDate.value = toYMD(firstDay)
    endDate.value = toYMD(now)
  } else if (presetId === 'CUSTOM') {
    customStartDate.value = startDate.value || toYMD(now)
    customEndDate.value = endDate.value || toYMD(now)
    return
  }

  loadAllReports()
}

const applyCustomDateRange = () => {
  if (!customStartDate.value || !customEndDate.value) {
    toast.error('Silakan tentukan tanggal mulai dan tanggal selesai.')
    return
  }
  if (customStartDate.value > customEndDate.value) {
    toast.error('Tanggal mulai tidak boleh lebih besar dari tanggal selesai.')
    return
  }
  startDate.value = customStartDate.value
  endDate.value = customEndDate.value
  loadAllReports()
}

const clearOrderSearch = () => {
  orderSearchQuery.value = ''
  pageOrders.value = 1
}

const clearExpenseSearch = () => {
  expenseSearchQuery.value = ''
  pageExpenses.value = 1
}

const clearCustomerSearch = () => {
  customerSearchQuery.value = ''
  pageCustomers.value = 1
}

// Master DataLoader
const loadAllReports = async () => {
  loading.value = true
  try {
    const queryParams = { start: startDate.value, end: endDate.value }

    const [finData, servData, expData, unpaidData, topCustData, ordersData, custData] =
      await Promise.all([
        api.get('/reports/financial', queryParams),
        api.get('/reports/services', queryParams),
        api.get('/expenses', queryParams),
        api.get('/reports/unpaid'),
        api.get('/reports/top-customers', { limit: 10 }),
        api.get('/pos/orders'),
        api.get('/customers'),
      ])

    financialReport.value = finData || {
      total_omset: 0,
      total_pendapatan: 0,
      total_pengeluaran: 0,
      laba_bersih: 0,
      daily_chart: [],
    }

    serviceBreakdown.value = servData || []
    expenses.value = expData || []
    unpaidOrders.value = unpaidData || []
    topCustomers.value = topCustData || []
    allOrders.value = Array.isArray(ordersData) ? ordersData : []
    allCustomers.value = Array.isArray(custData) ? custData : []

    // Reset page numbers
    pageOrders.value = 1
    pagePickedUp.value = 1
    pageReady.value = 1
    pageCancelled.value = 1
    pageOmset.value = 1
    pagePaid.value = 1
    pageExpenses.value = 1
    pageUnpaid.value = 1
    pageCustomers.value = 1
    pageGrowth.value = 1
  } catch (err) {
    console.error('Failed to load reports', err)
    toast.error('Gagal memuat beberapa data laporan.')
  } finally {
    loading.value = false
  }
}

// ==========================================
// Tab 1 Transaksi Filters & Pagination
// ==========================================
const periodOrders = computed(() => {
  if (!startDate.value || !endDate.value) return allOrders.value
  return allOrders.value.filter((o) => {
    if (!o.created_at) return false
    const d = o.created_at.slice(0, 10)
    return d >= startDate.value && d <= endDate.value
  })
})

const filteredPeriodOrders = computed(() => {
  let list = periodOrders.value
  if (orderSearchQuery.value.trim()) {
    const q = orderSearchQuery.value.toLowerCase()
    list = list.filter(
      (o) =>
        (o.order_number && o.order_number.toLowerCase().includes(q)) ||
        (o.customer?.full_name && o.customer.full_name.toLowerCase().includes(q)) ||
        (o.customer?.phone && o.customer.phone.includes(q))
    )
  }
  return list
})

const paginatedPeriodOrders = computed(() => {
  const start = (pageOrders.value - 1) * pageSizeOrders.value
  return filteredPeriodOrders.value.slice(start, start + pageSizeOrders.value)
})

const completedOrders = computed(() => {
  return periodOrders.value.filter((o) => o.status === 'SELESAI')
})

const paginatedCompletedOrders = computed(() => {
  const start = (pagePickedUp.value - 1) * pageSizePickedUp.value
  return completedOrders.value.slice(start, start + pageSizePickedUp.value)
})

const readyOrders = computed(() => {
  return allOrders.value.filter((o) => o.status === 'SIAP_AMBIL')
})

const paginatedReadyOrders = computed(() => {
  const start = (pageReady.value - 1) * pageSizeReady.value
  return readyOrders.value.slice(start, start + pageSizeReady.value)
})

const cancelledOrders = computed(() => {
  return periodOrders.value.filter((o) => o.status === 'BATAL')
})

const paginatedCancelledOrders = computed(() => {
  const start = (pageCancelled.value - 1) * pageSizeCancelled.value
  return cancelledOrders.value.slice(start, start + pageSizeCancelled.value)
})

// ==========================================
// Tab 2 Keuangan Filters & Pagination
// ==========================================
const omsetOrders = computed(() => {
  return periodOrders.value.filter((o) => o.status !== 'BATAL')
})

const paginatedOmsetOrders = computed(() => {
  const start = (pageOmset.value - 1) * pageSizeOmset.value
  return omsetOrders.value.slice(start, start + pageSizeOmset.value)
})

const paidOrders = computed(() => {
  return periodOrders.value.filter((o) => o.payment_status === 'PAID' || Number(o.paid_amount) > 0)
})

const paginatedPaidOrders = computed(() => {
  const start = (pagePaid.value - 1) * pageSizePaid.value
  return paidOrders.value.slice(start, start + pageSizePaid.value)
})

const filteredExpenses = computed(() => {
  let list = expenses.value
  if (selectedCategoryFilter.value !== 'ALL') {
    list = list.filter((e) => e.category === selectedCategoryFilter.value)
  }
  if (expenseSearchQuery.value.trim()) {
    const q = expenseSearchQuery.value.toLowerCase()
    list = list.filter(
      (e) =>
        (e.description && e.description.toLowerCase().includes(q)) ||
        (e.created_by && e.created_by.toLowerCase().includes(q))
    )
  }
  return list
})

const paginatedExpenses = computed(() => {
  const start = (pageExpenses.value - 1) * pageSizeExpenses.value
  return filteredExpenses.value.slice(start, start + pageSizeExpenses.value)
})

const categoryBreakdown = computed(() => {
  const map = {}
  let total = 0
  expenses.value.forEach((e) => {
    const cat = e.category || 'LAIN_LAIN'
    const amt = Number(e.amount) || 0
    total += amt
    if (!map[cat]) map[cat] = { category: cat, total_amount: 0, count: 0 }
    map[cat].total_amount += amt
    map[cat].count += 1
  })

  return Object.values(map)
    .map((item) => ({
      ...item,
      percent: total > 0 ? Math.round((item.total_amount / total) * 100) : 0,
    }))
    .sort((a, b) => b.total_amount - a.total_amount)
})

const totalPiutangNominal = computed(() => {
  return unpaidOrders.value.reduce((sum, o) => sum + (Number(o.sisa_tagihan) || 0), 0)
})

const paginatedUnpaidOrders = computed(() => {
  const start = (pageUnpaid.value - 1) * pageSizeUnpaid.value
  return unpaidOrders.value.slice(start, start + pageSizeUnpaid.value)
})

const profitMarginPercent = computed(() => {
  const rev = Number(financialReport.value.total_pendapatan) || 0
  const profit = Number(financialReport.value.laba_bersih) || 0
  if (rev <= 0) return 0
  return Math.round((profit / rev) * 100)
})

// ==========================================
// Tab 3 CRM Filters, Activity Recap & Pagination
// ==========================================
const newCustomersInPeriod = computed(() => {
  if (!startDate.value || !endDate.value) return []
  return allCustomers.value.filter((c) => {
    if (!c.created_at) return false
    const d = c.created_at.slice(0, 10)
    return d >= startDate.value && d <= endDate.value
  })
})

const paginatedGrowthCustomers = computed(() => {
  const start = (pageGrowth.value - 1) * pageSizeGrowth.value
  return newCustomersInPeriod.value.slice(start, start + pageSizeGrowth.value)
})

const activeCustomerCount = computed(() => {
  const ids = new Set()
  periodOrders.value.forEach((o) => {
    if (o.customer_id) ids.add(o.customer_id)
  })
  return ids.size
})

const avgOrderValue = computed(() => {
  const validOrders = periodOrders.value.filter((o) => o.status !== 'BATAL')
  if (validOrders.length === 0) return 0
  const sum = validOrders.reduce((acc, o) => acc + (Number(o.total_price) || 0), 0)
  return Math.round(sum / validOrders.length)
})

// Customer Activity Aggregation per Period
const customerPeriodMap = computed(() => {
  const map = {}
  periodOrders.value.forEach((ord) => {
    if (!ord.customer_id) return
    const cid = ord.customer_id
    if (!map[cid]) {
      map[cid] = {
        orderCount: 0,
        totalWeight: 0,
        totalSpend: 0,
        servicesUsed: new Set(),
        orders: [],
      }
    }
    map[cid].orders.push(ord)
    if (ord.status !== 'BATAL') {
      map[cid].orderCount += 1
      map[cid].totalWeight += Number(ord.total_weight) || 0
      map[cid].totalSpend += Number(ord.total_price) || 0
    }
    if (Array.isArray(ord.items)) {
      ord.items.forEach((item) => {
        if (item.category) map[cid].servicesUsed.add(item.category)
      })
    }
  })
  return map
})

const filteredCustomerWithStats = computed(() => {
  let list = allCustomers.value.map((c) => {
    const stats = customerPeriodMap.value[c.id] || {
      orderCount: 0,
      totalWeight: 0,
      totalSpend: 0,
      servicesUsed: new Set(),
      orders: [],
    }
    return {
      customer: c,
      periodOrderCount: stats.orderCount,
      periodWeight: stats.totalWeight,
      periodSpend: stats.totalSpend,
      servicesUsed: Array.from(stats.servicesUsed),
      orders: stats.orders,
    }
  })

  // Filter Service
  if (customerServiceFilter.value !== 'ALL') {
    list = list.filter((item) => item.servicesUsed.includes(customerServiceFilter.value))
  }

  // Filter Search
  if (customerSearchQuery.value.trim()) {
    const q = customerSearchQuery.value.toLowerCase()
    list = list.filter(
      (item) =>
        (item.customer.full_name && item.customer.full_name.toLowerCase().includes(q)) ||
        (item.customer.phone && item.customer.phone.includes(q)) ||
        (item.customer.customer_code && item.customer.customer_code.toLowerCase().includes(q)) ||
        (item.customer.address && item.customer.address.toLowerCase().includes(q))
    )
  }

  // Sort: Customers with transactions first, then alphabetically
  return list.sort((a, b) => {
    if (b.periodSpend !== a.periodSpend) return b.periodSpend - a.periodSpend
    return a.customer.full_name.localeCompare(b.customer.full_name)
  })
})

const paginatedCustomers = computed(() => {
  const start = (pageCustomers.value - 1) * pageSizeCustomers.value
  return filteredCustomerWithStats.value.slice(start, start + pageSizeCustomers.value)
})

const openCustomerCRM = (cust, orders = null) => {
  selectedCrmCustomerId.value = cust.id
  selectedCrmCustomer.value = cust
  if (orders) {
    selectedCrmCustomerOrders.value = orders
  } else {
    // Lookup from period orders
    selectedCrmCustomerOrders.value = periodOrders.value.filter((o) => o.customer_id === cust.id)
  }
  showCrmModal.value = true
}

// Helpers
const getOrderStatusBadgeClass = (status) => {
  switch (status) {
    case 'ANTRIAN':
      return 'bg-amber-100 text-amber-800'
    case 'CUCI':
    case 'KERING':
    case 'SETRIKA':
      return 'bg-blue-100 text-blue-800'
    case 'SIAP_AMBIL':
      return 'bg-purple-100 text-purple-800'
    case 'SELESAI':
      return 'bg-emerald-100 text-emerald-800'
    case 'BATAL':
      return 'bg-rose-100 text-rose-800'
    default:
      return 'bg-slate-100 text-slate-700'
  }
}

const formatCategoryName = (cat) => {
  switch (cat) {
    case 'KILOAN':
      return 'Cuci Kiloan'
    case 'SELF_SERVICE':
      return 'Self Service (Koin)'
    case 'SATUAN':
      return 'Cuci Satuan'
    case 'SEPATU_TAS':
      return 'Sepatu & Tas'
    case 'KARPET':
      return 'Cuci Karpet'
    default:
      return cat || 'Layanan'
  }
}

const getServiceCategoryIcon = (cat) => {
  switch (cat) {
    case 'KILOAN':
      return '🧺'
    case 'SELF_SERVICE':
      return '🪙'
    case 'SATUAN':
      return '👔'
    case 'SEPATU_TAS':
      return '👟'
    case 'KARPET':
      return '🧶'
    default:
      return '✨'
  }
}

const getCategoryLabel = (cat) => {
  switch (cat) {
    case 'BAHAN_BAKU':
      return 'Bahan Baku'
    case 'GAS':
      return 'Gas Dryer'
    case 'LISTRIK_AIR':
      return 'Listrik & Air'
    case 'GAJI':
      return 'Gaji Pegawai'
    case 'BONUS':
      return 'Uang Makan'
    case 'MAINTENANCE':
      return 'Servis Mesin'
    case 'SEWA':
      return 'Sewa Tempat'
    case 'LAIN_LAIN':
      return 'Lain-Lain'
    default:
      return cat || 'Operasional'
  }
}

const getCategoryIcon = (cat) => {
  switch (cat) {
    case 'BAHAN_BAKU':
      return '🧴'
    case 'GAS':
      return '🔥'
    case 'LISTRIK_AIR':
      return '⚡'
    case 'GAJI':
      return '👥'
    case 'BONUS':
      return '🍱'
    case 'MAINTENANCE':
      return '🔧'
    case 'SEWA':
      return '🏠'
    default:
      return '📦'
  }
}

const getCategoryBadgeClass = (cat) => {
  switch (cat) {
    case 'BAHAN_BAKU':
      return 'bg-blue-50 text-blue-700 border border-blue-200'
    case 'GAS':
      return 'bg-orange-50 text-orange-700 border border-orange-200'
    case 'LISTRIK_AIR':
      return 'bg-yellow-50 text-yellow-700 border border-yellow-200'
    case 'GAJI':
      return 'bg-purple-50 text-purple-700 border border-purple-200'
    case 'BONUS':
      return 'bg-pink-50 text-pink-700 border border-pink-200'
    case 'MAINTENANCE':
      return 'bg-amber-50 text-amber-700 border border-amber-200'
    case 'SEWA':
      return 'bg-indigo-50 text-indigo-700 border border-indigo-200'
    default:
      return 'bg-slate-100 text-slate-700 border border-slate-200'
  }
}

const promptDeleteExpense = (item) => {
  expenseToDelete.value = item
  showDeleteExpenseModal.value = true
}

const confirmDeleteExpense = async () => {
  if (!expenseToDelete.value) return
  deletingExpense.value = true
  try {
    await api.delete(`/expenses/${expenseToDelete.value.id}`)
    toast.success('Catatan pengeluaran berhasil dihapus.')
    showDeleteExpenseModal.value = false
    expenseToDelete.value = null
    loadAllReports()
  } catch (err) {
    toast.error(err.message || 'Gagal menghapus pengeluaran.')
  } finally {
    deletingExpense.value = false
  }
}

// WhatsApp Generator Helpers
const getReadyOrderWhatsappUrl = (order) => {
  const phone = order.customer?.phone || ''
  const cleanPhone = phone.replace(/^0/, '62').replace(/[^0-9]/g, '')
  const name = order.customer?.full_name || 'Kak'
  const rack = order.rack_number ? `di ${order.rack_number}` : 'di rak outlet'
  const isPaid = order.payment_status === 'PAID'
  const payInfo = isPaid
    ? 'Status: SUDAH LUNAS ✅'
    : `Status: BELUM LUNAS (Sisa: Rp ${formatRupiah(order.total_price - (order.paid_amount || 0))})`

  const msg = encodeURIComponent(
    `Halo Kak ${name}! 🧺 Cucian Kakak dengan nomor order #${order.order_number} sudah SELESAI, bersih, wangi, dan rapi disimpan ${rack} FAN Laundry Kalibaru.\n\n${payInfo}\n\nSilakan diambil saat outlet buka ya Kak. Terima kasih! 🙏`
  )
  return `https://wa.me/${cleanPhone}?text=${msg}`
}

const getUnpaidWhatsappUrl = (order) => {
  const phone = order.customer_phone || ''
  const cleanPhone = phone.replace(/^0/, '62').replace(/[^0-9]/g, '')
  const name = order.customer_name || 'Kak'
  const sisa = formatRupiah(order.sisa_tagihan)

  const msg = encodeURIComponent(
    `Halo Kak ${name}! 🙏 Mengingatkan dari FAN Laundry Kalibaru untuk pesanan #${order.order_number} masih terdapat sisa tagihan sebesar *Rp ${sisa}*.\n\nPembayaran dapat dilakukan via Tunai di kasir atau Transfer QRIS/BCA saat pengambilan. Terima kasih banyak ya Kak! ✨`
  )
  return `https://wa.me/${cleanPhone}?text=${msg}`
}

const getCustomerGreetingWhatsappUrl = (cust) => {
  const phone = cust.phone || ''
  const cleanPhone = phone.replace(/^0/, '62').replace(/[^0-9]/g, '')
  const name = cust.full_name || 'Kak'
  const msg = encodeURIComponent(
    `Halo Kak ${name}! ✨ Kami dari manajemen FAN Laundry Kalibaru ingin mengucapkan terima kasih banyak telah menjadi salah satu pelanggan paling setia kami 🙏 Semoga pelayanan kami selalu memuaskan ya Kak!`
  )
  return `https://wa.me/${cleanPhone}?text=${msg}`
}

// CSV Export
const exportCurrentTabCSV = () => {
  let filename = 'laporan-laundryos'
  let csvContent = 'data:text/csv;charset=utf-8,'

  if (activeTab.value === 'TRANSACTIONS') {
    filename = `laporan-transaksi-${txSubTab.value}-${startDate.value}-sd-${endDate.value}.csv`
    csvContent += 'No Order,Pelanggan,Telepon,Tanggal,Status,Status Bayar,Total\n'
    const targetOrders =
      txSubTab.value === 'PICKED_UP'
        ? completedOrders.value
        : txSubTab.value === 'READY_UNCOLLECTED'
          ? readyOrders.value
          : txSubTab.value === 'CANCELLED'
            ? cancelledOrders.value
            : periodOrders.value

    targetOrders.forEach((o) => {
      csvContent += `"${o.order_number}","${(o.customer?.full_name || '').replace(/"/g, '""')}","${o.customer?.phone || ''}","${o.created_at}","${o.status}","${o.payment_status}",${o.total_price}\n`
    })
  } else if (activeTab.value === 'FINANCIAL') {
    filename = `laporan-keuangan-${finSubTab.value}-${startDate.value}-sd-${endDate.value}.csv`
    if (finSubTab.value === 'PENGELUARAN') {
      csvContent += 'Tanggal,Kategori,Keterangan,Dicatat Oleh,Nominal\n'
      filteredExpenses.value.forEach((e) => {
        csvContent += `"${e.expense_date}","${e.category}","${(e.description || '').replace(/"/g, '""')}","${e.created_by || ''}",${e.amount}\n`
      })
    } else if (finSubTab.value === 'PIUTANG') {
      csvContent += 'No Order,Nama Pelanggan,Telepon,Total Tagihan,Sudah Dibayar,Sisa Piutang\n'
      unpaidOrders.value.forEach((u) => {
        csvContent += `"${u.order_number}","${u.customer_name}","${u.customer_phone}",${u.total_price},${u.paid_amount},${u.sisa_tagihan}\n`
      })
    } else {
      csvContent += 'No Order,Pelanggan,Tanggal,Status Bayar,Nominal\n'
      omsetOrders.value.forEach((o) => {
        csvContent += `"${o.order_number}","${o.customer?.full_name || ''}","${o.created_at}","${o.payment_status}",${o.total_price}\n`
      })
    }
  } else if (activeTab.value === 'CRM') {
    if (crmSubTab.value === 'DETAIL') {
      filename = `rekap-aktivitas-pelanggan-${startDate.value}-sd-${endDate.value}.csv`
      csvContent +=
        'Kode,Nama,Telepon,Total Order Periode,Total Kg Periode,Total Belanja Periode,Layanan Dipakai\n'
      filteredCustomerWithStats.value.forEach((item) => {
        csvContent += `"${item.customer.customer_code || ''}","${item.customer.full_name}","${item.customer.phone}",${item.periodOrderCount},${item.periodWeight},${item.periodSpend},"${item.servicesUsed.join('; ')}"\n`
      })
    } else {
      filename = `laporan-pelanggan-${crmSubTab.value}-${startDate.value}.csv`
      csvContent += 'Kode,Nama,Telepon,Alamat,Tanggal Daftar\n'
      allCustomers.value.forEach((c) => {
        csvContent += `"${c.customer_code}","${c.full_name}","${c.phone}","${(c.address || '').replace(/"/g, '""')}","${c.created_at}"\n`
      })
    }
  }

  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute('download', filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  toast.success(`Data ${filename} berhasil diunduh.`)
}

watch(
  () => route.query.tab,
  (newTab) => {
    if (newTab === 'transactions') activeTab.value = 'TRANSACTIONS'
    else if (newTab === 'financial' || newTab === 'expenses') activeTab.value = 'FINANCIAL'
    else if (newTab === 'crm') activeTab.value = 'CRM'
  },
  { immediate: true }
)

onMounted(() => {
  // Default filter: Hari Ini (TODAY) as requested by user
  applyPeriod('TODAY')
})
</script>
