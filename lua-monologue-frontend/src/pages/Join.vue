<script setup>
import {ref} from 'vue';
import { sendJoinInfo } from '@/api/join';
import router from '@/router';

const userId = ref('')
const password = ref('')
const error = ref('')

const sendJoinInforeq = async () => {
    if (!userId.value) return;
    if (!password.value) return;

    try {
        let response = await sendJoinInfo(userId.value, password.value);
        console.log(response);
        error.value = response.data;
        router.push('/login')
    } catch (err) {
        console.error('회원가입입 실패', err)
    }
}

</script>

<template>
    <div class="login-container">
        <h2>회원가입</h2>
        <form @submit.prevent="login">
            <input v-model="userId" type="text" placeholder="ID" required />
            <input v-model="password" type="password" placeholder="Password" required />
            <button type="submit" @click="sendJoinInforeq">회원가입</button>
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
