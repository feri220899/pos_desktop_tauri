import { createRouter, createWebHashHistory } from 'vue-router'
import { config, device, lisensi as lisensiTauri } from '../services/tauriApi'
import LisensiApi from '../services/LisensiApi'

import Aktivasi   from '../views/Aktivasi.vue'
import Setup      from '../views/base/Setup.vue'
import Login      from '../views/base/Login.vue'
import Pengaturan from '../views/base/Pengaturan.vue'
import Dashboard  from '../views/Dashboard.vue'
import Kasir      from '../views/Kasir.vue'
import Transaksi  from '../views/Transaksi.vue'
import Laporan    from '../views/Laporan.vue'
import Produk     from '../views/produk/MasterProduk.vue'
import Kategori   from '../views/produk/Kategori.vue'
import Pemasok    from '../views/produk/Pemasok.vue'

const routes = [
    { path: '/', redirect: '/aktivasi' },
    { path: '/setup',       component: Setup,      meta: { layout: false } },
    { path: '/aktivasi',    component: Aktivasi,   meta: { layout: false } },
    { path: '/login',       component: Login,      meta: { layout: false } },
    { path: '/dashboard',   component: Dashboard,  meta: { layout: true, auth: true } },
    { path: '/kasir',       component: Kasir,      meta: { layout: true, auth: true, permission: 'kasir' } },
    { path: '/master-produk', component: Produk,   meta: { layout: true, auth: true, permission: 'produk' } },
    { path: '/kategori',    component: Kategori,   meta: { layout: true, auth: true, permission: 'produk' } },
    { path: '/pemasok',     component: Pemasok,    meta: { layout: true, auth: true, permission: 'produk' } },
    { path: '/transaksi',   component: Transaksi,  meta: { layout: true, auth: true, permission: 'transaksi' } },
    { path: '/laporan',     component: Laporan,    meta: { layout: true, auth: true, permission: 'laporan' } },
    { path: '/pengaturan',  component: Pengaturan, meta: { layout: true, auth: true, permission: 'pengaturan' } },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

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

router.beforeEach(async (to) => {
    if (to.path === '/setup') return true

    const appMode = await config.get('app_mode')
    if (!appMode) return '/setup'

    if (to.path === '/aktivasi' || to.path === '/login') return true

    const licenseToken = await config.get('license_token')
    if (!licenseToken) return '/aktivasi'

    const verify = await lisensiTauri.verifyToken(licenseToken)

    if (verify.valid) {
        if ((verify.days_left ?? 99) < 2) renewInBackground()
    } else if (verify.expired) {
        const lastValidated = await config.get('last_validated_at')
        if (lastValidated && (Date.now() - Number(lastValidated)) / 86400000 < 3) {
            renewInBackground()
        } else {
            return '/aktivasi'
        }
    } else {
        return '/aktivasi'
    }

    if (!to.meta.auth) return true

    const authToken = localStorage.getItem('auth_token')
    if (!authToken) return '/login'

    if (to.meta.permission) {
        const user = JSON.parse(localStorage.getItem('auth_user') || 'null')
        if (!user?.permissions?.includes(to.meta.permission)) return '/dashboard'
    }

    return true
})

export default router
