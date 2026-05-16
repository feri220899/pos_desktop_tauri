<script setup>
import { ref, reactive, h } from 'vue'
import { FlexRender } from '@tanstack/vue-table'
import { useServerTable } from '../../composables/useServerTable.js'
import { useToast } from '../../composables/useToast.js'
import AppPagination from '../../components/AppPagination.vue'
import PemasokApi from '../../services/PemasokApi.js'

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
    { accessorKey: 'nama_pemasok', header: 'Nama Pemasok', meta: { cellClass: 'font-medium' } },
    {
        accessorKey: 'telepon', header: 'No. Telepon', enableSorting: false,
        meta: { cellClass: 'text-xs' },
        cell: info => info.getValue() || '—',
    },
    {
        accessorKey: 'alamat', header: 'Alamat', enableSorting: false,
        meta: { cellClass: 'text-xs text-base-content/70 max-w-56 truncate' },
        cell: info => info.getValue() || '—',
    },
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
    fetchFn: params => PemasokApi.getAll(params),
    pageSize: 10,
})

// ── Form Tambah ────────────────────────────────────────────
const formOpen  = ref(true)
const saving    = ref(false)
const form      = reactive({ nama_pemasok: '', telepon: '', alamat: '' })
const formErr   = reactive({ nama_pemasok: '' })

function validateForm() {
    formErr.nama_pemasok = form.nama_pemasok.trim() ? '' : 'Nama pemasok wajib diisi.'
    return !formErr.nama_pemasok
}

function resetForm() {
    Object.assign(form, { nama_pemasok: '', telepon: '', alamat: '' })
    formErr.nama_pemasok = ''
}

async function simpan() {
    if (!validateForm()) return
    saving.value = true
    try {
        await PemasokApi.create({
            nama_pemasok: form.nama_pemasok.trim(),
            telepon: form.telepon.trim() || null,
            alamat:  form.alamat.trim()  || null,
        })
        showToast('Pemasok berhasil disimpan.')
        resetForm()
        fetchData()
    } catch (e) { showToast(e.response?.data?.message || 'Gagal menyimpan pemasok.', 'error') }
    finally { saving.value = false }
}

// ── Edit ───────────────────────────────────────────────────
const modalEdit = ref(null)
const editId    = ref(null)
const updating  = ref(false)
const editForm  = reactive({ nama_pemasok: '', telepon: '', alamat: '' })
const editErr   = reactive({ nama_pemasok: '' })

function validateEdit() {
    editErr.nama_pemasok = editForm.nama_pemasok.trim() ? '' : 'Nama pemasok wajib diisi.'
    return !editErr.nama_pemasok
}

function openEdit(row) {
    editId.value = row.id
    editErr.nama_pemasok = ''
    Object.assign(editForm, {
        nama_pemasok: row.nama_pemasok ?? '',
        telepon:      row.telepon      ?? '',
        alamat:       row.alamat       ?? '',
    })
    modalEdit.value?.showModal()
}

async function update() {
    if (!validateEdit()) return
    updating.value = true
    try {
        await PemasokApi.update(editId.value, {
            nama_pemasok: editForm.nama_pemasok.trim(),
            telepon: editForm.telepon.trim() || null,
            alamat:  editForm.alamat.trim()  || null,
        })
        showToast('Pemasok berhasil diupdate.')
        modalEdit.value?.close()
        fetchData()
    } catch (e) { showToast(e.response?.data?.message || 'Gagal mengupdate pemasok.', 'error') }
    finally { updating.value = false }
}

// ── Hapus ──────────────────────────────────────────────────
async function hapus(id) {
    if (!confirm('Hapus pemasok ini?')) return
    try {
        await PemasokApi.remove(id)
        showToast('Pemasok berhasil dihapus.')
        fetchData()
    } catch { showToast('Gagal menghapus pemasok.', 'error') }
}
</script>

<template>
    <div class="flex-1 flex flex-col gap-4 min-h-0 overflow-y-auto pb-4">

        <!-- Header -->
        <div class="shrink-0">
            <h1 class="text-lg font-bold">Pemasok</h1>
            <p class="text-sm text-base-content/50">Kelola data pemasok / supplier produk</p>
        </div>

        <!-- Form Tambah -->
        <div class="card bg-base-100 border border-base-300 shadow-sm shrink-0">
            <div class="card-body p-0">
                <div class="flex items-center justify-between px-4 py-2.5 border-b border-base-300 bg-base-200/50 cursor-pointer"
                     @click="formOpen = !formOpen">
                    <span class="text-sm font-semibold">Tambah Pemasok</span>
                    <span class="text-xs text-base-content/40">{{ formOpen ? 'Sembunyikan ▲' : 'Tampilkan ▼' }}</span>
                </div>
                <div v-show="formOpen" class="p-4">
                    <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 max-w-2xl">
                        <div>
                            <label class="label-text text-xs font-medium mb-1 block">
                                Nama Pemasok <span class="text-error">*</span>
                            </label>
                            <input v-model="form.nama_pemasok" type="text"
                                   placeholder="Nama supplier / distributor"
                                   :class="['input input-sm input-bordered w-full', formErr.nama_pemasok ? 'input-error' : '']"
                                   @input="formErr.nama_pemasok = ''"
                                   @keyup.enter="simpan" />
                            <p v-if="formErr.nama_pemasok" class="text-error text-xs mt-1">{{ formErr.nama_pemasok }}</p>
                        </div>
                        <div>
                            <label class="label-text text-xs font-medium mb-1 block">
                                No. Telepon <span class="text-base-content/40">(opsional)</span>
                            </label>
                            <input v-model="form.telepon" type="text" placeholder="0812-xxxx-xxxx"
                                   class="input input-sm input-bordered w-full" />
                        </div>
                        <div>
                            <label class="label-text text-xs font-medium mb-1 block">
                                Alamat <span class="text-base-content/40">(opsional)</span>
                            </label>
                            <input v-model="form.alamat" type="text" placeholder="Alamat pemasok"
                                   class="input input-sm input-bordered w-full" />
                        </div>
                    </div>
                    <div class="mt-3">
                        <button class="btn btn-primary btn-sm gap-2" :disabled="saving" @click="simpan">
                            <span v-if="saving" class="loading loading-spinner loading-xs"></span>
                            Simpan Pemasok
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Tabel -->
        <div class="card bg-base-100 border border-base-300 shadow-sm flex-1 min-h-0 overflow-hidden">
            <div class="card-body flex flex-col overflow-hidden p-4">
                <div class="flex items-center justify-between mb-1 shrink-0">
                    <span class="text-sm font-semibold">Daftar Pemasok</span>
                    <span class="text-xs text-base-content/40">{{ table.getRowCount() }} pemasok</span>
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
                                <td colspan="6" class="text-center py-10 text-base-content/40">
                                    <span class="loading loading-spinner loading-sm"></span>
                                </td>
                            </tr>
                            <tr v-else-if="table.getRowModel().rows.length === 0">
                                <td colspan="6" class="text-center py-10 text-base-content/40">Belum ada pemasok</td>
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
            <h3 class="font-bold text-base mb-4">Edit Pemasok</h3>
            <div class="flex flex-col gap-3">
                <div>
                    <fieldset class="fieldset">
                        <legend class="fieldset-legend">
                            Nama Pemasok <span class="text-error">*</span>
                        </legend>
                        <input v-model="editForm.nama_pemasok" type="text"
                               :class="['input input-bordered w-full', editErr.nama_pemasok ? 'input-error' : '']"
                               @input="editErr.nama_pemasok = ''" />
                    </fieldset>
                    <p v-if="editErr.nama_pemasok" class="text-error text-xs mt-1 px-1">{{ editErr.nama_pemasok }}</p>
                </div>
                <fieldset class="fieldset">
                    <legend class="fieldset-legend">No. Telepon</legend>
                    <input v-model="editForm.telepon" type="text" placeholder="Opsional"
                           class="input input-bordered w-full" />
                </fieldset>
                <fieldset class="fieldset">
                    <legend class="fieldset-legend">Alamat</legend>
                    <input v-model="editForm.alamat" type="text" placeholder="Opsional"
                           class="input input-bordered w-full" />
                </fieldset>
            </div>
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
