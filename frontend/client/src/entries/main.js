import Vue from 'vue'
import Main from '@/apps/Main.vue'
import router from '@/routers/main/router'
import store from '@/stores/main/store'
import '@/components/_global'
import '@/components/main/home'

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(Main)
}).$mount('#app')
