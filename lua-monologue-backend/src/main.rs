mod pb;

use pb::chat_service_server::{ChatService, ChatServiceServer};
use pb::{ChatRequest, ChatResponse};

use tonic::{transport::Server, Request, Response, Status};

use sqlx::postgres::PgPoolOptions;
use std::env;

#[derive(Debug, Default)]
pub struct LuaChatService {}

#[tonic::async_trait]
impl ChatService for LuaChatService {
    async fn send_message(
        &self,
        request: Request<ChatRequest>,
    ) -> Result<Response<ChatResponse>, Status> {
        let req_msg = &request.get_ref().content;
        let req_role = &request.get_ref().role;
        let req_user_id = &request.get_ref().id;
        println!("📨 Rust 서버에서 받은 메시지: msg: {} role: {} user_id: {}", req_msg, req_role, req_user_id);

        dotenv::dotenv().ok();

        let db_url = env::var("DATABASE_URL").expect("DATABASE_URL not set");

        // DB 풀 연결
        let pool = PgPoolOptions::new()
            .max_connections(5)
            .connect(&db_url)
            .await.unwrap();

        // 테스트용 INSERT
        let user_id = req_user_id;
        let role = req_role; // 또는 "assistant"
        let content = req_msg;

        sqlx::query!(
            "INSERT INTO messages (user_id, role, content) VALUES ($1, $2, $3)",
            user_id,
            role,
            content
        )
        .execute(&pool)
        .await.unwrap();

        let reply = format!("안녕! [{}] 잘 받았어", req_msg);
        Ok(Response::new(ChatResponse { reply }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    dotenv::dotenv().ok();

    let db_url = env::var("DATABASE_URL").expect("DATABASE_URL not set");

    println!("📦 연결 주소: {}", db_url);

    // DB 풀 연결
    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&db_url)
        .await?;

    /* 
    // 테스트용 INSERT
    let user_id = "user-test";
    let role = "user"; // 또는 "assistant"
    let content = "오늘 정말 지쳤다...";

    sqlx::query!(
        "INSERT INTO messages (user_id, role, content) VALUES ($1, $2, $3)",
        user_id,
        role,
        content
    )
    .execute(&pool)
    .await?;
    
    let role = "assistant"; // 또는 "assistant"
    let content = "고생 많았어!";

    sqlx::query!(
        "INSERT INTO messages (user_id, role, content) VALUES ($1, $2, $3)",
        user_id,
        role,
        content
    )
    .execute(&pool)
    .await?;

    println!("✅ 메시지 저장 완료!");
    */
    //let addr = "[::1]:50051".parse()?;
    let addr = "127.0.0.1:50051".parse()?;
    let svc = LuaChatService::default();

    println!("🚀 Rust gRPC 서버 실행중: {}", addr);

    Server::builder()
        .add_service(ChatServiceServer::new(svc))
        .serve(addr)
        .await?;

    Ok(())
}