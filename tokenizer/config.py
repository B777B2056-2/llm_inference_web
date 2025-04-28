#!/usr/bin/env python
# -*- coding: UTF-8 -*-


# 服务配置
SERVER_PORT = 9000
MAX_WORKS = 10


# LLM配置
MODEL_NAME = "Qwen/Qwen2-0.5B"

# 日志配置
LOGGER = {
    "level": "debug",
    "outputPath": "output/logs",
    "maxSingleFileSizeMB": 100,
    "maxBackups": 1,
    "maxStorageAgeInDays": 7,
}
