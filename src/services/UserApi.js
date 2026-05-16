import apibackend from './api.js'

async function getAll(params = {})    { return (await apibackend.get('/api/users', { params })).data }
async function getRoles()             { return (await apibackend.get('/api/roles')).data }
async function create(data)           { return (await apibackend.post('/api/users', data)).data }
async function update(id, data)       { return (await apibackend.put(`/api/users/${id}`, data)).data }
async function remove(id)             { return (await apibackend.delete(`/api/users/${id}`)).data }
async function updateRole(id, data)   { return (await apibackend.put(`/api/roles/${id}`, data)).data }
async function createRole(data)       { return (await apibackend.post('/api/roles', data)).data }

export default { getAll, getRoles, create, update, remove, updateRole, createRole }
