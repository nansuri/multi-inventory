<template>
  <div class="login-container">
    <div class="logo-area">
      <h1>Multi Inventory</h1>
      <p>Manage your inventory with ease</p>
    </div>

    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field
          v-model="username"
          name="Username"
          label="Username"
          placeholder="Username"
          :rules="[{ required: true, message: 'Username is required' }]"
        />
        <van-field
          v-model="password"
          type="password"
          name="Password"
          label="Password"
          placeholder="Password"
          :rules="[{ required: true, message: 'Password is required' }]"
        />
      </van-cell-group>
      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          Login
        </van-button>
      </div>
      <div style="margin: 16px; text-align: center;">
        <van-button size="small" type="default" to="/register">
          Create Account
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { showToast, showSuccessToast, showFailToast } from 'vant';
const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

const username = ref('');
const password = ref('');
const loading = ref(false);
const router = useRouter();

const onSubmit = async (values) => {
  loading.value = true;
  try {
    const response = await fetch(`${apiBase}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
      }),
    });

    if (!response.ok) {
      throw new Error('Login failed');
    }

    const data = await response.json();
    localStorage.setItem('user', JSON.stringify(data));
    showSuccessToast('Login successful');
    router.push('/dashboard');
  } catch (error) {
    showFailToast('Invalid credentials');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background-color: #f7f8fa;
}

.logo-area {
  text-align: center;
  margin-bottom: 40px;
}

.logo-area h1 {
  color: #1989fa;
  margin-bottom: 8px;
}

.logo-area p {
  color: #969799;
  font-size: 14px;
}
</style>
