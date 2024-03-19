https://github.com/iot-salzburg/gpu-jupyter

```shell
docker run --gpus all -d -it -p 8848:8888 -v $(pwd)/jupyter:/home/jovyan/work -e GRANT_SUDO=yes -e JUPYTER_ENABLE_LAB=yes --user root cschranz/gpu-jupyter:v1.6_cuda-12.0_ubuntu-22.04
```