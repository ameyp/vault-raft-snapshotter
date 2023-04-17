allow_k8s_contexts('default')

docker_build('vault-snapshot', '.', dockerfile='local-development/Dockerfile')
k8s_yaml('local-development/job.yaml')
