from fastapi import FastAPI, Request
from pydantic import BaseModel
import httpx
import psycopg2
import httpx
import json
import numpy as np
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

app = FastAPI()
model = SentenceTransformer("hkunlp/instructor-xl")

async def update_demo_user_embeddings():
    conn = psycopg2.connect(
        host="localhost",
        port=5432,
        dbname="journal",
        user="myuser",
        password="mypass"
    )

    cur = conn.cursor()

    cur.execute("""
        SELECT id, content
        FROM messages
        WHERE role = 'user'
        AND embedding IS NULL OR embedding_cls IS NULL
    """, )

    rows = cur.fetchall()
    print(f"🔍 {len(rows)}개 메시지에 대해 embedding 생성 중...")

    print("📌 가져온 메시지:")
    for row in rows:
        msg_id, content = row
        print(f"[{msg_id}] {content}")
        input_text = ["Represent the semantic meaning of the sentence", content]

        token_embeddings = model.encode([input_text], output_value='token_embeddings')
        token_embeddings = token_embeddings[0]

        cls_embedding = token_embeddings[0].cpu().numpy()
        cls_str = f"[{','.join(map(str, cls_embedding))}]"

        mean_embedding = np.mean(token_embeddings.cpu().numpy(), axis=0)
        mean_str = f"[{','.join(map(str, mean_embedding))}]"

        cur.execute("""
            UPDATE messages
            SET embedding = %s,
                embedding_cls = %s
            WHERE id = %s         
        """, (mean_str, cls_str, msg_id))
        print(f"✅ [{msg_id}] embedding 저장 완료")

    conn.commit()
    cur.close()
    conn.close()

async def test_main():
    query_text = "요즘 너무 무기력하고 힘들어."

    input_text = ["Represent the semantic meaning of the sentence", query_text]

    query_token_embeddings = model.encode([input_text], output_value="token_embeddings")
    query_token_embeddings = query_token_embeddings[0]

    query_mean_embedding = query_token_embeddings.mean(dim=0).cpu().numpy()
    query_cls_embedding = query_token_embeddings[0].cpu().numpy()

    conn = psycopg2.connect(
        host="localhost",
        port=5432,
        dbname="journal",
        user="myuser",
        password="mypass"
    )

    cur = conn.cursor()

    query_vec_str = f"[{','.join(map(str, query_mean_embedding))}]"
    cur.execute("""
        SELECT id, content, embedding, created_at
        FROM messages
        WHERE user_id = %s
        ORDER BY embedding <-> %s
        LIMIT 5
    """, ("demo-user", query_vec_str))

    mean_results = [(id, content, emb, created_at, "mean") for id, content, emb, created_at in cur.fetchall()]

    query_cls_vec_str = f"[{','.join(map(str, query_cls_embedding))}]"
    cur.execute("""
        SELECT id, content, embedding_cls, created_at
        FROM messages
        WHERE user_id = %s
        ORDER BY embedding_cls <-> %s
        LIMIT 5
    """, ("demo-user", query_cls_vec_str))

    cls_results = [(id, content, emb, created_at, "cls") for id, content, emb, created_at in cur.fetchall()]

    cur.close()
    conn.close()

    candidates = []

    for id, content, emb_vector, created_at, source in mean_results + cls_results:
        candidates.append((id, content, emb_vector, created_at, source))

    rerank_scores = []

    for id, content, emb_vector, created_at, source in candidates:
        emb_np = np.array(eval(emb_vector))

        if source == 'mean':
            score = cosine_similarity([query_mean_embedding], [emb_np])[0][0]
        if source == 'cls':
            score = cosine_similarity([query_cls_embedding], [emb_np])[0][0]
        
        rerank_scores.append((id, content, score, created_at, source))

    rerank_scores.sort(key=lambda x: x[2], reverse=True)
    top_contexts = [(content, created_at) for _, content, _, created_at, _ in rerank_scores[:3]]
    print(rerank_scores)
    print(top_contexts)

    context_block = "\n\n".join(f"Context {i+1} ({d.strftime('%Y-%m-%d')}): {c}" for i, (c, d) in enumerate(top_contexts))

    final_prompt = f"""
    YOU are a helpful AI assistant.
    Based on the following past diary entries, answer the user's question.

    {context_block}

    User's Question:
    {query_text}

    그리고 한글로 대답해줘, Please always respond in Korean.
    
    Your Answer:
    """

    system_prompt =  "너는 감정에 공감 잘하고, 사용자와 친밀하게 대화하는 AI 친구야. 사용자의 일기 내용을 보고 따뜻하게 반응해줘. 그리고 한글로 대답해줘, Please always respond in Korean."
    async with httpx.AsyncClient(timeout=None) as client:
        response = await client.post(
            "http://localhost:11434/api/generate",
            json={
                "model": "llama3",
                "prompt": final_prompt,
                "system": system_prompt,
                "stream": False
            }
        )

        data = response.json()
        print("최종 응답: " + data["response"])


@app.on_event("startup")
async def startup_event():
    await update_demo_user_embeddings()
    #await test_main()


class PromptRequest(BaseModel):
    id: str
    prompt: str

@app.post("/generate")
async def generate_response(req: PromptRequest):
    req_id = req.id
    req_prompt = req.prompt
    print("🚨 RAW BODY:", req)
    input_text = ["Represent the semantic meaning of the sentence", req_prompt]

    query_token_embeddings = model.encode([input_text], output_value="token_embeddings")
    query_token_embeddings = query_token_embeddings[0]

    query_mean_embedding = query_token_embeddings.mean(dim=0).cpu().numpy()
    query_cls_embedding = query_token_embeddings[0].cpu().numpy()

    conn = psycopg2.connect(
        host="localhost",
        port=5432,
        dbname="journal",
        user="myuser",
        password="mypass"
    )

    cur = conn.cursor()

    query_vec_str = f"[{','.join(map(str, query_mean_embedding))}]"
    cur.execute("""
        SELECT id, content, embedding, created_at
        FROM messages
        WHERE user_id = %s
        ORDER BY embedding <-> %s
        LIMIT 5
    """, (req_id, query_vec_str))

    mean_results = [(id, content, emb, created_at, "mean") for id, content, emb, created_at in cur.fetchall()]

    query_cls_vec_str = f"[{','.join(map(str, query_cls_embedding))}]"
    cur.execute("""
        SELECT id, content, embedding_cls, created_at
        FROM messages
        WHERE user_id = %s
        ORDER BY embedding_cls <-> %s
        LIMIT 5
    """, (req_id, query_cls_vec_str))

    cls_results = [(id, content, emb, created_at, "cls") for id, content, emb, created_at in cur.fetchall()]

    cur.close()
    conn.close()

    candidates = []

    for id, content, emb_vector, created_at, source in mean_results + cls_results:
        candidates.append((id, content, emb_vector, created_at, source))

    rerank_scores = []

    for id, content, emb_vector, created_at, source in candidates:
        if emb_vector is None:
            print("emb_vector is not yet completed")
            continue
        
        emb_np = np.array(eval(emb_vector))

        if source == 'mean':
            score = cosine_similarity([query_mean_embedding], [emb_np])[0][0]
        if source == 'cls':
            score = cosine_similarity([query_cls_embedding], [emb_np])[0][0]
        
        rerank_scores.append((id, content, score, created_at, source))

    rerank_scores.sort(key=lambda x: x[2], reverse=True)
    top_contexts = [(content, created_at) for _, content, _, created_at, _ in rerank_scores[:3]]
    print(rerank_scores)
    print(top_contexts)

    context_block = "\n\n".join(f"Context {i+1} ({d.strftime('%Y-%m-%d')}): {c}" for i, (c, d) in enumerate(top_contexts))

    final_prompt = f"""
    YOU are a helpful AI assistant.
    Based on the following past diary entries, answer the user's question.

    {context_block}

    User's Question:
    {req_prompt}

    그리고 한글로 대답해줘, Please always respond in Korean.
    
    Your Answer:
    """

    print(final_prompt)

    system_prompt =  "너는 감정에 공감 잘하고, 사용자와 친밀하게 대화하는 AI 친구야. 사용자의 일기 내용을 보고 따뜻하게 반응해줘. 그리고 한글로 대답해줘, Please always respond in Korean."
    async with httpx.AsyncClient(timeout=None) as client:
        response = await client.post(
            "http://localhost:11434/api/generate",
            json={
                "model": "llama3",
                "prompt": final_prompt,
                "system": system_prompt,
                "stream": False
            }
        )

        data = response.json()
        print("최종 응답: " + data["response"])
        await update_demo_user_embeddings()
        return {"response": data["response"]}
    