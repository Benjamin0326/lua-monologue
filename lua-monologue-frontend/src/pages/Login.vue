<script setup>
import {ref} from 'vue';
import { sendLogInfo } from '@/api/login';
import { useRouter } from 'vue-router'
import router from '@/router';

const userId = ref('')
const password = ref('')
const error = ref('')

const sendLoginInfo = async () => {
    if (!userId.value) return;
    if (!password.value) return;

    try {
        let response = await sendLogInfo(userId.value, password.value);
        console.log(response);
        error.value = response.data;
        router.push('/chat')
    } catch (err) {
        console.error('로그인 실패', err)
    }
}

const gotoJoin = async () => {
    router.push('/join')
}

</script>

<template>
    <div class="login-container">
        <h2>로그인</h2>
        <form @submit.prevent="login">
            <input v-model="userId" type="text" placeholder="ID" required />
            <input v-model="password" type="password" placeholder="Password" required />
            <button type="submit" @click="sendLoginInfo">로그인</button>
            <button type="submit" @click="gotoJoin">회원가입</button>
            <p v-if="error">{{ error }}</p>
        </form>
    </div>
</template>

<style scoped>
.login-container {
    max-width: 300px;
    margin: auto;
}
input {
    display: block;
    margin-bottom: 10px;
    width: 100%;
}
</style>
