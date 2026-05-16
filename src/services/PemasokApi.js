import apibackend from './api.js'

async function getAll(params = {}) { return (await apibackend.get('/api/pemasok', { params })).data }
async function create(data)     { return (await apibackend.post('/api/pemasok', data)).data }
async function update(id, data) { return (await apibackend.put(`/api/pemasok/${id}`, data)).data }
async function remove(id)       { return (await apibackend.delete(`/api/pemasok/${id}`)).data }

export default { getAll, create, update, remove }
