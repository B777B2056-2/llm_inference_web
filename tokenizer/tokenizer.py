#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import config
from loguru import logger
import threading
from transformers import AutoTokenizer


class Qwen2TokenizerSingleton:
    """
    Qwen2 Tokenizer单例管理器
    """
    _instance = None
    _lock = threading.Lock()
    _is_initialized = False

    def __new__(cls):
        with cls._lock:
            if cls._instance is None:
                cls._instance = super().__new__(cls)
                # 延迟加载机制
                cls._instance._model_name = config.MODEL_NAME
                cls._instance._tokenizer = None
            return cls._instance

    def initialize(self):
        """显式初始化方法（线程安全）"""
        if self._tokenizer is None:
            with self.__class__._lock:
                if self._tokenizer is None:  # 双重检查锁定
                    try:
                        self._tokenizer = AutoTokenizer.from_pretrained(
                            self._model_name,
                            use_fast=True,
                            add_prefix_space=True,
                            truncation_side="left",
                        )
                        self.__class__._is_initialized = True
                    except Exception as e:
                        raise RuntimeError(f"Tokenizer初始化失败: {str(e)}")

    @property
    def tokenizer(self):
        """安全访问tokenizer属性"""
        if not self._is_initialized:
            self.initialize()
        return self._tokenizer

    @classmethod
    def get_tokenizer(cls):
        """获取全局tokenizer实例"""
        instance = cls()
        return instance.tokenizer
