<template>
    <div class="flex h-screen">

        <!-- Left Panel -->
        <div data-theme="night"
            class="hidden lg:flex w-2/6 bg-base-100 flex-col justify-between p-10 relative overflow-hidden">

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
                    Kelola bisnis<br />lebih mudah
                </h1>
                <p class="text-base-content/50 text-sm leading-relaxed mb-8">
                    Sistem kasir modern dengan manajemen produk, laporan real-time, dan dukungan multi kasir.
                </p>
                <div class="space-y-3">
                    <div v-for="item in features" :key="item" class="flex items-center gap-3">
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
        <div class="flex-1 flex items-center justify-center bg-base-200 p-8">
            <div class="w-full max-w-sm bg-base-100/30 p-8 rounded-2xl shadow-sm border transition-colors"
                :class="[ errors.api ? 'border-error' : 'border-transparent', shake ? 'animate-shake' : '' ]">

                <!-- Mobile logo -->
                <div class="lg:hidden flex items-center justify-center gap-2 mb-8">
                    <div class="btn btn-primary btn-square btn-sm no-animation pointer-events-none">
                        <LayoutGrid class="size-4" />
                    </div>
                    <span class="font-bold text-lg">POS Desktop</span>
                </div>

                <!-- Heading -->
                <div class="mb-7">
                    <h2 class="text-2xl font-bold text-base-content">Selamat datang!</h2>
                    <p class="text-base-content/50 text-sm mt-1">Masuk untuk melanjutkan ke aplikasi</p>
                </div>

                <!-- Form -->
                <form @submit.prevent="submit" class="space-y-4">

                    <div class="space-y-1.5">
                        <label class="text-sm font-medium text-base-content">Username</label>
                        <label class="input input-bordered flex items-center gap-2 w-full"
                            :class="{ 'input-error': errors.username || errors.api, 'input-disabled': loading }">
                            <User class="size-4 text-base-content/30 shrink-0" />
                            <input v-model="form.username" type="text" class="grow" placeholder="Masukkan username"
                                autocomplete="username" autofocus :disabled="loading"
                                @input="errors.username = ''; errors.api = ''" />
                        </label>
                        <p v-if="errors.username" class="text-xs text-error mt-1">{{ errors.username }}</p>
                    </div>

                    <div class="space-y-1.5">
                        <label class="text-sm font-medium text-base-content">Password</label>
                        <label class="input input-bordered flex items-center gap-2 w-full"
                            :class="{ 'input-error': errors.password || errors.api, 'input-disabled': loading }">
                            <Lock class="size-4 text-base-content/30 shrink-0" />
                            <input v-model="form.password" :type="showPass ? 'text' : 'password'" class="grow"
                                placeholder="Masukkan password" autocomplete="current-password" :disabled="loading"
                                @input="errors.password = ''; errors.api = ''" />
                            <button type="button" @click="showPass = !showPass"
                                class="text-base-content/30 hover:text-base-content transition-colors" tabindex="-1">
                                <EyeOff v-if="showPass" class="size-4" />
                                <Eye v-else class="size-4" />
                            </button>
                        </label>
                        <p v-if="errors.password" class="text-xs text-error mt-1">{{ errors.password }}</p>
                    </div>

                    <button type="submit" class="btn btn-primary w-full mt-2" :disabled="loading">
                        <span v-if="loading" class="loading loading-spinner loading-sm"></span>
                        {{ loading ? 'Memproses...' : 'Masuk' }}
                    </button>

                </form>

                <p class="text-center text-xs text-base-content/25 mt-10">v1.0.0</p>
            </div>
        </div>

    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { LayoutGrid, User, Lock, Eye, EyeOff, Check } from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = ref({ username: '', password: '' })
const loading = ref(false)
const showPass = ref(false)
const shake = ref(false)
const year = new Date().getFullYear()
const errors = ref({ username: '', password: '', api: '' })

const features = [
    'Manajemen produk & stok',
    'Laporan penjualan real-time',
    'Multi kasir & multi role',
    'Mode master & client (LAN)',
]

async function submit() {
    errors.value = { username: '', password: '', api: '' }

    if (!form.value.username.trim()) errors.value.username = 'Username tidak boleh kosong.'
    if (!form.value.password.trim()) errors.value.password = 'Password tidak boleh kosong.'
    if (errors.value.username || errors.value.password) return

    loading.value = true
    try {
        await authStore.login(form.value.username, form.value.password)
        router.push('/dashboard')
    } catch (e) {
        errors.value.api = e.response?.data?.message ?? 'Username atau password salah.'
        shake.value = true
        setTimeout(() => { shake.value = false }, 500)
    } finally {
        loading.value = false
    }
}
</script>

<style scoped>
@keyframes shake {

    0%,
    100% {
        transform: translateX(0);
    }

    15% {
        transform: translateX(-6px);
    }

    30% {
        transform: translateX(6px);
    }

    45% {
        transform: translateX(-4px);
    }

    60% {
        transform: translateX(4px);
    }

    75% {
        transform: translateX(-2px);
    }

    90% {
        transform: translateX(2px);
    }
}

.animate-shake {
    animation: shake 0.5s ease-in-out;
}
</style>
