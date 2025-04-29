FROM nvidia/cuda:12.4.1-cudnn-runtime-ubuntu22.04 as builder

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PYTHON_VERSION=3.12

RUN apt-get update && \
    apt upgrade -y && \
    apt install software-properties-common -y && \
    add-apt-repository ppa:deadsnakes/ppa && \
    apt-get install -y --no-install-recommends \
    python${PYTHON_VERSION} \
    python3-pip \
    python${PYTHON_VERSION}-dev \
    gcc \
    && rm -rf /var/lib/apt/lists/* \
    && ln -s /usr/bin/python${PYTHON_VERSION} /usr/bin/python3 \
    && ln -s /usr/bin/python${PYTHON_VERSION} /usr/bin/python

WORKDIR /app
COPY requirements.txt .
RUN python -m pip install --user -r requirements.txt

FROM nvidia/cuda:12.4.1-cudnn-runtime-ubuntu22.04

ENV PYTHON_VERSION=3.12 \
    PATH=/home/app/.local/bin:$PATH \
    NVIDIA_VISIBLE_DEVICES=all \
    NVIDIA_DRIVER_CAPABILITIES=compute,utility

# 安装Python运行时
RUN apt-get update && \
    apt upgrade -y && \
    apt install software-properties-common -y && \
    add-apt-repository ppa:deadsnakes/ppa && \
    apt-get install -y --no-install-recommends \
    python${PYTHON_VERSION} \
    python3-pip \
    && rm -rf /var/lib/apt/lists/* \
    && ln -s /usr/bin/python${PYTHON_VERSION} /usr/bin/python3 \
    && ln -s /usr/bin/python${PYTHON_VERSION} /usr/bin/python

RUN groupadd -r app && useradd -r -g app app
USER app

COPY --from=builder --chown=app:app /root/.local /home/app/.local

WORKDIR /home/app
COPY --chown=app:app . .

ENV PYTHONPATH=/home/app

EXPOSE 9001

CMD ["python", "app.py"]