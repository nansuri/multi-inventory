<template>
  <div class="order-checker-container">
    <van-nav-bar
      title="Order Details"
      left-text="Back"
      left-arrow
      @click-left="$router.back()"
    />

    <div v-if="order" class="order-info">
        <van-cell-group inset>
            <van-cell title="Order ID" :value="order.id" />
            <van-cell title="Status" :value="order.status" />
            <van-cell title="Total Price" :value="`$${order.total_price}`" />
            <van-cell title="Date" :value="new Date(order.created_at).toLocaleString()" />
        </van-cell-group>

        <div style="margin: 16px;">
            <h3>Items (Checker Mode)</h3>
            <p class="hint">Tap checkbox to mark as fulfilled</p>
        </div>

        <van-checkbox-group v-model="fulfilledItems">
            <van-cell-group inset>
                <van-cell
                    v-for="item in order.items"
                    :key="item.id"
                    :title="item.item_name"
                    :label="`Qty: ${item.quantity}`"
                    clickable
                    @click="toggle(item)"
                >
                    <template #right-icon>
                        <van-checkbox :name="item.id" @click.stop="toggle(item)" />
                    </template>
                </van-cell>
            </van-cell-group>
        </van-checkbox-group>
    </div>
    <van-loading v-else size="24px" vertical>Loading...</van-loading>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { showToast } from 'vant';

const route = useRoute();
const order = ref(null);
const fulfilledItems = ref([]);

const loadOrder = async () => {
    try {
        const response = await fetch(`http://localhost:8080/api/sales/${route.params.id}`);
        if (!response.ok) throw new Error('Failed to load order');
        const data = await response.json();
        order.value = data;
        
        // Set initial checkboxes
        fulfilledItems.value = data.items
            .filter(i => i.is_fulfilled)
            .map(i => i.id);
            
    } catch (error) {
        showToast.fail('Failed to load order');
    }
};

const toggle = async (item) => {
    const isNowFulfilled = !fulfilledItems.value.includes(item.id);
    // Optimistic update
    if (isNowFulfilled) {
        fulfilledItems.value.push(item.id);
    } else {
        fulfilledItems.value = fulfilledItems.value.filter(id => id !== item.id);
    }

    try {
        const response = await fetch(`http://localhost:8080/api/sales/items/${item.id}/fulfillment`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ is_fulfilled: isNowFulfilled }),
        });
        
        if (!response.ok) throw new Error('Update failed');
        showToast.success(isNowFulfilled ? 'Marked fulfilled' : 'Marked unfulfilled');
    } catch (error) {
        // Revert
        if (isNowFulfilled) {
            fulfilledItems.value = fulfilledItems.value.filter(id => id !== item.id);
        } else {
            fulfilledItems.value.push(item.id);
        }
        showToast.fail('Failed to update status');
    }
};

onMounted(() => {
    loadOrder();
});
</script>

<style scoped>
.order-checker-container {
  min-height: 100vh;
  background-color: #f7f8fa;
}
.order-info {
    padding-top: 16px;
}
.hint {
    font-size: 12px;
    color: #969799;
}
</style>
