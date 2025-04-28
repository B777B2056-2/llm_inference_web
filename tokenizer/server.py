#!/usr/bin/env python
# -*- coding: UTF-8 -*-
from logger import setup_logging
from loguru import logger
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
            input_ids = encoding.input_ids
            token_type_ids = []
            if "token_type_ids" in encoding.__dict__:
                token_type_ids = encoding.token_type_ids
            return pb.tokenizer_pb2.TokenizerResp(
                token_ids=input_ids,
                token_type_ids=token_type_ids,
            )

        except Exception as e:
            logger.error(f"[{request.trace_id}] TokenizerServicer.Tokenizer error: {e}", exc_info=True)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Tokenizer error: {str(e)}")
            return pb.tokenizer_pb2.TokenizerResp()

    def DeTokenizer(self, request, context):
        try:
            text = self.tokenizer.decode(request.token_ids)
            return pb.tokenizer_pb2.DeTokenizerResult(text=text)
        except Exception as e:
            logger.error(f"[{request.trace_id}] TokenizerServicer.DeTokenizer error: {e}", exc_info=True)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Tokenizer error: {str(e)}")
            return pb.tokenizer_pb2.DeTokenizerResult()

def serve():
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
    setup_logging() # 初始化日志
    try:
        Qwen2TokenizerSingleton().initialize()  # 初始化分词器
        serve() # 开始服务
    except Exception as e:
        logger.critical(f"TokenizerServicer fatal error: {e}", exc_info=True)
        raise
