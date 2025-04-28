#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import json
from kafka import KafkaConsumer
from multiprocessing import Process, Queue
from common import SamplingParams
import prefill
import decoder


def consume_kafka_messages(queue):
  """持续消费Kafka消息放入预处理队列"""
  from config import KAFKA_CONFIG
  consumer = KafkaConsumer(
    KAFKA_CONFIG["topic"],
    bootstrap_servers=KAFKA_CONFIG["bootstrap_servers"],
    group_id=KAFKA_CONFIG["group_id"],
    auto_offset_reset=KAFKA_CONFIG["auto_offset_reset"]
  )

  for msg in consumer:
    try:
      data = json.loads(msg.value)
      sampling_params = SamplingParams.from_dict(data["sampling_params"])
      queue.put({
        "id": data["batch_inference_id"],
        "name": data["batch_inference_name"],
        "user_id": data["user_id"],
        "sampling_params": sampling_params,
        "prompts": data["prompts"]
      })
    except Exception as e:
      print(f"无效消息: {msg.value}, 错误: {e}")


def main():
  prefill_queue = Queue()
  decode_queue = Queue()

  prefill_process = Process(target=prefill.run, args=(prefill_queue, decode_queue,))
  decode_process = Process(target=decoder.run, args=(decode_queue,))

  prefill_process.start()
  decode_process.start()

  consume_kafka_messages(prefill_queue)

  decode_process.join()
  prefill_process.terminate()


if __name__ == "__main__":
  main()
