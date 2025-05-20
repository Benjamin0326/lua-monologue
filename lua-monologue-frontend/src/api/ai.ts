import api from './index';

export const sendChatMessage = async(content: string) => {
    const response = await api.post('/sendchatmessage', { content });
    return response.data;
}

export const getChatMessages = async() => {
    const response = await api.get('/getchatmessages');
    return response.data;
}