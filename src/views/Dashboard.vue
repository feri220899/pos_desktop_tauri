<template>
    <div class="space-y-6">

        <!-- Stats -->
        <div class="grid grid-cols-3 gap-4">
            <div v-for="stat in stats" :key="stat.label" class="stat bg-base-100 rounded-box shadow-sm">
                <div class="stat-figure">
                    <div :class="['rounded-btn p-2', stat.color]">
                        <svg xmlns="http://www.w3.org/2000/svg" class="size-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="stat.icon"/>
                        </svg>
                    </div>
                </div>
                <div class="stat-title">{{ stat.label }}</div>
                <div class="stat-value text-2xl">{{ stat.value }}</div>
                <div class="stat-desc">{{ stat.desc }}</div>
            </div>
        </div>

        <!-- Info -->
        <div class="card bg-base-100 shadow-sm">
            <div class="card-body py-4">
                <div class="flex items-center gap-2 text-sm text-base-content/60">
                    <svg xmlns="http://www.w3.org/2000/svg" class="size-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                    </svg>
                    {{ tanggal }}
                </div>
            </div>
        </div>

    </div>
</template>

<script setup>
const stats = [
    {
        label: 'Transaksi Hari Ini',
        value: '0',
        desc:  'transaksi',
        color: 'bg-primary',
        icon:  'M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 11h.01M12 11h.01M15 11h.01M4 19h16a2 2 0 002-2V7a2 2 0 00-2-2H4a2 2 0 00-2 2v10a2 2 0 002 2z',
    },
    {
        label: 'Pendapatan Hari Ini',
        value: 'Rp 0',
        desc:  'total pendapatan',
        color: 'bg-success',
        icon:  'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    },
    {
        label: 'Total Produk',
        value: '0',
        desc:  'produk aktif',
        color: 'bg-secondary',
        icon:  'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10',
    },
]

const tanggal = new Date().toLocaleDateString('id-ID', {
    weekday: 'long', year: 'numeric', month: 'long', day: 'numeric'
})
</script>
