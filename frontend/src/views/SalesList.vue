<template>
  <div class="sales-list-container">
    <van-nav-bar
      title="Sales History"
      left-text="Back"
      left-arrow
      @click-left="$router.back()"
      right-text="New Sale"
      @click-right="$router.push('/sales/new')"
    />

    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="No more orders"
      @load="onLoad"
    >
      <van-cell
        v-for="order in list"
        :key="order.id"
        :title="`Order #${order.id}`"
        :label="formatDate(order.created_at)"
        :value="`$${order.total_price}`"
        is-link
        :to="`/sales/${order.id}`"
      >
        <template #right-icon>
            <van-tag :type="getStatusType(order.status)">{{ order.status }}</van-tag>
        </template>
      </van-cell>
    </van-list>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { showToast } from 'vant';
import { apiBase } from '../config/api';

const list = ref([]);
const loading = ref(false);
const finished = ref(false);

const formatDate = (dateStr) => {
    return new Date(dateStr).toLocaleString();
};

const getStatusType = (status) => {
    switch(status) {
        case 'completed': return 'success';
        case 'pending': return 'warning';
        case 'cancelled': return 'danger';
        default: return 'primary';
    }
};

const onLoad = async () => {
    if (list.value.length > 0) {
        loading.value = false;
        finished.value = true;
        return;
    }
    
    try {
        const response = await fetch(`${apiBase}/api/sales`);
        if (!response.ok) throw new Error('Failed to fetch orders');
        const data = await response.json();
        list.value = data || [];
    } catch (error) {
        showToast.fail('Failed to load sales history');
    } finally {
        loading.value = false;
        finished.value = true;
    }
};

onMounted(() => {
    loading.value = true;
    onLoad();
});
</script>

<style scoped>
.sales-list-container {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
