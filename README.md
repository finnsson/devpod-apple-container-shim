# devpod-apple-container-shim


```
devpod provider add ./provider-local.yaml
devpod provider set-options apple-container -o SHIM_PATH=$(pwd)/build/devpod-apple-container-shim
devpod up --provider apple-container ubuntu:latest
```