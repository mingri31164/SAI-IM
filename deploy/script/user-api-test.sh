#!/bin/bash
reso_addr='crpi-sosoingx8l87enoq.cn-shanghai.personal.cr.aliyuncs.com/my-zero-im/user-api-dev'
tag='latest'

container_name="sai-im-user-api-test"

docker stop ${container_name}

docker rm ${container_name}

docker rmi ${reso_addr}:${tag}

docker pull ${reso_addr}:${tag}


# 如果需要指定配置文件的
# docker run -p 10001:8080 --network imooc_sai-im -v /sai-im/config/user-rpc:/user/conf/ --name=${container_name} -d ${reso_addr}:${tag}
docker run -p 8888:8888 -p 7888:7888 -p 7999:7999 --name=${container_name} -d ${reso_addr}:${tag}