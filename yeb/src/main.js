import Vue from 'vue'
import App from './App.vue'
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import router from './router'
import * as echarts from 'echarts'

router.beforeEach((to,from,next)=>{
  console.log(to)
  console.log(from)
  next()
})
Vue.config.productionTip = false
Vue.prototype.$echarts = echarts
Vue.use(ElementUI);
new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
