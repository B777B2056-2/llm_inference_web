#!/usr/bin/env python
# -*- coding: UTF-8 -*-
# test_perf.py
import grpc
import time
import threading
from concurrent import futures
import pb.tokenizer_pb2 as pb
import pb.tokenizer_pb2_grpc
from config import SERVER_PORT

class PerformanceTester:
    def __init__(self, concurrency=10, total_requests=1000):
        self.concurrency = concurrency
        self.total_requests = total_requests
        self.latencies = []
        self.success_count = 0
        self.error_count = 0
        self.lock = threading.Lock()
        self.test_prompts = [
            "Language models are changing AI.",
            "预训练模型正在改变人工智能",
            "The quick brown fox jumps over the lazy dog"
        ]

    def _worker(self, stub):
        while True:
            with self.lock:
                if self.total_requests <= 0:
                    return
                self.total_requests -= 1

            prompt = self.test_prompts[self.success_count % len(self.test_prompts)]
            start_time = time.perf_counter()

            try:
                response = stub.Tokenizer(pb.tokenizer_pb2.TokenizerReq(
                    parent_message_id="perf_test",
                    prompt=prompt
                ))
                latency = (time.perf_counter() - start_time) * 1000  # 毫秒

                with self.lock:
                    self.latencies.append(latency)
                    self.success_count += 1

            except Exception as e:
                with self.lock:
                    self.error_count += 1
                print(f"Request failed: {str(e)}")

    def run_test(self):
        # 预热连接
        channel = grpc.insecure_channel(f'localhost:{SERVER_PORT}')
        stub = pb.tokenizer_pb2_grpc.TokenizerServiceStub(channel)
        stub.Tokenizer(pb.tokenizer_pb2.TokenizerReq(prompt="warmup"))

        start_time = time.time()

        with futures.ThreadPoolExecutor(max_workers=self.concurrency) as executor:
            for _ in range(self.concurrency):
                executor.submit(self._worker, stub)

        total_time = time.time() - start_time

        # 生成报告
        self._generate_report(total_time)

    def _generate_report(self, total_time):
        sorted_latencies = sorted(self.latencies)

        print("\n=== 性能测试报告 ===")
        print(f"总请求数:   {self.success_count + self.error_count}")
        print(f"成功请求:   {self.success_count}")
        print(f"失败请求:   {self.error_count}")
        print(f"总耗时:     {total_time:.2f}s")
        print(f"QPS:       {self.success_count / total_time:.2f}")
        print(f"平均延迟:   {sum(sorted_latencies)/len(sorted_latencies):.2f}ms")
        print(f"P50延迟:    {sorted_latencies[int(len(sorted_latencies)*0.5)]:.2f}ms")
        print(f"P90延迟:    {sorted_latencies[int(len(sorted_latencies)*0.9)]:.2f}ms")
        print(f"P99延迟:    {sorted_latencies[int(len(sorted_latencies)*0.99)]:.2f}ms")
        print(f"最大延迟:   {sorted_latencies[-1]:.2f}ms")

if __name__ == '__main__':
    # 测试参数配置
    tester = PerformanceTester(
        concurrency=1000,      # 并发线程数
        total_requests=50000  # 总请求数
    )
    tester.run_test()