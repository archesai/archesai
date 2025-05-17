// import RunForm from "#components/forms/run-form";
import { Handle, Position } from '@xyflow/react'

import type { PipelineStepEntity } from '@archesai/domain'

import { Card } from '@archesai/ui/components/shadcn/card'

function RunFormNode({ data }: { data: PipelineStepEntity }) {
  return (
    <div>
      {/* Include your RunForm component */}
      {/* <RunForm /> */}
      <Card className='flex items-center justify-center px-2 py-1'>
        {data.tool.name}
      </Card>
      {/* Add handles for connecting nodes */}
      <Handle
        position={Position.Left}
        type='target'
      />
      <Handle
        position={Position.Right}
        type='source'
      />
    </div>
  )
}

export default RunFormNode
