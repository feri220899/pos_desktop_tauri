import axios from 'axios'
import { config } from './tauriApi'

let client = null

async function getClient() {
    if (client) return client

    const mode      = await config.get('app_mode')
    const port      = (await config.get('server_port')) ?? 3001
    const masterIp  = await config.get('master_ip')
    const masterPort = (await config.get('master_port')) ?? 3001

    const baseURL = mode === 'client'
        ? `http://${masterIp}:${masterPort}`
        : `http://localhost:${port}`

    const token = localStorage.getItem('auth_token')

    client = axios.create({
        baseURL,
        timeout: 10000,
        headers: token ? { Authorization: `Bearer ${token}` } : {},
    })

    client.interceptors.response.use(
        res => res,
        err => {
            if (err.response?.status === 401) {
                localStorage.removeItem('auth_token')
                localStorage.removeItem('auth_user')
                client = null
                window.location.hash = '/login'
            }
            return Promise.reject(err)
        }
    )

    return client
}

function reset() { client = null }

export default {
    get:    async (url, cfg)         => (await getClient()).get(url, cfg),
    post:   async (url, data, cfg)   => (await getClient()).post(url, data, cfg),
    put:    async (url, data, cfg)   => (await getClient()).put(url, data, cfg),
    delete: async (url, cfg)         => (await getClient()).delete(url, cfg),
    reset,
}
