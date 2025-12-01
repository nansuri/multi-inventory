<template>
  <div class="scan-page">
    <van-nav-bar
      title="Scan Barcode"
      left-text="Back"
      left-arrow
      @click-left="$router.back()"
    />
    
    <div class="scanner-wrapper">
      <BarcodeScanner @scan="onScan" @close="$router.back()" />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';
import BarcodeScanner from '../components/BarcodeScanner.vue';
import { showSuccessToast, showFailToast } from 'vant';
import { apiBase } from '../config/api';

const router = useRouter();

const onScan = async (barcode) => {
  console.log('Scanned barcode:', barcode);
  
  try {
    // Check if item exists in inventory
    const token = JSON.parse(localStorage.getItem('user'))?.token;
    const response = await fetch(`/api/inventory/barcode/${barcode}`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    
    if (response.ok) {
      const item = await response.json();
      showSuccessToast('Item found!');
      // Navigate to edit page for the found item
      router.push(`/inventory/edit/${item.id}`);
    } else {
      showFailToast('Item not found. Add new item?');
      // Navigate to add page with barcode pre-filled
      router.push({
        name: 'inventory-add',
        query: { barcode }
      });
    }
  } catch (error) {
    console.error('Error searching for item:', error);
    showFailToast('Error scanning barcode');
  }
};
</script>

<style scoped>
.scan-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.scanner-wrapper {
  padding: 0;
}
</style>
