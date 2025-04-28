#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import json
from dataclasses import dataclass
import vllm


@dataclass
class SamplingParams:
  presence_penalty: float = 0.1
  frequency_penalty: float = 0.1
  repetition_penalty: float = 0.1
  temperature: float = 0.7
  top_p: float = 0.1
  top_k: int = 2

  @classmethod
  def from_dict(cls, d: dict) -> "SamplingParams":
    return cls(**d)

  def to_vllm(self) -> vllm.SamplingParams:
    return vllm.SamplingParams(
      presence_penalty=self.presence_penalty,
      frequency_penalty=self.frequency_penalty,
      repetition_penalty=self.repetition_penalty,
      temperature=self.temperature,
      top_p=self.top_p,
      top_k=self.top_k,
    )
