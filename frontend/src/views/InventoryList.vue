<template>
  <div class="inventory-list-container">
    <van-nav-bar
      title="Inventory"
      left-text="Back"
      left-arrow
      @click-left="$router.back()"
      right-text="Add"
      @click-right="$router.push('/inventory/add')"
    />

    <van-search v-model="search" placeholder="Search items..." />

    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="No more items"
      @load="onLoad"
    >
      <van-cell
        v-for="item in filteredItems"
        :key="item.id"
        :title="item.name"
        :label="`Price: $${item.price} | Qty: ${item.quantity}`"
        is-link
        :to="`/inventory/edit/${item.id}`"
      >
        <template #value>
            <van-tag type="primary" v-if="item.is_halal">Halal</van-tag>
            <van-tag type="danger" v-else>Non-Halal</van-tag>
        </template>
      </van-cell>
    </van-list>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { showToast } from 'vant';
import { apiBase } from '../config/api';

const list = ref([]);
const loading = ref(false);
const finished = ref(false);
const search = ref('');

const filteredItems = computed(() => {
  if (!search.value) return list.value;
  const lower = search.value.toLowerCase();
  return list.value.filter(item => 
    item.name.toLowerCase().includes(lower) || 
    item.barcode.includes(lower)
  );
});

const onLoad = async () => {
    // In a real app, we would support pagination.
    // Here we just load all items once.
    if (list.value.length > 0) {
        loading.value = false;
        finished.value = true;
        return;
    }
    
    try {
        const response = await fetch(`${apiBase}/api/inventory`);
        if (!response.ok) throw new Error('Failed to fetch items');
        const data = await response.json();
        list.value = data || [];
    } catch (error) {
        showToast.fail('Failed to load inventory');
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
.inventory-list-container {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
