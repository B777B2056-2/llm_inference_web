#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import sys
import config
from loguru import logger
import os
from datetime import timedelta


def setup_logging():
    os.makedirs(config.LOGGER["outputPath"], exist_ok=True)
    retention = timedelta(days=config.LOGGER["maxStorageAgeInDays"])
    log_level = config.LOGGER["level"].upper()
    log_format = "{time:YYYY-MM-DD HH:mm:ss.SSS} | {level: <8} | {name}:{function}:{line} - {message}"

    logger.add(sys.stderr, level=log_level, format=log_format, colorize=True, backtrace=True, diagnose=True)
    logger.add(
        sink=os.path.join(config.LOGGER["outputPath"], "tokenizer.log"),
        format="{time:YYYY-MM-DD HH:mm:ss.SSS} | {level: <8} | {name}:{function}:{line} - {message}",
        retention=retention,
        rotation=f'{config.LOGGER["maxSingleFileSizeMB"]} MB',
        compression="zip",
        level=config.LOGGER["level"].upper(),
        enqueue=True
    )
