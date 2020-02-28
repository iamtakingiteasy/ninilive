import Vue from 'vue'
import VueRouter from 'vue-router'
import Live from '../views/Live.vue'
import store from '../store'

Vue.use(VueRouter);

const routes = [
  {
    path: '',
    alias: ['/live'],
    name: 'root',
    component: Live,
    meta: {
      handler() {
        store.dispatch('channels/selectNone');
      }
    }
  },
  {
    path: '/live/:channel',
    name: 'live',
    component: Live,
    meta: {
      handler({params: {channel}}) {
        store.dispatch('channels/selectName', {name: channel});
      }
    }
  },
];

const router = new VueRouter({
  mode: 'history',
  routes
});

router.beforeEach((to, from, next) => {
  to.matched.forEach(record => record.meta.handler && record.meta.handler(to, from));
  next();
});

export default router
