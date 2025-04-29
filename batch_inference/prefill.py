#!/usr/bin/env python
# -*- coding: UTF-8 -*-
from loguru import logger
from multiprocessing import Queue
from typing import List
from config import VLLM_CONFIG
from common import SamplingParams
from vllm import LLM
from vllm.config import KVTransferConfig


class Prefill(object):
  def __init__(self, model: str):
    kv_conf_dict = {
      "kv_connector": VLLM_CONFIG["kv_cache"]["kv_connector"],
      "kv_role": "kv_producer",
      "kv_rank": VLLM_CONFIG["kv_cache"]["kv_rank"],
      "kv_parallel_size": VLLM_CONFIG["kv_cache"]["kv_parallel_size"]
    }
    ktc = KVTransferConfig.from_cli(str(kv_conf_dict))
    self.llm = LLM(model=model,
                   kv_transfer_config=ktc,
                   max_model_len=VLLM_CONFIG["max_model_len"],
                   gpu_memory_utilization=VLLM_CONFIG["gpu_memory_utilization"])

  def do_prefill(self, sampling_params: SamplingParams, prompts: List[str]) -> None:
    self.llm.generate(prompts, sampling_params.to_vllm())


def run(input_q: Queue, output_q: Queue) -> None:
  import os
  os.environ["CUDA_VISIBLE_DEVICES"] = "0"
  p = Prefill(model=VLLM_CONFIG["model"])
  while True:
    param = input_q.get(True)
    if param is None:
      continue
    if param["sampling_params"] is None or len(param["prompts"]) == 0:
      continue

    trace_id = "unknown"
    try:
      trace_id = param["trace_id"]
      p.do_prefill(sampling_params=param["sampling_params"], prompts=param["prompts"])
      output_q.put(param)
    except Exception as e:
      logger.error(f'[{trace_id}] Batch Inference Prefill Error: {e}')