import { createRouter, createWebHistory } from 'vue-router';
import Home from '@/pages/Home.vue';
import Diary from '@/pages/Diary.vue';
import Chat from '@/pages/Chat.vue';
import Login from '@/pages/Login.vue';
import Join from '@/pages/Join.vue';

const routes = [
    //{ path: '/', component: Home},
    { path: '/', component: Chat},
    { path: '/login', component: Login},
    { path: '/join', component: Join},
    { path: '/diary', component: Diary},
    { path: '/chat', component: Chat,
        meta: {requiresAuth: true},
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, from, next) => {
    const isLoggedIn = !!localStorage.getItem('access_token')

    if (to.meta.requiresAuth && !isLoggedIn) {
        return next('login')
    }

    next()
})

export default router;