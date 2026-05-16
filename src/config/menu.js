import { LayoutGrid, Monitor, ShoppingBag, ClipboardList, BarChart2, Settings, Tag, Truck } from 'lucide-vue-next'

export const allMenu = [
    {
        title: 'Utama',
        items: [
            { to: '/dashboard', label: 'Dashboard', icon: LayoutGrid, permission: 'dashboard' },
        ],
    },
    {
        title: 'Manajemen',
        items: [
            //  { to: '/produk', label: 'Produk', icon: ClipboardList, permission: 'produk' },
            {
                label: 'Produk',
                icon: ShoppingBag,
                permission: 'produk',
                children: [
                    { to: '/master-produk', label: 'Master Produk',    icon: Monitor },
                    { to: '/kategori',      label: 'Kategori Produk',  icon: Tag },
                    { to: '/pemasok',       label: 'Pemasok',          icon: Truck },
                ],
            },
            {
                label: 'Penjualan',
                icon: ShoppingBag,
                permission: 'transaksi',
                children: [
                    { to: '/kasir',     label: 'Kasir',     icon: Monitor,      permission: 'kasir' },
                    { to: '/transaksi', label: 'Transaksi', icon: ClipboardList, permission: 'transaksi' },
                    { to: '/laporan',   label: 'Laporan',   icon: BarChart2,     permission: 'laporan' },
                ],
            },
        ],
    },
]

export const bottomMenu = [
    { to: '/pengaturan', label: 'Pengaturan', icon: Settings, permission: 'pengaturan' },
]

export function buildBreadcrumbs(path) {
    for (const section of allMenu) {
        for (const item of section.items) {
            if (item.to === path) {
                return [
                    { label: item.label, icon: item.icon },
                ]
            }
            if (item.children) {
                const child = item.children.find(c => c.to === path)
                if (child) {
                    return [
                        { label: item.label, icon: item.icon },
                        { label: child.label, icon: child.icon },
                    ]
                }
            }
        }
    }
    for (const item of bottomMenu) {
        if (item.to === path) {
            return [{ label: item.label, icon: item.icon }]
        }
    }
    return [{ label: path.replace('/', '') }]
}
