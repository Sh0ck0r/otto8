package workflow

import (
	"context"

	"github.com/otto8-ai/nah/pkg/name"
	"github.com/otto8-ai/nah/pkg/router"
	"github.com/otto8-ai/otto8/pkg/create"
	v1 "github.com/otto8-ai/otto8/pkg/storage/apis/otto.otto8.ai/v1"
	"github.com/otto8-ai/otto8/pkg/system"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func createWorkspace(ctx context.Context, c kclient.Client, workflow *v1.Workflow) error {
	if workflow.Status.WorkspaceName != "" {
		return nil
	}

	ws := &v1.Workspace{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  workflow.Namespace,
			Name:       name.SafeConcatName(system.WorkspacePrefix, workflow.Name),
			Finalizers: []string{v1.WorkspaceFinalizer},
		},
		Spec: v1.WorkspaceSpec{
			AgentName: workflow.Name,
		},
	}
	if err := create.OrGet(ctx, c, ws); err != nil {
		return err
	}

	workflow.Status.WorkspaceName = ws.Name
	return c.Status().Update(ctx, workflow)
}

func createKnowledgeSet(ctx context.Context, c kclient.Client, workflow *v1.Workflow) error {
	if len(workflow.Status.KnowledgeSetNames) > 0 {
		return nil
	}

	ks := &v1.KnowledgeSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  workflow.Namespace,
			Name:       name.SafeConcatName(system.KnowledgeSetPrefix, workflow.Name),
			Finalizers: []string{v1.KnowledgeSetFinalizer},
		},
		Spec: v1.KnowledgeSetSpec{
			AgentName: workflow.Name,
		},
	}
	if err := create.OrGet(ctx, c, ks); err != nil {
		return err
	}

	workflow.Status.KnowledgeSetNames = append(workflow.Status.KnowledgeSetNames, ks.Name)
	return c.Status().Update(ctx, workflow)
}

func CreateWorkspaceAndKnowledgeSet(req router.Request, _ router.Response) error {
	workflow := req.Object.(*v1.Workflow)

	if err := createWorkspace(req.Ctx, req.Client, workflow); err != nil {
		return err
	}

	return createKnowledgeSet(req.Ctx, req.Client, workflow)
}
