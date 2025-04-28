#!/usr/bin/env python
# -*- coding: UTF-8 -*-
from multiprocessing import Queue
from typing import List, Tuple
from config import VLLM_CONFIG
from common import SamplingParams
from db import BatchInferenceDatabase
from vllm import LLM
from vllm.config import KVTransferConfig


class Decoder(object):
  def __init__(self, model: str):
    kv_conf_dict = {
      "kv_connector": VLLM_CONFIG["kv_connector"],
      "kv_role": "kv_consumer",
      "kv_rank": VLLM_CONFIG["kv_rank"],
      "kv_parallel_size": VLLM_CONFIG["kv_parallel_size"]
    }
    ktc = KVTransferConfig.from_cli(str(kv_conf_dict))
    self.llm = LLM(model=model,
                   kv_transfer_config=ktc,
                   max_model_len=VLLM_CONFIG["max_model_len"],
                   gpu_memory_utilization=VLLM_CONFIG["gpu_memory_utilization"])

  def do_decoder(self, sampling_params: SamplingParams, prompts: List[str]) -> List[Tuple[str, str]]:
    results = []
    outputs = self.llm.generate(prompts, sampling_params.to_vllm())
    for output in outputs:
      prompt = output.prompt
      generated_text = output.outputs[0].text
      results.append((prompt, generated_text))
    return results


def run(q: Queue) -> None:
  import os
  os.environ["CUDA_VISIBLE_DEVICES"] = "1"
  database = BatchInferenceDatabase()
  d = Decoder(model=VLLM_CONFIG["model"])
  while True:
    param = q.get(True)
    if param is None:
      continue
    if param["sampling_params"] is None or len(param["prompts"]) == 0:
      continue
    results = d.do_decoder(sampling_params=param["sampling_params"], prompts=param["prompts"])
    # 存入数据库
    database.insert_results(id=param["id"], user_id=param["user_id"], results=results)
