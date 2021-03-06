import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '../views/Login.vue'
import Home from '../views/Home.vue'
import ImageManager from '../views/servicePublish/ImageManager.vue'
import ServicePublish from '../views/servicePublish/ServicePublish.vue'
import ConfigCenter from '../views/servicePublish/ConfigCenter.vue'
import ServiceManager from '../views/serviceManager/ServiceManager.vue'
import Test from '../views/admin/Manager.vue'
import ServiceTrace from '../views/serviceManager/ServiceTrace.vue'
Vue.use(VueRouter)

const routes = [
  {
    path:'/da',
    name:'系统管理',
    component: ServiceManager,
    children:[
      
    ]
  },
  {
    path: '/',
    name: 'Login',
    component: Login,
    hidden: true
  },
  {
    path: '/',
    name: '服务发布',
    component: Home,
    children:[
    {
      path: '/img',
      name: '镜像管理',
      component: ImageManager
    },
    {
      path: '/svc',
      name: '服务发布',
      component: ServicePublish
    },
    {
      path:'/cfg',
      name: '配置中心',
      component: ConfigCenter
    },
    {
      path:'/test',
      name: 'Test',
      component: Test
    }
    ]
  },
  {
    path:'/',
    name:'服务治理',
    component: Home,
    children:[
      {
        path:'/trace',
        name: '服务链路追踪',
        component: ServiceTrace
      },
    ]
  },
  {
    path:'',
    name:'自动运维',
    component: ServiceManager,
    children:[
      
    ]
  }
]

const router = new VueRouter({
  routes
})

export default router
