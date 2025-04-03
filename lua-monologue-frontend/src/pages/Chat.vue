<script setup>
import {ref} from 'vue';
import { sendChatMessage } from '@/api/ai';

const userMessage = ref('');
const messages = ref([]);

const sendMessage = async () => {
    if (!userMessage.value) return;
    console.log(userMessage.value);

    await sendChatMessage(userMessage.value);
    userMessage.value = '';
}
</script>

<template>
    <div>
        <h1>AI와 대화</h1>
        <div v-for="(msg, index) in messages" :key="index">
            <p><strong>{{  msg.role === 'user' ? '나' : 'AI' }}:</strong>{{  msg.text }}</p>
        </div>
        <input v-model="userMessage" placeholder="메시지를 입력하세요" @keyup.enter="sendMessage" />
        
        <button @click="sendMessage">보내기</button>
    </div>
</template>