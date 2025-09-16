import type { Connection, Edge, Node } from "@xyflow/react";
import type { JSX } from "react";

// import '@xyflow/react/dist/style.css'

import { Button, Card } from "@archesai/ui";
import {
  addEdge,
  Background,
  Controls,
  Handle,
  MiniMap,
  Panel,
  Position,
  ReactFlow,
  useEdgesState,
  useNodesState,
} from "@xyflow/react";
import { useCallback, useMemo } from "react";

// PipelineStepEntity doesn't exist in generated types yet
interface PipelineStepEntity {
  id: string;
  prerequisites: string[];
  toolID: string;
}

export default function CreatePipelinePage(): JSX.Element {
  return <CreatePipelineContent />;
}

function RunFormNode({ data }: { data: PipelineStepEntity }) {
  return (
    <div>
      {/* Include your RunForm component */}
      {/* <RunForm /> */}
      <Card className="flex items-center justify-center px-2 py-1">
        {data.toolID}
      </Card>
      {/* Add handles for connecting nodes */}
      <Handle
        position={Position.Left}
        type="target"
      />
      <Handle
        position={Position.Right}
        type="source"
      />
    </div>
  );
}

const initialNodes: Node[] = [];
const initialEdges: Edge[] = [];

export const CreatePipelineContent = (): JSX.Element => {
  const [nodes, _setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const onConnect = useCallback(
    (params: Connection | Edge) => {
      setEdges((edges) => addEdge<Edge>(params, edges));
    },
    [setEdges],
  );

  // useEffect(() => {
  //   if (pipelines?.data[0]) {
  //     const pipelineSteps = pipelines.data[0]?.steps
  //     const nodes = pipelineSteps.map((step, index) => ({
  //       data: step,
  //       id: step.id,
  //       position: { x: 200 + index * 200, y: 100 },
  //       type: 'runFormNode'
  //     }))
  //     setNodes(nodes)
  //     for (const step of pipelineSteps) {
  //       step.prerequisites.forEach((prereq) => {
  //         onConnect({
  //           animated: true,
  //           id: `${prereq.pipelineStepID}-${step.id}`,
  //           label: 'depends on',
  //           source: step.id.toString(),
  //           target: step.id,
  //           type: 'smoothstep'
  //         })
  //       })
  //     }
  //   }
  // }, [pipelines, onConnect, setNodes])

  const nodeTypes = useMemo(() => ({ runFormNode: RunFormNode }), []);

  return (
    <ReactFlow
      attributionPosition="top-right"
      edges={edges}
      elementsSelectable
      fitView
      nodes={nodes}
      nodeTypes={nodeTypes}
      onConnect={onConnect}
      onEdgesChange={onEdgesChange}
      onNodesChange={onNodesChange}
    >
      <Panel position="top-right">
        <Button className="z-50">Useless Button</Button>
      </Panel>
      <Controls />
      <MiniMap
        pannable
        zoomable
      />
      <Background
        gap={12}
        size={1}
      />
    </ReactFlow>
  );
};
