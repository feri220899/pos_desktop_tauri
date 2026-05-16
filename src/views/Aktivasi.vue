<template>
    <div class="flex h-screen">

        <!-- Left Panel -->
        <div data-theme="night" class="hidden lg:flex max-w-2/6 bg-base-100 flex-col justify-between p-10 relative overflow-hidden">
            <div class="absolute -top-24 -left-24 w-72 h-72 rounded-full bg-primary/10 pointer-events-none"></div>
            <div class="absolute -bottom-24 -right-24 w-96 h-96 rounded-full bg-primary/5 pointer-events-none"></div>
            <div class="relative z-10 flex items-center gap-3">
                <div class="btn btn-primary btn-square btn-sm no-animation pointer-events-none">
                    <LayoutGrid class="size-4" />
                </div>
                <span class="font-bold text-base-content">POS Desktop</span>
            </div>
            <div class="relative z-10">
                <h1 class="text-3xl font-bold text-base-content leading-tight mb-4">
                    Satu lisensi,<br/>siap digunakan
                </h1>
                <p class="text-base-content/50 text-sm leading-relaxed mb-8">
                    Aktifkan aplikasi dengan license key yang Anda terima setelah pembelian.
                </p>
                <div class="space-y-3">
                    <div v-for="item in info" :key="item" class="flex items-center gap-3">
                        <div class="size-5 rounded-full bg-primary/20 flex items-center justify-center shrink-0">
                            <Check class="size-3 text-primary" />
                        </div>
                        <span class="text-sm text-base-content/60">{{ item }}</span>
                    </div>
                </div>
            </div>
            <div class="relative z-10 text-xs text-base-content/20">
                &copy; {{ year }} POS Desktop
            </div>
        </div>

        <!-- Right Panel -->
        <div class="flex-1 flex items-center justify-center bg-base-200 p-6">
            <div class="w-full max-w-lg bg-base-100/30 p-8 rounded-2xl shadow-sm">

                <div class="lg:hidden flex items-center justify-center gap-2 mb-8">
                    <div class="btn btn-primary btn-square btn-sm no-animation pointer-events-none">
                        <LayoutGrid class="size-4" />
                    </div>
                    <span class="font-bold text-lg">POS Desktop</span>
                </div>

                <div class="mb-7">
                    <h2 class="text-2xl font-bold text-base-content">Aktivasi Lisensi</h2>
                    <p class="text-base-content/50 text-sm mt-1">Masukkan license key untuk mulai menggunakan aplikasi</p>
                </div>

                <div class="space-y-5">
                    <div class="space-y-3">
                        <label class="text-sm font-medium text-base-content">License Key</label>
                        <div class="flex items-center gap-2">
                            <template v-for="(_, i) in parts" :key="i">
                                <input
                                    :ref="el => inputRefs[i] = el"
                                    v-model="parts[i]"
                                    type="text"
                                    maxlength="5"
                                    :disabled="loading"
                                    class="input input-bordered w-full text-center font-mono tracking-widest uppercase text-sm transition-colors"
                                    :class="{
                                        'input-error':   !!error,
                                        'input-success': isComplete && !error,
                                        'input-primary': parts[i].length === 5 && !isComplete && !error,
                                    }"
                                    @input="onInput(i)"
                                    @keydown="onKeydown($event, i)"
                                    @paste.prevent="onPaste($event)"
                                    @focus="$event.target.select()"
                                />
                                <span v-if="i < 3" class="text-base-content/30 font-mono font-bold shrink-0">—</span>
                            </template>
                        </div>
                        <p class="text-xs" :class="error ? 'text-error' : 'text-base-content/40'">
                            {{ error || 'License key dikirim ke email setelah pembelian' }}
                        </p>
                    </div>

                    <button class="btn btn-primary w-full" :disabled="loading || !isComplete" @click="aktivasi">
                        <span v-if="loading" class="loading loading-spinner loading-sm"></span>
                        {{ loading ? 'Memproses...' : 'Aktifkan Sekarang' }}
                    </button>

                    <div class="divider text-xs text-base-content/40 my-3">Belum punya lisensi?</div>

                    <div class="grid grid-cols-2 gap-2 mt-1">
                        <button class="btn btn-outline w-full gap-2" @click="bukaTrial">
                            <ExternalLink class="size-4" />
                            Coba Gratis 7 Hari
                        </button>
                        <button class="btn btn-ghost w-full gap-2" @click="bukaPaket">
                            <PackageSearch class="size-4" />
                            Cek Lisensi Saya
                        </button>
                    </div>
                </div>

                <p class="text-center text-xs text-base-content/25 mt-10">v1.0.0</p>
            </div>
        </div>

    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { openUrl } from '@tauri-apps/plugin-opener'
import { LayoutGrid, Check, ExternalLink, PackageSearch } from 'lucide-vue-next'
import { config, device, lisensi as lisensiTauri } from '../services/tauriApi'
import LisensiApi from '../services/LisensiApi'

const router = useRouter()
const year   = new Date().getFullYear()

const parts     = ref(['', '', '', ''])
const inputRefs = ref([])
const error     = ref('')
const loading   = ref(false)

const licenseKey = computed(() => parts.value.join('-'))
const isComplete = computed(() => parts.value.every(p => p.length === 5))

const info = [
    'Aktivasi sekali, gunakan selamanya',
    'Mendukung multi kasir dalam satu jaringan',
    'Update aplikasi gratis',
]

function bukaTrial()  { openUrl('https://lisansi.test') }
function bukaPaket()  { openUrl('https://lisansi.test/paket-saya') }

function onInput(i) {
    parts.value[i] = parts.value[i].toUpperCase().replace(/[^A-Z0-9]/g, '')
    error.value = ''
    if (parts.value[i].length === 5 && i < 3) {
        inputRefs.value[i + 1]?.focus()
    }
}

function onKeydown(e, i) {
    if (e.key === 'Backspace' && parts.value[i] === '' && i > 0) {
        inputRefs.value[i - 1]?.focus()
    }
    if (e.key === 'Enter' && isComplete.value) {
        aktivasi()
    }
}

function onPaste(e) {
    const text = e.clipboardData.getData('text').toUpperCase().replace(/[^A-Z0-9]/g, '')
    for (let i = 0; i < 4; i++) {
        parts.value[i] = text.slice(i * 5, i * 5 + 5)
    }
    inputRefs.value[3]?.focus()
}

onMounted(async () => {
    const savedKey = await config.get('license_key')
    const deviceId = await device.getId()

    if (savedKey) {
        const token = await config.get('license_token')
        if (token) {
            const verify = await lisensiTauri.verifyToken(token)
            if (verify.valid) {
                router.replace('/dashboard')
                return
            }
        }
        const result = await LisensiApi.validasi(savedKey, deviceId)
        if (result.success && result.token) {
            await config.set('license_token', result.token)
            await config.set('last_validated_at', Date.now())
            router.replace('/dashboard')
            return
        }
    }

    inputRefs.value[0]?.focus()
})

async function aktivasi() {
    if (!isComplete.value) {
        error.value = 'License key harus diisi lengkap.'
        return
    }
    loading.value = true
    error.value   = ''

    const deviceId = await device.getId()
    const result   = await LisensiApi.aktivasi(licenseKey.value, deviceId)

    if (result.success) {
        await config.set('license_key', licenseKey.value)
        if (result.token) {
            await config.set('license_token', result.token)
            await config.set('last_validated_at', Date.now())
        }
        router.replace('/dashboard')
    } else {
        error.value = result.pesan || result.message || 'Aktivasi gagal.'
    }
    loading.value = false
}
</script>
