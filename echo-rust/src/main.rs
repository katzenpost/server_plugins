
extern crate protobuf;
extern crate grpc;
extern crate futures;
extern crate futures_cpupool;
extern crate tls_api;
extern crate tls_api_native_tls;

use std::thread;
use std::env;

pub mod kaetzchen;
pub mod kaetzchen_grpc;

use kaetzchen::{Request, Response};
use kaetzchen_grpc::{KaetzchenServer, Kaetzchen};

use tls_api::TlsAcceptorBuilder;


struct Echo;

impl Kaetzchen for Echo {
    fn on_request(&self, _m: grpc::RequestOptions, req: Request) -> grpc::SingleResponse<Response> {
        let mut r = Response::new();
        grpc::SingleResponse::completed(r)
    }
}

/*fn test_tls_acceptor() -> tls_api_native_tls::TlsAcceptor {
    let pkcs12 = include_bytes!("../foobar.com.p12");
    let builder = tls_api_native_tls::TlsAcceptorBuilder::from_pkcs12(pkcs12, "mypass").unwrap();
    builder.build().unwrap()
}*/

fn is_tls() -> bool {
    env::args().any(|a| a == "--tls")
}

fn main() {
    let tls = is_tls();

    let port = if !tls { 50051 } else { 50052 };
    
    let mut server = grpc::ServerBuilder::new();
    server.http.set_port(port);
    server.add_service(KaetzchenServer::new_service_def(Echo));
    server.http.set_cpu_pool_threads(4);
    if tls {
        server.http.set_tls(test_tls_acceptor());
    }
    let _server = server.build().expect("server");

    println!("greeter server started on port {} {}",
        port, if tls { "with tls" } else { "without tls" });

    loop {
        thread::park();
    }
}
