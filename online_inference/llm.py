#!/usr/bin/env python
# -*- coding: UTF-8 -*-
from typing import List
from vllm.engine.async_llm_engine import AsyncLLMEngine
from vllm import SamplingParams, AsyncEngineArgs, TokensPrompt


class StreamLLMInference(object):
  def __init__(self, model:str):
    # 初始化引擎参数
    engine_args = AsyncEngineArgs(model=model, skip_tokenizer_init=True, tokenizer=None)
    # 创建异步引擎实例
    self.engine = AsyncLLMEngine.from_engine_args(engine_args)

  async def stream_generation(self, request_id: str, sampling_params: SamplingParams, input_ids:List[int],
                              token_type_ids: List[int]=None):
    """从引擎流式获取生成结果的异步生成器"""
    output_length = 0
    tokens = TokensPrompt(prompt_token_ids=input_ids)
    if token_type_ids is not None and len(token_type_ids) != 0:
      tokens.token_type_ids = token_type_ids
    generator = self.engine.generate(prompt=tokens, sampling_params=sampling_params, request_id=request_id)
    async for request_output in generator:
      current_token_ids = request_output.outputs[-1].token_ids
      delta = current_token_ids[output_length:]
      yield delta
      output_length = len(current_token_ids)
      if request_output.finished:
        break
