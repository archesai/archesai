// import RunForm from "@/components/forms/run-form";
import { Card } from '@/components/ui/card'
import { PipelineStepEntity } from '@/generated/archesApiSchemas'
import { Handle, Position } from '@xyflow/react'
import React from 'react'

function RunFormNode({ data }: { data: PipelineStepEntity }) {
  return (
    <div>
      {/* Include your RunForm component */}
      {/* <RunForm /> */}
      <Card className='flex items-center justify-center px-2 py-1'>{data.tool.name}</Card>
      {/* Add handles for connecting nodes */}
      <Handle position={Position.Left} type='target' />
      <Handle position={Position.Right} type='source' />
    </div>
  )
}

export default RunFormNode
