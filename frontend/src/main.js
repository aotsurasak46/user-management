import './style.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Vue3Toastify from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'
import App from './App.vue'
import router from './router'
import PrimeVue from 'primevue/config';

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Vue3Toastify, {
    position: 'top-right',
    transition: 'Vue-Toastify__bounce'
})
app.use(PrimeVue);
app.mount('#app')
