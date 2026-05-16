import { invoke } from '@tauri-apps/api/core'

export const config = {
    get: (key) => invoke('get_config', { key }),
    set: (key, value) => invoke('set_config', { key, value }),
}

export const device = {
    getId: () => invoke('get_device_id'),
}

export const lisensi = {
    verifyToken: (token) => invoke('verify_license_token', { token }),
}
