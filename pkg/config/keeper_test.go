package config

import (
	"testing"

	"github.com/jenkins-x/lighthouse-client/pkg/config/branchprotection"
	"github.com/jenkins-x/lighthouse-client/pkg/config/job"
	"github.com/jenkins-x/lighthouse-client/pkg/config/keeper"
	"github.com/stretchr/testify/assert"
)

func TestConfigGetKeeperContextPolicy(t *testing.T) {
	cases := []struct {
		name                 string
		bpOrgs               map[string]branchprotection.Org
		presubmits           []job.Presubmit
		skipUnknownContexts  bool
		fromBranchProtection bool

		expectedRequired  []string
		expectedOptional  []string
		expectedIfPresent []string
	}{
		{
			name: "basic",
			presubmits: []job.Presubmit{
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "always-run",
						SkipReport: false,
					},
				},
				{
					RegexpChangeMatcher: job.RegexpChangeMatcher{
						RunIfChanged: "foo",
					},
					AlwaysRun: false,
					Reporter: job.Reporter{
						Context:    "run-if-changed",
						SkipReport: false,
					},
				},
				{
					AlwaysRun: false,
					Reporter: job.Reporter{
						Context:    "not-always",
						SkipReport: false,
					},
				},
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "skip-report",
						SkipReport: true,
					},
					Brancher: job.Brancher{
						SkipBranches: []string{"master"},
					},
				},
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "optional",
						SkipReport: false,
					},
					Optional: true,
				},
			},
			expectedRequired:  []string{"always-run"},
			expectedIfPresent: []string{"run-if-changed", "not-always"},
			expectedOptional:  []string{"optional"},
		},
		{
			name: "from branch protection",
			presubmits: []job.Presubmit{
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "always-run",
						SkipReport: false,
					},
				},
				{
					RegexpChangeMatcher: job.RegexpChangeMatcher{
						RunIfChanged: "foo",
					},
					AlwaysRun: false,
					Reporter: job.Reporter{
						Context:    "run-if-changed",
						SkipReport: false,
					},
				},
				{
					AlwaysRun: false,
					Reporter: job.Reporter{
						Context:    "not-always",
						SkipReport: false,
					},
				},
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "skip-report",
						SkipReport: true,
					},
					Brancher: job.Brancher{
						SkipBranches: []string{"master"},
					},
				},
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "optional",
						SkipReport: false,
					},
					Optional: true,
				},
			},
			fromBranchProtection: true,
			bpOrgs: map[string]branchprotection.Org{
				"o": {
					Policy: branchprotection.Policy{},
					Repos: map[string]branchprotection.Repo{
						"r": {
							Policy: branchprotection.Policy{
								RequiredStatusChecks: &branchprotection.ContextPolicy{
									Contexts: []string{
										"always-run",
										"run-if-changed",
									},
								},
							},
						},
					},
				},
			},
			expectedRequired:  []string{"always-run"},
			expectedIfPresent: []string{"run-if-changed", "not-always"},
			expectedOptional:  []string{"optional"},
		},
		{
			name: "from branch protection with unknown context",
			presubmits: []job.Presubmit{
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "always-run",
						SkipReport: false,
					},
				},
				{
					RegexpChangeMatcher: job.RegexpChangeMatcher{
						RunIfChanged: "foo",
					},
					AlwaysRun: false,
					Reporter: job.Reporter{
						Context:    "run-if-changed",
						SkipReport: false,
					},
				},
				{
					AlwaysRun: false,
					Reporter: job.Reporter{
						Context:    "not-always",
						SkipReport: false,
					},
				},
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "skip-report",
						SkipReport: true,
					},
					Brancher: job.Brancher{
						SkipBranches: []string{"master"},
					},
				},
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "optional",
						SkipReport: false,
					},
					Optional: true,
				},
			},
			fromBranchProtection: true,
			bpOrgs: map[string]branchprotection.Org{
				"o": {
					Policy: branchprotection.Policy{},
					Repos: map[string]branchprotection.Repo{
						"r": {
							Policy: branchprotection.Policy{
								RequiredStatusChecks: &branchprotection.ContextPolicy{
									Contexts: []string{
										"always-run",
										"run-if-changed",
										"non-lighthouse-job",
									},
								},
							},
						},
					},
				},
			},
			expectedRequired:  []string{"always-run", "non-lighthouse-job"},
			expectedIfPresent: []string{"run-if-changed", "not-always"},
			expectedOptional:  []string{"optional"},
		},
		{
			name: "from branch protection skipping unknown context",
			presubmits: []job.Presubmit{
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "always-run",
						SkipReport: false,
					},
				},
				{
					RegexpChangeMatcher: job.RegexpChangeMatcher{
						RunIfChanged: "foo",
					},
					AlwaysRun: false,
					Reporter: job.Reporter{
						Context:    "run-if-changed",
						SkipReport: false,
					},
				},
				{
					AlwaysRun: false,
					Reporter: job.Reporter{
						Context:    "not-always",
						SkipReport: false,
					},
				},
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "skip-report",
						SkipReport: true,
					},
					Brancher: job.Brancher{
						SkipBranches: []string{"master"},
					},
				},
				{
					AlwaysRun: true,
					Reporter: job.Reporter{
						Context:    "optional",
						SkipReport: false,
					},
					Optional: true,
				},
			},
			fromBranchProtection: true,
			bpOrgs: map[string]branchprotection.Org{
				"o": {
					Policy: branchprotection.Policy{},
					Repos: map[string]branchprotection.Repo{
						"r": {
							Policy: branchprotection.Policy{
								RequiredStatusChecks: &branchprotection.ContextPolicy{
									Contexts: []string{
										"always-run",
										"run-if-changed",
										"non-lighthouse-job",
									},
								},
							},
						},
					},
				},
			},
			skipUnknownContexts: true,
			expectedRequired:    []string{"always-run", "non-lighthouse-job"},
			expectedIfPresent:   []string{"run-if-changed", "not-always"},
			expectedOptional:    []string{"optional"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for i := range tc.presubmits {
				if err := tc.presubmits[i].SetRegexes(); err != nil {
					t.Fatalf("could not set regexes: %v", err)
				}
			}
			presubmits := map[string][]job.Presubmit{
				"o/r": tc.presubmits,
			}
			cfg := Config{
				JobConfig: job.Config{
					Presubmits: presubmits,
				},
				ProwConfig: ProwConfig{
					Keeper: keeper.Config{
						ContextOptions: keeper.ContextPolicyOptions{
							ContextPolicy: keeper.ContextPolicy{
								SkipUnknownContexts:  &tc.skipUnknownContexts,
								FromBranchProtection: &tc.fromBranchProtection,
							},
						},
					},
				},
			}
			if tc.bpOrgs != nil {
				cfg.ProwConfig.BranchProtection = branchprotection.Config{
					ProtectTested: true,
					Orgs:          tc.bpOrgs,
				}
			}
			ctxPolicy, err := cfg.GetKeeperContextPolicy("o", "r", "master")
			assert.NoError(t, err)
			assert.NotNil(t, ctxPolicy)

			assert.ElementsMatch(t, tc.expectedRequired, ctxPolicy.RequiredContexts)
			assert.ElementsMatch(t, tc.expectedIfPresent, ctxPolicy.RequiredIfPresentContexts)
			assert.ElementsMatch(t, tc.expectedOptional, ctxPolicy.OptionalContexts)
		})
	}
}
