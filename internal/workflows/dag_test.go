package workflows

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDAG(t *testing.T) {
	tests := []struct {
		name         string
		steps        []PipelineStep
		dependencies map[uuid.UUID][]uuid.UUID
		wantErr      bool
		errContains  string
		wantRoots    int
		wantNodes    int
	}{
		{
			name: "simple linear pipeline",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Step 1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Step 2"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Step 3"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000002")},
			},
			wantRoots: 1,
			wantNodes: 3,
		},
		{
			name: "diamond pattern",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Start"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Parallel1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Parallel2"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Name: "End"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000004"): {
					uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				},
			},
			wantRoots: 1,
			wantNodes: 4,
		},
		{
			name: "multiple roots",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Root1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Root2"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Child"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {
					uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				},
			},
			wantRoots: 2,
			wantNodes: 3,
		},
		{
			name: "simple cycle",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Step 1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Step 2"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"): {uuid.MustParse("00000000-0000-0000-0000-000000000002")},
				uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
			},
			wantErr:     true,
			errContains: "cycle detected",
		},
		{
			name: "self-referencing cycle",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Step 1"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
			},
			wantErr:     true,
			errContains: "cycle detected",
		},
		{
			name:         "empty pipeline",
			steps:        []PipelineStep{},
			dependencies: map[uuid.UUID][]uuid.UUID{},
			wantRoots:    0,
			wantNodes:    0,
		},
		{
			name: "dependency to non-existent step",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Step 1"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"): {uuid.MustParse("00000000-0000-0000-0000-000000000999")},
			},
			wantErr:     true,
			errContains: "not found for node",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dag, err := NewDAG(tt.steps, tt.dependencies)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, dag)
			assert.Len(t, dag.RootNodes, tt.wantRoots)
			assert.Len(t, dag.Nodes, tt.wantNodes)
		})
	}
}

func TestDAG_TopologicalSort(t *testing.T) {
	tests := []struct {
		name         string
		steps        []PipelineStep
		dependencies map[uuid.UUID][]uuid.UUID
		validate     func(t *testing.T, sorted []*DAGNode)
	}{
		{
			name: "linear dependencies",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Step 1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Step 2"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Step 3"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000002")},
			},
			validate: func(t *testing.T, sorted []*DAGNode) {
				require.Len(t, sorted, 3)
				assert.Equal(t, "Step 1", sorted[0].Step.Name)
				assert.Equal(t, "Step 2", sorted[1].Step.Name)
				assert.Equal(t, "Step 3", sorted[2].Step.Name)
			},
		},
		{
			name: "parallel branches",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Root"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Branch1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Branch2"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
			},
			validate: func(t *testing.T, sorted []*DAGNode) {
				require.Len(t, sorted, 3)
				assert.Equal(t, "Root", sorted[0].Step.Name)
				// Branch1 and Branch2 can be in any order
				names := []string{sorted[1].Step.Name, sorted[2].Step.Name}
				assert.Contains(t, names, "Branch1")
				assert.Contains(t, names, "Branch2")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dag, err := NewDAG(tt.steps, tt.dependencies)
			require.NoError(t, err)

			sorted, err := dag.TopologicalSort()
			require.NoError(t, err)

			tt.validate(t, sorted)
		})
	}
}

func TestDAG_ExecutionPlanLevels(t *testing.T) {
	tests := []struct {
		name         string
		steps        []PipelineStep
		dependencies map[uuid.UUID][]uuid.UUID
		wantLevels   []int // number of nodes at each level
	}{
		{
			name: "linear pipeline",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Step 1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Step 2"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Step 3"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000002")},
			},
			wantLevels: []int{1, 1, 1},
		},
		{
			name: "diamond pattern",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Start"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Parallel1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Parallel2"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Name: "End"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000004"): {
					uuid.MustParse("00000000-0000-0000-0000-000000000002"),
					uuid.MustParse("00000000-0000-0000-0000-000000000003"),
				},
			},
			wantLevels: []int{1, 2, 1},
		},
		{
			name: "complex parallel",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "A"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "B"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "C"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Name: "D"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000005"), Name: "E"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000004"): {uuid.MustParse("00000000-0000-0000-0000-000000000002")},
				uuid.MustParse("00000000-0000-0000-0000-000000000005"): {
					uuid.MustParse("00000000-0000-0000-0000-000000000003"),
					uuid.MustParse("00000000-0000-0000-0000-000000000004"),
				},
			},
			wantLevels: []int{1, 2, 1, 1},
		},
		{
			name: "multiple independent roots",
			steps: []PipelineStep{
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Root1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Root2"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Child1"},
				{Id: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Name: "Child2"},
			},
			dependencies: map[uuid.UUID][]uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
				uuid.MustParse("00000000-0000-0000-0000-000000000004"): {uuid.MustParse("00000000-0000-0000-0000-000000000002")},
			},
			wantLevels: []int{2, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dag, err := NewDAG(tt.steps, tt.dependencies)
			require.NoError(t, err)

			levels, err := dag.GetExecutionPlan()
			require.NoError(t, err)
			require.Len(t, levels, len(tt.wantLevels))

			for i, wantCount := range tt.wantLevels {
				assert.Len(t, levels[i], wantCount, "Level %d should have %d nodes", i, wantCount)
			}
		})
	}
}

func TestDAG_GetExecutionPlan(t *testing.T) {
	t.Run("linear execution plan", func(t *testing.T) {
		steps := []PipelineStep{
			{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Step 1"},
			{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Step 2"},
			{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Step 3"},
		}
		dependencies := map[uuid.UUID][]uuid.UUID{
			uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
			uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000002")},
		}

		dag, err := NewDAG(steps, dependencies)
		require.NoError(t, err)

		plan, err := dag.GetExecutionPlan()
		require.NoError(t, err)

		// Should have 3 levels for linear execution
		assert.Len(t, plan, 3)
		assert.Len(t, plan[0], 1) // First level: Step 1
		assert.Len(t, plan[1], 1) // Second level: Step 2
		assert.Len(t, plan[2], 1) // Third level: Step 3
	})

	t.Run("parallel execution plan", func(t *testing.T) {
		steps := []PipelineStep{
			{Id: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Start"},
			{Id: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Parallel1"},
			{Id: uuid.MustParse("00000000-0000-0000-0000-000000000003"), Name: "Parallel2"},
			{Id: uuid.MustParse("00000000-0000-0000-0000-000000000004"), Name: "End"},
		}
		dependencies := map[uuid.UUID][]uuid.UUID{
			uuid.MustParse("00000000-0000-0000-0000-000000000002"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
			uuid.MustParse("00000000-0000-0000-0000-000000000003"): {uuid.MustParse("00000000-0000-0000-0000-000000000001")},
			uuid.MustParse("00000000-0000-0000-0000-000000000004"): {
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		}

		dag, err := NewDAG(steps, dependencies)
		require.NoError(t, err)

		plan, err := dag.GetExecutionPlan()
		require.NoError(t, err)

		// Should have 3 levels: Start -> (Parallel1, Parallel2) -> End
		assert.Len(t, plan, 3)
		assert.Len(t, plan[0], 1) // First level: Start
		assert.Len(t, plan[1], 2) // Second level: Parallel1, Parallel2
		assert.Len(t, plan[2], 1) // Third level: End
	})
}
