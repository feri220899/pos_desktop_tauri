import { defineStore } from 'pinia'
import { ref } from 'vue'
import { config, device, lisensi as lisensiApi } from '../services/tauriApi'
import LisensiApi from '../services/LisensiApi'

export const useLisensiStore = defineStore('lisensi', () => {
    const tokenValid = ref(false)
    const daysLeft   = ref(null)
    const checking   = ref(false)

    async function check() {
        checking.value = true
        try {
            const token = await config.get('license_token')
            if (!token) return false

            const result = await lisensiApi.verifyToken(token)
            tokenValid.value = result.valid === true
            daysLeft.value   = result.days_left ?? null

            if (result.valid && result.days_left <= 2) {
                renewInBackground()
            }

            if (!result.valid && result.expired) {
                const lastValidated = await config.get('last_validated_at')
                if (lastValidated) {
                    const diff = Date.now() - Number(lastValidated)
                    if (diff < 3 * 24 * 60 * 60 * 1000) {
                        tokenValid.value = true
                        renewInBackground()
                        return true
                    }
                }
                return false
            }

            return tokenValid.value
        } finally {
            checking.value = false
        }
    }

    async function aktivasi(licenseKey) {
        const deviceId = await device.getId()
        const result   = await LisensiApi.aktivasi(licenseKey, deviceId)
        if (result.success && result.token) {
            await config.set('license_key', licenseKey)
            await config.set('license_token', result.token)
            await config.set('last_validated_at', Date.now())
            tokenValid.value = true
        }
        return result
    }

    async function renewInBackground() {
        try {
            const [licenseKey, deviceId] = await Promise.all([
                config.get('license_key'),
                device.getId(),
            ])
            if (!licenseKey) return
            const result = await LisensiApi.validasi(licenseKey, deviceId)
            if (result.success && result.token) {
                await config.set('license_token', result.token)
                await config.set('last_validated_at', Date.now())
            }
        } catch { /* silent fail — offline */ }
    }

    return { tokenValid, daysLeft, checking, check, aktivasi }
})
