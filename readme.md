https://github.com/iot-salzburg/gpu-jupyter

```shell
docker run --gpus all -d -it -p 8848:8888 -v $(pwd)/jupyter:/home/jovyan/work -e GRANT_SUDO=yes -e JUPYTER_ENABLE_LAB=yes --user root cschranz/gpu-jupyter:v1.6_cuda-12.0_ubuntu-22.04
```

- [Белый список DNS](https://www.dnswl.org/?page_id=4)
- https://github.com/tb0hdan/domains
- https://www.team-cymru.com/ip-asn-mapping
- https://cert.pl/en/posts/2020/03/malicious_domains/
- https://github.com/stamparm/blackbook
- https://github.com/sefinek24/Sefinek-Blocklist-Collection