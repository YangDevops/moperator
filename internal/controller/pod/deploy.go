package pod

import (
	"context"
	"fmt"

	"github.com/infraboard/mpaas/apps/deploy"
	"github.com/infraboard/mpaas/common/format"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *PodReconciler) HandleDeploy(ctx context.Context, obj v1.Pod) error {
	// 获取日志对象
	l := log.FromContext(ctx)

	// 更新Deploy注解 更新Deploy
	deployId := obj.Annotations[deploy.ANNOTATION_DEPLOY_ID]
	if deployId != "" {
		l.Info(fmt.Sprintf("get mpaas deploy: %s", deployId))

		// 查询Deploy
		ins, err := r.mpaas.Deploy().DescribeDeployment(ctx, deploy.NewDescribeDeploymentRequest(deployId))
		if err != nil {
			return fmt.Errorf("get deploy error, %s", err)
		}

		// 更新Depoy
		updateReq := deploy.NewUpdateDeploymentStatusRequest(deployId)
		if ins.Credential != nil {
			updateReq.UpdateToken = ins.Credential.Token
		}

		updateReq.UpdatedK8SConfig.Pods = format.MustToYaml(obj)
		updateReq.UpdateBy = r.name
		_, err = r.mpaas.Deploy().UpdateDeploymentStatus(ctx, updateReq)
		if err != nil {
			return err
		}
		l.Info(fmt.Sprintf("deploy: %s, update success", deployId))
	}
	return nil
}
