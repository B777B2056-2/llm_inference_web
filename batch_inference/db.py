#!/usr/bin/env python
# -*- coding: UTF-8 -*-
from typing import List, Tuple
import pymongo
from config import MONGO_CONFIG


class BatchInferenceDatabase(object):
  def __init__(self):
    uri = f'mongodb://{MONGO_CONFIG["user"]}:{MONGO_CONFIG["password"]}@{MONGO_CONFIG["host"]}:{MONGO_CONFIG["port"]}/'
    self.clt = pymongo.MongoClient(uri)
    self.db = self.clt[MONGO_CONFIG["db"]]

  def insert_results(self, id: str, name: str, user_id: int, results: List[Tuple[str, str]]) -> None:
    collection = self.db[MONGO_CONFIG["collection"]]
    data = []
    for (prompt, text) in results:
      data.append(
        {
          "batch_inference_id": id,
          "batch_inference_name": name,
          "user_id":user_id,
          "prompt": prompt,
          "text": text
        }
      )
    collection.insert_many(data)
