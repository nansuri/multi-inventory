<template>
  <div class="dashboard-container">
    <van-nav-bar title="Dashboard" />
    
    <div class="content">
      <van-grid :column-num="2" :gutter="10">
        <van-grid-item icon="qr" text="Scan Item" to="/inventory/scan" />
        <van-grid-item icon="records" text="Inventory" to="/inventory" />
        <van-grid-item icon="cart-o" text="New Sale" to="/sales/new" />
        <van-grid-item icon="todo-list-o" text="Sales History" to="/sales" />
      </van-grid>

      <div style="margin-top: 20px; padding: 0 16px;">
        <h3>Recent Activity</h3>
        <van-empty description="No recent activity" />
      </div>

      <div style="margin: 20px 16px;">
        <van-button block type="danger" @click="logout">Logout</van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const logout = () => {
  localStorage.removeItem('user');
  router.push('/login');
};
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background-color: #f7f8fa;
}
.content {
  padding-top: 20px;
}
</style>
