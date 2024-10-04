import { IntersectionType, PickType } from "@nestjs/swagger";

import { AnimationEntity } from "../entities/animation.entity";

export class CreateAnimationDto extends IntersectionType(
  PickType(AnimationEntity, [
    "width",
    "height",
    "fps",
    "maxFrames",
    "animationPrompts",
    "audioUrl",
    "useAudio",
    "audioStart",
  ] as const)
) {}
