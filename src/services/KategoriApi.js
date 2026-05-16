import apibackend from './api.js'

async function getAll(params = {}) { return (await apibackend.get('/api/kategori', { params })).data }
async function create(data)     { return (await apibackend.post('/api/kategori', data)).data }
async function update(id, data) { return (await apibackend.put(`/api/kategori/${id}`, data)).data }
async function remove(id)       { return (await apibackend.delete(`/api/kategori/${id}`)).data }

export default { getAll, create, update, remove }
