<template>
    <div class="flex h-screen">

        <!-- Left Panel -->
        <div data-theme="night" class="hidden lg:flex max-w-2/6 bg-base-100 flex-col justify-between p-10 relative overflow-hidden">

            <!-- Decorative circles -->
            <div class="absolute -top-24 -left-24 w-72 h-72 rounded-full bg-primary/10 pointer-events-none"></div>
            <div class="absolute -bottom-24 -right-24 w-96 h-96 rounded-full bg-primary/5 pointer-events-none"></div>

            <!-- Logo -->
            <div class="relative z-10 flex items-center gap-3">
                <div class="btn btn-primary btn-square btn-sm no-animation pointer-events-none">
                    <LayoutGrid class="size-4" />
                </div>
                <span class="font-bold text-base-content">POS Desktop</span>
            </div>

            <!-- Center content -->
            <div class="relative z-10">
                <h1 class="text-3xl font-bold text-base-content leading-tight mb-4">
                    Siapkan kasir<br/>dalam hitungan menit
                </h1>
                <p class="text-base-content/50 text-sm leading-relaxed mb-8">
                    Tentukan peran komputer ini lalu sistem siap digunakan. Bisa diubah kapan saja melalui pengaturan.
                </p>
                <div class="space-y-3">
                    <div v-for="item in tips" :key="item" class="flex items-center gap-3">
                        <div class="size-5 rounded-full bg-primary/20 flex items-center justify-center shrink-0">
                            <Check class="size-3 text-primary" />
                        </div>
                        <span class="text-sm text-base-content/60">{{ item }}</span>
                    </div>
                </div>
            </div>

            <!-- Footer -->
            <div class="relative z-10 text-xs text-base-content/20">
                &copy; {{ year }} POS Desktop
            </div>
        </div>

        <!-- Right Panel -->
        <div class="flex-1 flex items-center justify-center bg-base-200 p-6 overflow-y-auto">
            <div class="w-full max-w-lg bg-base-100/30 p-8 rounded-2xl shadow-sm">

                <!-- Mobile logo -->
                <div class="lg:hidden flex items-center justify-center gap-2 mb-8">
                    <div class="btn btn-primary btn-square btn-sm no-animation pointer-events-none">
                        <LayoutGrid class="size-4" />
                    </div>
                    <span class="font-bold text-lg">POS Desktop</span>
                </div>

                <!-- Heading -->
                <div class="mb-6">
                    <h2 class="text-2xl font-bold text-base-content">Pengaturan Awal</h2>
                    <p class="text-base-content/50 text-sm mt-1">Tentukan peran komputer ini dalam sistem kasir Anda</p>
                </div>

                <!-- Steps -->
                <ul class="steps steps-horizontal w-full text-xs mb-6">
                    <li class="step" :class="mode ? 'step-primary' : ''">Pilih Peran</li>
                    <li class="step" :class="bisa ? 'step-primary' : ''">Konfigurasi</li>
                    <li class="step">Selesai</li>
                </ul>

                <div class="space-y-4">

                    <!-- Mode selection: radio cards -->
                    <div>
                        <!-- Prompt saat belum dipilih -->
                        <div v-if="!mode" class="flex items-center gap-2 mb-3">
                            <MousePointerClick class="size-4 text-primary shrink-0" />
                            <p class="text-sm font-medium text-base-content">Pilih salah satu peran untuk komputer ini:</p>
                        </div>
                        <p v-else class="text-sm font-medium text-base-content mb-3">Komputer ini berfungsi sebagai:</p>

                        <div class="grid grid-cols-2 gap-3">

                            <label
                                class="card card-border border-2 cursor-pointer transition-all has-checked:border-primary has-checked:bg-primary/5 has-checked:shadow-none"
                                :class="!mode ? 'hover:border-primary hover:bg-primary/5 hover:shadow-sm' : ''"
                            >
                                <div class="card-body items-center gap-2 py-5">
                                    <input type="radio" name="mode" class="hidden" value="master" v-model="mode" @change="pilihMode('master')" />
                                    <div class="size-14 rounded-2xl flex items-center justify-center transition-colors"
                                        :class="mode === 'master' ? 'bg-primary text-primary-content' : 'bg-base-200 text-base-content/50'">
                                        <Server class="size-7" />
                                    </div>
                                    <div class="text-center">
                                        <p class="font-bold text-sm">Kasir Utama</p>
                                        <p class="text-xs text-base-content/50 mt-0.5 leading-relaxed">Menyimpan semua data &amp; laporan</p>
                                    </div>
                                    <div v-if="mode === 'master'" class="badge badge-primary badge-sm">Dipilih</div>
                                    <div v-else class="badge badge-ghost badge-sm opacity-0 group-hover:opacity-100">Pilih</div>
                                </div>
                            </label>

                            <label
                                class="card card-border border-2 cursor-pointer transition-all has-checked:border-primary has-checked:bg-primary/5 has-checked:shadow-none"
                                :class="!mode ? 'hover:border-primary hover:bg-primary/5 hover:shadow-sm' : ''"
                            >
                                <div class="card-body items-center gap-2 py-5">
                                    <input type="radio" name="mode" class="hidden" value="client" v-model="mode" @change="pilihMode('client')" />
                                    <div class="size-14 rounded-2xl flex items-center justify-center transition-colors"
                                        :class="mode === 'client' ? 'bg-primary text-primary-content' : 'bg-base-200 text-base-content/50'">
                                        <Monitor class="size-7" />
                                    </div>
                                    <div class="text-center">
                                        <p class="font-bold text-sm">Kasir Tambahan</p>
                                        <p class="text-xs text-base-content/50 mt-0.5 leading-relaxed">Terhubung ke kasir utama via WiFi</p>
                                    </div>
                                    <div v-if="mode === 'client'" class="badge badge-primary badge-sm">Dipilih</div>
                                    <div v-else class="badge badge-ghost badge-sm opacity-0">Pilih</div>
                                </div>
                            </label>

                        </div>
                    </div>

                    <!-- Info alert -->
                    <div v-if="mode === 'master'" role="alert" class="alert alert-info alert-soft text-sm py-2.5">
                        <Info class="size-4 shrink-0 opacity-70" />
                        <span>Kasir utama harus <strong>selalu menyala</strong> agar kasir tambahan bisa beroperasi.</span>
                    </div>

                    <div v-if="mode === 'client'" role="alert" class="alert alert-warning alert-soft text-sm py-2.5">
                        <AlertTriangle class="size-4 shrink-0 opacity-70" />
                        <span>Pastikan kasir utama sudah menyala dan terhubung ke <strong>WiFi yang sama.</strong></span>
                    </div>

                    <!-- Master: port config -->
                    <div v-if="mode === 'master'" class="card card-border border-2">
                        <div class="card-body py-4 px-4 flex-row items-center gap-4">
                            <div class="flex-1">
                                <p class="font-semibold text-sm">Port Server</p>
                                <p class="text-xs text-base-content/50 mt-0.5">Nomor yang dipakai kasir tambahan untuk terhubung. Biarkan <strong>3001</strong> jika tidak yakin.</p>
                            </div>
                            <input
                                v-model="port"
                                type="number"
                                placeholder="3001"
                                class="input input-bordered w-28 shrink-0 text-center"
                            />
                        </div>
                    </div>

                    <!-- Client: discovery -->
                    <div v-if="mode === 'client'" class="flex flex-col gap-3">

                        <div class="flex items-center justify-between">
                            <div>
                                <p class="font-semibold text-sm">Cari Kasir Utama</p>
                                <p class="text-xs text-base-content/50">Kasir utama aktif di jaringan WiFi Anda</p>
                            </div>
                            <button class="btn btn-sm btn-ghost" :disabled="scanning" @click="startScan">
                                <span v-if="scanning" class="loading loading-spinner loading-xs"></span>
                                <RefreshCw v-else class="size-4" />
                                Cari Ulang
                            </button>
                        </div>

                        <!-- Scanning -->
                        <div v-if="scanning && masters.length === 0" class="flex flex-col items-center gap-2 py-8">
                            <span class="loading loading-dots loading-md text-primary"></span>
                            <p class="text-sm text-base-content/60">Sedang mencari kasir utama...</p>
                        </div>

                        <!-- Empty -->
                        <div v-if="!scanning && masters.length === 0" role="alert" class="alert">
                            <WifiOff class="size-5 text-base-content/40" />
                            <div>
                                <p class="font-semibold text-sm">Tidak ada kasir utama ditemukan</p>
                                <p class="text-xs text-base-content/50">Nyalakan kasir utama lalu klik "Cari Ulang"</p>
                            </div>
                        </div>

                        <!-- Master list -->
                        <div v-if="masters.length > 0" class="flex flex-col gap-2">
                            <p class="text-xs text-base-content/50 font-medium">Pilih kasir utama yang ingin dihubungkan:</p>
                            <label
                                v-for="m in masters"
                                :key="m.ip + m.port"
                                class="card card-border border-2 cursor-pointer transition-colors has-checked:border-primary has-checked:bg-primary/5"
                            >
                                <div class="card-body flex-row items-center gap-3 py-3 px-4">
                                    <input type="radio" name="master" class="radio radio-primary radio-sm" :value="m" v-model="selectedMaster" />
                                    <div class="flex-1">
                                        <p class="font-semibold text-sm">{{ m.name }}</p>
                                        <p class="text-xs text-base-content/50">{{ m.ip }} · Port {{ m.port }}</p>
                                    </div>
                                    <div class="badge badge-success badge-soft gap-1.5">
                                        <span class="size-1.5 rounded-full bg-success inline-block"></span>
                                        Online
                                    </div>
                                </div>
                            </label>
                        </div>

                    </div>

                    <!-- Error -->
                    <div v-if="error" role="alert" class="alert alert-error">
                        <AlertCircle class="size-4 shrink-0" />
                        <span class="text-sm">{{ error }}</span>
                    </div>

                    <!-- Submit -->
                    <button class="btn btn-primary w-full" :disabled="!bisa || loading" @click="simpan">
                        <span v-if="loading" class="loading loading-spinner loading-sm"></span>
                        {{ loading ? 'Menyimpan...' : 'Simpan & Lanjutkan' }}
                    </button>

                </div>
            </div>
        </div>

    </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { LayoutGrid, Check, Server, Monitor, Info, AlertTriangle, AlertCircle, RefreshCw, WifiOff, MousePointerClick } from 'lucide-vue-next'
import { config } from '../../services/tauriApi'

const router = useRouter()
const year   = new Date().getFullYear()

const tips = [
    'Kasir utama menyimpan semua data lokal',
    'Kasir tambahan terhubung via jaringan WiFi',
    'Satu jaringan bisa punya banyak kasir',
]

const mode           = ref(null)
const port           = ref(3001)
const masters        = ref([])
const selectedMaster = ref(null)
const scanning       = ref(false)
const loading        = ref(false)
const error          = ref(null)

const bisa = computed(() => {
    if (!mode.value) return false
    if (mode.value === 'client') return !!selectedMaster.value
    return true
})

function pilihMode(m) {
    if (m === 'client') startScan()
}

async function startScan() {
    masters.value        = []
    selectedMaster.value = null
    scanning.value       = true

    try {
        const res = await fetch('http://localhost:3001/api/internal/discovery/scan')
        if (res.ok) {
            const data = await res.json()
            masters.value = data.data ?? []
        }
    } catch { /* backend belum jalan / offline */ }

    scanning.value = false
}

async function simpan() {
    error.value = null
    loading.value = true
    try {
        await config.set('app_mode', mode.value)

        if (mode.value === 'master') {
            const serverPort = Number(port.value) || 3001
            await config.set('server_port', serverPort)
            fetch('http://localhost:3001/api/internal/discovery/advertise', {
                method:  'POST',
                headers: { 'Content-Type': 'application/json' },
                body:    JSON.stringify({ port: serverPort }),
            }).catch(() => {})
        } else {
            await config.set('master_ip',   selectedMaster.value.ip)
            await config.set('master_port', selectedMaster.value.port)
        }

        router.push('/aktivasi')
    } catch {
        error.value = 'Gagal menyimpan konfigurasi. Coba lagi.'
    } finally {
        loading.value = false
    }
}
</script>
