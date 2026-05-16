import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apibackend from '../services/api'

export const useAuthStore = defineStore('auth', () => {
    const token = ref(localStorage.getItem('auth_token') || null)
    const user  = ref(JSON.parse(localStorage.getItem('auth_user') || 'null'))

    const isLoggedIn  = computed(() => !!token.value)
    const permissions = computed(() => user.value?.permissions ?? [])
    const can         = (permission) => permissions.value.includes(permission)

    async function login(username, password) {
        const res = await apibackend.post('/api/auth/login', { username, password })
        token.value = res.data.data.token
        user.value  = res.data.data.user        
        localStorage.setItem('auth_token', token.value)
        localStorage.setItem('auth_user', JSON.stringify(user.value))
        apibackend.reset()
    }

    function logout() {
        token.value = null
        user.value  = null
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')
        apibackend.reset()
    }

    return { token, user, isLoggedIn, permissions, can, login, logout }
})
