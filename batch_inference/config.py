#!/usr/bin/env python
# -*- coding: UTF-8 -*-

# vllm配置
VLLM_CONFIG = {
    "model": "Qwen/Qwen2-0.5B",
    "max_model_len": 2000,
    "gpu_memory_utilization": 0.8,
    "kv_cache": {
        "kv_connector":"PyNcclConnector",
        "kv_rank":1,
        "kv_parallel_size":2
    }
}

# Kafka配置
KAFKA_CONFIG = {
    "bootstrap_servers": "localhost:9092",                 # Kafka服务器地址
    "group_id": "llm-batch-inference-workers",             # 消费组ID
    "auto_offset_reset": "latest",                         # 从最新位置开始消费
    "topic": "batch-inference-requests"                    # 订阅的主题
}

# MongoDB配置
MONGO_CONFIG = {
    "host": "localhost",
    "port": 27017,
    "user": "root",
    "password": "123456",
    "db": "batch_inference",
    "collection": "batch_inference_results",
}

# 日志配置
LOGGER = {
    "level": "debug",
    "outputPath": "output/logs",
    "maxSingleFileSizeMB": 100,
    "maxBackups": 1,
    "maxStorageAgeInDays": 7,
}