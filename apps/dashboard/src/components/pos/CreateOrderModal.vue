<template>
  <div
    v-if="modelValue"
    class="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 bg-slate-900/60 backdrop-blur-xs animate-in fade-in"
  >
    <div
      class="bg-white rounded-2xl max-w-4xl w-full p-4 sm:p-6 shadow-2xl border border-slate-100 flex flex-col max-h-[96vh] sm:max-h-[92vh]"
    >
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-100 shrink-0">
        <div class="flex items-center gap-2.5">
          <div
            class="w-10 h-10 rounded-xl bg-brand-50 text-brand-600 flex items-center justify-center shadow-xs"
          >
            <PlusCircle class="w-6 h-6" />
          </div>
          <div>
            <h3 class="font-black text-lg text-slate-900">Transaksi Kasir Walk-in</h3>
            <p class="text-xs text-slate-500">Pencatatan cucian masuk, timbangan & pembayaran</p>
          </div>
        </div>
        <button
          class="p-2 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition-all cursor-pointer"
          @click="close"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Main Form (Scrollable) -->
      <div class="flex-1 overflow-y-auto py-3 space-y-5 pr-1 text-slate-800">
        <!-- STEP 1: Pelanggan (Search / Add New) -->
        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <User class="w-4 h-4 text-brand-600" />
              <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
                1. Data Pelanggan
              </h4>
            </div>
            <button
              type="button"
              class="text-xs font-bold text-brand-600 hover:underline cursor-pointer"
              @click="isNewCustomer = !isNewCustomer"
            >
              {{ isNewCustomer ? 'Pilih Pelanggan Terdaftar' : '+ Pelanggan Baru' }}
            </button>
          </div>

          <!-- Existing Customer Selected Pill -->
          <div
            v-if="!isNewCustomer && selectedCustomer"
            class="p-3 rounded-xl bg-white border border-brand-200 flex items-center justify-between shadow-xs"
          >
            <div>
              <p class="font-bold text-sm text-slate-900">{{ selectedCustomer.full_name }}</p>
              <p class="text-xs text-slate-500">
                {{ selectedCustomer.phone }} • {{ selectedCustomer.address || 'Kalibaru' }}
              </p>
            </div>
            <button
              class="text-xs font-bold text-rose-600 hover:underline cursor-pointer"
              @click="selectedCustomer = null"
            >
              Ganti
            </button>
          </div>

          <!-- Existing Customer Autocomplete Input -->
          <div v-else-if="!isNewCustomer" class="relative">
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-3 pointer-events-none" />
            <input
              v-model="customerSearchQuery"
              type="text"
              placeholder="Ketik nama atau nomor WhatsApp pelanggan (cth: Maya, 0852...)"
              class="w-full pl-9 pr-8 py-2.5 rounded-xl border border-slate-200 bg-white text-xs font-medium focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              @input="onCustomerSearch"
            />
            <button
              v-if="customerSearchQuery"
              class="absolute right-3 top-3 text-slate-400 hover:text-slate-600"
              @click="clearCustomerSearch"
            >
              <X class="w-3.5 h-3.5" />
            </button>

            <!-- Search Results Dropdown -->
            <div
              v-if="customerSearchQuery && customerResults.length > 0"
              class="absolute top-full left-0 right-0 mt-1 bg-white rounded-xl border border-slate-200 shadow-xl z-20 max-h-48 overflow-y-auto divide-y divide-slate-100 animate-in fade-in slide-in-from-top-1"
            >
              <div
                v-for="c in customerResults"
                :key="c.id"
                class="p-2.5 hover:bg-brand-50/70 cursor-pointer text-xs flex justify-between items-center transition-colors"
                @click="selectCustomer(c)"
              >
                <div>
                  <span class="font-bold text-slate-900">{{ c.full_name }}</span>
                  <span class="text-brand-600 font-semibold ml-1.5">({{ c.phone }})</span>
                </div>
                <span class="text-[10px] text-slate-500 truncate max-w-[150px]">{{
                  c.address || 'Kalibaru'
                }}</span>
              </div>
            </div>

            <!-- Empty Search State -->
            <div
              v-else-if="customerSearchQuery && customerResults.length === 0 && !searchingCustomer"
              class="absolute top-full left-0 right-0 mt-1 bg-white rounded-xl border border-slate-200 shadow-xl z-20 p-3 text-center text-xs text-slate-500"
            >
              Pelanggan tidak ditemukan untuk "{{ customerSearchQuery }}".
              <button
                type="button"
                class="block mx-auto mt-1.5 text-xs font-bold text-brand-600 hover:underline cursor-pointer"
                @click="quickAddNewCustomer"
              >
                + Tambah Sebagai Pelanggan Baru
              </button>
            </div>
          </div>

          <!-- Inline Form New Customer -->
          <div v-else class="grid grid-cols-1 sm:grid-cols-3 gap-2.5 pt-1">
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">Nama Lengkap *</label>
              <input
                v-model="newCustomer.full_name"
                type="text"
                required
                placeholder="Nama pelanggan"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Nomor WhatsApp *</label
              >
              <input
                v-model="newCustomer.phone"
                type="tel"
                required
                placeholder="0852... / 628..."
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Alamat (Opsional)</label
              >
              <input
                v-model="newCustomer.address"
                type="text"
                placeholder="Rt/Rw, patokan rumah"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>
          </div>
        </div>

        <!-- STEP 2: Pilih Layanan & Keranjang Multi-Layanan -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <ShoppingBag class="w-4 h-4 text-brand-600" />
              <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
                2. Pilih Layanan Laundry
              </h4>
            </div>
            <span class="text-[11px] text-slate-500 font-semibold"
              >{{ cartItems.length }} Layanan di keranjang</span
            >
          </div>

          <!-- Quick Popular Service Chips (1-Tap for Kasir) -->
          <div>
            <span class="text-[11px] font-semibold text-slate-500 block mb-1.5"
              >Layanan Terpopuler (1-Klik Cepat):</span
            >
            <div class="flex flex-wrap gap-1.5">
              <button
                v-for="fav in favoriteServices"
                :key="fav.id"
                type="button"
                class="px-3 py-1.5 rounded-xl border border-slate-200 bg-slate-50 hover:bg-brand-50 hover:border-brand-300 hover:text-brand-700 text-xs font-bold text-slate-700 transition-all flex items-center gap-1.5 shadow-2xs cursor-pointer active:scale-95"
                @click="openQuantityModal(fav)"
              >
                <Plus class="w-3.5 h-3.5 text-brand-600" />
                <span>{{ fav.service_name }}</span>
                <span class="text-[10px] text-slate-400 font-normal">
                  Rp {{ formatRupiah(fav.price) }} / {{ fav.unit }}
                </span>
              </button>
            </div>
          </div>

          <!-- Trigger Buka Modal Katalog Layanan -->
          <div>
            <button
              type="button"
              class="w-full py-3.5 px-4 rounded-xl border-2 border-dashed border-brand-300 hover:border-brand-500 bg-brand-50/50 hover:bg-brand-50 text-brand-800 font-bold text-xs flex items-center justify-between transition-all group cursor-pointer shadow-2xs"
              @click="showServicePickerModal = true"
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-9 h-9 rounded-xl bg-brand-600 text-white flex items-center justify-center shadow-xs group-hover:scale-105 transition-transform"
                >
                  <Plus class="w-5 h-5" />
                </div>
                <div class="text-left">
                  <span
                    class="block font-black text-xs sm:text-sm text-slate-900 group-hover:text-brand-600"
                  >
                    + Buka Katalog Layanan Laundry (Semua Kategori)
                  </span>
                  <span class="text-[11px] text-slate-500 font-normal">
                    Pilih Self Service (Cuci Basah/Kering), Kiloan Reguler/Kilat, Satuan, & Karpet
                  </span>
                </div>
              </div>
              <div class="flex items-center gap-1 text-brand-600 font-bold text-xs">
                <span>Pilih</span>
                <ChevronRight class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </div>
            </button>
          </div>

          <!-- Active Cart Table -->
          <div
            v-if="cartItems.length > 0"
            class="mt-3 rounded-xl border border-slate-200 overflow-hidden shadow-xs"
          >
            <div
              class="bg-slate-100 px-3.5 py-2 flex justify-between text-xs font-bold text-slate-700"
            >
              <span>Keranjang Cucian ({{ cartItems.length }} item)</span>
              <span>Subtotal</span>
            </div>
            <div class="divide-y divide-slate-100 bg-white text-xs">
              <div
                v-for="(item, idx) in cartItems"
                :key="idx"
                class="px-3.5 py-2.5 flex items-center justify-between"
              >
                <div class="flex items-center gap-2.5">
                  <button
                    class="p-1 text-slate-300 hover:text-rose-600 transition-colors cursor-pointer"
                    title="Hapus"
                    @click="removeItem(idx)"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                  <div>
                    <div class="flex items-center gap-1.5">
                      <span class="font-bold text-slate-900">{{ item.service_name }}</span>
                      <span
                        v-if="item.category === 'TABUNG'"
                        class="px-1.5 py-0.5 rounded text-[10px] font-extrabold bg-blue-100 text-blue-700"
                      >
                        Self Service
                      </span>
                    </div>
                    <p class="text-[11px] text-slate-500">
                      {{ item.quantity }} {{ item.unit }} x Rp
                      {{ formatRupiah(item.price_per_unit) }}
                      <span v-if="item.notes" class="text-amber-700 font-semibold ml-1">
                        • {{ item.notes }}
                      </span>
                    </p>
                  </div>
                </div>
                <span class="font-black text-slate-900">Rp {{ formatRupiah(item.subtotal) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- STEP 3: Konfigurasi Order (Parfum, Diskon, Rak, SLA) -->
        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-3">
          <div class="flex items-center gap-2">
            <Sparkles class="w-4 h-4 text-brand-600" />
            <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
              3. Pilihan Tambahan & Lokasi Rak
            </h4>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-4 gap-2.5">
            <!-- Parfum -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">Aroma Parfum</label>
              <select
                v-model="orderConfig.perfume_name"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-medium"
              >
                <template v-if="perfumesList.length > 0">
                  <option v-for="p in perfumesList" :key="p.id" :value="p.name">
                    {{ p.name }}
                  </option>
                </template>
                <template v-else>
                  <option value="Akasia">Akasia (Segar Favorit)</option>
                  <option value="Sakura">Sakura (Manis Lembut)</option>
                  <option value="Ocean Fresh">Ocean Fresh (Semilir Laut)</option>
                  <option value="Snappy">Snappy (Elegan Segar)</option>
                  <option value="Downy Black">Downy Black (Mewah Modern)</option>
                </template>
              </select>
            </div>

            <!-- Diskon / Promo -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">Diskon / Promo</label>
              <select
                v-model="orderConfig.discount_type"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-medium"
                @change="applyDiscount"
              >
                <option value="NONE">Tanpa Diskon</option>
                <option value="PROMO_10">Promo Opening (10%)</option>
                <option value="LOYALTY">Loyalitas (Gratis 1x Cuci)</option>
                <option value="CUSTOM">Custom Nominal Rp</option>
              </select>
            </div>

            <!-- Custom Diskon input if selected -->
            <div v-if="orderConfig.discount_type === 'CUSTOM'">
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Potongan Diskon (Rp)</label
              >
              <input
                v-model.number="orderConfig.custom_discount"
                type="number"
                placeholder="5000"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
                @input="applyDiscount"
              />
            </div>

            <!-- Nomor Rak Simpan -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Nomor Rak Simpan</label
              >
              <input
                v-model="orderConfig.rack_number"
                type="text"
                placeholder="Contoh: A-01"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-bold"
              />
            </div>

            <!-- Jumlah Potong Pakaian (Pcs) Opsional -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">
                Jumlah Potong Pcs <span class="text-[10px] text-slate-400 font-normal">(Opsional)</span>
              </label>
              <input
                v-model="orderConfig.piece_count"
                type="number"
                min="0"
                placeholder="Cth: 14"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-bold"
              />
            </div>

            <!-- Estimasi Selesai (SLA) -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Estimasi Selesai (SLA)</label
              >
              <input
                v-model="orderConfig.estimated_completion_at"
                type="datetime-local"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-medium"
              />
            </div>
          </div>

          <!-- Catatan Order -->
          <div>
            <label class="block text-[11px] font-bold text-slate-600 mb-1"
              >Catatan Khusus (Opsional)</label
            >
            <input
              v-model="orderConfig.notes"
              type="text"
              placeholder="Contoh: Ada noda oli di lengan kiri, baju luntur dipisah..."
              class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
            />
          </div>
        </div>

        <!-- STEP 4: Opsi Pembayaran Awal -->
        <div class="p-4 rounded-xl border border-slate-200 bg-slate-50/50 space-y-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <CreditCard class="w-4 h-4 text-brand-600" />
              <h4 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
                4. Pembayaran Awal
              </h4>
            </div>
            <!-- Toggle Bayar Sekarang vs Bayar Nanti -->
            <div
              class="flex items-center gap-1.5 p-1 bg-white rounded-xl border border-slate-200 text-xs"
            >
              <button
                type="button"
                class="px-2.5 py-1 rounded-lg font-bold transition-all cursor-pointer"
                :class="
                  paymentOption === 'SEKARANG'
                    ? 'bg-brand-600 text-white'
                    : 'text-slate-500 hover:text-slate-800'
                "
                @click="paymentOption = 'SEKARANG'"
              >
                Bayar Sekarang
              </button>
              <button
                type="button"
                class="px-2.5 py-1 rounded-lg font-bold transition-all cursor-pointer"
                :class="
                  paymentOption === 'NANTI'
                    ? 'bg-slate-900 text-white'
                    : 'text-slate-500 hover:text-slate-800'
                "
                @click="paymentOption = 'NANTI'"
              >
                Bayar Saat Ambil
              </button>
            </div>
          </div>

          <div
            v-if="paymentOption === 'SEKARANG'"
            class="grid grid-cols-1 sm:grid-cols-3 gap-2.5 pt-1"
          >
            <!-- Pilihan Metode -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1">Metode Bayar</label>
              <select
                v-model="paymentMethod"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs focus:ring-2 focus:ring-brand-500 focus:outline-hidden font-bold"
              >
                <option value="TUNAI">Tunai / Cash</option>
                <option value="BCA">Transfer BCA (4140500068)</option>
                <option value="QRIS">QRIS AstraPay</option>
              </select>
            </div>

            <!-- Nominal Uang Diterima -->
            <div>
              <label class="block text-[11px] font-bold text-slate-600 mb-1"
                >Nominal Diterima (Rp)</label
              >
              <input
                v-model.number="paidAmount"
                type="number"
                min="0"
                class="w-full px-3 py-2 rounded-xl border border-slate-200 bg-white text-xs font-bold focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
            </div>

            <!-- Kembalian Info -->
            <div
              v-if="paymentMethod === 'TUNAI' && paidAmount > grandTotal"
              class="flex flex-col justify-end"
            >
              <div
                class="p-2 rounded-xl bg-emerald-50 border border-emerald-200 text-xs text-emerald-800"
              >
                <span class="text-[10px] text-emerald-600 block">Kembalian:</span>
                <span class="font-extrabold text-sm text-emerald-700"
                  >Rp {{ formatRupiah(paidAmount - grandTotal) }}</span
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer Bar (Total & Submit) -->
      <div
        class="pt-3 border-t border-slate-100 flex flex-col sm:flex-row items-center justify-between gap-3 shrink-0"
      >
        <div class="flex items-center gap-4 w-full sm:w-auto justify-between sm:justify-start">
          <div>
            <span class="text-[10px] text-slate-400 block font-semibold">TOTAL TAGIHAN</span>
            <span class="font-black text-xl text-slate-900">Rp {{ formatRupiah(grandTotal) }}</span>
          </div>
          <div v-if="discountAmount > 0" class="text-xs text-emerald-600 font-bold">
            Hemat: Rp {{ formatRupiah(discountAmount) }}
          </div>
        </div>

        <div class="flex items-center gap-2 w-full sm:w-auto">
          <button
            type="button"
            class="w-1/3 sm:w-auto px-4 py-2.5 rounded-xl border border-slate-200 text-slate-700 font-bold text-xs hover:bg-slate-50 cursor-pointer"
            @click="close"
          >
            Batal
          </button>
          <button
            type="button"
            :disabled="creating || cartItems.length === 0"
            class="w-2/3 sm:w-auto px-6 py-2.5 rounded-xl bg-gradient-to-r from-brand-600 to-sky-500 hover:from-brand-700 hover:to-sky-600 text-white font-extrabold text-xs shadow-md shadow-brand-500/20 flex items-center justify-center gap-1.5 disabled:opacity-50 transition-all cursor-pointer"
            @click="submitOrder"
          >
            <span v-if="creating">Menyimpan...</span>
            <span v-else>BUAT ORDERAN KASIR</span>
          </button>
        </div>
      </div>
    </div>

    <!-- SUB-MODAL 1: Service Picker Modal (Grouped by Category) -->
    <div
      v-if="showServicePickerModal"
      class="fixed inset-0 z-60 flex items-center justify-center p-3 sm:p-4 bg-slate-950/70 backdrop-blur-xs animate-in fade-in"
    >
      <div
        class="bg-white rounded-2xl max-w-2xl w-full p-4 sm:p-5 shadow-2xl border border-slate-100 flex flex-col max-h-[88vh]"
      >
        <!-- Service Picker Header -->
        <div class="flex items-center justify-between pb-3 border-b border-slate-100 shrink-0">
          <div>
            <h3 class="font-black text-base text-slate-900">Katalog Layanan FAN Laundry</h3>
            <p class="text-xs text-slate-500">
              Pilih layanan laundry untuk dimasukkan ke keranjang
            </p>
          </div>
          <button
            class="p-1.5 rounded-xl text-slate-400 hover:text-slate-600 hover:bg-slate-100 cursor-pointer"
            @click="showServicePickerModal = false"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Filter & Search Inside Service Modal -->
        <div class="pt-3 pb-2 space-y-2 shrink-0">
          <div class="relative">
            <Search class="w-4 h-4 text-slate-400 absolute left-3 top-3 pointer-events-none" />
            <input
              v-model="serviceModalSearch"
              type="text"
              placeholder="Cari layanan (cth: Cuci Basah, Setrika, Bed Cover, Karpet)..."
              class="w-full pl-9 pr-4 py-2 rounded-xl border border-slate-200 bg-slate-50 focus:bg-white text-xs font-medium focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
            />
          </div>

          <!-- Category Filter Pills -->
          <div class="flex items-center gap-1.5 overflow-x-auto pb-1 text-xs">
            <button
              v-for="tab in categoryTabs"
              :key="tab.id"
              type="button"
              class="px-2.5 py-1 rounded-lg font-bold shrink-0 transition-all cursor-pointer flex items-center gap-1.5"
              :class="
                activeServiceTab === tab.id
                  ? 'bg-slate-900 text-white shadow-xs'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              "
              @click="activeServiceTab = tab.id"
            >
              <span>{{ tab.icon }}</span>
              <span>{{ tab.name }}</span>
              <span
                class="text-[10px] px-1.5 py-0.2 rounded-full font-semibold"
                :class="
                  activeServiceTab === tab.id
                    ? 'bg-slate-800 text-white'
                    : 'bg-slate-200 text-slate-700'
                "
              >
                {{ tab.count }}
              </span>
            </button>
          </div>
        </div>

        <!-- Service List (Grouped Sections, Scrollable) -->
        <div class="flex-1 overflow-y-auto py-2 space-y-4 pr-1">
          <div v-for="grp in filteredServiceGroups" :key="grp.categoryId" class="space-y-2">
            <div class="flex items-center gap-2 border-b border-slate-100 pb-1">
              <span class="text-xs font-black text-slate-800">{{ grp.icon }} {{ grp.categoryName }}</span>
              <span class="text-[10px] text-slate-400 font-semibold"
                >({{ grp.items.length }} layanan)</span
              >
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
              <div
                v-for="s in grp.items"
                :key="s.id"
                class="p-3 rounded-xl border border-slate-200 hover:border-brand-300 hover:bg-brand-50/30 transition-all flex items-center justify-between cursor-pointer group"
                @click="selectServiceFromModal(s)"
              >
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="font-bold text-xs text-slate-900 group-hover:text-brand-700">
                      {{ s.service_name }}
                    </span>
                    <span
                      v-if="s.service_code"
                      class="px-1.5 py-0.2 rounded text-[9px] font-mono font-bold bg-slate-100 text-slate-600 border border-slate-200"
                    >
                      {{ s.service_code }}
                    </span>
                    <span
                      v-if="s.unit === 'Tabung' || (s.category && s.category.toLowerCase().includes('self service'))"
                      class="px-1.5 py-0.2 rounded text-[9px] font-black bg-blue-100 text-blue-700"
                    >
                      Kapasitas 8 Kg
                    </span>
                  </div>
                  <div class="flex items-center gap-2 mt-0.5">
                    <span class="text-xs font-black text-brand-600">
                      Rp {{ formatRupiah(s.price) }}
                    </span>
                    <span class="text-[10px] text-slate-400">/ {{ s.unit }}</span>
                    <span v-if="s.duration_day > 0" class="text-[10px] text-slate-400">
                      • {{ s.duration_day }} Hari
                    </span>
                    <span v-else class="text-[10px] text-emerald-600 font-semibold">
                      • Siap ~3 Jam
                    </span>
                  </div>
                </div>

                <button
                  type="button"
                  class="px-2.5 py-1.5 rounded-lg bg-brand-50 group-hover:bg-brand-600 group-hover:text-white text-brand-700 font-bold text-xs flex items-center gap-1 transition-all cursor-pointer"
                >
                  <Plus class="w-3.5 h-3.5" />
                  <span>Pilih</span>
                </button>
              </div>
            </div>
          </div>

          <div
            v-if="filteredServiceGroups.length === 0"
            class="py-12 text-center text-xs text-slate-400"
          >
            Layanan tidak ditemukan untuk kata kunci tersebut.
          </div>
        </div>
      </div>
    </div>

    <!-- SUB-MODAL 2: Quantity & Timbangan Calculator Sub-Dialog -->
    <div
      v-if="showQtyModal"
      class="fixed inset-0 z-70 flex items-center justify-center p-3 sm:p-4 bg-slate-950/70 backdrop-blur-xs animate-in fade-in"
    >
      <div
        class="bg-white rounded-2xl max-w-md w-full p-5 shadow-2xl border border-slate-100 space-y-4 max-h-[92vh] overflow-y-auto"
      >
        <!-- Header -->
        <div class="flex items-center justify-between pb-2 border-b border-slate-100">
          <div>
            <div class="flex items-center gap-1.5">
              <h4 class="font-black text-base text-slate-900">
                {{ currentService?.service_name }}
              </h4>
              <span
                v-if="isSelfService"
                class="px-2 py-0.5 rounded text-[10px] font-black bg-blue-100 text-blue-700 uppercase"
              >
                Self Service
              </span>
            </div>
            <p v-if="isSelfService" class="text-xs text-brand-600 font-extrabold">
              Rp {{ formatRupiah(currentService?.price) }} / Tabung
              <span class="text-[10px] text-slate-500 font-normal">
                (Kapasitas 8 Kg + toleransi 0.5 Kg)
              </span>
            </p>
            <p v-else class="text-xs text-brand-600 font-extrabold">
              Rp {{ formatRupiah(currentService?.price) }} / {{ currentService?.unit }}
            </p>
          </div>
          <button
            class="p-1 rounded-lg text-slate-400 hover:text-slate-600 cursor-pointer"
            @click="showQtyModal = false"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- SPECIFIC FORM: SELF SERVICE (Cuci Basah / Cuci Kering) -->
        <div v-if="isSelfService" class="space-y-3.5">
          <div
            class="p-3 rounded-xl bg-blue-50/70 border border-blue-200 text-xs text-blue-900 space-y-1"
          >
            <div class="font-extrabold flex items-center gap-1.5 text-blue-950">
              <Scale class="w-4 h-4 text-blue-600" />
              <span>Aturan Kapasitas Tabung Mesin Cuci:</span>
            </div>
            <p class="text-[11px] text-blue-800">
              Kapasitas normal 1 tabung adalah <strong>{{ tubeCapacity }} Kg</strong> dengan batas
              toleransi <strong>+{{ tubeTolerance }} Kg</strong> (Maksimal aman
              <strong>{{ maxSafeWeight }} Kg</strong>).
            </p>
          </div>

          <!-- Weight Input (Kg) -->
          <div>
            <label class="block text-xs font-bold text-slate-700 mb-1">
              Berat Timbangan Cucian (Kg):
            </label>
            <div class="flex items-center gap-2">
              <button
                type="button"
                class="w-11 h-11 rounded-xl bg-slate-100 font-black text-slate-700 hover:bg-slate-200 text-base cursor-pointer"
                @click="
                  selfServiceWeight = Math.max(0.5, Number((selfServiceWeight - 0.5).toFixed(2)))
                "
              >
                -
              </button>
              <div class="relative flex-1">
                <input
                  v-model.number="selfServiceWeight"
                  type="number"
                  step="0.1"
                  min="0.5"
                  class="w-full py-2.5 text-center text-xl font-black rounded-xl border border-slate-200 focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
                />
                <span class="absolute right-3.5 top-3 text-xs font-bold text-slate-400">Kg</span>
              </div>
              <button
                type="button"
                class="w-11 h-11 rounded-xl bg-slate-100 font-black text-slate-700 hover:bg-slate-200 text-base cursor-pointer"
                @click="selfServiceWeight = Number((selfServiceWeight + 0.5).toFixed(2))"
              >
                +
              </button>
            </div>
          </div>

          <!-- Quick Presets for Self Service -->
          <div class="flex flex-wrap gap-1.5 justify-center">
            <button
              v-for="w in [6.0, 7.5, 8.0, 8.5, 9.0, 10.0, 12.0]"
              :key="w"
              type="button"
              class="px-2.5 py-1 rounded-lg text-xs font-bold transition-all cursor-pointer"
              :class="
                selfServiceWeight === w
                  ? 'bg-brand-600 text-white shadow-2xs'
                  : 'bg-slate-100 text-slate-700 hover:bg-slate-200'
              "
              @click="selfServiceWeight = w"
            >
              {{ w }} Kg
            </button>
          </div>

          <!-- CAPACITY TOLERANCE STATUS CHECK -->
          <!-- Case A: Within normal + tolerance (<= maxSafeWeight Kg) -->
          <div
            v-if="selfServiceWeight <= maxSafeWeight"
            class="p-3 rounded-xl bg-emerald-50 border border-emerald-200 flex items-start gap-2.5"
          >
            <CheckCircle2 class="w-5 h-5 text-emerald-600 shrink-0 mt-0.5" />
            <div class="text-xs">
              <span class="font-extrabold text-emerald-900 block">Kapasitas Aman (1 Tabung)</span>
              <p class="text-[11px] text-emerald-700">
                Berat cucian {{ selfServiceWeight }} Kg cukup dicuci dalam
                <strong>1 Tabung</strong> (Maks {{ maxSafeWeight }} Kg termasuk toleransi).
              </p>
            </div>
          </div>

          <!-- Case B: Exceeds maxSafeWeight Kg (> maxSafeWeight Kg) -->
          <div v-else class="p-3.5 rounded-xl bg-amber-50 border-2 border-amber-300 space-y-2.5">
            <div class="flex items-start gap-2">
              <AlertTriangle class="w-5 h-5 text-amber-600 shrink-0 mt-0.5" />
              <div class="text-xs">
                <span class="font-black text-amber-950 block">
                  Melebihi Toleransi 1 Tabung ({{ selfServiceWeight }} Kg > {{ maxSafeWeight }} Kg)
                </span>
                <p class="text-[11px] text-amber-800">
                  Terdapat kelebihan muatan sekitar <strong>{{ excessWeightKg }} Kg</strong>.
                  Silakan konfirmasi opsi dengan pelanggan:
                </p>
              </div>
            </div>

            <!-- 2 Options for Kasir / Customer -->
            <div class="grid grid-cols-1 gap-2 pt-1">
              <button
                type="button"
                class="p-2.5 rounded-xl border text-left transition-all cursor-pointer flex items-center justify-between"
                :class="
                  selfServiceAction === 'TWO_TUBES'
                    ? 'border-brand-500 bg-brand-50/80 ring-2 ring-brand-500/30'
                    : 'border-slate-200 bg-white hover:bg-slate-50'
                "
                @click="selfServiceAction = 'TWO_TUBES'"
              >
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="font-bold text-xs text-slate-900">Opsi 1: Jadikan 2 Tabung</span>
                    <span
                      class="px-1.5 py-0.2 rounded text-[9px] font-black bg-emerald-100 text-emerald-800"
                    >
                      Rekomendasi
                    </span>
                  </div>
                  <p class="text-[10px] text-slate-500">
                    Cucian dibagi 2 mesin, pakaian bersih & bilas maksimal (Total: 2 Tabung).
                  </p>
                </div>
                <div class="text-right pl-2 shrink-0">
                  <span class="text-xs font-black text-brand-700">
                    Rp {{ formatRupiah(currentService.price * 2) }}
                  </span>
                  <span class="block text-[9px] text-slate-400">2 Tabung</span>
                </div>
              </button>

              <button
                type="button"
                class="p-2.5 rounded-xl border text-left transition-all cursor-pointer flex items-center justify-between"
                :class="
                  selfServiceAction === 'SPLIT_EXCESS'
                    ? 'border-amber-500 bg-amber-50/80 ring-2 ring-amber-500/30'
                    : 'border-slate-200 bg-white hover:bg-slate-50'
                "
                @click="selfServiceAction = 'SPLIT_EXCESS'"
              >
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="font-bold text-xs text-slate-900"
                      >Opsi 2: Tetap 1 Tabung & Sisihkan</span
                    >
                  </div>
                  <p class="text-[10px] text-slate-500">
                    8.0 Kg dicuci sekarang, kelebihan
                    <strong>~{{ excessWeightKg }} Kg disisihkan</strong>/dibawa pulang.
                  </p>
                </div>
                <div class="text-right pl-2 shrink-0">
                  <span class="text-xs font-black text-amber-700">
                    Rp {{ formatRupiah(currentService.price * 1) }}
                  </span>
                  <span class="block text-[9px] text-slate-400">1 Tabung</span>
                </div>
              </button>
            </div>
          </div>

          <!-- Self Service Subtotal Box -->
          <div class="p-3 rounded-xl bg-slate-100 flex justify-between items-center text-xs">
            <div>
              <span class="text-slate-500 block">Total Tabung Dihitung:</span>
              <span class="font-bold text-slate-800 text-xs">
                {{ calculatedTubeCount }} Tabung
                <span class="text-slate-400 font-normal"
                  >({{ selfServiceWeight }} Kg ditimbang)</span
                >
              </span>
            </div>
            <div class="text-right">
              <span class="text-slate-500 block">Subtotal:</span>
              <span class="font-black text-base text-brand-700">
                Rp {{ formatRupiah(calculatedTubeCount * (currentService?.price || 0)) }}
              </span>
            </div>
          </div>
        </div>

        <!-- STANDARD FORM: KILOAN / SATUAN / METERAN -->
        <div v-else class="space-y-3.5">
          <!-- Quantity Decimal Input -->
          <div>
            <label class="block text-xs font-bold text-slate-600 mb-1">
              Jumlah / Timbangan ({{ currentService?.unit }}):
            </label>
            <div class="flex items-center gap-2">
              <button
                type="button"
                class="w-10 h-10 rounded-xl bg-slate-100 font-bold text-slate-700 hover:bg-slate-200 cursor-pointer"
                @click="
                  qtyInput = Math.max(
                    0.1,
                    Number(
                      (qtyInput - (currentService?.unit === 'Kg' ? 0.5 : 1)).toFixed(2)
                    )
                  )
                "
              >
                -
              </button>
              <input
                v-model.number="qtyInput"
                type="number"
                :step="currentService?.unit === 'Kg' ? '0.1' : '1'"
                min="0.1"
                class="flex-1 py-2 text-center text-lg font-black rounded-xl border border-slate-200 focus:ring-2 focus:ring-brand-500 focus:outline-hidden"
              />
              <button
                type="button"
                class="w-10 h-10 rounded-xl bg-slate-100 font-bold text-slate-700 hover:bg-slate-200 cursor-pointer"
                @click="
                  qtyInput = Number(
                    (qtyInput + (currentService?.unit === 'Kg' ? 0.5 : 1)).toFixed(2)
                  )
                "
              >
                +
              </button>
            </div>
          </div>

          <!-- Quick Presets -->
          <div class="flex gap-1.5 justify-center flex-wrap">
            <button
              v-for="p in presetOptions"
              :key="p"
              type="button"
              class="px-2.5 py-1 rounded-lg bg-slate-50 border border-slate-200 hover:bg-brand-50 text-xs font-bold text-slate-700 cursor-pointer"
              @click="qtyInput = p"
            >
              {{ p }} {{ currentService?.unit }}
            </button>
          </div>

          <!-- Subtotal Calculation -->
          <div class="p-2.5 rounded-xl bg-slate-50 flex justify-between items-center text-xs">
            <span class="text-slate-500">Subtotal:</span>
            <span class="font-extrabold text-sm text-slate-900">
              Rp {{ formatRupiah(qtyInput * (currentService?.price || 0)) }}
            </span>
          </div>
        </div>

        <button
          type="button"
          class="w-full py-3 rounded-xl bg-brand-600 hover:bg-brand-700 text-white font-extrabold text-xs shadow-md transition-all cursor-pointer"
          @click="addCurrentServiceToCart"
        >
          Tambahkan ke Keranjang
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import {
  PlusCircle,
  User,
  ShoppingBag,
  Plus,
  Trash2,
  Sparkles,
  CreditCard,
  Search,
  X,
  ChevronRight,
  Scale,
  CheckCircle2,
  AlertTriangle,
} from 'lucide-vue-next'
import { api, formatRupiah } from '../../api/client'
import { useToast } from '../../composables/useToast'
import { useAppSettings } from '../../composables/useAppSettings'

defineProps({
  modelValue: Boolean,
})

const emit = defineEmits(['update:modelValue', 'order-created'])
const toast = useToast()
const { settings } = useAppSettings()

const tubeCapacity = computed(() => Number(settings.machine?.tubeCapacityKg) || 8.0)
const tubeTolerance = computed(() => Number(settings.machine?.toleranceKg) || 0.5)
const maxSafeWeight = computed(() => Number((tubeCapacity.value + tubeTolerance.value).toFixed(2)))

const creating = ref(false)
const isNewCustomer = ref(false)
const customerSearchQuery = ref('')
const customerResults = ref([])
const selectedCustomer = ref(null)
const searchingCustomer = ref(false)

const newCustomer = reactive({
  full_name: '',
  phone: '',
  address: '',
})

const services = ref([])
const serviceCategories = ref([])
const perfumesList = ref([])
const discountsList = ref([])
const cartItems = ref([])

// Service Picker Modal State
const showServicePickerModal = ref(false)
const serviceModalSearch = ref('')
const activeServiceTab = ref('ALL')

const getCategoryIcon = (identifier) => {
  const idStr = (identifier || '').toUpperCase()
  if (idStr.includes('SELF') || idStr.includes('TABUNG')) return '⚡'
  if (idStr.includes('SETRIKA') && !idStr.includes('CUCI')) return '👔'
  if (idStr.includes('CUCI_SETRIKA') || idStr.includes('KILOAN')) return '🧺'
  if (idStr.includes('SATUAN')) return '🧥'
  if (idStr.includes('KARPET')) return '🧼'
  if (idStr.includes('BONEKA')) return '🧸'
  if (idStr.includes('SEPATU') || idStr.includes('TAS')) return '👟'
  if (idStr.includes('BED_COVER') || idStr.includes('SELIMUT')) return '🛏️'
  return '🏷️'
}

// Category Tabs for Service Picker Modal
const categoryTabs = computed(() => {
  const list = []
  const activeServices = services.value.filter((s) => s.active !== false)
  list.push({ id: 'ALL', name: 'Semua', count: activeServices.length, icon: '✨' })

  if (serviceCategories.value.length > 0) {
    serviceCategories.value.forEach((cat) => {
      const count = activeServices.filter((s) => {
        return (
          s.category_id === cat.id ||
          s.category === cat.name ||
          s.category === cat.code
        )
      }).length
      list.push({
        id: cat.id,
        code: cat.code,
        name: cat.name,
        count: count,
        icon: getCategoryIcon(cat.code || cat.name),
      })
    })
  } else {
    // Dynamic fallback from services
    const uniqueCats = [...new Set(activeServices.map((s) => s.category).filter(Boolean))]
    uniqueCats.forEach((cname) => {
      const count = activeServices.filter((s) => s.category === cname).length
      list.push({
        id: cname,
        code: cname,
        name: cname,
        count: count,
        icon: getCategoryIcon(cname),
      })
    })
  }
  return list
})

// Quantity Sub-Modal State
const showQtyModal = ref(false)
const currentService = ref(null)
const qtyInput = ref(1.0)

// Self Service specific calculation state
const selfServiceWeight = ref(8.0)
const selfServiceAction = ref('TWO_TUBES') // 'TWO_TUBES' or 'SPLIT_EXCESS'

const orderConfig = reactive({
  piece_count: '',
  perfume_name: 'Akasia',
  discount_type: 'NONE',
  custom_discount: 0,
  rack_number: 'A-01',
  estimated_completion_at: '',
  notes: '',
})

const paymentOption = ref('SEKARANG')
const paymentMethod = ref('TUNAI')
const paidAmount = ref(0)

const isSelfService = computed(() => {
  if (!currentService.value) return false
  const name = (currentService.value.service_name || '').toLowerCase()
  const code = (currentService.value.service_code || '').toLowerCase()
  const cat = (currentService.value.category || '').toLowerCase()
  return (
    currentService.value.unit === 'Tabung' ||
    cat.includes('self service') ||
    cat === 'tabung' ||
    name.includes('self service') ||
    name.includes('cuci basah') ||
    name.includes('cuci kering') ||
    code === 'svc018' ||
    code === 'svc019'
  )
})

const excessWeightKg = computed(() => {
  if (selfServiceWeight.value <= maxSafeWeight.value) return '0.0'
  return (selfServiceWeight.value - tubeCapacity.value).toFixed(1)
})

const calculatedTubeCount = computed(() => {
  if (selfServiceWeight.value <= maxSafeWeight.value) {
    return 1
  }
  // If > maxSafeWeight
  if (selfServiceAction.value === 'TWO_TUBES') {
    return 2
  }
  return 1 // 'SPLIT_EXCESS' -> 1 tabung
})

const presetOptions = computed(() => {
  if (!currentService.value) return [1, 2, 3, 5]
  const unit = currentService.value.unit || ''
  if (unit === 'Kg') {
    return [3, 4, 5, 7, 10]
  }
  if (unit === 'Meter') {
    return [2, 3, 4, 6]
  }
  return [1, 2, 3, 4]
})

const favoriteServices = computed(() => {
  const favCodes = ['SVC018', 'SVC001', 'SVC005', 'SVC019', 'SVC006']
  const matched = services.value.filter((s) => favCodes.includes(s.service_code) && s.active !== false)
  if (matched.length > 0) return matched
  return services.value.filter((s) => s.active !== false).slice(0, 4)
})

const groupedServices = computed(() => {
  const activeServices = services.value.filter((s) => s.active !== false)

  if (serviceCategories.value.length > 0) {
    return serviceCategories.value
      .map((cat) => {
        const items = activeServices.filter((s) => {
          return (
            s.category_id === cat.id ||
            s.category === cat.name ||
            s.category === cat.code
          )
        })
        return {
          categoryId: cat.id,
          categoryCode: cat.code,
          categoryName: cat.name,
          icon: getCategoryIcon(cat.code || cat.name),
          items,
        }
      })
      .filter((grp) => grp.items.length > 0)
  }

  // Fallback grouping by category string
  const groupsMap = new Map()
  activeServices.forEach((s) => {
    const cname = s.category || 'Lainnya'
    if (!groupsMap.has(cname)) {
      groupsMap.set(cname, {
        categoryId: cname,
        categoryCode: cname,
        categoryName: cname,
        icon: getCategoryIcon(cname),
        items: [],
      })
    }
    groupsMap.get(cname).items.push(s)
  })
  return Array.from(groupsMap.values())
})

const filteredServiceGroups = computed(() => {
  const search = serviceModalSearch.value.trim().toLowerCase()
  return groupedServices.value
    .filter((grp) => {
      if (activeServiceTab.value === 'ALL') return true
      return (
        grp.categoryId === activeServiceTab.value ||
        grp.categoryCode === activeServiceTab.value ||
        grp.categoryName === activeServiceTab.value
      )
    })
    .map((grp) => {
      if (!search) return grp
      const items = grp.items.filter(
        (s) =>
          s.service_name.toLowerCase().includes(search) ||
          (s.service_code && s.service_code.toLowerCase().includes(search)) ||
          (s.category && s.category.toLowerCase().includes(search))
      )
      return { ...grp, items }
    })
    .filter((grp) => grp.items.length > 0)
})

const selectServiceFromModal = (srv) => {
  showServicePickerModal.value = false
  openQuantityModal(srv)
}

const subtotalCart = computed(() => {
  return cartItems.value.reduce((acc, it) => acc + (it.subtotal || 0), 0)
})

const discountAmount = computed(() => {
  if (orderConfig.discount_type === 'PROMO_10') {
    return Math.round(subtotalCart.value * 0.1)
  }
  if (orderConfig.discount_type === 'LOYALTY') {
    return Math.min(subtotalCart.value, 20000)
  }
  if (orderConfig.discount_type === 'CUSTOM') {
    return orderConfig.custom_discount || 0
  }
  return 0
})

const grandTotal = computed(() => {
  return Math.max(0, subtotalCart.value - discountAmount.value)
})

const close = () => {
  emit('update:modelValue', false)
  cartItems.value = []
  customerResults.value = []
  customerSearchQuery.value = ''
  selectedCustomer.value = null
  isNewCustomer.value = false
}

const onCustomerSearch = async () => {
  const q = (customerSearchQuery.value || '').trim().toLowerCase()
  if (q.length < 1) {
    customerResults.value = []
    return
  }
  searchingCustomer.value = true
  try {
    const data = await api.get('/customers', { search: q })
    // Double filter client-side for immediate precision
    customerResults.value = (data || []).filter(
      (c) =>
        (c.full_name && c.full_name.toLowerCase().includes(q)) ||
        (c.phone && c.phone.includes(q)) ||
        (c.address && c.address.toLowerCase().includes(q))
    )
  } catch {
    customerResults.value = []
  } finally {
    searchingCustomer.value = false
  }
}

const selectCustomer = (c) => {
  selectedCustomer.value = c
  customerResults.value = []
  customerSearchQuery.value = ''
}

const clearCustomerSearch = () => {
  customerSearchQuery.value = ''
  customerResults.value = []
}

const quickAddNewCustomer = () => {
  isNewCustomer.value = true
  newCustomer.full_name = customerSearchQuery.value
  customerSearchQuery.value = ''
  customerResults.value = []
}

const openQuantityModal = (srv) => {
  currentService.value = srv
  const name = (srv.service_name || '').toLowerCase()
  const code = (srv.service_code || '').toLowerCase()
  const isSS =
    srv.category === 'TABUNG' ||
    srv.unit === 'Tabung' ||
    name.includes('self service') ||
    name.includes('cuci basah') ||
    name.includes('cuci kering') ||
    code === 'svc018' ||
    code === 'svc019'

  if (isSS) {
    selfServiceWeight.value = 8.0
    selfServiceAction.value = 'TWO_TUBES'
  } else if (srv.category === 'KILOAN') {
    qtyInput.value = 4.0
  } else {
    qtyInput.value = 1.0
  }
  showQtyModal.value = true
}

const addCurrentServiceToCart = () => {
  if (!currentService.value) return

  if (isSelfService.value) {
    const tubeQty = calculatedTubeCount.value
    const price = currentService.value.price
    const itemSubtotal = tubeQty * price

    let notes = `Ditimbang ${selfServiceWeight.value} Kg`
    if (selfServiceWeight.value > maxSafeWeight.value) {
      if (selfServiceAction.value === 'TWO_TUBES') {
        notes = `Ditimbang ${selfServiceWeight.value} Kg (Diproses 2 Tabung)`
      } else {
        notes = `Ditimbang ${selfServiceWeight.value} Kg (1 Tabung, Sisihkan ~${excessWeightKg.value} Kg)`
      }
    }

    cartItems.value.push({
      service_id: currentService.value.id,
      service_name: currentService.value.service_name,
      category: currentService.value.category,
      quantity: tubeQty,
      weighed_weight: selfServiceWeight.value,
      unit: currentService.value.unit,
      price_per_unit: price,
      subtotal: itemSubtotal,
      notes: notes,
    })
  } else {
    if (qtyInput.value <= 0) return
    const itemSubtotal = qtyInput.value * currentService.value.price
    cartItems.value.push({
      service_id: currentService.value.id,
      service_name: currentService.value.service_name,
      category: currentService.value.category,
      quantity: qtyInput.value,
      unit: currentService.value.unit,
      price_per_unit: currentService.value.price,
      subtotal: itemSubtotal,
      notes: '',
    })
  }

  // Auto calculate SLA if not manually set
  recalculateSLA()
  showQtyModal.value = false
}

const removeItem = (idx) => {
  cartItems.value.splice(idx, 1)
  recalculateSLA()
}

const recalculateSLA = () => {
  const now = new Date()
  let hasSelfService = cartItems.value.some((it) => it.category === 'TABUNG')

  if (hasSelfService) {
    // 3 hours SLA
    now.setHours(now.getHours() + 3)
  } else {
    // default 3 days
    now.setDate(now.getDate() + 3)
  }

  // Format for datetime-local: YYYY-MM-DDTHH:MM
  const pad = (n) => String(n).padStart(2, '0')
  orderConfig.estimated_completion_at = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())}T${pad(now.getHours())}:${pad(now.getMinutes())}`

  // Set paid amount default if bayar sekarang
  if (paymentOption.value === 'SEKARANG') {
    paidAmount.value = grandTotal.value
  }
}

const applyDiscount = () => {
  if (paymentOption.value === 'SEKARANG') {
    paidAmount.value = grandTotal.value
  }
}

const loadMasterServices = async () => {
  try {
    const [srvData, catData, prfData, discData] = await Promise.all([
      api.get('/master/services').catch(() => []),
      api.get('/master/service-categories').catch(() => []),
      api.get('/master/perfumes').catch(() => []),
      api.get('/master/discounts').catch(() => []),
    ])
    services.value = Array.isArray(srvData) ? srvData : []
    serviceCategories.value = Array.isArray(catData) ? catData : []
    if (Array.isArray(prfData) && prfData.length > 0) {
      perfumesList.value = prfData
    }
    if (Array.isArray(discData) && discData.length > 0) {
      discountsList.value = discData
    }
  } catch (err) {
    console.error('Failed to load master services:', err)
    services.value = []
    serviceCategories.value = []
  }
}

const submitOrder = async () => {
  // Validate customer
  let custID = ''
  let custName = ''
  let custPhone = ''
  let custAddress = ''

  if (isNewCustomer.value) {
    if (!newCustomer.full_name || !newCustomer.phone) {
      toast.warning('Nama pelanggan dan nomor WhatsApp wajib diisi!')
      return
    }
    custName = newCustomer.full_name
    custPhone = newCustomer.phone
    custAddress = newCustomer.address
  } else if (selectedCustomer.value) {
    custID = selectedCustomer.value.id
    custName = selectedCustomer.value.full_name
    custPhone = selectedCustomer.value.phone
    custAddress = selectedCustomer.value.address
  } else {
    toast.warning('Silakan pilih atau tambahkan pelanggan terlebih dahulu!')
    return
  }

  if (cartItems.value.length === 0) {
    toast.warning('Keranjang cucian masih kosong! Silakan pilih layanan.')
    return
  }

  creating.value = true
  try {
    const totalWeight = cartItems.value.reduce((acc, it) => {
      if (it.unit === 'Kg') return acc + it.quantity
      if (it.weighed_weight) return acc + it.weighed_weight
      return acc
    }, 0)

    let estDate = null
    if (orderConfig.estimated_completion_at) {
      estDate = new Date(orderConfig.estimated_completion_at).toISOString()
    }

    let finalNotes = orderConfig.notes ? orderConfig.notes.trim() : ''
    if (orderConfig.piece_count) {
      finalNotes = `[${orderConfig.piece_count} Pcs] ${finalNotes}`.trim()
    }

    const payload = {
      customer_id: custID,
      customer_name: custName,
      customer_phone: custPhone,
      customer_address: custAddress,
      items: cartItems.value,
      perfume_name: orderConfig.perfume_name,
      discount_amount: discountAmount.value,
      total_weight: totalWeight,
      total_price: grandTotal.value,
      paid_amount: paymentOption.value === 'SEKARANG' ? paidAmount.value : 0,
      payment_method_name: paymentOption.value === 'SEKARANG' ? paymentMethod.value : '',
      rack_number: orderConfig.rack_number,
      estimated_completion_at: estDate,
      notes: finalNotes,
      piece_count: Number(orderConfig.piece_count) || null,
      pickup_required: false,
      created_by: 'Kasir Utama',
    }

    const res = await api.post('/pos/orders', payload)
    toast.success('Transaksi kasir walk-in berhasil dibuat!')
    emit('order-created', res)
    close()
  } catch (err) {
    toast.error(err.message || 'Gagal membuat order.')
  } finally {
    creating.value = false
  }
}

onMounted(() => {
  loadMasterServices()
  recalculateSLA()
})
</script>
