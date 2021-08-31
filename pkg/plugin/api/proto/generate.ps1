$dir = $PSScriptRoot
$octantRoot = (get-item $dir).parent.parent.parent.parent.FullName
$module = "github.com/vmware-tanzu/octant/pkg/plugin/api/proto"

protoc -I ${octantRoot}/vendor -I ${octantRoot} -I ${dir} --go_out=plugins=grpc:${dir} --go_opt=module=${module} ${dir}/dashboard_api.proto