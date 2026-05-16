<script setup>
import { ref, reactive, h } from 'vue'
import { FlexRender } from '@tanstack/vue-table'
import { useServerTable } from '../../composables/useServerTable.js'
import { useToast } from '../../composables/useToast.js'
import AppPagination from '../../components/AppPagination.vue'
import KategoriApi from '../../services/KategoriApi.js'

const { showToast } = useToast()

// ── Tabel ──────────────────────────────────────────────────
const columns = [
    {
        id: 'index', header: 'No', enableSorting: false,
        meta: { headerClass: 'w-10', cellClass: 'text-center text-xs text-base-content/50' },
        cell: info => {
            const { pageIndex, pageSize } = info.table.getState().pagination
            return pageIndex * pageSize + info.row.index + 1
        },
    },
    { accessorKey: 'nama_kategori', header: 'Nama Kategori', meta: { cellClass: 'font-medium' } },
    {
        accessorKey: 'created_at', header: 'Dibuat', enableSorting: false,
        meta: { cellClass: 'text-xs text-base-content/50' },
        cell: info => info.getValue() ? new Date(info.getValue()).toLocaleDateString('id-ID') : '—',
    },
    {
        id: 'aksi', header: 'Aksi', enableSorting: false,
        meta: { headerClass: 'text-center', cellClass: 'text-center' },
        cell: info => h('div', { class: 'flex gap-1 justify-center' }, [
            h('button', { class: 'btn btn-ghost btn-xs', onClick: () => openEdit(info.row.original) }, 'Edit'),
            h('button', { class: 'btn btn-ghost btn-xs text-error', onClick: () => hapus(info.row.original.id) }, 'Hapus'),
        ]),
    },
]

const { table, loading, search, fetchData } = useServerTable({
    columns,
    fetchFn: params => KategoriApi.getAll(params),
    pageSize: 10,
})

// ── Form Tambah ────────────────────────────────────────────
const formOpen = ref(true)
const saving   = ref(false)
const form     = reactive({ nama_kategori: '' })

async function simpan() {
    if (!form.nama_kategori.trim()) { showToast('Nama kategori wajib diisi.', 'error'); return }
    saving.value = true
    try {
        await KategoriApi.create({ nama_kategori: form.nama_kategori.trim() })
        showToast('Kategori berhasil disimpan.')
        form.nama_kategori = ''
        fetchData()
    } catch (e) { showToast(e.response?.data?.message || 'Gagal simpan.', 'error') }
    finally { saving.value = false }
}

// ── Edit ───────────────────────────────────────────────────
const modalEdit = ref(null)
const editId    = ref(null)
const updating  = ref(false)
const editForm  = reactive({ nama_kategori: '' })

function openEdit(row) {
    editId.value = row.id
    editForm.nama_kategori = row.nama_kategori
    modalEdit.value?.showModal()
}

async function update() {
    if (!editForm.nama_kategori.trim()) { showToast('Nama kategori wajib diisi.', 'error'); return }
    updating.value = true
    try {
        await KategoriApi.update(editId.value, { nama_kategori: editForm.nama_kategori.trim() })
        showToast('Kategori berhasil diupdate.')
        modalEdit.value?.close()
        fetchData()
    } catch (e) { showToast(e.response?.data?.message || 'Gagal update.', 'error') }
    finally { updating.value = false }
}

// ── Hapus ──────────────────────────────────────────────────
async function hapus(id) {
    if (!confirm('Hapus kategori ini?')) return
    try {
        await KategoriApi.remove(id)
        showToast('Kategori berhasil dihapus.')
        fetchData()
    } catch { showToast('Gagal hapus kategori.', 'error') }
}
</script>

<template>
    <div class="flex-1 flex flex-col gap-4 min-h-0 overflow-y-auto pb-4">

        <!-- Header -->
        <div class="shrink-0">
            <h1 class="text-lg font-bold">Kategori Produk</h1>
            <p class="text-sm text-base-content/50">Kelola kategori untuk pengelompokan produk</p>
        </div>

        <!-- Form Tambah -->
        <div class="card bg-base-100 border border-base-300 shadow-sm shrink-0">
            <div class="card-body p-0">
                <div class="flex items-center justify-between px-4 py-2.5 border-b border-base-300 bg-base-200/50 cursor-pointer"
                     @click="formOpen = !formOpen">
                    <span class="text-sm font-semibold">Tambah Kategori</span>
                    <span class="text-xs text-base-content/40">{{ formOpen ? 'Sembunyikan ▲' : 'Tampilkan ▼' }}</span>
                </div>
                <div v-show="formOpen" class="p-4">
                    <div class="flex gap-2 items-end max-w-sm">
                        <div class="flex-1">
                            <label class="label-text text-xs font-medium mb-1 block">Nama Kategori <span class="text-error">*</span></label>
                            <input v-model="form.nama_kategori" type="text" placeholder="Contoh: Minuman, Snack, dll"
                                   class="input input-sm input-bordered w-full"
                                   @keyup.enter="simpan" />
                        </div>
                        <button class="btn btn-primary btn-sm gap-2" :disabled="saving" @click="simpan">
                            <span v-if="saving" class="loading loading-spinner loading-xs"></span>
                            Simpan
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Tabel -->
        <div class="card bg-base-100 border border-base-300 shadow-sm flex-1 min-h-0 overflow-hidden">
            <div class="card-body flex flex-col overflow-hidden p-4">
                <div class="flex items-center justify-between mb-1 shrink-0">
                    <span class="text-sm font-semibold">Daftar Kategori</span>
                    <span class="text-xs text-base-content/40">{{ table.getRowCount() }} kategori</span>
                </div>
                <AppPagination :table="table" v-model:search="search" class="flex-1 min-h-0">
                    <table class="table table-xs">
                        <thead class="sticky top-0 z-10 bg-base-100">
                            <tr>
                                <th v-for="header in table.getFlatHeaders()" :key="header.id"
                                    :class="[header.column.columnDef.meta?.headerClass, header.column.getCanSort() ? 'cursor-pointer select-none' : '']"
                                    @click="header.column.getToggleSortingHandler()?.($event)">
                                    <div class="flex items-center gap-1">
                                        <FlexRender :render="header.column.columnDef.header" :props="header.getContext()" />
                                        <span v-if="header.column.getIsSorted() === 'asc'">↑</span>
                                        <span v-else-if="header.column.getIsSorted() === 'desc'">↓</span>
                                    </div>
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-if="loading">
                                <td colspan="4" class="text-center py-10 text-base-content/40">
                                    <span class="loading loading-spinner loading-sm"></span>
                                </td>
                            </tr>
                            <tr v-else-if="table.getRowModel().rows.length === 0">
                                <td colspan="4" class="text-center py-10 text-base-content/40">Belum ada kategori</td>
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

    </div>

    <!-- Modal Edit -->
    <dialog ref="modalEdit" class="modal">
        <div class="modal-box max-w-sm">
            <h3 class="font-bold text-base mb-4">Edit Kategori</h3>
            <fieldset class="fieldset">
                <legend class="fieldset-legend">Nama Kategori</legend>
                <input v-model="editForm.nama_kategori" type="text"
                       class="input input-bordered w-full"
                       @keyup.enter="update" />
            </fieldset>
            <div class="modal-action mt-4">
                <button class="btn btn-ghost btn-sm" @click="modalEdit?.close()">Batal</button>
                <button class="btn btn-primary btn-sm gap-2" :disabled="updating" @click="update">
                    <span v-if="updating" class="loading loading-spinner loading-xs"></span>
                    Update
                </button>
            </div>
        </div>
        <form method="dialog" class="modal-backdrop"><button>close</button></form>
    </dialog>

</template>
