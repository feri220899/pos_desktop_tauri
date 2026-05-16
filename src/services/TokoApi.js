import apibackend from './api.js'

async function get()         { return (await apibackend.get('/api/toko')).data }
async function update(data)  { return (await apibackend.put('/api/toko', data)).data }

export default { get, update }
