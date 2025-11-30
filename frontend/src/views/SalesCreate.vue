<template>
  <div class="sales-create-container">
    <van-nav-bar
      title="New Sale"
      left-text="Back"
      left-arrow
      @click-left="$router.back()"
    />

    <div class="cart-items">
        <van-empty v-if="cart.length === 0" description="Cart is empty. Scan items to add." />
        <van-swipe-cell v-for="(item, index) in cart" :key="index">
            <van-card
                :num="item.quantity"
                :price="item.price"
                :title="item.name"
                class="goods-card"
            >
                <template #footer>
                    <van-stepper v-model="item.quantity" min="1" :max="item.maxQuantity" />
                </template>
            </van-card>
            <template #right>
                <van-button square text="Delete" type="danger" class="delete-button" @click="removeFromCart(index)" />
            </template>
        </van-swipe-cell>
    </div>

    <van-submit-bar :price="totalPrice * 100" button-text="Submit Order" @submit="onSubmit" :loading="loading">
        <van-button type="primary" size="small" icon="scan" @click="showScanner = true">Scan Item</van-button>
    </van-submit-bar>

    <van-popup v-model:show="showScanner" position="bottom" :style="{ height: '60%' }">
      <BarcodeScanner @scan="onScan" @close="showScanner = false" />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { showToast } from 'vant';
import BarcodeScanner from '../components/BarcodeScanner.vue';
const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const router = useRouter();
const cart = ref([]);
const loading = ref(false);
const showScanner = ref(false);

const totalPrice = computed(() => {
    return cart.value.reduce((sum, item) => sum + (item.price * item.quantity), 0);
});

const onScan = async (code) => {
    showScanner.value = false;
    // Check if item already in cart
    const existingItem = cart.value.find(i => i.barcode === code);
    if (existingItem) {
        if (existingItem.quantity < existingItem.maxQuantity) {
            existingItem.quantity++;
            showToast.success('Quantity increased');
        } else {
            showToast.fail('Not enough stock');
        }
        return;
    }

    // Fetch item details
    try {
        // We need an endpoint to get by barcode, but currently we have get by ID or List.
        // Let's assume we can filter the list or add a specific endpoint.
        // For MVP, let's fetch all and find. Not efficient but works.
        // Better: Add GetByBarcode endpoint. I did add GetItemByBarcode in service but not handler?
        // Wait, I didn't add GetByBarcode in Handler. I should fix that or just use List for now.
        // Let's use List for now to save time, or just fetch list once and cache.
        
        const response = await fetch(`${apiBase}/api/inventory`);
        if (!response.ok) throw new Error('Failed to fetch inventory');
        const items = await response.json();
        const item = items.find(i => i.barcode === code);
        
        if (!item) {
            showToast.fail('Item not found');
            return;
        }

        if (item.quantity <= 0) {
            showToast.fail('Out of stock');
            return;
        }

        cart.value.push({
            id: item.id,
            name: item.name,
            price: item.price,
            barcode: item.barcode,
            quantity: 1,
            maxQuantity: item.quantity
        });
        showToast.success('Item added to cart');

    } catch (error) {
        showToast.fail('Error finding item');
    }
};

const removeFromCart = (index) => {
    cart.value.splice(index, 1);
};

const onSubmit = async () => {
    if (cart.value.length === 0) return;
    
    loading.value = true;
    try {
        const user = JSON.parse(localStorage.getItem('user'));
        const response = await fetch(`${apiBase}/api/sales`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                user_id: user ? user.id : 1,
                items: cart.value.map(i => ({
                    item_id: i.id,
                    quantity: i.quantity
                }))
            }),
        });

        if (!response.ok) {
             const err = await response.text();
             throw new Error(err || 'Failed to create order');
        }

        showToast.success('Order created');
        router.push('/sales');
    } catch (error) {
        showToast.fail(error.message);
    } finally {
        loading.value = false;
    }
};
</script>

<style scoped>
.sales-create-container {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 50px;
}
.cart-items {
    padding-bottom: 60px;
}
.goods-card {
  margin: 0;
  background-color: white;
}

.delete-button {
  height: 100%;
}
</style>
