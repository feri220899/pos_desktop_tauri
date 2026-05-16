import apibackend from './api.js'

async function getAll() { return (await apibackend.get('/api/permissions')).data }

export default { getAll }
