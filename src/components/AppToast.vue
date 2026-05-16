<template>
    <div class="toast toast-end z-9999">
        <transition-group name="toast">
            <div
                v-for="toast in toasts"
                :key="toast.id"
                class="alert shadow-lg"
                :class="toast.type === 'success' ? 'alert-success' : 'alert-error'"
            >
                <CheckCircle v-if="toast.type === 'success'" class="size-4 shrink-0" />
                <XCircle v-else class="size-4 shrink-0" />
                <span class="text-sm">{{ toast.message }}</span>
            </div>
        </transition-group>
    </div>
</template>

<script setup>
import { CheckCircle, XCircle } from 'lucide-vue-next'
import { useToast } from '../composables/useToast'

const { toasts } = useToast()
</script>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all .25s ease; }
.toast-enter-from, .toast-leave-to      { opacity: 0; transform: translateX(1rem); }
</style>
