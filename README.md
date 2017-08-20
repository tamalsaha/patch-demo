# patch-demo
Using Kubernetes patch apis

I am exploring an idea to auto generate json patch and apply the JSON patch type. This is an working demo of that.

## Background
https://kubernetes.io/docs/tasks/run-application/update-api-object-kubectl-patch/

## Kubernetes Client-go Util
I have put together a util library for client-go to easily use PATCH apis. You can find it here: https://github.com/appscode/kutil
