fn main() {
    tonic_build::configure()
        .build_server(true)
        .out_dir(std::env::var("OUT_DIR").unwrap())  // 기본적으로 OUT_DIR에 생성
        .compile(&["proto/chat.proto"], &["proto"])
        .expect("Failed to compile proto files");
}
