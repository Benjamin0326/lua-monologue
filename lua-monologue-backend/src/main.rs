mod pb;

use pb::chat_service_server::{ChatService, ChatServiceServer};
use pb::{ChatRequest, ChatResponse};

use tonic::{transport::Server, Request, Response, Status};

#[derive(Debug, Default)]
pub struct LuaChatService {}

#[tonic::async_trait]
impl ChatService for LuaChatService {
    async fn send_message(
        &self,
        request: Request<ChatRequest>,
    ) -> Result<Response<ChatResponse>, Status> {
        let msg = request.into_inner().content;
        println!("📨 Rust 서버에서 받은 메시지: {}", msg);

        let reply = format!("안녕! [{}] 잘 받았어", msg);
        Ok(Response::new(ChatResponse { reply }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let svc = LuaChatService::default();

    println!("🚀 Rust gRPC 서버 실행중: {}", addr);

    Server::builder()
        .add_service(ChatServiceServer::new(svc))
        .serve(addr)
        .await?;

    Ok(())
}