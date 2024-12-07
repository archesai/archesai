'use client'
import { Button } from '@/components/ui/button'
import { usePipelinesControllerFindAll } from '@/generated/archesApiComponents'
import { PipelineStepEntity } from '@/generated/archesApiSchemas'
import { useAuth } from '@/hooks/use-auth'
import {
  addEdge,
  Background,
  Connection,
  Controls,
  Edge,
  MiniMap,
  Node,
  Panel,
  ReactFlow,
  useEdgesState,
  useNodesState
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import React, { useCallback, useEffect, useMemo } from 'react'

import RunFormNode from './node'

const initialNodes: Node<PipelineStepEntity>[] = []
const initialEdges: Edge[] = []

export default function App() {
  const { defaultOrgname } = useAuth()
  const { data: pipelines } = usePipelinesControllerFindAll({
    pathParams: {
      orgname: defaultOrgname
    }
  })

  console.log(pipelines)
  useEffect(() => {
    if (pipelines && pipelines.results[0]) {
      const pipelineSteps = pipelines.results[0].pipelineSteps
      const nodes = pipelineSteps.map((step, index) => ({
        data: step,
        id: step.id,
        position: { x: 200 + index * 200, y: 100 },
        type: 'runFormNode'
      }))
      setNodes(nodes)
      for (const step of pipelineSteps) {
        step.dependsOn?.forEach((id) =>
          onConnect({
            animated: true,
            id: `${id.id}-${step.id}`,
            label: 'depends on',
            source: id.id,
            target: step.id,
            type: 'smoothstep'
          })
        )
      }
    }
  }, [pipelines])

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges)

  const onConnect = useCallback(
    (params: Connection | Edge) =>
      setEdges((edges) => addEdge<Edge>(params, edges)),
    [setEdges]
  )

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
        <Button
          className='z-50'
          onClick={() => console.log(nodes, edges)}
        >
          Log
        </Button>
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
