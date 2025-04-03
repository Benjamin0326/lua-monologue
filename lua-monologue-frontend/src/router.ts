import { createRouter, createWebHistory } from 'vue-router';
import Home from '@/pages/Home.vue';
import Diary from '@/pages/Diary.vue';
import Chat from '@/pages/Chat.vue';

const routes = [
    { path: '/', component: Home},
    { path: '/diary', component: Diary},
    { path: '/chat', component: Chat},
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;