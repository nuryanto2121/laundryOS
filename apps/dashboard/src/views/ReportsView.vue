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
            title="Download CSV"
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

      <!-- 2. Global Date Filter Toolbar (Default: Hari Ini) -->
      <div
        class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 bg-white p-3 sm:p-4 rounded-2xl border border-slate-100 shadow-xs"
      >
        <!-- Quick Date Range Pills -->
        <div class="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0 text-xs scrollbar-none">
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

        <!-- Date Range Indicators -->
        <div class="flex items-center gap-2 text-xs font-medium text-slate-500">
          <Calendar class="w-4 h-4 text-brand-600 shrink-0" />
          <span>{{ formatDate(startDate) }}</span>
          <span>&mdash;</span>
          <span>{{ formatDate(endDate) }}</span>
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

      <!-- 4. TIGA KATEGORI LAPORAN (Tabs) -->
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
            {{ serviceBreakdown.length }} Layanan
          </span>
        </button>

        <!-- Tab 2: Laporan Keuangan -->
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
            Top 10
          </span>
        </button>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 1: LAPORAN TRANSAKSI & OPERASIONAL                          -->
      <!-- ============================================================== -->
      <div v-if="activeTab === 'TRANSACTIONS'" class="space-y-6">
        <!-- Rekap Volume per Satuan & Kategori Layanan -->
        <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4">
          <div class="flex items-center justify-between">
            <h3
              class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
            >
              <Package class="w-4 h-4 text-brand-600" />
              <span>Rekapitulasi Volume & Penjualan Layanan</span>
            </h3>
            <span class="text-xs text-slate-400"
              >Periode: {{ formatDate(startDate) }} - {{ formatDate(endDate) }}</span
            >
          </div>

          <div
            v-if="serviceBreakdown.length === 0"
            class="py-10 text-center text-slate-400 text-xs"
          >
            Belum ada transaksi layanan pada periode ini.
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

        <!-- Monitoring Cucian Siap Ambil di Rak & Reminder WA -->
        <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <CheckCircle2 class="w-4 h-4 text-emerald-600" />
                <span>Monitoring Cucian Siap Ambil di Rak</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Pakaian yang sudah selesai diproses dan menunggu diambil atau diantar
              </p>
            </div>
            <span
              class="px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-50 text-emerald-700 border border-emerald-200 self-start sm:self-auto"
            >
              {{ readyOrders.length }} Order di Rak
            </span>
          </div>

          <div v-if="readyOrders.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Tidak ada cucian yang sedang menunggu di rak. Semua sudah diambil atau sedang diproses.
          </div>

          <div v-else class="overflow-x-auto">
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
                  v-for="order in readyOrders"
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
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 2: LAPORAN KEUANGAN & KAS KECIL                            -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'FINANCIAL'" class="space-y-6">
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

        <!-- Buku Kas Kecil Operasional (CRUD Milestone 4) -->
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
                Rincian belanja bahan baku, operasional mesin & kasbon harian
              </p>
            </div>

            <!-- Filters -->
            <div class="flex items-center gap-2 flex-wrap">
              <select
                v-model="selectedCategoryFilter"
                class="px-3 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-700 focus:outline-none focus:border-brand-500 cursor-pointer"
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
                />
                <button
                  v-if="expenseSearchQuery"
                  class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                  @click="expenseSearchQuery = ''"
                >
                  <X class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>
          </div>

          <div v-if="filteredExpenses.length === 0" class="py-8 text-center text-slate-400 text-xs">
            Tidak ada catatan pengeluaran yang cocok dengan filter.
          </div>

          <div v-else class="overflow-x-auto">
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
                  v-for="item in filteredExpenses"
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
        </div>

        <!-- Laporan Piutang Tagihan Belum Lunas -->
        <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4">
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

          <div v-else class="overflow-x-auto">
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
                  v-for="item in unpaidOrders"
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
                  <td class="py-3.5 px-4 text-right font-semibold text-slate-700 whitespace-nowrap">
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
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- TAB 3: LAPORAN PELANGGAN (CRM) & LEADERBOARD TOP 10            -->
      <!-- ============================================================== -->
      <div v-else-if="activeTab === 'CRM'" class="space-y-6">
        <!-- Top 10 Leaderboard Card -->
        <div class="bg-white p-5 rounded-2xl border border-slate-100 shadow-xs space-y-4">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div>
              <h3
                class="text-sm font-bold text-slate-800 uppercase tracking-wider flex items-center gap-2"
              >
                <Award class="w-4 h-4 text-purple-600" />
                <span>Leaderboard 10 Pelanggan Paling Loyal</span>
              </h3>
              <p class="text-xs text-slate-400 mt-0.5">
                Berdasarkan akumulasi omset belanja & total volume cucian di FAN Laundry
              </p>
            </div>

            <!-- Customer Search -->
            <div class="relative min-w-[200px] max-w-xs">
              <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
              <input
                v-model="customerSearchQuery"
                type="text"
                placeholder="Cari pelanggan..."
                class="w-full pl-9 pr-7 py-2 rounded-xl bg-slate-50 border border-slate-200 text-xs font-medium text-slate-800 placeholder-slate-400 focus:outline-none focus:border-brand-500 focus:bg-white"
              />
              <button
                v-if="customerSearchQuery"
                class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600"
                @click="customerSearchQuery = ''"
              >
                <X class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <div
            v-if="filteredTopCustomers.length === 0"
            class="py-10 text-center text-slate-400 text-xs"
          >
            Belum ada data pelanggan tercatat.
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
                  v-for="(cust, idx) in filteredTopCustomers"
                  :key="cust.id"
                  class="hover:bg-slate-50 transition-colors"
                >
                  <!-- Rank Badge -->
                  <td class="py-3.5 px-4 text-center font-black">
                    <span
                      v-if="idx === 0"
                      class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-amber-100 text-amber-700 text-sm shadow-xs"
                      title="Juara 1"
                    >
                      🥇
                    </span>
                    <span
                      v-else-if="idx === 1"
                      class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-slate-200 text-slate-700 text-sm shadow-xs"
                      title="Juara 2"
                    >
                      🥈
                    </span>
                    <span
                      v-else-if="idx === 2"
                      class="inline-flex items-center justify-center w-7 h-7 rounded-full bg-amber-50 text-amber-800 text-sm shadow-xs"
                      title="Juara 3"
                    >
                      🥉
                    </span>
                    <span v-else class="text-slate-400 font-bold text-xs"> #{{ idx + 1 }} </span>
                  </td>

                  <!-- Customer Name -->
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

                  <!-- Total Orders -->
                  <td class="py-3.5 px-4 text-center font-black text-slate-800 text-sm">
                    {{ cust.total_orders }}x
                  </td>

                  <!-- Total Kg -->
                  <td class="py-3.5 px-4 text-right font-semibold text-slate-600">
                    {{ Number(cust.total_kg || 0).toFixed(1) }} Kg
                  </td>

                  <!-- Total Spend -->
                  <td
                    class="py-3.5 px-4 text-right font-black text-emerald-600 text-sm whitespace-nowrap"
                  >
                    Rp {{ formatRupiah(cust.total_spend) }}
                  </td>

                  <!-- Actions -->
                  <td class="py-3.5 px-4 text-center whitespace-nowrap">
                    <div class="flex items-center justify-center gap-1.5">
                      <button
                        class="px-3 py-1.5 rounded-xl bg-purple-50 hover:bg-purple-100 text-purple-700 font-bold text-xs border border-purple-200 transition-all cursor-pointer"
                        @click="openCustomerCRM(cust)"
                      >
                        Lihat CRM
                      </button>
                      <a
                        :href="getCustomerGreetingWhatsappUrl(cust)"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="p-1.5 rounded-xl bg-emerald-50 hover:bg-emerald-100 text-emerald-700 border border-emerald-200 transition-all cursor-pointer"
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
} from 'lucide-vue-next'
import { api, formatRupiah, formatDate } from '../api/client'
import { useToast } from '../composables/useToast'

const route = useRoute()
const toast = useToast()

const activeTab = ref('TRANSACTIONS')
const loading = ref(false)

// Date Range (Default: Hari Ini / TODAY)
const selectedPeriod = ref('TODAY')
const startDate = ref('')
const endDate = ref('')

const periodPresets = [
  { id: 'TODAY', label: 'Hari Ini' },
  { id: 'LAST_7_DAYS', label: '7 Hari Terakhir' },
  { id: 'THIS_MONTH', label: 'Bulan Ini' },
  { id: 'ALL', label: 'Semua' },
]

// Top Financial Report Data
const financialReport = ref({
  total_omset: 0,
  total_pendapatan: 0,
  total_pengeluaran: 0,
  laba_bersih: 0,
  daily_chart: [],
})

// Tab 1: Services Breakdown & Ready Orders
const serviceBreakdown = ref([])
const readyOrders = ref([])

// Tab 2: Expenses & Unpaid Orders
const expenses = ref([])
const unpaidOrders = ref([])
const expenseSearchQuery = ref('')
const selectedCategoryFilter = ref('ALL')
const showCreateExpenseModal = ref(false)
const showDeleteExpenseModal = ref(false)
const expenseToDelete = ref(null)
const deletingExpense = ref(false)

// Tab 3: CRM Top Customers
const topCustomers = ref([])
const customerSearchQuery = ref('')
const showCrmModal = ref(false)
const selectedCrmCustomerId = ref('')
const selectedCrmCustomer = ref(null)

// Formatters for date range
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
  } else if (presetId === 'THIS_MONTH') {
    const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
    startDate.value = toYMD(firstDay)
    endDate.value = toYMD(now)
  } else {
    // ALL
    startDate.value = '2020-01-01'
    endDate.value = toYMD(now)
  }

  loadAllReports()
}

// Master DataLoader
const loadAllReports = async () => {
  loading.value = true
  try {
    const queryParams = { start: startDate.value, end: endDate.value }

    // Fetch Financial & Service Reports in parallel
    const [finData, servData, expData, unpaidData, topCustData, ordersData] = await Promise.all([
      api.get('/reports/financial', queryParams),
      api.get('/reports/services', queryParams),
      api.get('/expenses', queryParams),
      api.get('/reports/unpaid'),
      api.get('/reports/top-customers', { limit: 10 }),
      api.get('/pos/orders'),
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

    // Filter ready orders for Tab 1
    if (Array.isArray(ordersData)) {
      readyOrders.value = ordersData.filter((o) => o.status === 'SIAP_AMBIL')
    }
  } catch (err) {
    console.error('Failed to load reports', err)
    toast.error('Gagal memuat beberapa data laporan.')
  } finally {
    loading.value = false
  }
}

// Category formatting helpers
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

// Tab 2 Expenses Helpers
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

// Tab 3 CRM Helpers
const filteredTopCustomers = computed(() => {
  if (!customerSearchQuery.value.trim()) return topCustomers.value
  const q = customerSearchQuery.value.toLowerCase()
  return topCustomers.value.filter(
    (c) =>
      (c.full_name && c.full_name.toLowerCase().includes(q)) ||
      (c.phone && c.phone.includes(q)) ||
      (c.customer_code && c.customer_code.toLowerCase().includes(q))
  )
})

const openCustomerCRM = (cust) => {
  selectedCrmCustomerId.value = cust.id
  selectedCrmCustomer.value = cust
  showCrmModal.value = true
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
    filename = `laporan-transaksi-${startDate.value}-sd-${endDate.value}.csv`
    csvContent += 'Kategori,Satuan,Volume Total,Total Order,Total Pendapatan\n'
    serviceBreakdown.value.forEach((s) => {
      csvContent += `"${s.category}","${s.unit}",${s.total_volume},${s.total_order},${s.total_pendapatan}\n`
    })
  } else if (activeTab.value === 'FINANCIAL') {
    filename = `laporan-pengeluaran-${startDate.value}-sd-${endDate.value}.csv`
    csvContent += 'Tanggal,Kategori,Keterangan,Dicatat Oleh,Nominal\n'
    filteredExpenses.value.forEach((e) => {
      csvContent += `"${e.expense_date}","${e.category}","${(e.description || '').replace(/"/g, '""')}","${e.created_by || ''}",${e.amount}\n`
    })
  } else if (activeTab.value === 'CRM') {
    filename = `laporan-pelanggan-loyal-${new Date().toISOString().slice(0, 10)}.csv`
    csvContent += 'Ranking,Kode,Nama,Telepon,Total Orders,Total Kg,Total Belanja\n'
    topCustomers.value.forEach((c, idx) => {
      csvContent += `${idx + 1},"${c.customer_code}","${c.full_name}","${c.phone}",${c.total_orders},${c.total_kg},${c.total_spend}\n`
    })
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
