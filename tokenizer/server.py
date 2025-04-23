#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import grpc
import config
from concurrent import futures
import pb.tokenizer_pb2
import pb.tokenizer_pb2_grpc
from tokenizer import Qwen2TokenizerSingleton


class TokenizerServicer(pb.tokenizer_pb2_grpc.TokenizerServiceServicer):
    def __init__(self):
        self.tokenizer = Qwen2TokenizerSingleton.get_tokenizer()

    def Tokenizer(self, request, context):
        try:
            encoding = self.tokenizer(request.prompt)
            return pb.tokenizer_pb2.TokenizerResp(
                input_ids=encoding.input_ids,
                attention_mask=encoding.attention_mask,
                current_tokens_cnt=len(encoding.input_ids)
            )

        except Exception as e:
            # 错误处理
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Tokenizer error: {str(e)}")
            return pb.tokenizer_pb2.TokenizerResp()

def serve():
    # 初始化分词器
    Qwen2TokenizerSingleton().initialize()

    # 启动服务
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=config.MAX_WORKS))
    pb.tokenizer_pb2_grpc.add_TokenizerServiceServicer_to_server(
        TokenizerServicer(), server
    )
    server.add_insecure_port(f'[::]:{config.SERVER_PORT}')
    print(f"Server started on port {config.SERVER_PORT}")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
