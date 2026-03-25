import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/main.css'
import { useTheme } from './composables/useTheme'
import { i18n } from './i18n'
import { initLocale } from './composables/useLocale'

useTheme().initTheme()

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(i18n)
app.mount('#app')

// Initialize locale after app is mounted
initLocale(i18n)
