<template>
    <div class="flex-1 flex flex-col min-h-0">

        <!-- Header -->
        <div class="flex items-center justify-between mb-4 shrink-0">
            <div>
                <h1 class="text-2xl font-bold tracking-tight">Master Produk</h1>
                <p class="text-sm text-base-content/60 mt-0.5">Kelola data produk, stok, satuan, dan harga</p>
            </div>
            <button class="btn btn-outline btn-sm gap-2" @click="openTrashed">
                <Trash2 class="size-4" />
                Produk Terhapus
            </button>
        </div>

        <!-- Tab Navigation — segmented control style -->
        <div class="flex bg-base-200 rounded-xl p-1 w-fit mb-2 shrink-0 gap-0.5">
            <button :class="[ 'px-5 py-2 rounded-lg text-sm font-medium transition-all duration-150 flex items-center gap-2',
                activeTab === 'tambah'
                    ? 'bg-base-100 shadow-sm text-base-content'
                    : 'text-base-content/50 hover:text-base-content' ]" @click="activeTab = 'tambah'">
                <Plus class="size-4" />
                Tambah Produk
            </button>
            <button :class="[ 'px-5 py-2 rounded-lg text-sm font-medium transition-all duration-150 flex items-center gap-2',
                activeTab === 'list'
                    ? 'bg-base-100 shadow-sm text-base-content'
                    : 'text-base-content/50 hover:text-base-content' ]" @click="activeTab = 'list'">
                <List class="size-4" />
                Daftar Produk
                <span
                    :class="[ 'badge badge-xs p-2 pb-1.5 mb-0.5', activeTab === 'list' ? 'badge-primary' : 'badge-neutral' ]">
                    {{ table.getRowCount() }}
                </span>
            </button>
        </div>

        <!-- Tab: Daftar Produk -->
        <div v-show="activeTab === 'list'" class="flex-1 min-h-0 overflow-hidden">
            <div
                class="bg-base-100 rounded-2xl border border-base-200 shadow-sm h-full flex flex-col overflow-hidden px-4 py-3">
                <AppPagination :table="table" v-model:search="search" class="flex-1 min-h-0">
                    <table class="table table-bordered table-hover">
                        <thead class="sticky top-0 z-10">
                            <tr class="bg-base-200 border-b-2 border-base-300">
                                <th v-for="header in table.getFlatHeaders()" :key="header.id"
                                    :class="[ 'text-sm font-medium cemalecase tracking-widest py-2',
                                        header.column.columnDef.meta?.headerClass,
                                        header.column.getCanSort() ? 'cursor-pointer select-none hover:text-primary hover:bg-base-300/50 transition-colors' : '' ]"
                                    @click="header.column.getToggleSortingHandler()?.($event)">
                                    <div :class="['flex items-center gap-1', header.column.columnDef.meta?.headerClass?.includes('text-center') ? 'justify-center' : '']">
                                        <FlexRender :render="header.column.columnDef.header" :props="header.getContext()" />
                                        <span v-if="header.column.getIsSorted() === 'asc'" class="text-primary">↑</span>
                                        <span v-else-if="header.column.getIsSorted() === 'desc'" class="text-primary">↓</span>
                                    </div>
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="loading">
                                <td colspan="9" class="py-24 text-center">
                                    <div class="flex flex-col items-center gap-3">
                                        <span class="loading loading-spinner loading-md text-primary"></span>
                                        <span class="text-sm text-base-content/50 font-medium">Memuat data...</span>
                                    </div>
                                </td>
                            </tr>
                            <tr v-else-if="table.getRowModel().rows.length === 0">
                                <td colspan="9" class="py-24 text-center">
                                    <div class="flex flex-col items-center gap-3">
                                        <div class="size-14 rounded-2xl bg-base-200 flex items-center justify-center">
                                            <Package class="size-7 text-base-content/30" />
                                        </div>
                                        <div>
                                            <p class="font-semibold text-base-content/60">Belum ada data produk</p>
                                            <p class="text-sm text-base-content/40 mt-0.5">
                                                Klik tab
                                                <button class="text-primary font-semibold hover:underline"
                                                    @click="activeTab = 'tambah'">Tambah Produk</button>
                                                untuk menambahkan produk baru
                                            </p>
                                        </div>
                                    </div>
                                </td>
                            </tr>
                            <tr v-else v-for="row in table.getRowModel().rows" :key="row.id"
                                class="border-b border-base-200 hover:bg-primary/5 transition-colors duration-100">
                                <td v-for="cell in row.getVisibleCells()" :key="cell.id"
                                    :class="[ 'py-1', cell.column.columnDef.meta?.cellClass ]">
                                    <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </AppPagination>
            </div>
        </div>

        <!-- Tab: Tambah Produk -->
        <div v-show="activeTab === 'tambah'" class="flex-1 overflow-y-auto pb-6">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">

                <!-- Kiri: Informasi Produk -->
                <div class="bg-base-100 rounded-2xl border border-base-200 shadow-sm overflow-hidden">
                    <div class="px-5 py-3.5 border-b border-base-200 flex items-center gap-3">
                        <div class="w-1 h-5 bg-primary rounded-full"></div>
                        <h3 class="font-semibold text-base-content">Informasi Produk</h3>
                    </div>
                    <div class="p-5 grid grid-cols-2 gap-4">
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Kode Produk</label>
                            <div class="flex gap-2">
                                <input v-model="form.kode_produk" type="text" placeholder="Contoh: A000001"
                                    class="input input-bordered flex-1" />
                                <button class="btn btn-outline gap-1.5" @click="generateKode(form)">
                                    <RefreshCw class="size-4" />
                                    Generate
                                </button>
                            </div>
                            <p class="text-xs text-base-content/40 mt-1">Kosongkan untuk auto-generate saat simpan</p>
                        </div>
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">
                                Nama Produk <span class="text-error">*</span>
                            </label>
                            <input v-model="form.nama_produk" type="text" placeholder="Masukkan nama produk"
                                class="input input-bordered w-full" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Kategori</label>
                            <Multiselect v-model="form.kategori" :options="fetchKategoriOptions" :object="true"
                                value-prop="id" label="nama_kategori" :searchable="true" :can-clear="true"
                                :can-deselect="true" :min-chars="0" :delay="300" placeholder="Pilih kategori..."
                                no-options-text="Belum ada kategori" no-results-text="Tidak ditemukan" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Tipe
                                Kepemilikan</label>
                            <select v-model="form.tipe_kepemilikan" class="select select-bordered w-full">
                                <option value="MILIK_SENDIRI">Milik Sendiri</option>
                                <option value="TITIPAN">Titipan / Konsinyasi</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Pemasok</label>
                            <Multiselect v-model="form.pemasok" :options="fetchPemasokOptions" :object="true"
                                value-prop="id" label="nama_pemasok" :searchable="true" :can-clear="true"
                                :can-deselect="true" :min-chars="0" :delay="300" placeholder="Pilih pemasok..."
                                no-options-text="Belum ada pemasok" no-results-text="Tidak ditemukan" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">
                                Kode Barcode
                                <span class="text-xs font-normal text-base-content/40 ml-1">(opsional)</span>
                            </label>
                            <input v-model="form.kode_barcode" type="text" placeholder="Scan atau ketik manual"
                                class="input input-bordered w-full" />
                        </div>
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">
                                Deskripsi
                                <span class="text-xs font-normal text-base-content/40 ml-1">(opsional)</span>
                            </label>
                            <textarea v-model="form.deskripsi" rows="2"
                                placeholder="Catatan atau keterangan tambahan..."
                                class="textarea textarea-bordered w-full resize-none"></textarea>
                        </div>
                    </div>
                </div>

                <!-- Kanan: Stok, Satuan & Harga -->
                <div class="bg-base-100 rounded-2xl border border-base-200 shadow-sm overflow-hidden">
                    <div class="px-5 py-3.5 border-b border-base-200 flex items-center gap-3">
                        <div class="w-1 h-5 bg-success rounded-full"></div>
                        <h3 class="font-semibold text-base-content">Stok, Satuan & Harga</h3>
                    </div>
                    <div class="p-5 grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Nama Satuan</label>
                            <input v-model="form.nama_satuan" type="text" placeholder="PCS / KG / Lusin"
                                class="input input-bordered w-full" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Jenis Input</label>
                            <select v-model="form.jenis_input" class="select select-bordered w-full">
                                <option value="PCS">PCS (satuan)</option>
                                <option value="BERAT">Berat (gram/kg)</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Nilai Konversi</label>
                            <input v-model="form.nilai_konversi" type="number" step="0.0001" min="0.0001"
                                class="input input-bordered w-full" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Stok Awal</label>
                            <input v-model="form.stok_awal" type="number" step="0.0001" min="0"
                                class="input input-bordered w-full" />
                        </div>

                        <div class="col-span-2">
                            <div class="flex items-center gap-3 py-1">
                                <div class="h-px flex-1 bg-base-200"></div>
                                <span
                                    class="text-xs font-semibold uppercase tracking-wider text-base-content/40">Harga</span>
                                <div class="h-px flex-1 bg-base-200"></div>
                            </div>
                        </div>

                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Harga Beli</label>
                            <AppCurrencyInput v-model="form.harga_beli"
                                :disabled="form.tipe_kepemilikan === 'TITIPAN'" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">
                                Harga Jual <span class="text-error">*</span>
                            </label>
                            <AppCurrencyInput v-model="form.harga_jual" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Harga Titip</label>
                            <AppCurrencyInput v-model="form.harga_titip"
                                :disabled="form.tipe_kepemilikan === 'MILIK_SENDIRI'" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Status</label>
                            <select v-model="form.status_aktif" class="select select-bordered w-full">
                                <option :value="true">Aktif</option>
                                <option :value="false">Nonaktif</option>
                            </select>
                        </div>
                    </div>
                </div>

            </div>

            <div class="mt-5 flex items-center gap-3">
                <button class="btn btn-primary px-6 gap-2" :disabled="saving" @click="simpanProduk">
                    <span v-if="saving" class="loading loading-spinner loading-sm"></span>
                    <Check v-else class="size-4" />
                    Simpan Produk
                </button>
                <button class="btn btn-ghost" @click="resetForm">Batal / Reset</button>
            </div>
        </div>

    </div>

    <!-- Modal: Edit Produk -->
    <dialog ref="modalEdit" class="modal">
        <div class="modal-box max-w-4xl w-full p-0 overflow-hidden">
            <div class="px-6 py-4 border-b border-base-200 flex items-center justify-between">
                <h3 class="text-lg font-bold">Edit Produk</h3>
                <button class="btn btn-ghost btn-sm btn-circle" @click="modalEdit?.close()">✕</button>
            </div>
            <div class="p-6 grid grid-cols-1 lg:grid-cols-2 gap-5 overflow-y-auto max-h-[70vh]">

                <div class="bg-base-100 rounded-xl border border-base-200 overflow-hidden">
                    <div class="px-4 py-3 border-b border-base-200 flex items-center gap-3">
                        <div class="w-1 h-5 bg-primary rounded-full"></div>
                        <h4 class="font-semibold">Informasi Produk</h4>
                    </div>
                    <div class="p-4 grid grid-cols-2 gap-3">
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Kode Produk</label>
                            <div class="flex gap-2">
                                <input v-model="editForm.kode_produk" type="text" class="input input-bordered flex-1" />
                                <button class="btn btn-outline" @click="generateKode(editForm)">Generate</button>
                            </div>
                        </div>
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Nama Produk <span
                                    class="text-error">*</span></label>
                            <input v-model="editForm.nama_produk" type="text" class="input input-bordered w-full" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Kategori</label>
                            <Multiselect v-model="editForm.kategori" :options="fetchKategoriOptions" :object="true"
                                value-prop="id" label="nama_kategori" :searchable="true" :can-clear="true"
                                :can-deselect="true" :min-chars="0" :delay="300" placeholder="Pilih kategori..."
                                no-options-text="Belum ada" no-results-text="Tidak ditemukan" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Tipe
                                Kepemilikan</label>
                            <select v-model="editForm.tipe_kepemilikan" class="select select-bordered w-full">
                                <option value="MILIK_SENDIRI">Milik Sendiri</option>
                                <option value="TITIPAN">Titipan / Konsinyasi</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Pemasok</label>
                            <Multiselect v-model="editForm.pemasok" :options="fetchPemasokOptions" :object="true"
                                value-prop="id" label="nama_pemasok" :searchable="true" :can-clear="true"
                                :can-deselect="true" :min-chars="0" :delay="300" placeholder="Pilih pemasok..."
                                no-options-text="Belum ada" no-results-text="Tidak ditemukan" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Status</label>
                            <select v-model="editForm.status_aktif" class="select select-bordered w-full">
                                <option :value="true">Aktif</option>
                                <option :value="false">Nonaktif</option>
                            </select>
                        </div>
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Kode Barcode</label>
                            <input v-model="editForm.kode_barcode" type="text" class="input input-bordered w-full" />
                        </div>
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Deskripsi</label>
                            <textarea v-model="editForm.deskripsi" rows="2"
                                class="textarea textarea-bordered w-full resize-none"></textarea>
                        </div>
                    </div>
                </div>

                <div class="bg-base-100 rounded-xl border border-base-200 overflow-hidden">
                    <div class="px-4 py-3 border-b border-base-200 flex items-center gap-3">
                        <div class="w-1 h-5 bg-success rounded-full"></div>
                        <h4 class="font-semibold">Stok, Satuan & Harga</h4>
                    </div>
                    <div class="p-4 grid grid-cols-2 gap-3">
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Nama Satuan</label>
                            <input v-model="editForm.nama_satuan" type="text" class="input input-bordered w-full" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Jenis Input</label>
                            <select v-model="editForm.jenis_input" class="select select-bordered w-full">
                                <option value="PCS">PCS (satuan)</option>
                                <option value="BERAT">Berat (gram/kg)</option>
                            </select>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Nilai Konversi</label>
                            <input v-model="editForm.nilai_konversi" type="number" step="0.0001" min="0.0001"
                                class="input input-bordered w-full" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Jumlah Stok</label>
                            <input v-model="editForm.jumlah_stok" type="number" step="0.0001" min="0"
                                class="input input-bordered w-full" />
                        </div>
                        <div class="col-span-2">
                            <div class="flex items-center gap-3 py-1">
                                <div class="h-px flex-1 bg-base-200"></div>
                                <span
                                    class="text-xs font-semibold uppercase tracking-wider text-base-content/40">Harga</span>
                                <div class="h-px flex-1 bg-base-200"></div>
                            </div>
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Harga Beli</label>
                            <AppCurrencyInput v-model="editForm.harga_beli"
                                :disabled="editForm.tipe_kepemilikan === 'TITIPAN'" />
                        </div>
                        <div>
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Harga Jual <span
                                    class="text-error">*</span></label>
                            <AppCurrencyInput v-model="editForm.harga_jual" />
                        </div>
                        <div class="col-span-2">
                            <label class="block text-sm font-medium text-base-content/80 mb-1.5">Harga Titip</label>
                            <AppCurrencyInput v-model="editForm.harga_titip"
                                :disabled="editForm.tipe_kepemilikan === 'MILIK_SENDIRI'" />
                        </div>
                    </div>
                </div>

            </div>
            <div class="px-6 py-4 border-t border-base-200 flex justify-end gap-3">
                <button class="btn btn-ghost" @click="modalEdit?.close()">Batal</button>
                <button class="btn btn-primary gap-2" :disabled="updating" @click="updateProduk">
                    <span v-if="updating" class="loading loading-spinner loading-sm"></span>
                    <Check v-else class="size-4" />
                    Simpan Perubahan
                </button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop"><button>close</button></form>
    </dialog>

    <!-- Modal: Multi-Satuan -->
    <dialog ref="modalSatuan" class="modal">
        <div class="modal-box max-w-3xl w-full">
            <h3 class="font-bold text-base mb-1">
                Multi-Satuan: <span class="text-primary">{{ satuanProdukNama }}</span>
            </h3>
            <p class="text-sm text-base-content/50 mb-4">
                Nilai konversi terhadap satuan utama. Contoh: satuan utama PCS → BOX = 12 (1 BOX = 12 PCS).
            </p>

            <div class="border border-base-300 rounded-xl overflow-hidden mb-4">
                <div class="px-3 py-2 bg-green-50 border-b border-base-300">
                    <span class="text-sm font-semibold text-green-700">Tambah Satuan</span>
                </div>
                <div class="p-3 grid grid-cols-2 sm:grid-cols-4 gap-2">
                    <div>
                        <label class="label-text text-sm font-medium mb-1 block">Nama Satuan</label>
                        <input v-model="satuanForm.nama_satuan" type="text" placeholder="BOX / PCS / KG"
                            class="input input-bordered w-full" />
                    </div>
                    <div>
                        <label class="label-text text-sm font-medium mb-1 block">Jenis</label>
                        <select v-model="satuanForm.jenis_input" class="select select-bordered w-full">
                            <option value="PCS">PCS</option>
                            <option value="BERAT">BERAT</option>
                        </select>
                    </div>
                    <div>
                        <label class="label-text text-sm font-medium mb-1 block">Konversi</label>
                        <input v-model="satuanForm.nilai_konversi" type="number" step="0.0001" min="0.0001"
                            class="input input-bordered w-full" />
                    </div>
                    <div>
                        <label class="label-text text-sm font-medium mb-1 block">Keterangan</label>
                        <select v-model="satuanForm.satuan_utama" class="select select-bordered w-full">
                            <option :value="false">Bukan Utama</option>
                            <option :value="true">Satuan Utama</option>
                        </select>
                    </div>
                    <div>
                        <label class="label-text text-sm font-medium mb-1 block">Harga Beli</label>
                        <AppCurrencyInput v-model="satuanForm.harga_beli" />
                    </div>
                    <div>
                        <label class="label-text text-sm font-medium mb-1 block">Harga Jual <span
                                class="text-error">*</span></label>
                        <AppCurrencyInput v-model="satuanForm.harga_jual" />
                    </div>
                    <div>
                        <label class="label-text text-sm font-medium mb-1 block">Harga Titip</label>
                        <AppCurrencyInput v-model="satuanForm.harga_titip" />
                    </div>
                    <div>
                        <label class="label-text text-sm font-medium mb-1 block">Barcode <span
                                class="text-base-content/40">(opsional)</span></label>
                        <input v-model="satuanForm.kode_barcode" type="text" class="input input-bordered w-full" />
                    </div>
                    <div class="col-span-2 sm:col-span-4 flex justify-end mt-1">
                        <button class="btn btn-success btn-sm gap-2" :disabled="satuanSaving" @click="tambahSatuan">
                            <span v-if="satuanSaving" class="loading loading-spinner loading-xs"></span>
                            Tambah Satuan
                        </button>
                    </div>
                </div>
            </div>

            <div class="overflow-x-auto">
                <div v-if="satuanLoading" class="text-center py-6 text-base-content/40">
                    <span class="loading loading-spinner loading-sm"></span>
                </div>
                <table v-else class="table table-sm">
                    <thead>
                        <tr>
                            <th>Satuan</th>
                            <th>Jenis</th>
                            <th>Konversi</th>
                            <th>Harga Jual</th>
                            <th>Barcode</th>
                            <th>Utama</th>
                            <th class="text-center">Hapus</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="satuanList.length === 0">
                            <td colspan="7" class="text-center py-4 text-base-content/40">Belum ada satuan</td>
                        </tr>
                        <tr v-for="s in satuanList" :key="s.id" class="hover:bg-base-200">
                            <td class="font-semibold">{{ s.nama_satuan }}</td>
                            <td>{{ s.jenis_input }}</td>
                            <td class="font-mono">{{ s.nilai_konversi }}</td>
                            <td>Rp {{ formatRp(s.harga_jual) }}</td>
                            <td class="font-mono text-base-content/50">{{ s.kode_barcode || '—' }}</td>
                            <td>
                                <span v-if="s.satuan_utama" class="badge badge-sm badge-info">Utama</span>
                                <span v-else class="text-base-content/30">—</span>
                            </td>
                            <td class="text-center">
                                <button class="btn btn-ghost btn-xs text-error"
                                    @click="hapusSatuan(s.id)">Hapus</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <div class="modal-action mt-3">
                <button class="btn btn-ghost btn-sm" @click="modalSatuan?.close()">Tutup</button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop"><button>close</button></form>
    </dialog>

    <!-- Modal: Produk Terhapus -->
    <dialog ref="modalTrashed" class="modal">
        <div class="modal-box max-w-4xl w-full">
            <h3 class="font-bold text-base mb-4 text-error">Produk Terhapus</h3>
            <div v-if="trashedLoading" class="text-center py-8 text-base-content/40">
                <span class="loading loading-spinner"></span>
            </div>
            <div v-else class="overflow-x-auto">
                <table class="table table-sm">
                    <thead>
                        <tr>
                            <th>Nama Produk</th>
                            <th>Kategori</th>
                            <th>Tipe</th>
                            <th>Satuan</th>
                            <th>Barcode</th>
                            <th>Dihapus</th>
                            <th class="text-center">Aksi</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="trashedList.length === 0">
                            <td colspan="7" class="text-center py-6 text-base-content/40">Tidak ada produk terhapus</td>
                        </tr>
                        <tr v-for="p in trashedList" :key="p.id" class="hover:bg-base-200">
                            <td class="font-medium">{{ p.nama_produk }}</td>
                            <td>{{ p.nama_kategori || '—' }}</td>
                            <td>
                                <span
                                    :class="p.tipe_kepemilikan === 'TITIPAN' ? 'badge badge-sm badge-warning' : 'badge badge-sm badge-info'">
                                    {{ p.tipe_kepemilikan === 'TITIPAN' ? 'Titipan' : 'Milik' }}
                                </span>
                            </td>
                            <td>{{ p.nama_satuan || 'PCS' }}</td>
                            <td class="font-mono text-base-content/50">{{ p.kode_barcode || '—' }}</td>
                            <td class="text-sm text-base-content/50">
                                {{ p.deleted_at ? new Date(p.deleted_at).toLocaleDateString('id-ID') : '—' }}
                            </td>
                            <td class="text-center">
                                <button class="btn btn-ghost btn-xs text-success"
                                    @click="restoreProduk(p.id)">Pulihkan</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="modal-action">
                <button class="btn btn-ghost btn-sm" @click="modalTrashed?.close()">Tutup</button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop"><button>close</button></form>
    </dialog>

</template>
<script setup>
import { ref, reactive, h } from 'vue'
import { FlexRender } from '@tanstack/vue-table'
import Multiselect from '@vueform/multiselect'
import { Trash2, Plus, List, Package, RefreshCw, Check, Edit2, Layers } from 'lucide-vue-next'
import { useServerTable } from '../../composables/useServerTable.js'
import { useToast } from '../../composables/useToast.js'
import AppPagination from '../../components/AppPagination.vue'
import AppBarcode from '../../components/AppBarcode.vue'
import AppCurrencyInput from '../../components/AppCurrencyInput.vue'
import ProdukApi from '../../services/ProdukApi.js'
import KategoriApi from '../../services/KategoriApi.js'
import PemasokApi from '../../services/PemasokApi.js'

const { showToast } = useToast()

// ── Async options untuk Multiselect ───────────────────────
async function fetchKategoriOptions(query) {
    try {
        const res = await KategoriApi.getAll({ search: query, pageSize: 20 })
        return res.data.data ?? []
    } catch { return [] }
}

async function fetchPemasokOptions(query) {
    try {
        const res = await PemasokApi.getAll({ search: query, pageSize: 20 })
        return res.data.data ?? []
    } catch { return [] }
}

// ── Tab ────────────────────────────────────────────────────
const activeTab = ref('list')

// ── Form Tambah ────────────────────────────────────────────
const saving = ref(false)

const form = reactive({
    kode_produk: '', nama_produk: '', kategori: null, tipe_kepemilikan: 'MILIK_SENDIRI',
    pemasok: null, kode_barcode: '', deskripsi: '', nama_satuan: 'PCS',
    jenis_input: 'PCS', nilai_konversi: 1, stok_awal: 0,
    harga_beli: '', harga_jual: '', harga_titip: '', status_aktif: true,
})

function resetForm() {
    Object.assign(form, {
        kode_produk: '', nama_produk: '', kategori: null, tipe_kepemilikan: 'MILIK_SENDIRI',
        pemasok: null, kode_barcode: '', deskripsi: '', nama_satuan: 'PCS',
        jenis_input: 'PCS', nilai_konversi: 1, stok_awal: 0,
        harga_beli: '', harga_jual: '', harga_titip: '', status_aktif: true,
    })
}

async function generateKode(target) {
    const prefix = /[A-Za-z]/.test(target.kode_produk?.charAt(0))
        ? target.kode_produk.charAt(0).toUpperCase()
        : 'A'
    const res = await ProdukApi.nextKode(prefix)
    if (res.success) target.kode_produk = res.data.kode_produk
}

async function simpanProduk() {
    if (!form.nama_produk || form.harga_jual === '') {
        showToast('Nama produk dan harga jual wajib diisi.', 'error'); return
    }
    if (!form.kode_produk.trim()) await generateKode(form)
    saving.value = true
    try {
        const { kategori, pemasok, ...rest } = form
        await ProdukApi.create({ ...rest, kategori_id: kategori?.id ?? null, pemasok_id: pemasok?.id ?? null })
        showToast('Produk berhasil disimpan.')
        resetForm()
        fetchData()
        activeTab.value = 'list'
    } catch (e) {
        showToast(e.response?.data?.message || 'Gagal simpan produk.', 'error')
    } finally { saving.value = false }
}

// ── Tabel ──────────────────────────────────────────────────
const columns = [
    {
        id: 'index', header: '#', enableSorting: false,
        meta: { headerClass: 'w-10 text-center', cellClass: 'text-center text-sm text-base-content/40 tabular-nums' },
        cell: info => {
            const { pageIndex, pageSize } = info.table.getState().pagination
            return pageIndex * pageSize + info.row.index + 1
        },
    },
    {
        accessorKey: 'kode_barcode', header: 'Barcode', enableSorting: false,
        cell: info => info.getValue()
            ? h('div', { class: 'flex flex-col items-start gap-0.5' }, [
                h(AppBarcode, { value: info.getValue(), width: 1.2, height: 18 }),
            ])
            : h('span', { class: 'text-base-content/30' }, '—'),
    },
    {
        accessorKey: 'nama_produk', header: 'Nama Produk',
        meta: { cellClass: 'font-semibold text-base-content' },
    },
    {
        accessorKey: 'nama_kategori', header: 'Kategori',
        meta: { cellClass: 'text-base-content/70' },
        cell: info => info.getValue()
            ? h('span', { class: 'bg-base-200/60 text-base-content/80 text-xs px-2.5 py-1 rounded-full font-medium' }, info.getValue())
            : h('span', { class: 'text-base-content/30' }, '—'),
    },
    {
        accessorKey: 'tipe_kepemilikan', header: 'Tipe', enableSorting: false,
        cell: info => h('span', {
            class: info.getValue() === 'TITIPAN'
                ? 'badge badge-warning badge-soft badge-sm font-medium'
                : 'badge badge-primary badge-soft badge-sm font-medium',
        }, info.getValue() === 'TITIPAN' ? 'Titipan' : 'Milik'),
    },
    {
        accessorKey: 'nama_satuan', header: 'Satuan', enableSorting: false,
        meta: { cellClass: 'text-sm text-base-content/60 font-medium' },
        cell: info => info.getValue() || 'PCS',
    },
    {
        accessorKey: 'jumlah_stok', header: 'Stok',
        meta: { headerClass: 'text-center', cellClass: 'text-center font-bold tabular-nums' },
        cell: info => Number(info.getValue() || 0).toLocaleString('id-ID'),
    },
    {
        accessorKey: 'harga_jual', header: 'Harga Jual',
        meta: { headerClass: 'text-right', cellClass: 'text-right tabular-nums font-semibold text-base-content' },
        cell: info => h('span', {}, `Rp ${Number(info.getValue() || 0).toLocaleString('id-ID')}`),
    },
    {
        accessorKey: 'status_aktif', header: 'Status', enableSorting: false,
        meta: { headerClass: 'text-center', cellClass: 'text-center' },
        cell: info => h('span', {
            class: info.getValue()
                ? 'badge badge-success badge-soft badge-sm font-semibold'
                : 'badge badge-error badge-soft badge-sm font-semibold',
        }, info.getValue() ? 'Aktif' : 'Nonaktif'),
    },
    {
        id: 'aksi', header: 'ACT', enableSorting: false,
        meta: { headerClass: 'w-40 text-center', cellClass: 'text-right' },
        cell: info => h('div', { class: 'flex gap-0.5 justify-end' }, [
            h('button', { class: 'btn btn-sm px-3 font-medium gap-1.5', onClick: () => openEdit(info.row.original) }, [h(Edit2, { class: 'size-3.5' })]),
            h('button', { class: 'btn btn-sm px-3 font-medium text-info gap-1.5', onClick: () => openSatuan(info.row.original) }, [h(Layers, { class: 'size-3.5' })]),
            h('button', { class: 'btn btn-sm px-3 font-medium text-error gap-1.5', onClick: () => hapusProduk(info.row.original.id) }, [h(Trash2, { class: 'size-3.5' })]),
        ]),
    },
]

const { table, loading, search, fetchData } = useServerTable({
    columns,
    fetchFn: params => ProdukApi.getAll(params),
    pageSize: 20,
})

async function hapusProduk(id) {
    if (!confirm('Hapus produk ini? Bisa dipulihkan dari Produk Terhapus.')) return
    try {
        await ProdukApi.remove(id)
        showToast('Produk berhasil dihapus.')
        fetchData()
    } catch { showToast('Gagal hapus produk.', 'error') }
}

// ── Modal Edit ─────────────────────────────────────────────
const modalEdit = ref(null)
const editId = ref(null)
const updating = ref(false)

const editForm = reactive({
    kode_produk: '', nama_produk: '', kategori: null, tipe_kepemilikan: 'MILIK_SENDIRI',
    pemasok: null, kode_barcode: '', deskripsi: '', nama_satuan: 'PCS',
    jenis_input: 'PCS', nilai_konversi: 1, jumlah_stok: 0,
    harga_beli: '', harga_jual: '', harga_titip: '', status_aktif: true,
})

function openEdit(row) {
    editId.value = row.id
    Object.assign(editForm, {
        kode_produk: row.kode_produk ?? '', nama_produk: row.nama_produk ?? '',
        kategori: row.kategori_id ? { id: row.kategori_id, nama_kategori: row.nama_kategori ?? '' } : null,
        tipe_kepemilikan: row.tipe_kepemilikan ?? 'MILIK_SENDIRI',
        pemasok: row.pemasok_id ? { id: row.pemasok_id, nama_pemasok: row.nama_pemasok ?? '' } : null,
        kode_barcode: row.kode_barcode ?? '', deskripsi: row.deskripsi ?? '',
        nama_satuan: row.nama_satuan ?? 'PCS', jenis_input: row.jenis_input ?? 'PCS',
        nilai_konversi: row.nilai_konversi ?? 1, jumlah_stok: row.jumlah_stok ?? 0,
        harga_beli: row.harga_beli ?? '', harga_jual: row.harga_jual ?? '',
        harga_titip: row.harga_titip ?? '', status_aktif: !!row.status_aktif,
    })
    modalEdit.value?.showModal()
}

async function updateProduk() {
    if (!editForm.nama_produk || editForm.harga_jual === '') {
        showToast('Nama produk dan harga jual wajib diisi.', 'error'); return
    }
    if (!editForm.kode_produk.trim()) await generateKode(editForm)
    updating.value = true
    try {
        const { kategori, pemasok, ...rest } = editForm
        await ProdukApi.update(editId.value, { ...rest, kategori_id: kategori?.id ?? null, pemasok_id: pemasok?.id ?? null })
        showToast('Produk berhasil diupdate.')
        modalEdit.value?.close()
        fetchData()
    } catch (e) {
        showToast(e.response?.data?.message || 'Gagal update produk.', 'error')
    } finally { updating.value = false }
}

// ── Modal Multi-Satuan ─────────────────────────────────────
const modalSatuan = ref(null)
const satuanProdukId = ref(null)
const satuanProdukNama = ref('')
const satuanList = ref([])
const satuanLoading = ref(false)
const satuanSaving = ref(false)

const satuanForm = reactive({
    nama_satuan: '', jenis_input: 'PCS', nilai_konversi: 1, satuan_utama: false,
    kode_barcode: '', harga_beli: '', harga_jual: '', harga_titip: '',
})

function resetSatuanForm() {
    Object.assign(satuanForm, { nama_satuan: '', jenis_input: 'PCS', nilai_konversi: 1, satuan_utama: false, kode_barcode: '', harga_beli: '', harga_jual: '', harga_titip: '' })
}

async function loadSatuan(id) {
    satuanLoading.value = true
    try {
        const res = await ProdukApi.getSatuan(id)
        satuanList.value = res.data ?? []
    } catch { satuanList.value = [] } finally { satuanLoading.value = false }
}

function openSatuan(row) {
    satuanProdukId.value = row.id
    satuanProdukNama.value = row.nama_produk
    resetSatuanForm()
    loadSatuan(row.id)
    modalSatuan.value?.showModal()
}

async function tambahSatuan() {
    if (!satuanForm.nama_satuan || satuanForm.harga_jual === '') {
        showToast('Nama satuan dan harga jual wajib diisi.', 'error'); return
    }
    satuanSaving.value = true
    try {
        const res = await ProdukApi.storeSatuan(satuanProdukId.value, satuanForm)
        if (res.success) { showToast('Satuan ditambahkan.'); resetSatuanForm(); loadSatuan(satuanProdukId.value); fetchData() }
    } catch (e) {
        showToast(e.response?.data?.message || 'Gagal tambah satuan.', 'error')
    } finally { satuanSaving.value = false }
}

async function hapusSatuan(satuanId) {
    if (!confirm('Hapus satuan ini?')) return
    try {
        await ProdukApi.removeSatuan(satuanProdukId.value, satuanId)
        showToast('Satuan dihapus.')
        loadSatuan(satuanProdukId.value)
        fetchData()
    } catch (e) { showToast(e.response?.data?.message || 'Gagal hapus satuan.', 'error') }
}

// ── Modal Terhapus ─────────────────────────────────────────
const modalTrashed = ref(null)
const trashedList = ref([])
const trashedLoading = ref(false)

async function loadTrashed() {
    trashedLoading.value = true
    try {
        const res = await ProdukApi.getTrashed({ page: 1, pageSize: 100 })
        trashedList.value = res.data?.data ?? []
    } catch { trashedList.value = [] } finally { trashedLoading.value = false }
}

function openTrashed() { loadTrashed(); modalTrashed.value?.showModal() }

async function restoreProduk(id) {
    if (!confirm('Pulihkan produk ini?')) return
    try {
        await ProdukApi.restore(id)
        showToast('Produk berhasil dipulihkan.')
        loadTrashed()
        fetchData()
    } catch { showToast('Gagal pulihkan produk.', 'error') }
}

const formatRp = n => Number(n || 0).toLocaleString('id-ID')
</script>