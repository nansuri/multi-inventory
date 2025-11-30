import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// Import Vant styles (optional if using unplugin-vue-components but good for base styles)
import 'vant/lib/index.css';

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
