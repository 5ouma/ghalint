package policy

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/config"
	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type GitHubAppShouldLimitPermissionsPolicy struct{}

func (p *GitHubAppShouldLimitPermissionsPolicy) Name() string {
	return "github_app_should_limit_permissions"
}

func (p *GitHubAppShouldLimitPermissionsPolicy) ID() string {
	return "010"
}

func (p *GitHubAppShouldLimitPermissionsPolicy) ApplyStep(_ *logrus.Entry, _ *config.Config, _ *StepContext, step *workflow.Step) (ge error) { //nolint:cyclop
	action := p.checkUses(step.Uses)
	if action == "" {
		return nil
	}
	defer func() {
		if ge != nil {
			ge = logerr.WithFields(ge, logrus.Fields{
				"action": action,
			})
		}
	}()

	switch action {
	case "tibdex/github-app-token":
		if step.With == nil {
			return errPermissionsIsRequired
		}
		if _, ok := step.With["permissions"]; !ok {
			return errPermissionsIsRequired
		}
	case "actions/create-github-app-token":
		if step.With == nil {
			return errPermissionsIsRequired
		}
		err := errPermissionHyphenIsRequired
		for k := range step.With {
			if strings.HasPrefix(k, "permission-") {
				err = nil
				break
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *GitHubAppShouldLimitPermissionsPolicy) checkUses(uses string) string {
	if uses == "" {
		return ""
	}
	action, _, _ := strings.Cut(uses, "@")
	return action
}
