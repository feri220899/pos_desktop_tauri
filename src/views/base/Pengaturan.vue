<template>
    <div class="flex-1 flex flex-col min-h-0">
        <div ref="splitContainerRef" class="flex flex-row flex-1 min-h-0"
             :class="isResizing ? 'select-none cursor-col-resize' : ''">

            <!-- Card Kiri: Manajemen Pengguna -->
            <div class="card bg-base-100 border border-base-300 shadow-sm overflow-hidden min-w-0"
                 :style="{ width: leftWidth + '%' }">
                <div class="card-body overflow-hidden flex flex-col">

                    <div class="flex items-center justify-between">
                        <div>
                            <h2 class="card-title text-base">Pengguna</h2>
                            <p class="text-base-content/50 mt-0.5">Kelola akun yang dapat mengakses aplikasi</p>
                        </div>
                        <div class="flex gap-2">
                            <button class="btn btn-ghost btn-sm gap-2" @click="openPermModal">
                                <ShieldCheck class="size-4" />
                                Hak Akses
                            </button>
                            <button class="btn btn-primary btn-sm gap-2" @click="openCreate">
                                <UserPlus class="size-4" />
                                Tambah
                            </button>
                        </div>
                    </div>

                    <!-- Tabel User dengan toolbar & pagination -->
                    <AppPagination :table="table" v-model:search="search" class="mt-4 flex-1 min-h-0">
                        <table class="table">
                            <thead class="sticky top-0 z-10 bg-base-100">
                                <tr>
                                    <th v-for="header in table.getFlatHeaders()" :key="header.id"
                                        @click="header.column.getToggleSortingHandler()?.($event)"
                                        :style="header.column.getCanSort() ? 'cursor:pointer;user-select:none' : ''"
                                    >
                                        <div class="flex items-center gap-1" :class="header.column.columnDef.meta?.headerClass">
                                            <FlexRender :render="header.column.columnDef.header" :props="header.getContext()" />
                                            <span v-if="header.column.getIsSorted() === 'asc'">↑</span>
                                            <span v-else-if="header.column.getIsSorted() === 'desc'">↓</span>
                                        </div>
                                    </th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-if="loading">
                                    <td colspan="4" class="text-center py-8 text-base-content/40">
                                        <span class="loading loading-spinner loading-sm"></span>
                                    </td>
                                </tr>
                                <tr v-else-if="table.getRowModel().rows.length === 0">
                                    <td colspan="4" class="text-center py-8 text-base-content/40">Belum ada pengguna</td>
                                </tr>
                                <tr v-else v-for="row in table.getRowModel().rows" :key="row.id" class="hover:bg-base-200">
                                    <td v-for="cell in row.getVisibleCells()" :key="cell.id"
                                        :class="cell.column.columnDef.meta?.cellClass">
                                        <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </AppPagination>

                </div>
            </div>

            <!-- Divider -->
            <div class="w-2.5 flex items-center justify-center cursor-col-resize group"
                 @mousedown.prevent="startResize">
                <div class="w-1 h-12 rounded-full bg-base-300 group-hover:bg-primary/60 transition-colors duration-150 flex flex-col items-center justify-center gap-1">
                    <span v-for="i in 4" :key="i" class="block size-0.5 rounded-full bg-base-content/25 group-hover:bg-white/70 transition-colors"></span>
                </div>
            </div>

            <!-- Card Kanan: Tabbed -->
            <div class="card bg-base-100 border border-base-300 shadow-sm overflow-hidden flex-1 min-w-0">
                <div class="card-body flex flex-col gap-0 p-0 overflow-hidden">

                    <!-- Tab bar -->
                    <div class="flex border-b border-base-300 px-4 pt-3 gap-1 shrink-0">
                        <button v-for="tab in rightTabs" :key="tab.key"
                            class="flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-t-box border-b-2 transition-colors cursor-pointer"
                            :class="rightTab === tab.key
                                ? 'border-primary text-primary'
                                : 'border-transparent text-base-content/50 hover:text-base-content'"
                            @click="rightTab = tab.key">
                            <component :is="tab.icon" class="size-4" />
                            {{ tab.label }}
                        </button>
                    </div>

                    <!-- Tab: Informasi Toko -->
                    <div v-show="rightTab === 'toko'" class="flex flex-col flex-1 overflow-y-auto px-6 py-5 gap-4">

                        <div class="flex items-center gap-3">
                            <div class="bg-success/10 text-success rounded-box p-2.5 shrink-0">
                                <Store class="size-5" />
                            </div>
                            <div>
                                <h2 class="font-bold text-base leading-tight">Informasi Toko</h2>
                                <p class="text-base-content/50 text-sm mt-0.5">Tampil di struk dan laporan</p>
                            </div>
                        </div>

                        <div class="space-y-3">
                            <fieldset class="fieldset">
                                <legend class="fieldset-legend">Nama Toko</legend>
                                <label class="input w-full" :class="tokoErrors.nama ? 'input-error' : ''">
                                    <Store class="size-4 opacity-50" />
                                    <input v-model="tokoForm.nama" type="text" placeholder="Contoh: Toko Maju Jaya" />
                                </label>
                                <p v-if="tokoErrors.nama" class="fieldset-label text-error">{{ tokoErrors.nama }}</p>
                            </fieldset>

                            <fieldset class="fieldset">
                                <legend class="fieldset-legend">No. Telepon</legend>
                                <label class="input w-full">
                                    <Phone class="size-4 opacity-50" />
                                    <input v-model="tokoForm.telepon" type="text" placeholder="Contoh: 0812-3456-7890" />
                                </label>
                            </fieldset>

                            <fieldset class="fieldset">
                                <legend class="fieldset-legend">Email <span class="font-normal text-base-content/40">(opsional)</span></legend>
                                <label class="input w-full">
                                    <Mail class="size-4 opacity-50" />
                                    <input v-model="tokoForm.email" type="email" placeholder="Contoh: toko@email.com" />
                                </label>
                            </fieldset>

                            <fieldset class="fieldset">
                                <legend class="fieldset-legend">Alamat</legend>
                                <textarea v-model="tokoForm.alamat" class="textarea w-full" rows="3"
                                    placeholder="Jl. Contoh No. 1, Kota"></textarea>
                            </fieldset>
                        </div>

                        <div class="mt-auto pt-2">
                            <button class="btn btn-primary btn-sm gap-2 w-full" @click="saveToko" :disabled="tokoSaving">
                                <span v-if="tokoSaving" class="loading loading-spinner loading-sm"></span>
                                <Save v-else class="size-4" />
                                Simpan Informasi
                            </button>
                        </div>

                    </div>

                    <!-- Tab: Lisensi -->
                    <div v-show="rightTab === 'lisensi'" class="flex flex-col flex-1 overflow-y-auto px-6 py-5 gap-4">

                        <div class="flex items-center gap-3">
                            <div class="bg-warning/10 text-warning rounded-box p-2.5 shrink-0">
                                <Key class="size-5" />
                            </div>
                            <div>
                                <h2 class="font-bold text-base leading-tight">Detail Lisensi</h2>
                                <p class="text-base-content/50 text-sm mt-0.5">Informasi lisensi aplikasi ini</p>
                            </div>
                        </div>

                        <div class="space-y-2.5">

                            <!-- Status -->
                            <div class="flex items-center justify-between px-3 py-2.5 bg-base-200 rounded-box">
                                <span class="text-sm text-base-content/60">Status</span>
                                <span v-if="lisensiInfo.status === 'active'"   class="badge badge-success badge-sm">Aktif</span>
                                <span v-else-if="lisensiInfo.status === 'expiring'" class="badge badge-warning badge-sm">Hampir Habis</span>
                                <span v-else-if="lisensiInfo.status === 'expired'"  class="badge badge-error badge-sm">Kadaluarsa</span>
                                <span v-else class="badge badge-ghost badge-sm">Tidak diketahui</span>
                            </div>

                            <!-- License Key -->
                            <div class="px-3 py-2.5 bg-base-200 rounded-box">
                                <p class="text-xs text-base-content/40 mb-1 flex items-center gap-1.5">
                                    <Key class="size-3" /> License Key
                                </p>
                                <p class="font-mono text-sm font-semibold tracking-widest">{{ lisensiInfo.key || '-' }}</p>
                            </div>

                            <!-- Paket & Tipe -->
                            <div class="grid grid-cols-2 gap-px bg-base-300 rounded-box overflow-hidden">
                                <div class="bg-base-100 flex flex-col gap-1 p-3">
                                    <span class="text-xs text-base-content/40 flex items-center gap-1.5">
                                        <ShieldCheck class="size-3" /> Paket
                                    </span>
                                    <span class="font-semibold text-sm capitalize">{{ lisensiInfo.paket || '-' }}</span>
                                </div>
                                <div class="bg-base-100 flex flex-col gap-1 p-3">
                                    <span class="text-xs text-base-content/40 flex items-center gap-1.5">
                                        <Key class="size-3" /> Tipe
                                    </span>
                                    <span class="font-semibold text-sm capitalize">{{ lisensiInfo.tipe || '-' }}</span>
                                </div>
                            </div>

                            <!-- Expiry & Sisa Hari -->
                            <div class="grid grid-cols-2 gap-px bg-base-300 rounded-box overflow-hidden">
                                <div class="bg-base-100 flex flex-col gap-1 p-3">
                                    <span class="text-xs text-base-content/40 flex items-center gap-1.5">
                                        <Calendar class="size-3" /> Berlaku Hingga
                                    </span>
                                    <span class="font-semibold text-sm">{{ lisensiInfo.expiryDate || '-' }}</span>
                                </div>
                                <div class="bg-base-100 flex flex-col gap-1 p-3">
                                    <span class="text-xs text-base-content/40 flex items-center gap-1.5">
                                        <Clock class="size-3" /> Sisa Hari
                                    </span>
                                    <span v-if="lisensiInfo.tipe === 'lifetime'" class="font-semibold text-sm text-success">∞</span>
                                    <span v-else class="font-semibold text-sm"
                                        :class="lisensiInfo.daysLeft !== null && lisensiInfo.daysLeft <= 0 ? 'text-error' : lisensiInfo.daysLeft !== null && lisensiInfo.daysLeft <= 7 ? 'text-warning' : ''">
                                        {{ lisensiInfo.daysLeft !== null ? lisensiInfo.daysLeft + ' hari' : '-' }}
                                    </span>
                                </div>
                            </div>

                            <!-- Terakhir Divalidasi -->
                            <div class="px-3 py-2.5 bg-base-200 rounded-box">
                                <p class="text-xs text-base-content/40 mb-1 flex items-center gap-1.5">
                                    <Clock class="size-3" /> Terakhir Divalidasi
                                </p>
                                <p class="text-sm">{{ lisensiInfo.lastValidated || '-' }}</p>
                            </div>

                            <!-- Device ID -->
                            <div class="px-3 py-2.5 bg-base-200 rounded-box">
                                <p class="text-xs text-base-content/40 mb-1 flex items-center gap-1.5">
                                    <HardDrive class="size-3" /> Device ID
                                </p>
                                <p class="font-mono text-xs text-base-content/60 break-all leading-relaxed">{{ lisensiInfo.deviceId || '-' }}</p>
                            </div>

                        </div>

                        <!-- Konfigurasi Ulang -->
                        <div class="mt-auto pt-3 space-y-2.5">
                            <div class="divider text-xs text-base-content/30 my-0">Konfigurasi Perangkat</div>
                            <div class="flex items-center justify-between px-3 py-2.5 bg-base-200 rounded-box">
                                <span class="text-sm text-base-content/60">Mode Saat Ini</span>
                                <span class="badge badge-ghost badge-sm capitalize">
                                    {{ appMode === 'master' ? 'Kasir Utama' : appMode === 'client' ? 'Kasir Tambahan' : '-' }}
                                </span>
                            </div>
                            <button class="btn btn-outline btn-sm w-full gap-2" @click="openReconfigModal()">
                                <Settings class="size-4" />
                                Konfigurasi Ulang Perangkat
                            </button>
                        </div>

                    </div>

                    <!-- Tab: Database -->
                    <div v-show="rightTab === 'database' && isMaster" class="flex flex-col flex-1 overflow-y-auto">

                        <div class="flex items-center gap-3 px-6 pt-5 pb-4">
                            <div class="bg-primary/10 text-primary rounded-box p-2.5 shrink-0">
                                <DatabaseIcon class="size-5" />
                            </div>
                            <div>
                                <h2 class="font-bold text-base leading-tight">Database</h2>
                                <p class="text-base-content/50 text-sm mt-0.5">Kelola dan backup data aplikasi</p>
                            </div>
                        </div>

                        <div class="divider my-0 mx-6"></div>

                        <div class="px-6 pt-4 space-y-3 flex-1">
                            <div class="grid grid-cols-2 gap-px bg-base-300 rounded-box overflow-hidden">
                                <div class="bg-base-100 flex flex-col gap-1 p-3">
                                    <span class="text-xs text-base-content/40 flex items-center gap-1.5">
                                        <DatabaseIcon class="size-3" /> Ukuran
                                    </span>
                                    <span class="font-semibold text-sm">{{ formatDbSize(dbInfo.size) }}</span>
                                </div>
                                <div class="bg-base-100 flex flex-col gap-1 p-3">
                                    <span class="text-xs text-base-content/40 flex items-center gap-1.5">
                                        <Clock class="size-3" /> Backup Terakhir
                                    </span>
                                    <span class="font-semibold text-sm" :class="lastBackup ? '' : 'text-base-content/30'">
                                        {{ lastBackup || 'Belum pernah' }}
                                    </span>
                                </div>
                            </div>

                            <div class="flex items-start gap-2 bg-base-200 rounded-box px-3 py-2.5">
                                <HardDrive class="size-3.5 shrink-0 mt-0.5 text-base-content/40" />
                                <p class="font-mono text-xs text-base-content/50 break-all leading-relaxed">{{ dbInfo.path || '-' }}</p>
                            </div>
                        </div>

                        <div class="grid grid-cols-2 gap-3 px-6 py-5 mt-auto">
                            <button class="btn btn-primary btn-sm gap-2" @click="doBackup" :disabled="dbBusy">
                                <span v-if="dbBusy && dbAction === 'backup'" class="loading loading-spinner loading-sm"></span>
                                <Download v-else class="size-4" />
                                Backup
                            </button>
                            <button class="btn btn-outline btn-warning btn-sm gap-2" @click="restoreModalRef.showModal()" :disabled="dbBusy">
                                <span v-if="dbBusy && dbAction === 'restore'" class="loading loading-spinner loading-sm"></span>
                                <Upload v-else class="size-4" />
                                Restore
                            </button>
                        </div>

                    </div>

                    <!-- Tab: Database — client mode placeholder -->
                    <div v-show="rightTab === 'database' && !isMaster"
                         class="flex-1 flex flex-col items-center justify-center gap-2 text-base-content/30">
                        <DatabaseIcon class="size-10" />
                        <p class="text-sm">Database dikelola di perangkat master</p>
                    </div>

                    <!-- Tab: Printer -->
                    <div v-show="rightTab === 'printer'" class="flex flex-col flex-1 overflow-y-auto px-6 py-5 gap-4">

                        <div class="flex items-center gap-3">
                            <div class="bg-secondary/10 text-secondary rounded-box p-2.5 shrink-0">
                                <Printer class="size-5" />
                            </div>
                            <div>
                                <h2 class="font-bold text-base leading-tight">Pengaturan Printer</h2>
                                <p class="text-base-content/50 text-sm mt-0.5">Konfigurasi printer struk kasir</p>
                            </div>
                        </div>

                        <div class="space-y-3">

                            <!-- Pilih Printer -->
                            <fieldset class="fieldset">
                                <legend class="fieldset-legend">Printer</legend>
                                <div class="flex gap-2">
                                    <select v-model="printerName" class="select flex-1">
                                        <option value="">-- Pilih Printer --</option>
                                        <option v-for="p in printerList" :key="p.name" :value="p.name">
                                            {{ p.displayName || p.name }}{{ p.isDefault ? ' (Default)' : '' }}
                                        </option>
                                    </select>
                                    <button class="btn btn-ghost btn-square" :disabled="printerLoading"
                                        title="Muat ulang daftar printer" @click="loadPrinterList">
                                        <span v-if="printerLoading" class="loading loading-spinner loading-sm"></span>
                                        <RefreshCw v-else class="size-4" />
                                    </button>
                                </div>
                                <p v-if="!printerList.length && !printerLoading" class="fieldset-label text-base-content/40">
                                    Tidak ada printer terdeteksi
                                </p>
                            </fieldset>

                            <!-- Lebar Kertas -->
                            <fieldset class="fieldset">
                                <legend class="fieldset-legend">Lebar Kertas</legend>
                                <div class="flex gap-2">
                                    <label v-for="opt in [{ val: '58mm', label: '58 mm' }, { val: '80mm', label: '80 mm' }]"
                                        :key="opt.val"
                                        class="flex items-center gap-3 flex-1 border rounded-box px-3 py-2.5 cursor-pointer transition-colors"
                                        :class="paperWidth === opt.val
                                            ? 'border-primary/60 bg-primary/10 text-primary'
                                            : 'border-base-300 hover:bg-base-200'">
                                        <input type="radio" v-model="paperWidth" :value="opt.val"
                                            class="radio radio-primary radio-sm" />
                                        <span class="font-medium text-sm">{{ opt.label }}</span>
                                    </label>
                                </div>
                            </fieldset>

                            <!-- Auto Print -->
                            <div class="flex items-center justify-between px-3 py-2.5 bg-base-200 rounded-box cursor-pointer"
                                @click="printerAuto = !printerAuto">
                                <div>
                                    <p class="text-sm font-medium">Auto Print Struk</p>
                                    <p class="text-xs text-base-content/40 mt-0.5">Cetak otomatis setelah transaksi selesai</p>
                                </div>
                                <input v-model="printerAuto" type="checkbox" class="toggle toggle-primary toggle-sm" />
                            </div>

                        </div>

                        <div class="mt-auto pt-2 flex gap-2">
                            <button class="btn btn-ghost btn-sm gap-2 flex-1"
                                :disabled="!printerName || testPrinting" @click="testPrint">
                                <span v-if="testPrinting" class="loading loading-spinner loading-sm"></span>
                                <Printer v-else class="size-4" />
                                Test Print
                            </button>
                            <button class="btn btn-primary btn-sm gap-2 flex-1"
                                :disabled="printerSaving" @click="savePrinterConfig">
                                <span v-if="printerSaving" class="loading loading-spinner loading-sm"></span>
                                <Save v-else class="size-4" />
                                Simpan
                            </button>
                        </div>

                    </div>

                </div>
            </div>

        </div>
    </div>

    <!-- Modal Konfirmasi Hapus -->
    <dialog ref="deleteModalRef" class="modal modal-bottom sm:modal-middle">
        <div class="modal-box max-w-sm">
            <div class="flex items-center gap-3 mb-4">
                <div class="bg-error/10 text-error rounded-box p-2.5">
                    <Trash class="size-5" />
                </div>
                <div>
                    <h3 class="font-bold text-lg leading-tight">Hapus Pengguna</h3>
                    <p class="text-sm text-base-content/50">Tindakan ini tidak dapat dibatalkan</p>
                </div>
            </div>
            <p class="text-sm">
                Yakin ingin menghapus pengguna
                <span class="font-semibold">{{ userToDelete?.username }}</span>?
            </p>
            <div class="modal-action">
                <button class="btn btn-ghost" @click="deleteModalRef.close()">Batal</button>
                <button class="btn btn-error" :disabled="saving" @click="deleteUser">
                    <span v-if="saving" class="loading loading-spinner loading-sm"></span>
                    Hapus
                </button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button>close</button>
        </form>
    </dialog>

    <!-- Modal Tambah / Edit Pengguna -->
    <dialog ref="formModalRef" class="modal modal-bottom sm:modal-middle">
        <div class="modal-box max-w-sm">

            <!-- Header -->
            <div class="flex items-center gap-3 mb-6">
                <div class="bg-primary/10 text-primary rounded-box p-2.5">
                    <component :is="isEdit ? UserCog : UserPlus" class="size-5" />
                </div>
                <div>
                    <h3 class="font-bold text-lg leading-tight">{{ isEdit ? 'Edit Pengguna' : 'Tambah Pengguna' }}</h3>
                    <p class="text-sm text-base-content/50">
                        {{ isEdit ? 'Perbarui informasi akun' : 'Buat akun pengguna baru' }}
                    </p>
                </div>
            </div>

            <!-- Form fields -->
            <div class="space-y-3">
                <fieldset class="fieldset">
                    <legend class="fieldset-legend">Username</legend>
                    <label class="input w-full" :class="errors.username ? 'input-error' : ''">
                        <User class="size-4 opacity-50" />
                        <input v-model="form.username" type="text" placeholder="Masukkan username" />
                    </label>
                    <p v-if="errors.username" class="fieldset-label text-error">{{ errors.username }}</p>
                </fieldset>

                <fieldset class="fieldset">
                    <legend class="fieldset-legend">
                        Password
                        <span v-if="isEdit" class="font-normal text-base-content/40">(kosongkan jika tidak
                            diubah)</span>
                    </legend>
                    <label class="input w-full" :class="errors.password ? 'input-error' : ''">
                        <Lock class="size-4 opacity-50" />
                        <input v-model="form.password" type="password"
                            :placeholder="isEdit ? 'Password baru (opsional)' : 'Masukkan password'" />
                    </label>
                    <p v-if="errors.password" class="fieldset-label text-error">{{ errors.password }}</p>
                </fieldset>

                <div class="grid grid-cols-2 gap-3">
                    <fieldset class="fieldset">
                        <legend class="fieldset-legend">Role</legend>
                        <select v-model="form.role_id" class="select w-full"
                            :class="{ 'select-error': errors.role_id }">
                            <option value="" disabled>Pilih role</option>
                            <option v-for="role in roles" :key="role.id" :value="role.id">{{ role.name }}</option>
                        </select>
                        <p v-if="errors.role_id" class="fieldset-label text-error">{{ errors.role_id }}</p>
                    </fieldset>

                    <fieldset class="fieldset">
                        <legend class="fieldset-legend">Status</legend>
                        <div class="flex items-center gap-3 border border-base-300 rounded-box px-3 h-10 cursor-pointer"
                            @click="form.active = !form.active">
                            <input v-model="form.active" type="checkbox" class="toggle toggle-primary toggle-sm" />
                            <span class="text-sm font-medium"
                                :class="form.active ? 'text-success' : 'text-base-content/50'">
                                {{ form.active ? 'Aktif' : 'Nonaktif' }}
                            </span>
                        </div>
                    </fieldset>
                </div>
            </div>

            <div v-if="formError" class="alert alert-error mt-4">
                <span class="text-sm">{{ formError }}</span>
            </div>

            <div class="modal-action">
                <button class="btn btn-ghost" @click="closeForm">Batal</button>
                <button class="btn btn-primary" :disabled="saving" @click="save">
                    <span v-if="saving" class="loading loading-spinner loading-sm"></span>
                    {{ isEdit ? 'Simpan Perubahan' : 'Buat Akun' }}
                </button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button @click="closeForm">close</button>
        </form>
    </dialog>

    <!-- Modal Hak Akses -->
    <dialog ref="permModalRef" class="modal modal-bottom sm:modal-middle">
        <div class="modal-box max-w-lg">

            <!-- Header -->
            <div class="flex items-start justify-between mb-5">
                <div>
                    <h3 class="font-bold text-lg">Hak Akses</h3>
                    <p class="text-sm text-base-content/50 mt-0.5">Atur permission untuk setiap role</p>
                </div>
                <button class="btn btn-ghost btn-sm btn-square -mr-1 -mt-1" title="Tambah role baru"
                    @click="showNewRole = !showNewRole; newRoleName = ''; newRoleError = ''">
                    <Plus class="size-4" />
                </button>
            </div>

            <!-- Form tambah role baru -->
            <div v-if="showNewRole" class="bg-base-200 border border-base-300 rounded-box p-4 mb-5">
                <p class="font-semibold mb-3 flex items-center gap-2">
                    <ShieldCheck class="size-4 text-primary" />
                    Role Baru
                </p>
                <div class="flex gap-2">
                    <input v-model="newRoleName" type="text" class="input flex-1"
                        :class="{ 'input-error': newRoleError }" placeholder="Nama role, contoh: supervisor"
                        @keyup.enter="submitNewRole" />
                    <button class="btn btn-primary" :disabled="creatingRole" @click="submitNewRole">
                        <span v-if="creatingRole" class="loading loading-spinner loading-sm"></span>
                        Buat
                    </button>
                </div>
                <p v-if="newRoleError" class="text-error text-sm mt-2">{{ newRoleError }}</p>
            </div>

            <!-- Tabs role -->
            <div v-if="!showNewRole">
                <div role="tablist" class="tabs font-bold tabs-box mb-5">
                    <button v-for="role in roles" :key="role.id" role="tab" class="tab capitalize"
                        :class="{ 'tab-active': activeTab === role.id }" @click="activeTab = role.id">
                        {{ role.name }}
                    </button>
                </div>

                <!-- Permission cards untuk role aktif -->
                <div v-for="role in roles" :key="role.id" v-show="activeTab === role.id">
                    <div class="grid grid-cols-2 gap-2">
                        <label v-for="perm in allPermissions" :key="perm.key"
                            class="flex items-center gap-3 p-2 rounded-box border cursor-pointer transition-colors select-none"
                            :class="editPerms[ role.id ]?.includes(perm.key)
                                ? 'border-primary/60 bg-primary/20 text-primary'
                                : 'border-base-300 hover:bg-base-200 text-base-content'">
                            <component :is="perm.icon" class="size-4 shrink-0 transition-colors"
                                :class="editPerms[ role.id ]?.includes(perm.key) ? 'text-primary' : 'text-base-content/40'" />
                            <span class="flex-1 font-medium">{{ perm.label }}</span>
                            <input type="checkbox" class="checkbox checkbox-primary checkbox-sm"
                                :checked="editPerms[ role.id ]?.includes(perm.key)"
                                @change="togglePerm(role.id, perm.key)" />
                        </label>
                    </div>
                    <p class="text-xs text-base-content/40 mt-4">
                        Berlaku untuk semua pengguna dengan role
                        <span class="font-semibold capitalize">{{ role.name }}</span>.
                    </p>
                </div>
            </div>

            <div class="modal-action">
                <button class="btn btn-ghost" @click="permModalRef.close()">Tutup</button>
                <button v-if="!showNewRole" class="btn btn-primary"
                    :disabled="!isRoleChanged(activeTab) || savingRoleId === activeTab" @click="saveActiveRole">
                    <span v-if="savingRoleId === activeTab" class="loading loading-spinner loading-sm"></span>
                    Simpan
                </button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button>close</button>
        </form>
    </dialog>

    <!-- Modal Konfirmasi Restore -->
    <dialog ref="restoreModalRef" class="modal modal-bottom sm:modal-middle">
        <div class="modal-box max-w-sm">
            <div class="flex items-center gap-3 mb-4">
                <div class="bg-warning/10 text-warning rounded-box p-2.5">
                    <AlertTriangle class="size-5" />
                </div>
                <div>
                    <h3 class="font-bold text-lg leading-tight">Restore Database</h3>
                    <p class="text-sm text-base-content/50">Tindakan ini tidak dapat dibatalkan</p>
                </div>
            </div>
            <p class="text-sm">
                Restore akan <span class="font-semibold text-warning">mengganti seluruh data</span> dengan isi file backup yang dipilih.
                Pastikan sudah melakukan backup terlebih dahulu.
            </p>
            <div class="modal-action">
                <button class="btn btn-ghost" @click="restoreModalRef.close()">Batal</button>
                <button class="btn btn-warning" @click="doRestore">Pilih File &amp; Restore</button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button>close</button>
        </form>
    </dialog>

    <!-- Modal Konfirmasi Konfigurasi Ulang -->
    <dialog ref="reconfigModalRef" class="modal modal-bottom sm:modal-middle">
        <div class="modal-box max-w-sm">
            <div class="flex items-center gap-3 mb-4">
                <div class="bg-warning/10 text-warning rounded-box p-2.5">
                    <AlertTriangle class="size-5" />
                </div>
                <div>
                    <h3 class="font-bold text-lg leading-tight">Konfigurasi Ulang</h3>
                    <p class="text-sm text-base-content/50">Ubah peran komputer ini</p>
                </div>
            </div>

            <p class="text-sm mb-4">
                Mode <span class="font-semibold">{{ appMode === 'master' ? 'Kasir Utama' : 'Kasir Tambahan' }}</span>
                akan direset. Aplikasi akan memuat ulang dan menampilkan halaman pengaturan awal untuk memilih peran baru.
            </p>

            <!-- Wajib backup dulu jika master -->
            <template v-if="isMaster">
                <div v-if="!reconfigBackupDone" class="bg-error/10 border border-error/30 rounded-box px-3 py-3 flex flex-col gap-2.5">
                    <p class="text-sm font-semibold text-error flex items-center gap-2">
                        <AlertTriangle class="size-4 shrink-0" />
                        Backup database diperlukan
                    </p>
                    <p class="text-xs text-base-content/60">Komputer ini menyimpan seluruh data. Lakukan backup terlebih dahulu sebelum konfigurasi ulang agar data tidak hilang.</p>
                    <button class="btn btn-sm btn-error gap-2 w-full" :disabled="dbBusy" @click="doBackup">
                        <span v-if="dbBusy && dbAction === 'backup'" class="loading loading-spinner loading-sm"></span>
                        <Download v-else class="size-4" />
                        Backup Sekarang
                    </button>
                </div>
                <div v-else class="bg-success/10 border border-success/30 rounded-box px-3 py-2.5 flex items-center gap-2">
                    <Save class="size-4 text-success shrink-0" />
                    <p class="text-sm text-success font-medium">Backup berhasil — siap konfigurasi ulang</p>
                </div>
            </template>

            <div class="modal-action">
                <button class="btn btn-ghost" @click="reconfigModalRef.close()">Batal</button>
                <button class="btn btn-warning" :disabled="isMaster && !reconfigBackupDone" @click="doReconfig">Lanjutkan</button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop">
            <button>close</button>
        </form>
    </dialog>

</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, h } from 'vue'
import { UserPlus, UserCog, User, Lock, ShieldCheck, Settings, Plus, LayoutDashboard, ShoppingCart, Package, ArrowLeftRight, BarChart2, Trash, Edit, HardDrive, Database as DatabaseIcon, Clock, Download, Upload, AlertTriangle, Unlock, LockIcon, Store, Phone, Mail, Save, Key, Calendar, Printer, RefreshCw } from 'lucide-vue-next'
import { FlexRender } from '@tanstack/vue-table'
import { useServerTable }  from '../../composables/useServerTable'
import AppPagination       from '../../components/AppPagination.vue'
import { config, device } from '../../services/tauriApi'
import UserApi       from '../../services/UserApi'
import PermissionApi from '../../services/PermissionApi'
import TokoApi       from '../../services/TokoApi'
import { useToast } from '../../composables/useToast'

const { showToast } = useToast()

// --- Tab kanan ---
const rightTab  = ref('toko')
const rightTabs = [
    { key: 'toko',     label: 'Informasi Toko', icon: Store        },
    { key: 'lisensi',  label: 'Lisensi',         icon: Key          },
    { key: 'database', label: 'Database',         icon: DatabaseIcon },
    { key: 'printer',  label: 'Printer',          icon: Printer      },
]
// --- End Tab kanan ---

// --- Informasi Toko ---
const tokoForm   = reactive({ nama: '', telepon: '', email: '', alamat: '' })
const tokoErrors = reactive({ nama: '' })
const tokoSaving = ref(false)

async function loadToko() {
    try {
        const res = await TokoApi.get()
        const data = res.data ?? {}
        tokoForm.nama    = data.nama    ?? ''
        tokoForm.telepon = data.telepon ?? ''
        tokoForm.email   = data.email   ?? ''
        tokoForm.alamat  = data.alamat  ?? ''
    } catch { }
}

async function saveToko() {
    tokoErrors.nama = tokoForm.nama.trim() ? '' : 'Nama toko wajib diisi'
    if (tokoErrors.nama) return

    tokoSaving.value = true
    try {
        await TokoApi.update({
            nama:    tokoForm.nama.trim(),
            telepon: tokoForm.telepon.trim(),
            email:   tokoForm.email.trim(),
            alamat:  tokoForm.alamat.trim(),
        })
        showToast('Informasi toko berhasil disimpan', 'success')
    } catch {
        showToast('Gagal menyimpan informasi toko', 'error')
    } finally {
        tokoSaving.value = false
    }
}
// --- End Informasi Toko ---

// --- Lisensi ---
const lisensiInfo = reactive({
    key: '',
    deviceId: '',
    lastValidated: '',
    expiryDate: '',
    daysLeft: null,
    tipe: '',
    paket: '',
    status: 'unknown',
})

async function loadLisensi() {
    const key             = await config.get('license_key')
    const token           = await config.get('license_token')
    const deviceId        = await config.get('device_id')
    const lastValidatedAt = await config.get('last_validated_at')

    lisensiInfo.key      = key      || ''
    lisensiInfo.deviceId = deviceId || ''

    if (lastValidatedAt) {
        lisensiInfo.lastValidated = new Date(lastValidatedAt).toLocaleString('id-ID')
    }

    if (token) {
        try {
            const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')))
            lisensiInfo.tipe  = payload.tipe  || ''
            lisensiInfo.paket = payload.paket || ''

            if (payload.tipe === 'lifetime') {
                lisensiInfo.expiryDate = 'Seumur Hidup'
                lisensiInfo.daysLeft   = null
                lisensiInfo.status     = 'active'
            } else if (payload.expired_at) {
                lisensiInfo.expiryDate = new Date(payload.expired_at * 1000).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
                const daysLeft = Math.ceil((payload.expired_at * 1000 - Date.now()) / 86400000)
                lisensiInfo.daysLeft = daysLeft
                lisensiInfo.status   = daysLeft <= 0 ? 'expired' : daysLeft <= 7 ? 'expiring' : 'active'
            }
        } catch { }
    }
}
// --- End Lisensi ---

// --- Split pane ---
const splitContainerRef = ref(null)
const leftWidth         = ref(50)
const isResizing        = ref(false)

function startResize() {
    isResizing.value = true
    document.addEventListener('mousemove', doResize)
    document.addEventListener('mouseup', stopResize)
}

function doResize(e) {
    if (!isResizing.value || !splitContainerRef.value) return
    const rect = splitContainerRef.value.getBoundingClientRect()
    const pct  = ((e.clientX - rect.left) / rect.width) * 100
    leftWidth.value = Math.min(Math.max(pct, 25), 75)
}

function stopResize() {
    isResizing.value = false
    document.removeEventListener('mousemove', doResize)
    document.removeEventListener('mouseup', stopResize)
    config.set('pengaturan_split_width', leftWidth.value)
}

onUnmounted(() => {
    document.removeEventListener('mousemove', doResize)
    document.removeEventListener('mouseup', stopResize)
})
// --- End Split pane ---

// --- Database ---
const dbInfo       = reactive({ path: '', size: 0 })
const lastBackup   = ref('')
const dbBusy       = ref(false)
const dbAction     = ref('')
const restoreModalRef     = ref(null)
const reconfigModalRef    = ref(null)
const reconfigBackupDone  = ref(false)

const isMaster  = ref(false)
const appMode   = ref('')

async function loadDbInfo() {
    const mode = await config.get('app_mode')
    isMaster.value = mode === 'master'
    appMode.value  = mode || ''
    if (!isMaster.value) return
    const info = Promise.resolve({ path: "(Tauri - lihat app data dir)", size: 0 })
    dbInfo.path = info.path
    dbInfo.size = info.size
    const lb = await config.get('last_backup_at')
    if (lb) lastBackup.value = new Date(lb).toLocaleString('id-ID')
}

function formatDbSize(bytes) {
    if (!bytes) return '-'
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB'
}

async function doBackup() {
    dbBusy.value   = true
    dbAction.value = 'backup'
    try {
        const res = Promise.resolve({ success: false, error: "Backup belum tersedia di versi Tauri" })
        if (res.success) {
            await config.set('last_backup_at', Date.now())
            lastBackup.value        = new Date().toLocaleString('id-ID')
            reconfigBackupDone.value = true
            showToast('Backup berhasil disimpan', 'success')
        }
    } catch {
        showToast('Backup gagal', 'error')
    } finally {
        dbBusy.value   = false
        dbAction.value = ''
    }
}

function openReconfigModal() {
    reconfigBackupDone.value = false
    reconfigModalRef.value.showModal()
}

async function doReconfig() {
    reconfigModalRef.value.close()
    await config.set('app_mode', '')
    window.location.reload()
}

async function doRestore() {
    restoreModalRef.value.close()
    dbBusy.value  = true
    dbAction.value = 'restore'
    try {
        const res = Promise.resolve({ success: false, error: "Restore belum tersedia di versi Tauri" })
        console.log('[restore result]', res)
        if (res.success) {
            showToast('Restore berhasil. Memuat ulang...', 'success')
            setTimeout(() => window.location.reload(), 1200)
        } else if (res.error) {
            showToast('Restore gagal: ' + res.error, 'error')
        }
    } catch {
        showToast('Restore gagal', 'error')
    } finally {
        dbBusy.value   = false
        dbAction.value = ''
    }
}
// --- End Database ---

// --- Printer ---
const printerList    = ref([])
const printerName    = ref('')
const paperWidth     = ref('80mm')
const printerAuto    = ref(false)
const printerLoading = ref(false)
const printerSaving  = ref(false)
const testPrinting   = ref(false)

async function loadPrinterList() {
    printerLoading.value = true
    try {
        printerList.value = Promise.resolve([])
    } catch {
        printerList.value = []
    } finally {
        printerLoading.value = false
    }
}

async function loadPrinterConfig() {
    printerName.value = (await config.get('printer_name'))  || ''
    paperWidth.value  = (await config.get('printer_paper')) || '80mm'
    printerAuto.value = !!(await config.get('printer_auto'))
    await loadPrinterList()
}

async function savePrinterConfig() {
    printerSaving.value = true
    try {
        await config.set('printer_name',  printerName.value)
        await config.set('printer_paper', paperWidth.value)
        await config.set('printer_auto',  printerAuto.value)
        showToast('Pengaturan printer berhasil disimpan', 'success')
    } catch {
        showToast('Gagal menyimpan pengaturan printer', 'error')
    } finally {
        printerSaving.value = false
    }
}

async function testPrint() {
    if (!printerName.value) return
    testPrinting.value = true
    try {
        const res = Promise.resolve({ success: false, error: "Print belum tersedia di versi Tauri" })
        if (res.success) {
            showToast('Test print berhasil', 'success')
        } else {
            showToast('Test print gagal: ' + (res.failureReason || res.error || 'error'), 'error')
        }
    } catch {
        showToast('Test print gagal', 'error')
    } finally {
        testPrinting.value = false
    }
}
// --- End Printer ---

const ICON_MAP = { LayoutDashboard, ShoppingCart, Package, ArrowLeftRight, BarChart2, Settings }

const allPermissions = ref([])
const roles          = ref([])

const userColumns = [
    // nomor urut, disembunyikan tapi diperlukan untuk sorting
    {
        id: 'index',
        header: 'No.',
        meta: { headerClass: 'text-center' },
        cell: info => {
            const pageIndex = info.table.getState().pagination.pageIndex
            const pageSize  = info.table.getState().pagination.pageSize
            return pageIndex * pageSize + info.row.index + 1
        },
    },
    {
        accessorKey: 'username',
        header: 'Username',
        cell: info => h('span', { class: 'font-medium' }, info.getValue()),
    },
    {
        accessorKey: 'role_name',
        header: 'Role',
        cell: info => h('span', { class: 'badge badge-ghost badge-sm' }, info.getValue()),
    },
    {
        accessorKey: 'active',
        header: 'Status',
        cell: info => h('span', {
            class: `badge badge-sm ${info.getValue() ? 'badge-success' : 'badge-error'}`
        }, info.getValue() ? 'Aktif' : 'Nonaktif'),
    },
    {
        id: 'aksi',
        header: 'Aksi',
        enableSorting: false,
        meta: { headerClass: 'justify-center', cellClass: 'text-center' },
        cell: ({ row }) => {
            const user = row.original
            return h('div', { class: 'flex justify-center' }, [
                h('div', { class: 'tooltip tooltip-bottom', 'data-tip': 'Edit pengguna' },
                    h('button', { class: 'btn btn-ghost btn-sm', onClick: () => openEdit(user) },
                        h(Edit, { class: 'size-4' })
                    )
                ),
                h('div', { class: 'tooltip tooltip-bottom', 'data-tip': user.active ? 'Nonaktifkan' : 'Aktifkan' },
                    h('button', { class: 'btn btn-ghost btn-sm', onClick: () => toggleActive(user) },
                        h(user.active ? LockIcon : Unlock, { class: 'size-4' })
                    )
                ),
                h('div', { class: 'tooltip tooltip-bottom', 'data-tip': 'Hapus pengguna' },
                    h('button', { class: 'btn btn-ghost btn-sm text-error', onClick: () => openDeleteModal(user) },
                        h(Trash, { class: 'size-4' })
                    )
                ),
            ])
        },
    },
]

const { table, loading, search, fetchData: reloadUsers } = useServerTable({
    columns:  userColumns,
    fetchFn:  params => UserApi.getAll(params),
    pageSize: 10,
})

const saving = ref(false)
const formError = ref('')
const isEdit       = ref(false)
const editId       = ref(null)
const userToDelete = ref(null)

const editPerms = ref({})
const origPerms = ref({})
const savingRoleId = ref(null)
const activeTab = ref(null)
const showNewRole = ref(false)
const newRoleName = ref('')
const newRoleError = ref('')
const creatingRole = ref(false)

const formModalRef   = ref(null)
const permModalRef   = ref(null)
const deleteModalRef = ref(null)

const form = reactive({ username: '', password: '', role_id: '', active: true })
const errors = reactive({ username: '', password: '', role_id: '' })

onMounted(async () => {
    const saved = await config.get('pengaturan_split_width')
    if (saved) leftWidth.value = saved
    await Promise.all([ loadRoles(), loadPermissions(), loadDbInfo(), loadToko(), loadLisensi(), loadPrinterConfig() ])
})

async function loadPermissions() {
    const res = await PermissionApi.getAll()
    allPermissions.value = (res.data ?? []).map(p => ({
        ...p,
        icon: ICON_MAP[p.icon_name] ?? Settings,
    }))
}

async function loadRoles() {
    const res = await UserApi.getRoles()
    roles.value = res.data ?? []
    initRolePerms()
}

function initRolePerms() {
    const ep = {}
    const op = {}
    for (const role of roles.value) {
        ep[ role.id ] = [ ...role.permissions ]
        op[ role.id ] = [ ...role.permissions ]
    }
    editPerms.value = ep
    origPerms.value = op
}

function openCreate() {
    isEdit.value = false
    editId.value = null
    resetForm()
    formModalRef.value.showModal()
}

function openEdit(user) {
    isEdit.value = true
    editId.value = user.id
    resetForm()
    form.username = user.username
    form.role_id = user.role_id
    form.active = !!user.active
    formModalRef.value.showModal()
}

function closeForm() {
    formModalRef.value.close()
    editId.value = null
    resetForm()
}

function resetForm() {
    form.username = ''
    form.password = ''
    form.role_id = ''
    form.active = true
    errors.username = ''
    errors.password = ''
    errors.role_id = ''
    formError.value = ''
}

function validate() {
    let ok = true
    errors.username = ''
    errors.password = ''
    errors.role_id = ''

    if (!form.username.trim()) { errors.username = 'Username wajib diisi'; ok = false }
    if (!isEdit.value && !form.password) { errors.password = 'Password wajib diisi'; ok = false }
    if (!form.role_id) { errors.role_id = 'Role wajib dipilih'; ok = false }

    return ok
}

async function save() {
    if (!validate()) return

    saving.value = true
    formError.value = ''

    try {
        const payload = {
            username: form.username.trim(),
            role_id: form.role_id,
            active: form.active ? 1 : 0,
        }
        if (form.password) payload.password = form.password

        const res = isEdit.value
            ? await UserApi.update(editId.value, payload)
            : await UserApi.create(payload)

        if (res.success) {
            closeForm()
            await reloadUsers()
        } else {
            formError.value = res.message || 'Gagal menyimpan data'
        }
    } catch (err) {
        formError.value = err.response?.data?.message || 'Terjadi kesalahan'
    } finally {
        saving.value = false
    }
}

async function toggleActive(user) {
    try {
        await UserApi.update(user.id, { active: user.active ? 0 : 1 })
        await reloadUsers()
    } catch { }
}

function openDeleteModal(user) {
    userToDelete.value = user
    deleteModalRef.value.showModal()
}

async function deleteUser() {
    if (!userToDelete.value) return
    saving.value = true
    try {
        await UserApi.remove(userToDelete.value.id)
        deleteModalRef.value.close()
        userToDelete.value = null
        await reloadUsers()
        showToast('Pengguna berhasil dihapus')
    } catch {
        showToast('Gagal menghapus pengguna', 'error')
    } finally {
        saving.value = false
    }
}

function openPermModal() {
    if (roles.value.length) activeTab.value = roles.value[ 0 ].id
    showNewRole.value = false
    newRoleName.value = ''
    newRoleError.value = ''
    permModalRef.value.showModal()
}

async function submitNewRole() {
    const name = newRoleName.value.trim()
    if (!name) { newRoleError.value = 'Nama role wajib diisi'; return }

    creatingRole.value = true
    newRoleError.value = ''
    try {
        const res = await UserApi.createRole({ name, permissions: [] })
        if (res.success) {
            await loadRoles()
            activeTab.value = roles.value.find(r => r.name === name)?.id ?? roles.value.at(-1)?.id
            showNewRole.value = false
            newRoleName.value = ''
        } else {
            newRoleError.value = res.message || 'Gagal membuat role'
        }
    } catch (err) {
        newRoleError.value = err.response?.data?.message || 'Terjadi kesalahan'
    } finally {
        creatingRole.value = false
    }
}

function togglePerm(roleId, permKey) {
    const perms = editPerms.value[ roleId ]
    if (!perms) return
    const idx = perms.indexOf(permKey)
    if (idx === -1) perms.push(permKey)
    else perms.splice(idx, 1)
}

function isRoleChanged(roleId) {
    const curr = editPerms.value[ roleId ] ?? []
    const orig = origPerms.value[ roleId ] ?? []
    return curr.length !== orig.length || curr.some(p => !orig.includes(p))
}

async function saveActiveRole() {
    const role = roles.value.find(r => r.id === activeTab.value)
    if (!role) return
    savingRoleId.value = role.id
    try {
        const res = await UserApi.updateRole(role.id, { permissions: editPerms.value[ role.id ] })
        if (res.success) {
            origPerms.value[ role.id ] = [ ...editPerms.value[ role.id ] ]
            permModalRef.value.close()
            showToast(`Hak akses role "${role.name}" berhasil disimpan`, 'success')
        } else {
            showToast(res.message || 'Gagal menyimpan hak akses', 'error')
        }
    } catch {
        showToast('Terjadi kesalahan saat menyimpan', 'error')
    } finally {
        savingRoleId.value = null
    }
}

</script>
