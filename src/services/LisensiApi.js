import axios from 'axios'

const client = axios.create({
    baseURL: 'http://lisansi.test/api',
    timeout: 10000,
})

async function aktivasi(licenseKey, deviceId) {
    try {
        const { data } = await client.post('/lisensi/aktivasi', { license_key: licenseKey, device_id: deviceId })
        return { success: true, ...data }
    } catch (error) {
        const resData = error.response?.data
        return { success: false, message: resData?.pesan ?? resData?.message ?? 'Gagal terhubung ke server lisensi.' }
    }
}

async function validasi(licenseKey, deviceId) {
    try {
        const { data } = await client.post('/lisensi/validasi', { license_key: licenseKey, device_id: deviceId })
        return { success: true, ...data }
    } catch (error) {
        const resData = error.response?.data
        return { success: false, message: resData?.pesan ?? resData?.message ?? 'Gagal terhubung ke server lisensi.' }
    }
}

async function deaktivasi(licenseKey, deviceId) {
    try {
        const { data } = await client.post('/lisensi/deaktivasi', { license_key: licenseKey, device_id: deviceId })
        return { success: true, ...data }
    } catch (error) {
        const resData = error.response?.data
        return { success: false, message: resData?.pesan ?? resData?.message ?? 'Gagal terhubung ke server lisensi.' }
    }
}

export default { aktivasi, validasi, deaktivasi }
