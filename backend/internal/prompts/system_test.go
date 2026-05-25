package prompts

import (
	"strings"
	"testing"

	"github.com/NeevBhandari13/leetcoach/internal/session"
)

func TestGetSystemPrompt(t *testing.T) {
	tests := []struct {
		name        string
		state       session.State
		problemText string
		wantContain []string
	}{
		{
			name:        "includes problem text in base instructions",
			state:       session.IntroState,
			problemText: "Given an array of integers",
			wantContain: []string{"Given an array of integers"},
		},
		{
			name:        "intro state includes intro and present_problem instructions",
			state:       session.IntroState,
			problemText: "some problem",
			wantContain: []string{"intro", "present_problem"},
		},
		{
			name:        "present_problem state includes clarify and initial_solution transitions",
			state:       session.PresentProblemState,
			problemText: "some problem",
			wantContain: []string{"clarify", "initial_solution"},
		},
		{
			name:        "wrap_up state does not include further transitions",
			state:       session.WrapUpState,
			problemText: "some problem",
			wantContain: []string{"wrap_up"},
		},
		{
			name:        "result is non-empty for every state",
			state:       session.OptimisationState,
			problemText: "some problem",
			wantContain: []string{"optimisation"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSystemPrompt(tt.state, tt.problemText)

			if got == "" {
				t.Fatal("expected non-empty system prompt")
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf("expected prompt to contain %q\nfull prompt:\n%s", want, got)
				}
			}
		})
	}
}
