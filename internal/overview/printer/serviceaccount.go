package printer

import (
	"context"
	"fmt"

	"github.com/heptio/developer-dash/internal/cache"
	cacheutil "github.com/heptio/developer-dash/internal/cache/util"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func ServiceAccountListHandler(ctx context.Context, list *corev1.ServiceAccountList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("service account list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Secrets", "Age")
	table := component.NewTable("Service Accounts", cols)

	for _, serviceAccount := range list.Items {
		row := component.TableRow{}
		row["Name"] = link.ForObject(&serviceAccount, serviceAccount.Name)
		row["Labels"] = component.NewLabels(serviceAccount.Labels)
		row["Secrets"] = component.NewText(fmt.Sprint(len(serviceAccount.Secrets)))
		row["Age"] = component.NewTimestamp(serviceAccount.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

func ServiceAccountHandler(ctx context.Context, serviceAccount *corev1.ServiceAccount, options Options) (component.ViewComponent, error) {
	o := NewObject(serviceAccount)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return printServiceAccountConfig(ctx, serviceAccount, options.Cache)
	}, 16)

	o.EnableEvents()

	return o.ToComponent(ctx, options)
}

func printServiceAccountConfig(ctx context.Context, serviceAccount *corev1.ServiceAccount, c cache.Cache) (component.ViewComponent, error) {
	if serviceAccount == nil {
		return nil, errors.New("service account is nil")
	}

	var sections component.SummarySections

	var pullSecrets []string

	for _, s := range serviceAccount.ImagePullSecrets {
		pullSecrets = append(pullSecrets, s.Name)
	}

	if len(pullSecrets) > 0 {
		sections.Add("Image Pull Secrets",
			generateServiceAccountSecretsList(serviceAccount.Namespace, pullSecrets))
	}

	var mountSecrets []string
	for _, s := range serviceAccount.Secrets {
		mountSecrets = append(mountSecrets, s.Name)
	}

	if len(mountSecrets) > 0 {
		sections.Add("Mountable Secrets",
			generateServiceAccountSecretsList(serviceAccount.Namespace, mountSecrets))
	}

	tokens, err := serviceAccountTokens(ctx, *serviceAccount, c)
	if err != nil {
		return nil, errors.Wrap(err, "get tokens for service account")
	}

	if len(tokens) > 0 {
		sections.Add("Tokens",
			generateServiceAccountSecretsList(serviceAccount.Namespace, tokens))
	}

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func generateServiceAccountSecretsList(namespace string, secretNames []string) *component.List {
	var items []component.ViewComponent
	for _, name := range secretNames {
		items = append(items, link.ForGVK(namespace, "v1", "Secret", name, name))
	}
	return component.NewList("", items)
}

func serviceAccountTokens(ctx context.Context, serviceAccount corev1.ServiceAccount, c cache.Cache) ([]string, error) {
	key := cacheutil.Key{
		Namespace:  serviceAccount.Namespace,
		APIVersion: "v1",
		Kind:       "Secret",
	}
	secretList, err := c.List(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "find secrets for service account")
	}

	var tokens []string

	for _, u := range secretList {
		secret := &corev1.Secret{}

		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, secret); err != nil {
			return nil, errors.Wrap(err, "convert unstructured secret to structured")
		}

		if err := copyObjectMeta(secret, u); err != nil {
			return nil, errors.Wrap(err, "copy object metadata to secret")
		}

		if secret.Type == corev1.SecretTypeServiceAccountToken {
			name := secret.Annotations[corev1.ServiceAccountNameKey]
			uid := secret.Annotations[corev1.ServiceAccountUIDKey]

			if name == serviceAccount.Name && uid == string(serviceAccount.UID) {
				tokens = append(tokens, secret.Name)
			}
		}
	}

	return tokens, nil
}
