#!/usr/bin/env python
# -*- coding: UTF-8 -*-
SERVER_PORT = 9001
MODEL = "Qwen/Qwen2-0.5B"
MAX_TOKENS =4096

# 日志配置
LOGGER = {
    "level": "debug",
    "outputPath": "output/logs",
    "maxSingleFileSizeMB": 100,
    "maxBackups": 1,
    "maxStorageAgeInDays": 7,
}