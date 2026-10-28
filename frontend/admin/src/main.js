import { createApp } from 'vue'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'

import App from './App.vue'
import router from './router'
import { pinia } from './store'
import directives from './directives'
import { setupDateTimeLocale } from './utils/datetime'
import './assets/styles/main.less'

setupDateTimeLocale()

const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(Antd)
app.use(directives)
app.mount('#app')
