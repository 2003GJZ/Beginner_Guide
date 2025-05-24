import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import axios from 'axios'

// 使用ElementUI
Vue.use(ElementUI)

// 配置axios
// 使用相对路径，通过vue.config.js中的代理配置转发到后端
axios.defaults.baseURL = ''
console.log('前端已连接到后端服务器(通过代理)')
// 请求拦截器，添加token到请求头
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
// 响应拦截器，处理未授权错误
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response && error.response.status === 401) {
      // 清除本地存储的token
      localStorage.removeItem('token')
      // 跳转到登录页
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

Vue.prototype.$http = axios

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
