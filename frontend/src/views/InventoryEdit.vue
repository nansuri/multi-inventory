<template>
  <div class="inventory-edit-container">
    <van-nav-bar
      :title="isEdit ? 'Edit Item' : 'Add Item'"
      left-text="Back"
      left-arrow
      @click-left="$router.back()"
    />

    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field
          v-model="form.name"
          name="Name"
          label="Name"
          placeholder="Item Name"
          :rules="[{ required: true, message: 'Name is required' }]"
        />
        <van-field
          v-model="form.barcode"
          name="Barcode"
          label="Barcode"
          placeholder="Scan or enter barcode"
          right-icon="scan"
          @click-right-icon="showScanner = true"
          :rules="[{ required: true, message: 'Barcode is required' }]"
        />
        <van-field
          v-model.number="form.price"
          type="number"
          name="Price"
          label="Price"
          placeholder="Price"
          :rules="[{ required: true, message: 'Price is required' }]"
        />
        <van-field
          v-model.number="form.quantity"
          type="digit"
          name="Quantity"
          label="Quantity"
          placeholder="Quantity"
          :rules="[{ required: true, message: 'Quantity is required' }]"
        />
        <van-field
          v-model="form.location"
          name="Location"
          label="Location"
          placeholder="Location"
        />
        <van-field name="is_halal" label="Halal Status">
          <template #input>
            <van-switch v-model="form.is_halal" />
          </template>
        </van-field>
      </van-cell-group>

      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          {{ isEdit ? 'Update' : 'Create' }}
        </van-button>
      </div>
      
      <div v-if="isEdit" style="margin: 16px;">
        <van-button round block type="danger" @click="onDelete">
          Delete Item
        </van-button>
      </div>
    </van-form>

    <van-popup v-model:show="showScanner" position="bottom" :style="{ height: '60%' }">
      <BarcodeScanner @scan="onScan" @close="showScanner = false" />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { showToast, showConfirmDialog } from 'vant';
import BarcodeScanner from '../components/BarcodeScanner.vue';

const route = useRoute();
const router = useRouter();
const isEdit = computed(() => route.params.id !== undefined);
const loading = ref(false);
const showScanner = ref(false);

const form = ref({
  name: '',
  barcode: '',
  price: '',
  quantity: '',
  location: '',
  is_halal: true,
});

const onScan = (code) => {
  form.value.barcode = code;
  showScanner.value = false;
  showToast.success('Barcode scanned');
};

const loadItem = async () => {
  if (!isEdit.value) return;
  try {
    const response = await fetch(`http://localhost:8080/api/inventory/${route.params.id}`);
    if (!response.ok) throw new Error('Failed to load item');
    const data = await response.json();
    form.value = data;
  } catch (error) {
    showToast.fail('Failed to load item');
  }
};

const onSubmit = async () => {
  loading.value = true;
  try {
    const url = isEdit.value 
      ? `http://localhost:8080/api/inventory/${route.params.id}`
      : 'http://localhost:8080/api/inventory';
    
    const method = isEdit.value ? 'PUT' : 'POST';

    const response = await fetch(url, {
      method: method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(form.value),
    });

    if (!response.ok) {
        const err = await response.text();
        throw new Error(err || 'Operation failed');
    }

    showToast.success(isEdit.value ? 'Item updated' : 'Item created');
    router.back();
  } catch (error) {
    showToast.fail(error.message);
  } finally {
    loading.value = false;
  }
};

const onDelete = () => {
    showConfirmDialog({
        title: 'Delete Item',
        message: 'Are you sure you want to delete this item?',
    })
    .then(async () => {
        try {
            const response = await fetch(`http://localhost:8080/api/inventory/${route.params.id}`, {
                method: 'DELETE',
            });
            if (!response.ok) throw new Error('Failed to delete');
            showToast.success('Item deleted');
            router.back();
        } catch (error) {
            showToast.fail('Failed to delete item');
        }
    })
    .catch(() => {
        // on cancel
    });
};

onMounted(() => {
  loadItem();
});
</script>

<style scoped>
.inventory-edit-container {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
