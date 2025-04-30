FROM nvidia/cuda:12.4.1-cudnn-runtime-ubuntu22.04

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    PIP_NO_CACHE_DIR=1 \
    PYTHON_VERSION=3.10 \
    NVIDIA_VISIBLE_DEVICES=all \
    NVIDIA_DRIVER_CAPABILITIES=compute,utility

RUN echo "">sources.list && \
    echo "deb http://ftp2.cn.debian.org/debian/ buster main">>sources.list && \
    echo "deb http://ftp2.cn.debian.org/debian/debian-security buster/updates main">>sources.list && \
    echo "deb http://ftp2.cn.debian.org/debian/debian buster-updates main">>sources.list

RUN apt update && \
    apt install -y --no-install-recommends && \
    apt install -y python${PYTHON_VERSION} && \
    apt-get install -y python3-pip && \
    rm -rf /var/lib/apt/lists/*

RUN python3 --version

WORKDIR /app
COPY . .

RUN python3 -m pip install --user -r requirements.txt --use-deprecated=legacy-resolver

EXPOSE 9001

CMD ["python3", "app.py"]