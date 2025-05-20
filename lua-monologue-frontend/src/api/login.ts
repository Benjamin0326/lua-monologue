import api from './index';

export const sendLogInfo = async(id: string, pw: string) => {
    console.log("Test: sendLogInfo api called")
    console.log(id)
    console.log(pw)

    const response = await api.post('/sendloginfo', { username: id, password: pw }, 
    { withCredentials: true  // ✅ 쿠키도 받아오도록
    })

    // access token 저장
    localStorage.setItem('access_token', response.data.accessToken)
    console.log(response.data.accessToken)

    return response.data;
}

export const sendLogOut = async() => {
    console.log("Test: sendLogOut api called")
    const response = await api.post('/sendlogout')
    localStorage.removeItem('access_token')   // access token 제거

    return response.data;
}