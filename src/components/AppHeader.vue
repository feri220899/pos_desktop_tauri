<template>
    <header class="navbar bg-base-100 border-b border-base-200 min-h-16 px-6 shrink-0">
        <div class="navbar-start">
            <div class="breadcrumbs text-sm py-0">
                <ul>
                    <li v-for="(crumb, i) in breadcrumbs" :key="i">
                        <router-link v-if="crumb.to" :to="crumb.to"
                            class="text-base-content/50 hover:text-base-content gap-1.5">
                            <component v-if="crumb.icon" :is="crumb.icon" class="size-3.5 shrink-0" />
                            {{ crumb.label }}
                        </router-link>
                        <span v-else class="font-semibold text-base-content gap-1.5 flex items-center">
                            <component v-if="crumb.icon" :is="crumb.icon" class="size-3.5 shrink-0" />
                            {{ crumb.label }}
                        </span>
                    </li>
                </ul>
            </div>
        </div>
        <div class="navbar-end gap-4">
            <span class="text-sm text-base-content/60 font-mono">{{ time }}</span>
            <div class="divider divider-horizontal mx-0 h-5"></div>
            <div class="flex items-center gap-2">
                <span :class="[ 'badge badge-xs', online ? 'badge-success' : 'badge-error' ]"></span>
                <span class="text-sm text-base-content/60">{{ online ? 'Online' : 'Offline' }}</span>
            </div>
            <div class="divider divider-horizontal mx-0 h-5"></div>
            <div class="dropdown dropdown-end">
                <button tabindex="0" class="btn btn-ghost btn-sm gap-2">
                    <div class="avatar avatar-placeholder">
                        <div class="bg-primary text-primary-content rounded-full w-7 text-xs">
                            <span>{{ initials }}</span>
                        </div>
                    </div>
                    <span class="text-sm">{{ authStore.user?.username }}</span>
                    <ChevronDown class="size-3 text-base-content/50" />
                </button>
                <ul tabindex="0"
                    class="dropdown-content menu bg-base-100 rounded-box shadow-lg border border-base-200 w-44 p-1 mt-1 z-10">
                    <li class="menu-title text-xs">
                        <span class="badge badge-ghost badge-sm capitalize">{{ authStore.user?.role }}</span>
                    </li>
                    <li>
                        <button @click="logout" class="text-error">
                            <LogOut class="size-4" />
                            Keluar
                        </button>
                    </li>
                </ul>
            </div>
        </div>
    </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { LogOut, ChevronDown } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { buildBreadcrumbs } from '../config/menu'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const time = ref('')
const online = ref(navigator.onLine)

const breadcrumbs = computed(() => buildBreadcrumbs(route.path))

const initials = computed(() => {
    const name = authStore.user?.username ?? ''
    return name.slice(0, 2).toUpperCase()
})

function logout() {
    authStore.logout()
    router.push('/login')
}

let timer = null

function updateTime() {
    time.value = new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

onMounted(() => {
    updateTime()
    timer = setInterval(updateTime, 1000)
    window.addEventListener('online', () => online.value = true)
    window.addEventListener('offline', () => online.value = false)
})

onUnmounted(() => clearInterval(timer))
</script>
