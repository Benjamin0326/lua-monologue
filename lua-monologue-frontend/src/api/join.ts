import api from './index';

export const sendJoinInfo = async(id: string, pw: string) => {
    console.log("Test: sendJoinInfo api called")
    console.log(id)
    console.log(pw)

    const response = await api.post('/sendjoininfo', { username: id, password: pw }, 
    { withCredentials: true  // ✅ 쿠키도 받아오도록
    })

    // access token 저장
    localStorage.setItem('access_token', response.data.accessToken)
    
    return response.data;
}