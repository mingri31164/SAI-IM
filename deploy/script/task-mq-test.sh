#!/bin/bash
reso_addr='crpi-sosoingx8l87enoq.cn-shanghai.personal.cr.aliyuncs.com/my-zero-im/task-mq-dev'
tag='latest'

pod_idb="106.14.194.111"

container_name="sai-im-task-mq-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_sai-im -v /sai-im/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -e POD_IP=${pod_idb} --name=${container_name} -d ${reso_addr}:${tag}