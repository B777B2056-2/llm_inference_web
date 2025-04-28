#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import asyncio
from grpc import aio
from pb.model_server_pb2_grpc import add_ModelServerServiceServicer_to_server
from server import OnlineInferenceServerServicer


async def serve(model_path: str, port: int, max_tokens: int) -> None:
    server = aio.server()
    add_ModelServerServiceServicer_to_server(OnlineInferenceServerServicer(model_path, max_tokens), server)
    server.add_insecure_port(f'[::]:{port}')
    await server.start()
    print(f"服务已启动，监听端口 {port}")
    try:
        await server.wait_for_termination()
    except asyncio.CancelledError:
        await server.stop(5)


if __name__ == '__main__':
    import config
    asyncio.run(serve(model_path=config.MODEL, port=config.SERVER_PORT, max_tokens=config.MAX_TOKENS))