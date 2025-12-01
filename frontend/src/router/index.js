import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import InventoryList from '../views/InventoryList.vue'
import InventoryEdit from '../views/InventoryEdit.vue'
import InventoryScan from '../views/InventoryScan.vue'
import SalesList from '../views/SalesList.vue'
import SalesCreate from '../views/SalesCreate.vue'
import OrderChecker from '../views/OrderChecker.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/login',
            name: 'login',
            component: Login
        },
        {
            path: '/register',
            name: 'register',
            component: Register
        },
        {
            path: '/',
            redirect: '/dashboard'
        },
        {
            path: '/dashboard',
            name: 'dashboard',
            component: () => import('../views/Dashboard.vue'),
            meta: { requiresAuth: true }
        },
        {
            path: '/inventory',
            name: 'inventory-list',
            component: InventoryList,
            meta: { requiresAuth: true }
        },
        {
            path: '/inventory/add',
            name: 'inventory-add',
            component: InventoryEdit,
            meta: { requiresAuth: true }
        },
        {
            path: '/inventory/edit/:id',
            name: 'inventory-edit',
            component: InventoryEdit,
            meta: { requiresAuth: true }
        },
        {
            path: '/inventory/scan',
            name: 'inventory-scan',
            component: InventoryScan,
            meta: { requiresAuth: true }
        },
        {
            path: '/sales',
            name: 'sales-list',
            component: SalesList,
            meta: { requiresAuth: true }
        },
        {
            path: '/sales/create',
            name: 'sales-create',
            component: SalesCreate,
            meta: { requiresAuth: true }
        },
        {
            path: '/order-checker',
            name: 'order-checker',
            component: OrderChecker,
            meta: { requiresAuth: true }
        }
    ]
})

// Navigation guard
router.beforeEach((to, from, next) => {
    const isAuthenticated = localStorage.getItem('user')
    if (to.meta.requiresAuth && !isAuthenticated) {
        next('/login')
    } else {
        next()
    }
})

export default router
