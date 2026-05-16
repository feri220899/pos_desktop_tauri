import apibackend from './api.js'

async function getAll(params = {})      { return (await apibackend.get('/api/produk', { params })).data }
async function getById(id)              { return (await apibackend.get(`/api/produk/${id}`)).data }
async function getTrashed(params = {})  { return (await apibackend.get('/api/produk/trashed', { params })).data }
async function nextKode(prefix = 'A')   { return (await apibackend.get('/api/produk/next-kode', { params: { prefix } })).data }

async function create(data)             { return (await apibackend.post('/api/produk', data)).data }
async function update(id, data)         { return (await apibackend.put(`/api/produk/${id}`, data)).data }
async function remove(id)               { return (await apibackend.delete(`/api/produk/${id}`)).data }
async function restore(id)              { return (await apibackend.post(`/api/produk/${id}/restore`)).data }

async function getSatuan(produkId)             { return (await apibackend.get(`/api/produk/${produkId}/satuan`)).data }
async function storeSatuan(produkId, data)     { return (await apibackend.post(`/api/produk/${produkId}/satuan`, data)).data }
async function updateSatuan(produkId, satuanId, data) { return (await apibackend.put(`/api/produk/${produkId}/satuan/${satuanId}`, data)).data }
async function removeSatuan(produkId, satuanId)       { return (await apibackend.delete(`/api/produk/${produkId}/satuan/${satuanId}`)).data }

export default { getAll, getById, getTrashed, nextKode, create, update, remove, restore, getSatuan, storeSatuan, updateSatuan, removeSatuan }
