<template>
  <div class="min-h-screen bg-slate-900 text-slate-800 flex flex-col items-center justify-center p-2 sm:p-6 font-sans">
    
    <!-- Top Presentation Header Bar -->
    <div class="w-full max-w-[420px] flex items-center justify-between px-3 py-2 text-slate-300 text-xs mb-3">
      <div class="flex items-center gap-2">
        <router-link to="/" class="p-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 flex items-center gap-1 text-[11px] font-bold">
          <ArrowLeft class="w-3.5 h-3.5" />
          <span>Kembali ke Web POS</span>
        </router-link>
        <span class="px-2 py-0.5 rounded-md bg-brand-500/20 text-[10px] text-sky-400 font-mono font-bold">Mobile Native LMS</span>
      </div>
      <div class="flex items-center gap-1.5 text-[11px] text-slate-400">
        <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
        <span>Simulasi Mobile App</span>
      </div>
    </div>

    <!-- Smartphone Physical Canvas Frame -->
    <div class="relative w-full max-w-[400px] h-[100dvh] md:h-[844px] bg-slate-50 md:rounded-[48px] overflow-hidden flex flex-col shadow-2xl border-0 md:border-8 md:border-slate-800">

      <!-- 1. Dynamic Island & Hardware Status Bar -->
      <div class="pt-3 px-6 pb-2 bg-white flex items-center justify-between shrink-0 select-none border-b border-slate-100 z-30">
        <div class="text-xs font-black text-slate-900 tracking-tight">09:41</div>
        <!-- Dynamic Island Capsule -->
        <div class="w-24 h-5 bg-slate-950 rounded-full flex items-center justify-center gap-1.5 shadow-inner">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
          <span class="text-[9px] font-mono text-emerald-400 font-bold">BLE READY</span>
        </div>
        <div class="flex items-center gap-1.5 text-slate-900">
          <Wifi class="w-3.5 h-3.5" />
          <BatteryCharging class="w-4 h-4 text-emerald-600" />
        </div>
      </div>

      <!-- 2. Mobile App Header -->
      <div class="px-4 py-2.5 bg-white border-b border-slate-100 flex items-center justify-between shrink-0 z-20">
        <div class="flex items-center gap-2">
          <div class="w-8 h-8 rounded-xl bg-brand-600 text-white flex items-center justify-center font-black text-sm shadow-xs">
            FL
          </div>
          <div>
            <div class="flex items-center gap-1">
              <h1 class="text-xs font-extrabold text-slate-900">FAN Laundry Kalibaru</h1>
              <ChevronDown class="w-3 h-3 text-slate-400" />
            </div>
            <div class="flex items-center gap-1.5 mt-0.5">
              <button 
                type="button" 
                @click="toggleRole"
                class="px-2 py-0.5 rounded-full text-[9px] font-black transition-all cursor-pointer flex items-center gap-1 active:scale-95 shadow-2xs"
                :class="currentRole === 'OWNER' ? 'bg-amber-100 text-amber-900 border border-amber-300' : 'bg-slate-100 text-slate-700 border border-slate-200'"
                title="Klik untuk ganti mode Kasir / Owner"
              >
                <span>{{ currentRole === 'OWNER' ? '👑 Owner (Akses Penuh)' : '👤 Kasir: Nuryanto' }}</span>
                <span class="text-[8px] opacity-60 underline">Ganti</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Quick Bluetooth & Alert Header Actions -->
        <div class="flex items-center gap-1.5">
          <button 
            type="button"
            class="p-2 rounded-xl border text-xs font-bold flex items-center gap-1 transition-all cursor-pointer"
            :class="bluetoothConnected ? 'bg-sky-50 border-sky-200 text-sky-700' : 'bg-slate-100 border-slate-200 text-slate-400'"
            @click="toggleBluetooth"
            title="Status Bluetooth Thermal"
          >
            <Printer class="w-3.5 h-3.5" />
            <span class="text-[10px]">58mm</span>
          </button>
          <button 
            type="button"
            class="p-2 rounded-xl bg-slate-100 hover:bg-slate-200 text-slate-700 relative cursor-pointer"
            @click="showNotificationToast('Ada 2 cucian yang mendekati batas SLA!')"
          >
            <Bell class="w-3.5 h-3.5" />
            <span class="absolute top-1 right-1 w-2 h-2 rounded-full bg-rose-500"></span>
          </button>
        </div>
      </div>

      <!-- 3. Screen View Container (Scrollable Content) -->
      <div class="flex-1 overflow-y-auto relative bg-slate-50">

        <!-- ============================================================ -->
        <!-- SCREEN 1: BERANDA / QUICK RADAR                              -->
        <!-- ============================================================ -->
        <div v-if="activeNav === 'home'" class="p-4 space-y-4 pb-24">
          <!-- Live Quick Stats (Thumb friendly 2x2 grid) -->
          <div class="grid grid-cols-2 gap-2.5">
            <div class="bg-white p-3.5 rounded-2xl border border-slate-200/80 shadow-xs space-y-1">
              <div class="flex items-center justify-between">
                <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Omset Hari Ini</span>
                <div class="w-6 h-6 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center">
                  <Banknote class="w-3.5 h-3.5" />
                </div>
              </div>
              <div class="text-base font-black text-slate-900 tracking-tight">Rp 842.000</div>
              <div class="text-[10px] text-emerald-600 font-semibold flex items-center gap-0.5">
                <TrendingUp class="w-3 h-3" /> +18% vs kemarin
              </div>
            </div>

            <div class="bg-white p-3.5 rounded-2xl border border-slate-200/80 shadow-xs space-y-1">
              <div class="flex items-center justify-between">
                <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider">Antrian Aktif</span>
                <div class="w-6 h-6 rounded-lg bg-sky-50 text-sky-600 flex items-center justify-center">
                  <Clock class="w-3.5 h-3.5" />
                </div>
              </div>
              <div class="text-base font-black text-slate-900 tracking-tight">14 Order</div>
              <div class="text-[10px] text-amber-600 font-semibold flex items-center gap-0.5">
                <AlertCircle class="w-3 h-3" /> 2 cucian urgent SLA
              </div>
            </div>
          </div>

          <!-- 1-Tap Quick Action Row (Ergonomic Thumb Area) -->
          <div class="grid grid-cols-4 gap-2 text-center">
            <button 
              type="button"
              @click="activeNav = 'pos'"
              class="bg-brand-600 hover:bg-brand-700 text-white p-3 rounded-2xl shadow-sm flex flex-col items-center gap-1.5 transition-transform active:scale-95 cursor-pointer"
            >
              <Scale class="w-5 h-5" />
              <span class="text-[10px] font-extrabold leading-tight">Timbang (+POS)</span>
            </button>
            <button 
              type="button"
              @click="activeNav = 'scan'"
              class="bg-white hover:bg-slate-100 text-slate-800 p-3 rounded-2xl border border-slate-200/90 shadow-2xs flex flex-col items-center gap-1.5 transition-transform active:scale-95 cursor-pointer"
            >
              <QrCode class="w-5 h-5 text-indigo-600" />
              <span class="text-[10px] font-bold leading-tight">Scan QR</span>
            </button>
            <button 
              type="button"
              @click="activeNav = 'queue'"
              class="bg-white hover:bg-slate-100 text-slate-800 p-3 rounded-2xl border border-slate-200/90 shadow-2xs flex flex-col items-center gap-1.5 transition-transform active:scale-95 cursor-pointer"
            >
              <Layers class="w-5 h-5 text-amber-600" />
              <span class="text-[10px] font-bold leading-tight">Antrian Rak</span>
            </button>
            <button 
              type="button"
              @click="showNotificationToast('Kas kecil cepat siap dicatat.')"
              class="bg-white hover:bg-slate-100 text-slate-800 p-3 rounded-2xl border border-slate-200/90 shadow-2xs flex flex-col items-center gap-1.5 transition-transform active:scale-95 cursor-pointer"
            >
              <Wallet class="w-5 h-5 text-rose-600" />
              <span class="text-[10px] font-bold leading-tight">Kas Kecil</span>
            </button>
          </div>

          <!-- Section: Cucian Butuh Perhatian (SLA Warning) -->
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <h2 class="text-xs font-black text-slate-900 uppercase tracking-wider flex items-center gap-1.5">
                <Flame class="w-3.5 h-3.5 text-rose-500" />
                Prioritas Pengerjaan (Target SLA)
              </h2>
              <button type="button" @click="activeNav = 'queue'" class="text-[10px] font-bold text-brand-600 hover:underline cursor-pointer">
                Lihat Semua
              </button>
            </div>

            <div class="space-y-2">
              <div 
                v-for="order in urgentOrders" 
                :key="order.id"
                class="bg-white p-3.5 rounded-2xl border border-slate-200/90 shadow-xs flex items-center justify-between gap-3"
              >
                <div class="space-y-1">
                  <div class="flex items-center gap-1.5">
                    <span class="font-mono text-xs font-black text-slate-900">#{{ order.num }}</span>
                    <span class="px-2 py-0.5 rounded-md text-[10px] font-extrabold" :class="order.badgeClass">
                      {{ order.status }}
                    </span>
                    <span class="px-1.5 py-0.5 rounded-md text-[10px] font-black bg-slate-900 text-white font-mono">
                      Rak: {{ order.rack }}
                    </span>
                  </div>
                  <p class="text-xs font-bold text-slate-800">{{ order.customer }} ({{ order.weight }} Kg)</p>
                  <p class="text-[10px] text-rose-600 font-semibold flex items-center gap-1">
                    <Clock class="w-3 h-3" /> Batas SLA: {{ order.due }}
                  </p>
                </div>

                <!-- Quick Progress Action -->
                <button 
                  type="button"
                  @click="progressOrder(order)"
                  class="px-3 py-2 rounded-xl bg-brand-50 hover:bg-brand-100 text-brand-700 text-xs font-black flex flex-col items-center gap-0.5 border border-brand-200 shrink-0 cursor-pointer"
                >
                  <ArrowRightCircle class="w-4 h-4" />
                  <span class="text-[9px]">Lanjut</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- ============================================================ -->
        <!-- SCREEN 2: MOBILE NATIVE POS (TIMBANG & CATAT INSTAN)         -->
        <!-- ============================================================ -->
        <div v-if="activeNav === 'pos'" class="p-4 space-y-3.5 pb-28">
          <!-- 1. Customer Select Bar -->
          <div class="bg-white p-3 rounded-2xl border border-slate-200/90 shadow-xs flex items-center justify-between gap-2">
            <div class="flex items-center gap-2.5 flex-1">
              <div class="w-9 h-9 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center font-bold">
                <User class="w-4 h-4" />
              </div>
              <div class="flex-1 truncate">
                <span class="text-[10px] font-bold text-slate-400 block uppercase">Pelanggan</span>
                <span class="text-xs font-extrabold text-slate-900 block truncate">{{ selectedCustomer.name }}</span>
              </div>
            </div>
            <button 
              type="button"
              @click="showNotificationToast('Daftar pelanggan terbuka')"
              class="px-2.5 py-1.5 rounded-xl bg-slate-100 text-slate-700 font-bold text-xs hover:bg-slate-200 cursor-pointer"
            >
              Ubah
            </button>
          </div>

          <!-- 2. Layanan Chips Horizontal Scroll (Native Mobile) -->
          <div>
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block mb-1.5">
              Pilih Kategori Layanan
            </span>
            <div class="flex items-center gap-1.5 overflow-x-auto pb-1">
              <button 
                type="button"
                v-for="cat in serviceCategories"
                :key="cat.id"
                class="px-3 py-1.5 rounded-xl text-xs font-bold whitespace-nowrap shrink-0 transition-all cursor-pointer"
                :class="activeCategory === cat.id ? 'bg-slate-900 text-white shadow-xs' : 'bg-white text-slate-600 border border-slate-200'"
                @click="activeCategory = cat.id"
              >
                {{ cat.name }}
              </button>
            </div>
          </div>

          <!-- 3. Service Cards for Selected Category -->
          <div class="grid grid-cols-2 gap-2">
            <button 
              type="button"
              v-for="srv in filteredServices"
              :key="srv.id"
              @click="selectService(srv)"
              class="p-3 rounded-2xl text-left transition-all border relative cursor-pointer"
              :class="selectedService.id === srv.id ? 'bg-sky-50/60 border-brand-500 ring-2 ring-brand-500/20' : 'bg-white border-slate-200'"
            >
              <div class="text-xs font-extrabold text-slate-900 leading-tight mb-1">{{ srv.name }}</div>
              <div class="text-xs font-black text-brand-600">Rp {{ formatNumber(srv.price) }}/{{ srv.unit }}</div>
              <div class="text-[10px] text-slate-400 mt-1">SLA {{ srv.duration }}</div>
            </button>
          </div>

          <!-- 4. Big-Touch Scale Weight / Qty Stepper (Digital Scale Interface) -->
          <div class="bg-white p-4 rounded-2xl border border-slate-200/90 shadow-xs space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-xs font-extrabold text-slate-800 flex items-center gap-1.5">
                <Scale class="w-4 h-4 text-brand-600" />
                Input Timbangan Cucian
              </span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-slate-100 text-slate-700">
                Satuan: {{ selectedService.unit || 'Kg' }}
              </span>
            </div>

            <!-- Digital Scale Display -->
            <div class="bg-slate-900 rounded-2xl p-4 text-center text-white relative overflow-hidden">
              <div class="text-[10px] font-mono uppercase tracking-widest text-slate-400 mb-0.5">Berat Bersih Terukur</div>
              <div class="text-4xl font-mono font-black tracking-tight text-sky-400">
                {{ currentWeight.toFixed(1) }} <span class="text-lg font-bold text-white">{{ selectedService.unit }}</span>
              </div>
              <!-- Self Service Capacity Meter (8.0 Kg Limit indicator) -->
              <div v-if="selectedService.category === 'Self Service'" class="mt-2 pt-2 border-t border-slate-800 text-[10px] flex items-center justify-between">
                <span class="text-slate-400">Kapasitas Tabung (8 Kg):</span>
                <span :class="currentWeight > 8.5 ? 'text-rose-400 font-bold' : 'text-emerald-400 font-bold'">
                  {{ currentWeight > 8.5 ? '⚠️ Melebihi Tabung (+1 Tabung)' : '✓ Pas 1 Tabung' }}
                </span>
              </div>
            </div>

            <!-- Large Touch Stepper Buttons -->
            <div class="grid grid-cols-4 gap-2">
              <button type="button" @click="adjustWeight(-1)" class="py-2.5 rounded-xl bg-slate-100 hover:bg-slate-200 font-black text-sm text-slate-800 cursor-pointer">-1.0</button>
              <button type="button" @click="adjustWeight(-0.5)" class="py-2.5 rounded-xl bg-slate-100 hover:bg-slate-200 font-black text-sm text-slate-800 cursor-pointer">-0.5</button>
              <button type="button" @click="adjustWeight(0.5)" class="py-2.5 rounded-xl bg-slate-100 hover:bg-slate-200 font-black text-sm text-slate-800 cursor-pointer">+0.5</button>
              <button type="button" @click="adjustWeight(1.0)" class="py-2.5 rounded-xl bg-brand-50 hover:bg-brand-100 font-black text-sm text-brand-700 cursor-pointer">+1.0</button>
            </div>
          </div>

          <!-- 5. Perfume & Rack Selection -->
          <div class="grid grid-cols-2 gap-2.5">
            <div class="bg-white p-3 rounded-2xl border border-slate-200/90 shadow-xs">
              <label class="text-[10px] font-bold text-slate-400 uppercase block mb-1">Aroma Parfum</label>
              <select v-model="selectedPerfume" class="w-full text-xs font-bold bg-slate-50 border border-slate-200 rounded-xl p-2 focus:outline-hidden">
                <option value="Akasia">Akasia (Favorit)</option>
                <option value="Downy Black">Downy Black</option>
                <option value="Ocean Fresh">Ocean Fresh</option>
                <option value="Sakura">Sakura</option>
              </select>
            </div>
            <div class="bg-white p-3 rounded-2xl border border-slate-200/90 shadow-xs">
              <label class="text-[10px] font-bold text-slate-400 uppercase block mb-1">Nomor Rak</label>
              <input v-model="rackNumber" type="text" placeholder="Misal: A-01" class="w-full text-xs font-bold bg-slate-50 border border-slate-200 rounded-xl p-2 uppercase focus:outline-hidden">
            </div>
          </div>
        </div>

        <!-- ============================================================ -->
                <!-- ============================================================ -->
        <!-- SCREEN 3: ANTRIAN (FIXED STATUS FILTER & DIRECT TAP ACTIONS) -->
        <!-- ============================================================ -->
        <div v-if="activeNav === 'queue'" class="p-4 space-y-3 pb-24">
          <!-- Fixed Status Filter Tabs (100% Fit Width, No Scroll, No Cutoff) -->
          <div class="grid grid-cols-4 gap-1 p-1 bg-slate-100/90 rounded-2xl w-full border border-slate-200/60 select-none">
            <button 
              type="button" 
              v-for="tab in queueTabs" 
              :key="tab.id"
              class="py-2 px-1 rounded-xl text-center transition-all cursor-pointer flex flex-col items-center justify-center leading-tight truncate active:scale-95"
              :class="activeQueueTab === tab.id ? 'bg-slate-900 text-white shadow-xs' : 'text-slate-600 hover:text-slate-900 hover:bg-slate-200/60'"
              @click="activeQueueTab = tab.id"
            >
              <span class="text-[11px] font-extrabold truncate">{{ tab.label }}</span>
              <span class="text-[10px] font-semibold mt-0.5" :class="activeQueueTab === tab.id ? 'text-slate-300' : 'text-slate-400'">({{ tab.count }})</span>
            </button>
          </div>

          <!-- Order Cards with Direct 1-Tap Action Button -->
          <div class="space-y-2.5">
            <div 
              v-for="item in filteredQueueList" 
              :key="item.id"
              class="bg-white rounded-2xl border border-slate-200/90 p-3.5 shadow-xs space-y-2 relative overflow-hidden transition-all active:scale-[0.99] cursor-pointer"
              @click="openOrderDetailSheet(item)"
            >
              <!-- Card Header -->
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-1.5">
                  <span class="font-mono text-xs font-black text-slate-900">#{{ item.orderNumber }}</span>
                  <span class="px-2 py-0.5 rounded-md text-[10px] font-extrabold" :class="item.badgeClass">
                    {{ item.status }}
                  </span>
                </div>
                <span class="px-2 py-0.5 rounded-md text-[10px] font-black bg-slate-900 text-white font-mono">
                  Rak: {{ item.rack }}
                </span>
              </div>

              <!-- Customer & Service Info -->
              <div class="flex items-center justify-between text-xs">
                <div class="font-bold text-slate-800">{{ item.customer }}</div>
                <div class="font-black text-slate-900">Rp {{ formatNumber(item.total) }}</div>
              </div>

              <!-- Detail Meta Row -->
              <div class="flex items-center justify-between text-[11px] text-slate-500 pt-1 border-t border-slate-100">
                <span>{{ item.service }} ({{ item.weight }} Kg)</span>
                <span :class="item.paid ? 'text-emerald-600 font-bold' : 'text-rose-600 font-bold'">
                  {{ item.paid ? 'LUNAS' : 'BELUM LUNAS' }}
                </span>
              </div>

              <!-- Direct 1-Tap Quick Action Buttons (No Swipe Needed) -->
              <div class="flex items-center justify-between pt-1 text-[11px]">
                <button 
                  type="button" 
                  class="flex items-center gap-1 text-slate-500 hover:text-slate-900 font-bold py-1 px-2.5 rounded-lg bg-slate-50 border border-slate-200 hover:bg-slate-100 transition-all cursor-pointer active:scale-95"
                  @click.stop="printViaBluetooth"
                >
                  <Printer class="w-3.5 h-3.5 text-brand-600" />
                  <span>Tap Cetak</span>
                </button>

                <button 
                  type="button" 
                  class="flex items-center gap-1 px-3 py-1.5 rounded-xl font-bold transition-all cursor-pointer active:scale-95 shadow-2xs"
                  :class="item.status === 'SELESAI' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-sky-50 text-sky-700 hover:bg-sky-600 hover:text-white border border-sky-200'"
                  @click.stop="quickAdvanceQueue(item)"
                >
                  <span v-if="item.status === 'ANTRIAN'">Mulai Proses</span>
                  <span v-else-if="item.status === 'PROSES'">Siap Ambil</span>
                  <span v-else-if="item.status === 'SIAP AMBIL'">Selesaikan</span>
                  <span v-else>Selesai ✓</span>
                  <ChevronRight v-if="item.status !== 'SELESAI'" class="w-3.5 h-3.5" />
                </button>
              </div>
            </div>

            <div v-if="filteredQueueList.length === 0" class="py-12 text-center text-xs text-slate-400">
              Tidak ada cucian di status ini.
            </div>
          </div>
        </div>

        <!-- ============================================================ -->
        <!-- SCREEN 4: INSTANT CAMERA QR/BARCODE SCANNER                  -->
        <!-- ============================================================ -->
        <div v-if="activeNav === 'scan'" class="h-full min-h-[500px] flex flex-col justify-between p-4 bg-slate-950 text-white relative">
          <!-- Viewfinder Overlay -->
          <div class="text-center pt-2 space-y-1">
            <div class="text-sm font-extrabold tracking-wide">Pindai Kode QR / Barcode</div>
            <p class="text-xs text-slate-400">Arahkan kamera ke struk nota atau label rak cucian</p>
          </div>

          <!-- Camera Reticle Frame -->
          <div class="relative w-64 h-64 mx-auto my-auto border-2 border-dashed border-sky-400/80 rounded-3xl flex items-center justify-center overflow-hidden">
            <QrCode class="w-16 h-16 text-sky-400/30" />
          </div>

          <!-- Bottom Simulation Trigger -->
          <div class="space-y-2 pb-20">
            <button 
              type="button"
              @click="simulateScanSuccess"
              class="w-full py-3.5 rounded-2xl bg-sky-500 hover:bg-sky-600 text-white font-extrabold text-xs shadow-lg transition-transform active:scale-95 flex items-center justify-center gap-2 cursor-pointer"
            >
              <Zap class="w-4 h-4" />
              <span>Simulasi Scan Sukses (Nota #261120000002)</span>
            </button>
            <p class="text-[10px] text-center text-slate-500">Kamera aktif langsung via native hardware API</p>
          </div>
        </div>

        <!-- ============================================================ -->
        <!-- SCREEN 5: MENU & PENGATURAN (ROLE BASED: OWNER vs KASIR)     -->
        <!-- ============================================================ -->
        <div v-if="activeNav === 'menu'" class="p-4 space-y-4 pb-24">
          <!-- Role Indicator Card -->
          <div class="p-3.5 rounded-2xl border flex items-center justify-between"
            :class="currentRole === 'OWNER' ? 'bg-amber-50/80 border-amber-200' : 'bg-slate-100 border-slate-200'"
          >
            <div class="flex items-center gap-2.5">
              <div class="w-9 h-9 rounded-xl flex items-center justify-center font-bold text-base"
                :class="currentRole === 'OWNER' ? 'bg-amber-500 text-white shadow-xs' : 'bg-slate-800 text-white'"
              >
                {{ currentRole === 'OWNER' ? '👑' : '👤' }}
              </div>
              <div>
                <div class="font-extrabold text-xs text-slate-900">
                  {{ currentRole === 'OWNER' ? 'Akun Owner (Akses Penuh)' : 'Akun Kasir (Operasional Toko)' }}
                </div>
                <div class="text-[10px] text-slate-500">
                  {{ currentRole === 'OWNER' ? 'Hak akses Laporan, Omset, Settings, & POS' : 'Shift Pagi: 07:00 - 15:00 WIB' }}
                </div>
              </div>
            </div>
            <button 
              type="button" 
              @click="toggleRole"
              class="px-2.5 py-1 rounded-xl text-[10px] font-black border transition-all cursor-pointer active:scale-95"
              :class="currentRole === 'OWNER' ? 'bg-white text-amber-800 border-amber-300 shadow-2xs' : 'bg-white text-slate-700 border-slate-300 shadow-2xs'"
            >
              Switch Role
            </button>
          </div>

          <!-- OWNER EXCLUSIVE VIEW -->
          <div v-if="currentRole === 'OWNER'" class="space-y-3.5">
            <div class="text-xs font-black text-slate-800 uppercase tracking-wider flex items-center gap-1.5">
              <Banknote class="w-4 h-4 text-emerald-600" />
              <span>Ringkasan Eksekutif Keuangan</span>
            </div>

            <!-- Omset & Kas Grid -->
            <div class="grid grid-cols-2 gap-2">
              <div class="bg-white p-3 rounded-2xl border border-slate-200 shadow-xs space-y-1">
                <span class="text-[10px] font-bold text-slate-400 uppercase">Omset Hari Ini</span>
                <div class="text-sm font-black text-slate-900">Rp 842.000</div>
                <div class="text-[9px] text-emerald-600 font-bold">+18% dari target</div>
              </div>
              <div class="bg-white p-3 rounded-2xl border border-slate-200 shadow-xs space-y-1">
                <span class="text-[10px] font-bold text-slate-400 uppercase">Kas Kecil & Saldo</span>
                <div class="text-sm font-black text-slate-900">Rp 350.000</div>
                <div class="text-[9px] text-slate-500">Uang tunai kasir</div>
              </div>
            </div>

            <!-- Breakdown Pembayaran -->
            <div class="bg-white p-3.5 rounded-2xl border border-slate-200 shadow-xs space-y-2">
              <div class="text-xs font-bold text-slate-800 flex justify-between">
                <span>Metode Pembayaran</span>
                <span class="text-[10px] text-slate-400">Total 14 Transaksi</span>
              </div>
              <div class="space-y-1.5 text-[11px]">
                <div class="flex justify-between items-center text-slate-600">
                  <span class="flex items-center gap-1.5">🟣 QRIS (Statik & Dinamis)</span>
                  <span class="font-bold text-slate-900">Rp 450.000 (53%)</span>
                </div>
                <div class="flex justify-between items-center text-slate-600">
                  <span class="flex items-center gap-1.5">💵 Tunai Kasir</span>
                  <span class="font-bold text-slate-900">Rp 312.000 (37%)</span>
                </div>
                <div class="flex justify-between items-center text-slate-600">
                  <span class="flex items-center gap-1.5">🏦 Transfer Bank BCA/Mandiri</span>
                  <span class="font-bold text-slate-900">Rp 80.000 (10%)</span>
                </div>
              </div>
            </div>

            <!-- Master Data & Settings Menu List -->
            <div class="text-xs font-black text-slate-800 uppercase tracking-wider pt-1 flex items-center gap-1.5">
              <Layers class="w-4 h-4 text-brand-600" />
              <span>Menu Pengaturan & Master Data</span>
            </div>

            <div class="bg-white rounded-2xl border border-slate-200 shadow-xs divide-y divide-slate-100 overflow-hidden text-xs">
              <button 
                type="button" 
                @click="showNotificationToast('Membuka Master Layanan & Tarif (8 Kategori)')"
                class="w-full p-3.5 flex items-center justify-between hover:bg-slate-50 transition-colors cursor-pointer text-left"
              >
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center font-bold">🏷️</div>
                  <div>
                    <div class="font-bold text-slate-900">Master Layanan & Tarif</div>
                    <div class="text-[10px] text-slate-400">Atur harga 8 kategori layanan & SLA</div>
                  </div>
                </div>
                <ChevronRight class="w-4 h-4 text-slate-400" />
              </button>

              <button 
                type="button" 
                @click="showNotificationToast('Membuka Manajemen Pegawai & Hak Akses')"
                class="w-full p-3.5 flex items-center justify-between hover:bg-slate-50 transition-colors cursor-pointer text-left"
              >
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center font-bold">👥</div>
                  <div>
                    <div class="font-bold text-slate-900">Pegawai & Hak Akses (RBAC)</div>
                    <div class="text-[10px] text-slate-400">Kelola akun kasir, kurir & PIN kasir</div>
                  </div>
                </div>
                <ChevronRight class="w-4 h-4 text-slate-400" />
              </button>

              <button 
                type="button" 
                @click="showNotificationToast('Membuka Master Parfum & Diskon')"
                class="w-full p-3.5 flex items-center justify-between hover:bg-slate-50 transition-colors cursor-pointer text-left"
              >
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-lg bg-purple-50 text-purple-600 flex items-center justify-center font-bold">🌸</div>
                  <div>
                    <div class="font-bold text-slate-900">Master Parfum & Diskon Promosi</div>
                    <div class="text-[10px] text-slate-400">Varian aroma dan kupon diskon loyalitas</div>
                  </div>
                </div>
                <ChevronRight class="w-4 h-4 text-slate-400" />
              </button>

              <button 
                type="button" 
                @click="showNotificationToast('Membuka Pengeluaran Kas Kecil & Approval')"
                class="w-full p-3.5 flex items-center justify-between hover:bg-slate-50 transition-colors cursor-pointer text-left"
              >
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-lg bg-rose-50 text-rose-600 flex items-center justify-center font-bold">💵</div>
                  <div>
                    <div class="font-bold text-slate-900">Koreksi & Approval Kas Kecil</div>
                    <div class="text-[10px] text-slate-400">Persetujuan pengeluaran operasional outlet</div>
                  </div>
                </div>
                <ChevronRight class="w-4 h-4 text-slate-400" />
              </button>

              <button 
                type="button" 
                @click="showNotificationToast('Membuka Profil Outlet & Konfigurasi Printer')"
                class="w-full p-3.5 flex items-center justify-between hover:bg-slate-50 transition-colors cursor-pointer text-left"
              >
                <div class="flex items-center gap-2.5">
                  <div class="w-7 h-7 rounded-lg bg-amber-50 text-amber-600 flex items-center justify-center font-bold">🖨️</div>
                  <div>
                    <div class="font-bold text-slate-900">Profil Outlet & Template Nota</div>
                    <div class="text-[10px] text-slate-400">Nama outlet, WhatsApp hotline & footer nota</div>
                  </div>
                </div>
                <ChevronRight class="w-4 h-4 text-slate-400" />
              </button>
            </div>
          </div>

          <!-- CASHIER EXCLUSIVE VIEW -->
          <div v-else class="space-y-3.5">
            <div class="text-xs font-black text-slate-800 uppercase tracking-wider flex items-center gap-1.5">
              <Clock class="w-4 h-4 text-sky-600" />
              <span>Shift Kasir & Operasional</span>
            </div>

            <div class="bg-white p-3.5 rounded-2xl border border-slate-200 shadow-xs space-y-2.5">
              <div class="flex justify-between items-center text-xs">
                <span class="text-slate-500">Saldo Awal Laci</span>
                <span class="font-bold text-slate-900">Rp 100.000</span>
              </div>
              <div class="flex justify-between items-center text-xs">
                <span class="text-slate-500">Penerimaan Tunai Shift Ini</span>
                <span class="font-bold text-emerald-600">+Rp 312.000</span>
              </div>
              <div class="flex justify-between items-center text-xs pt-1 border-t border-slate-100">
                <span class="font-black text-slate-800">Total Uang di Laci</span>
                <span class="font-black text-slate-900 text-sm">Rp 412.000</span>
              </div>

              <button 
                type="button" 
                @click="showNotificationToast('Shift berhasil ditutup & rekap terkirim ke Owner!')"
                class="w-full mt-2 py-2.5 rounded-xl bg-slate-900 hover:bg-slate-800 text-white text-xs font-bold transition-all cursor-pointer shadow-xs active:scale-95"
              >
                Serah Terima & Tutup Shift
              </button>
            </div>

            <!-- Printer & Hardware Status -->
            <div class="bg-white p-3.5 rounded-2xl border border-slate-200 shadow-xs space-y-2">
              <div class="text-xs font-bold text-slate-800">Perangkat Kasir</div>
              <div class="flex items-center justify-between text-xs py-1">
                <span class="flex items-center gap-1.5 text-slate-600">
                  <Printer class="w-4 h-4 text-sky-600" /> Mini Thermal BT (58mm)
                </span>
                <span class="text-[10px] font-bold text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full">Terhubung</span>
              </div>
              <div class="flex items-center justify-between text-xs py-1">
                <span class="flex items-center gap-1.5 text-slate-600">
                  <Scan class="w-4 h-4 text-purple-600" /> Kamera Scanner Barcode
                </span>
                <span class="text-[10px] font-bold text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-full">Siap Digunakan</span>
              </div>
            </div>
          </div>
        </div>

      </div>

      <!-- 4. Mobile Sticky Bottom Action Bar (Contextual for POS Tab) -->
      <div v-if="activeNav === 'pos'" class="absolute bottom-[64px] left-0 right-0 p-3 bg-white/95 backdrop-blur-md border-t border-slate-200/90 z-20 flex items-center justify-between gap-3 shadow-lg">
        <div>
          <span class="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">Total Tagihan</span>
          <span class="text-base font-black text-slate-900 tracking-tight">
            Rp {{ formatNumber(calculateTotal) }}
          </span>
        </div>
        <button 
          type="button"
          @click="processNewOrder"
          class="px-5 py-3 rounded-2xl bg-brand-600 hover:bg-brand-700 text-white font-extrabold text-xs flex items-center gap-2 shadow-md transition-transform active:scale-95 cursor-pointer"
        >
          <span>Lanjut Bayar & Cetak</span>
          <ArrowRight class="w-4 h-4" />
        </button>
      </div>

      <!-- 5. Native Bottom Navigation Bar (Thumb Zone) -->
      <div class="h-16 bg-white border-t border-slate-200/90 flex items-center justify-around px-2 shrink-0 z-30 select-none">
        <button 
          type="button"
          @click="activeNav = 'home'"
          class="flex-1 py-1 flex flex-col items-center gap-1 transition-all cursor-pointer"
          :class="activeNav === 'home' ? 'text-brand-600 font-extrabold' : 'text-slate-400 font-medium'"
        >
          <Home class="w-5 h-5" />
          <span class="text-[10px]">Beranda</span>
        </button>

        <button 
          type="button"
          @click="activeNav = 'queue'"
          class="flex-1 py-1 flex flex-col items-center gap-1 transition-all relative cursor-pointer"
          :class="activeNav === 'queue' ? 'text-brand-600 font-extrabold' : 'text-slate-400 font-medium'"
        >
          <Layers class="w-5 h-5" />
          <span class="text-[10px]">Antrian</span>
          <span class="absolute top-0 right-5 w-2 h-2 rounded-full bg-brand-600"></span>
        </button>

        <!-- Center Floating Action POS Button -->
        <button 
          type="button"
          @click="activeNav = 'pos'"
          class="w-12 h-12 -mt-5 rounded-full bg-brand-600 text-white flex items-center justify-center shadow-lg transition-transform active:scale-95 border-2 border-white cursor-pointer"
          title="Input Order Cepat"
        >
          <Plus class="w-6 h-6 stroke-[3]" />
        </button>

        <button 
          type="button"
          @click="activeNav = 'scan'"
          class="flex-1 py-1 flex flex-col items-center gap-1 transition-all cursor-pointer"
          :class="activeNav === 'scan' ? 'text-brand-600 font-extrabold' : 'text-slate-400 font-medium'"
        >
          <Scan class="w-5 h-5" />
          <span class="text-[10px]">Scan QR</span>
        </button>

        <button 
          type="button"
          @click="activeNav = 'menu'"
          class="flex-1 py-1 flex flex-col items-center gap-1 transition-all cursor-pointer"
          :class="activeNav === 'menu' ? 'text-brand-600 font-extrabold' : 'text-slate-400 font-medium hover:text-slate-800'"
        >
          <Menu class="w-5 h-5" />
          <span class="text-[10px]">Menu</span>
        </button>
      </div>

      <!-- 6. Mobile Bottom Sheet: Order Quick Action (Draggable Simulation) -->
      <div 
        v-if="selectedOrderSheet" 
        class="absolute inset-0 bg-slate-900/60 backdrop-blur-xs z-50 flex flex-col justify-end"
        @click.self="selectedOrderSheet = null"
      >
        <div class="bg-white rounded-t-[32px] p-5 space-y-4 shadow-2xl animate-in slide-in-from-bottom duration-200">
          <!-- Sheet Drag Handle -->
          <div class="w-12 h-1.5 bg-slate-200 rounded-full mx-auto"></div>

          <div class="flex items-center justify-between pb-2 border-b border-slate-100">
            <div>
              <h3 class="font-mono text-sm font-black text-slate-900">#{{ selectedOrderSheet.orderNumber }}</h3>
              <p class="text-xs font-bold text-slate-600">{{ selectedOrderSheet.customer }} • Rak: {{ selectedOrderSheet.rack }}</p>
            </div>
            <span class="px-2.5 py-1 rounded-xl text-xs font-black" :class="selectedOrderSheet.badgeClass">
              {{ selectedOrderSheet.status }}
            </span>
          </div>

          <!-- 3 Core One-Handed Actions -->
          <div class="grid grid-cols-2 gap-2.5">
            <button 
              type="button"
              @click="printViaBluetooth"
              class="py-3 rounded-2xl bg-slate-900 text-white font-extrabold text-xs flex items-center justify-center gap-2 shadow-xs active:scale-95 cursor-pointer"
            >
              <Printer class="w-4 h-4 text-sky-400" />
              <span>Cetak BT 58mm</span>
            </button>
            <button 
              type="button"
              @click="progressSheetStatus"
              class="py-3 rounded-2xl bg-brand-600 text-white font-extrabold text-xs flex items-center justify-center gap-2 shadow-xs active:scale-95 cursor-pointer"
            >
              <CheckCircle class="w-4 h-4" />
              <span>Majukan Status</span>
            </button>
          </div>

          <button 
            type="button"
            @click="selectedOrderSheet = null"
            class="w-full py-2.5 rounded-xl bg-slate-100 font-bold text-xs text-slate-600 cursor-pointer"
          >
            Tutup Lembar Aksi
          </button>
        </div>
      </div>

      <!-- Notification Toast Popup -->
      <div 
        v-if="toastMessage" 
        class="absolute top-16 left-4 right-4 bg-slate-900 text-white px-4 py-3 rounded-2xl shadow-xl flex items-center gap-2 text-xs font-bold z-50 animate-in fade-in slide-in-from-top-2"
      >
        <Check class="w-4 h-4 text-emerald-400 shrink-0" />
        <span>{{ toastMessage }}</span>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  ArrowLeft,
  Wifi,
  BatteryCharging,
  ChevronDown,
  Printer,
  Bell,
  Banknote,
  TrendingUp,
  Clock,
  AlertCircle,
  Scale,
  QrCode,
  Layers,
  Wallet,
  Flame,
  ArrowRightCircle,
  User,
  Zap,
  ArrowRight,
  Home,
  Plus,
  Scan,
  Menu,
  CheckCircle,
  Check
} from 'lucide-vue-next'

const activeNav = ref('home')
const bluetoothConnected = ref(true)
const toastMessage = ref('')
const activeCategory = ref('cat-1')
const currentWeight = ref(8.0)
const selectedPerfume = ref('Akasia')
const rackNumber = ref('A-02')

const selectedCustomer = ref({
  name: 'Nuryanto (0812-3456-7890)',
  phone: '081234567890'
})

const serviceCategories = [
  { id: 'cat-1', name: 'Self Service' },
  { id: 'cat-2', name: 'Cuci Setrika' },
  { id: 'cat-3', name: 'Setrika Saja' },
  { id: 'cat-4', name: 'Bed Cover' },
]

const services = [
  { id: 's1', categoryId: 'cat-1', category: 'Self Service', name: 'Cuci Basah (Self Service)', price: 8000, unit: 'Tabung', duration: 'Instant' },
  { id: 's2', categoryId: 'cat-1', category: 'Self Service', name: 'Cuci Kering (Self Service)', price: 20000, unit: 'Tabung', duration: 'Instant' },
  { id: 's3', categoryId: 'cat-2', category: 'Cuci Setrika', name: 'Cuci Setrika Reguler', price: 7000, unit: 'Kg', duration: '3 Hari' },
  { id: 's4', categoryId: 'cat-2', category: 'Cuci Setrika', name: 'Cuci Setrika Kilat', price: 11000, unit: 'Kg', duration: '24 Jam' },
  { id: 's5', categoryId: 'cat-3', category: 'Setrika Saja', name: 'Setrika Saja Reguler', price: 6000, unit: 'Kg', duration: '3 Hari' },
  { id: 's6', categoryId: 'cat-4', category: 'Bed Cover', name: 'Bed Cover Sedang', price: 30000, unit: 'Pcs', duration: '3 Hari' },
]

const selectedService = ref(services[0])

const filteredServices = computed(() => {
  return services.filter(s => s.categoryId === activeCategory.value)
})

const urgentOrders = ref([
  { id: '1', num: '261120000002', status: 'PROSES', customer: 'Nuryanto', weight: 8.0, rack: 'A-01', due: '14:00 (2 jam lagi)', badgeClass: 'bg-amber-100 text-amber-800' },
  { id: '2', num: '261120000005', status: 'ANTRIAN', customer: 'Bu Sarah', weight: 4.5, rack: 'B-03', due: '17:00 (Hari ini)', badgeClass: 'bg-sky-100 text-sky-800' },
])

const currentRole = ref('CASHIER') // 'CASHIER' or 'OWNER'

const toggleRole = () => {
  currentRole.value = currentRole.value === 'CASHIER' ? 'OWNER' : 'CASHIER'
  showNotificationToast(`Mode beralih ke: ${currentRole.value === 'OWNER' ? '👑 Owner (Akses Penuh)' : '👤 Kasir (Operasional)'}`)
}

const queueTabs = computed(() => [
  { id: 'ALL', label: 'Semua', count: queueList.value.length },
  { id: 'ANTRIAN', label: 'Antrian', count: queueList.value.filter((q) => q.status === 'ANTRIAN').length },
  { id: 'PROSES', label: 'Proses', count: queueList.value.filter((q) => q.status === 'PROSES').length },
  { id: 'SIAP AMBIL', label: 'Siap', count: queueList.value.filter((q) => q.status === 'SIAP AMBIL').length },
])
const activeQueueTab = ref('ALL')

const filteredQueueList = computed(() => {
  if (activeQueueTab.value === 'ALL') return queueList.value
  return queueList.value.filter((q) => q.status === activeQueueTab.value)
})

const quickAdvanceQueue = (item) => {
  if (item.status === 'ANTRIAN') {
    item.status = 'PROSES'
    item.badgeClass = 'bg-amber-100 text-amber-800'
  } else if (item.status === 'PROSES') {
    item.status = 'SIAP AMBIL'
    item.badgeClass = 'bg-purple-100 text-purple-800'
  } else if (item.status === 'SIAP AMBIL') {
    item.status = 'SELESAI'
    item.badgeClass = 'bg-emerald-100 text-emerald-800'
  }
  showNotificationToast(`Status #${item.orderNumber} dimajukan ke ${item.status}!`)
}

const queueList = ref([
  { id: '1', orderNumber: '261120000002', status: 'PROSES', customer: 'Nuryanto', total: 28000, service: 'Cuci Basah (Self Service)', weight: 8.0, rack: 'A-01', paid: true, badgeClass: 'bg-amber-100 text-amber-800' },
  { id: '2', orderNumber: '261120000001', status: 'SIAP AMBIL', customer: 'Thiya', total: 56000, service: 'Cuci Setrika Reguler', weight: 8.0, rack: 'A-02', paid: true, badgeClass: 'bg-purple-100 text-purple-800' },
  { id: '3', orderNumber: '261119000007', status: 'ANTRIAN', customer: 'Doni Pratama', total: 152000, service: 'Cuci Setrika Turbo', weight: 9.5, rack: 'C-01', paid: false, badgeClass: 'bg-sky-100 text-sky-800' },
])

const selectedOrderSheet = ref(null)

const calculateTotal = computed(() => {
  return Math.round(selectedService.value.price * (currentWeight.value || 1))
})

const showNotificationToast = (msg) => {
  toastMessage.value = msg
  setTimeout(() => { toastMessage.value = '' }, 2500)
}

const selectService = (srv) => {
  selectedService.value = srv
  if (srv.unit === 'Tabung' && currentWeight.value < 1) currentWeight.value = 1
}

const adjustWeight = (delta) => {
  const next = Math.max(0.5, currentWeight.value + delta)
  currentWeight.value = Math.round(next * 10) / 10
}

const toggleBluetooth = () => {
  bluetoothConnected.value = !bluetoothConnected.value
  showNotificationToast(bluetoothConnected.value ? 'Printer Bluetooth 58mm Terhubung!' : 'Printer Bluetooth Diputus')
}

const processNewOrder = () => {
  showNotificationToast(`Order baru berhasil dibuat & dikirim via WA! Total Rp ${formatNumber(calculateTotal.value)}`)
  activeNav.value = 'queue'
}

const openOrderDetailSheet = (order) => {
  selectedOrderSheet.value = order
}

const printViaBluetooth = () => {
  showNotificationToast('Struk berhasil dicetak ke Mini Printer BT 58mm!')
}

const progressSheetStatus = () => {
  if (!selectedOrderSheet.value) return
  if (selectedOrderSheet.value.status === 'ANTRIAN') selectedOrderSheet.value.status = 'PROSES'
  else if (selectedOrderSheet.value.status === 'PROSES') selectedOrderSheet.value.status = 'SIAP AMBIL'
  else if (selectedOrderSheet.value.status === 'SIAP AMBIL') selectedOrderSheet.value.status = 'SELESAI'
  showNotificationToast(`Status #${selectedOrderSheet.value.orderNumber} diubah ke ${selectedOrderSheet.value.status}`)
  setTimeout(() => { selectedOrderSheet.value = null }, 800)
}

const progressOrder = (order) => {
  order.status = 'SIAP AMBIL'
  order.badgeClass = 'bg-purple-100 text-purple-800'
  showNotificationToast(`Status #${order.num} dimajukan ke SIAP AMBIL!`)
}

const simulateScanSuccess = () => {
  selectedOrderSheet.value = queueList.value[0]
  showNotificationToast('QR Code Terbaca! Membuka detail cucian...')
}

const formatNumber = (val) => {
  return (val || 0).toLocaleString('id-ID')
}
</script>
