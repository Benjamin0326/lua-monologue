from fastapi import FastAPI, Request
from pydantic import BaseModel
import httpx

app = FastAPI()

class PromptRequest(BaseModel):
    prompt: str


@app.post("/generate")
async def generate_response(req: PromptRequest):
    system_prompt =  "너는 감정에 공감 잘하고, 사용자와 친밀하게 대화하는 AI 친구야. 사용자의 일기 내용을 보고 따뜻하게 반응해줘. 그리고 한글로 대답해줘, Please always respond in Korean."
    async with httpx.AsyncClient(timeout=None) as client:
        response = await client.post(
            "http://localhost:11434/api/generate",
            json={
                "model": "llama3",
                "prompt": req.prompt,
                "system": system_prompt,
                "stream": False
            }
        )

        data = response.json()
        return {"response": data["response"]}
    