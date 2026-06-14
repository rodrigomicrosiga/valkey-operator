# Seta a versão mínima do Tilt
version_settings(constraint='>=0.30.0')

# Libera o deploy para o cluster da TCloud (Safeguard Bypass)
allow_k8s_contexts('francis')

# Diz ao Tilt para usar a sua conta no Docker Hub
default_registry('docker.io/rodrigomicrosiga')

# ====================================================================
# CORREÇÃO CRÍTICA: Configura as regras de exclusão do monitoramento.
# Argumento corrigido de "ignores" para "ignore"
# ====================================================================
watch_settings(ignore=[
    '**/zz_generated.deepcopy.go',  # Ignora o código Go gerado pelo controller-gen
    'config/crd/bases/**',          # Ignora os YAMLs de CRD gerados automaticamente
    'bin/**',                       # Ignora os binários locais de compilação
    'charts/**'                     # Ignora a pasta do Helm durante o dev com Tilt
])

# 1. Tarefa Local: Gera as CRDs e o código deepcopy sempre que a API mudar
local_resource(
    name='manifests',
    cmd='make manifests generate',
    deps=['api/', 'internal/', 'cmd/']
)

# 2. Infraestrutura: Carrega os manifestos de desenvolvimento nativos do Kubebuilder
k8s_yaml(kustomize('config/default'))

# 3. Build Automatizado: Constrói a imagem Docker toda vez que o código Go for salvo
docker_build(
    'valkey-operator', # <-- ALTERADO AQUI
    context='.',
    only=['api', 'cmd', 'internal', 'Dockerfile', 'go.mod', 'go.sum'],
)