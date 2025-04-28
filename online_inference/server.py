#!/usr/bin/env python
# -*- coding: UTF-8 -*-
from llm import StreamLLMInference
from vllm import SamplingParams
import grpc
import pb.model_server_pb2
from pb.model_server_pb2_grpc import ModelServerServiceServicer


class OnlineInferenceServerServicer(ModelServerServiceServicer):
  def __init__(self, model_path: str, max_tokens):
    self.llm = StreamLLMInference(model_path)
    self.max_tokens = max_tokens

  async def ChatCompletion(self, request, context):
    try:
      request_id = request.chat_session_id
      sampling_params = SamplingParams(
        presence_penalty=request.presence_penalty,
        frequency_penalty=request.frequency_penalty,
        repetition_penalty=request.repetition_penalty,
        temperature=request.temperature,
        top_p=request.top_p,
        top_k=request.top_k,
        max_tokens=self.max_tokens,
        detokenize=False,
      )

      generator = self.llm.stream_generation(request_id, sampling_params, request.tokens, request.token_type_ids)
      async for delta in generator:
        yield pb.model_server_pb2.ChatCompletionResult(token_ids=delta)
    except Exception as e:
      await context.abort(grpc.StatusCode.INTERNAL, f"生成错误: {str(e)}")