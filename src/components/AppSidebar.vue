<template>
    <aside data-theme="night" class="w-64 flex flex-col h-screen shrink-0 bg-base-100 border-r border-base-300">

        <!-- Logo -->
        <div class="flex items-center gap-3 px-4 h-16 shrink-0 border-b border-base-300">
            <div class="btn btn-primary btn-square btn-sm no-animation pointer-events-none">
                <LayoutGrid class="size-4" />
            </div>
            <div>
                <div class="text-sm font-bold text-base-content leading-tight">POS Desktop</div>
                <div class="text-xs text-green-400">{{ mode === 'master' ? '● Master' : '○ Client' }}</div>
            </div>
        </div>

        <!-- Nav -->
        <nav class="flex-1 overflow-y-auto py-3 px-2 space-y-4">
            <div v-for="section in visibleMenu" :key="section.title">
                <p class="px-3 mb-1 text-xs font-semibold uppercase tracking-wider text-base-content/50">
                    {{ section.title }}
                </p>
                <ul class="space-y-0.5">
                    <li v-for="item in section.items" :key="item.label">

                        <!-- Item dengan children -->
                        <template v-if="item.children?.length">
                            <button
                                class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-base font-medium transition-colors"
                                :class="isParentActive(item)
                                    ? 'bg-primary/10 text-primary'
                                    : 'text-white/70 hover:bg-base-200 hover:text-white'"
                                @click="toggleOpen(item.label)">
                                <component :is="item.icon" class="size-4 shrink-0" />
                                <span class="flex-1 text-left">{{ item.label }}</span>
                                <ChevronDown class="size-3.5 shrink-0 transition-transform duration-200"
                                    :class="openItems.has(item.label) ? 'rotate-180' : ''" />
                            </button>

                            <!-- Sub-menu -->
                            <ul v-if="openItems.has(item.label)"
                                class="mt-0.5 ml-2 pl-3 border-l border-base-300 space-y-0.5">
                                <li v-for="child in item.children" :key="child.to">
                                    <RouterLink :to="child.to" :class="childNavClass(child.to)">
                                        <component v-if="child.icon" :is="child.icon" class="size-3.5 shrink-0" />
                                        <span class="size-1.5 rounded-full bg-current shrink-0 opacity-50"
                                            v-else></span>
                                        {{ child.label }}
                                    </RouterLink>
                                </li>
                            </ul>
                        </template>

                        <!-- Item biasa -->
                        <RouterLink v-else :to="item.to" :class="navClass(item.to)">
                            <component :is="item.icon" class="size-4 shrink-0" />
                            {{ item.label }}
                        </RouterLink>

                    </li>
                </ul>
            </div>
        </nav>

        <!-- Bottom -->
        
        <div class="px-2 py-3 border-t border-base-300" v-if="authStore.can('pengaturan')">
            <RouterLink to="/pengaturan" :class="navClass('/pengaturan')">
                <Settings class="size-4 shrink-0" />
                Pengaturan
            </RouterLink>
        </div>

    </aside>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ChevronDown, Settings, LayoutGrid } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { allMenu, bottomMenu } from '../config/menu'
import { config } from '../services/tauriApi'

const route = useRoute()
const authStore = useAuthStore()
const mode = ref('master')
const openItems = ref(new Set())

function navClass(path) {
    const base = 'flex items-center gap-3 px-3 py-2.5 rounded-lg text-base font-medium transition-colors'
    return route.path === path
        ? `${base} bg-primary text-white`
        : `${base} text-white/70 hover:bg-base-200 hover:text-white`
}

function childNavClass(path) {
    const base = 'flex items-center gap-2.5 px-3 py-2 rounded-lg text-base font-medium transition-colors'
    return route.path === path
        ? `${base} bg-primary text-white`
        : `${base} text-white/60 hover:bg-base-200 hover:text-white hover:text-white/80`
}

function toggleOpen(label) {
    const s = new Set(openItems.value)
    s.has(label) ? s.delete(label) : s.add(label)
    openItems.value = s
}

function isParentActive(item) {
    return item.children?.some(c => route.path === c.to)
}

// Auto-buka parent jika child sedang aktif
function syncOpenFromRoute() {
    allMenu.forEach(section => {
        section.items.forEach(item => {
            if (item.children?.some(c => route.path === c.to)) {
                const s = new Set(openItems.value)
                s.add(item.label)
                openItems.value = s
            }
        })
    })
}

watch(() => route.path, syncOpenFromRoute)


const visibleMenu = computed(() =>
    allMenu
        .map(section => ({
            ...section,
            items: section.items
                .map(item => {
                    if (!item.children) return item
                    return {
                        ...item,
                        children: item.children.filter(c =>
                            authStore.can(c.permission ?? item.permission)
                        ),
                    }
                })
                .filter(item =>
                    item.children
                        ? item.children.length > 0
                        : authStore.can(item.permission)
                ),
        }))
        .filter(section => section.items.length > 0)
)

onMounted(async () => {
    mode.value = (await config.get('app_mode')) ?? 'master'
    syncOpenFromRoute()
})
</script>
