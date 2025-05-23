'use client'

import type { Connection, Edge, Node } from '@xyflow/react'

// import '@xyflow/react/dist/style.css'

import { useCallback, useMemo } from 'react'
import {
  addEdge,
  Background,
  Controls,
  MiniMap,
  Panel,
  ReactFlow,
  useEdgesState,
  useNodesState
} from '@xyflow/react'

import { Button } from '@archesai/ui/components/shadcn/button'

// import { useAuth } from '@archesai/ui/hooks/use-auth'

import RunFormNode from './node'

const initialNodes: Node[] = []
const initialEdges: Edge[] = []

export const CreatePipelineContent = () => {
  const [nodes, _setNodes, onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges)

  const onConnect = useCallback(
    (params: Connection | Edge) => {
      setEdges((edges) => addEdge<Edge>(params, edges))
    },
    [setEdges]
  )

  // useEffect(() => {
  //   if (pipelines?.data[0]) {
  //     const pipelineSteps = pipelines.data[0]?.attributes.steps
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
  //           id: `${prereq.pipelineStepId}-${step.id}`,
  //           label: 'depends on',
  //           source: step.id.toString(),
  //           target: step.id,
  //           type: 'smoothstep'
  //         })
  //       })
  //     }
  //   }
  // }, [pipelines, onConnect, setNodes])

  const nodeTypes = useMemo(() => ({ runFormNode: RunFormNode }), [])

  return (
    <ReactFlow
      attributionPosition='top-right'
      edges={edges}
      elementsSelectable
      fitView
      nodes={nodes}
      nodeTypes={nodeTypes}
      onConnect={onConnect}
      onEdgesChange={onEdgesChange}
      onNodesChange={onNodesChange}
    >
      <Panel position='top-right'>
        <Button className='z-50'>Useless Button</Button>
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
  )
}
