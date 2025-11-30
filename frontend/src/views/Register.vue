<template>
  <div class="register-container">
    <div class="logo-area">
      <h1>Register</h1>
      <p>Create a new account</p>
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
        <van-field
          v-model="confirmPassword"
          type="password"
          name="ConfirmPassword"
          label="Confirm"
          placeholder="Confirm Password"
          :rules="[{ validator: validateConfirm, message: 'Passwords do not match' }]"
        />
      </van-cell-group>
      <div style="margin: 16px;">
        <van-button round block type="primary" native-type="submit" :loading="loading">
          Register
        </van-button>
      </div>
      <div style="margin: 16px; text-align: center;">
        <van-button size="small" type="default" to="/login">
          Back to Login
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { showToast, showSuccessToast, showFailToast } from 'vant';

const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const loading = ref(false);
const router = useRouter();

const validateConfirm = (val) => val === password.value;

const onSubmit = async (values) => {
  loading.value = true;
  try {
    const response = await fetch('http://localhost:8080/api/auth/register', {
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
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || 'Registration failed');
    }

    showSuccessToast('Registration successful');
    router.push('/login');
  } catch (error) {
    showFailToast(error.message);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.register-container {
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
