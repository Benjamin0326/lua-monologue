<script setup>
import {onMounted, ref} from 'vue';
import { sendChatMessage } from '@/api/ai';
import { getChatMessages } from '@/api/ai';
import { sendLogOut } from '@/api/login';
import router from '@/router';
import api from '@/api';

const userMessage = ref('');
const messages = ref([]);
//const responseMessage = ref('');

const sendMessage = async () => {
    if (!userMessage.value) return;

    console.log(userMessage.value);

    const content = userMessage.value.trim()

    const userMsg = {
        id: Date.now(), // 임시 id
        role: 'user',
        content,
        created_at: new Date().toISOString()
    }
    messages.value.push(userMsg)

    let response = await sendChatMessage(userMessage.value);

    const reply = response.data

    // ✅ 3. assistant 응답 추가
    const aiMsg = {
      id: Date.now() + 1, // 임시 id
      role: 'assistant',
      content: reply,
      created_at: new Date().toISOString()
    }

    messages.value.push(aiMsg)

    userMessage.value = '';
    console.log(response);
    //responseMessage.value = response.data;
    //onMounted()
}

onMounted(async () => {
    const res = await getChatMessages()
    messages.value = res ?? []
    console.log(messages.value)
})

const sendLogout = async() => {
    try {
        let response = await sendLogOut();
        console.log(response);
        router.push('/login')
    } catch (err) {
        console.error('로그인 실패', err)
    }
    
}
</script>

<template>
    <div>
        <h1>AI와 대화</h1>
        <!--<div v-for="(msg, index) in messages" :key="index">
            <p><strong>{{  msg.role === 'user' ? '나' : 'AI' }}:</strong>{{  msg.content }}</p>
        </div>-->
        <div v-for="msg in messages" :key="msg.id" class="mb-2">
            <div
                :class="msg.role === 'user' ? 'text-right' : 'text-left'"
            >
                <div :class="[
                    'inline-block px-4 py-2 rounded-lg',
                    msg.role === 'user' ? 'bg-blue-500 text-white' : 'bg-gray-200 text-black'
                ]">
                {{ msg.content }}
                </div>
            </div>
        </div>
        <textarea v-model="userMessage" placeholder="메시지를 입력하세요" @keyup.enter="sendMessage" />
        <button @click="sendMessage">보내기</button>

        <!--<textarea v-model="responseMessage" rows="5" cols="40"></textarea>-->

        <button @click="sendLogout">로그아웃</button>
    </div>
</template>


<style scoped>
input {
    display: block;
    margin-bottom: 10px;
    width: 100%;
}

button {
    display: block;
    margin-bottom: 10px;
    width: 100%;
}

textarea {
    display: block;
    margin-bottom: 10px;
    width: 100%;
}
</style>