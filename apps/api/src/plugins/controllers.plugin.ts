import type { FastifyPluginAsync } from 'fastify'

import {
  accountsPlugin,
  authPlugin,
  emailChangeController,
  emailVerificationController,
  invitationsPlugin,
  membersPlugin,
  organizationsPlugin,
  passwordResetController,
  sessionsController,
  usersPlugin
} from '@archesai/auth'
// import {
//   callbacksController,
//   paymentMethodsController,
//   plansController,
//   stripeController,
//   subscriptionsController
// } from '@archesai/billing'
import { configController } from '@archesai/core'
import {
  artifactsController,
  labelsController,
  pipelinesController,
  runsController,
  toolsController
} from '@archesai/orchestration'

import type { Container } from '#utils/container'

export interface ControllersPluginOptions {
  container: Container
}

export const controllersPlugin: FastifyPluginAsync<
  ControllersPluginOptions
> = async (app, { container }) => {
  // Auth controllers
  await app.register(authPlugin, {
    authService: container.authService
  })

  await app.register(emailVerificationController)
  await app.register(emailChangeController)
  await app.register(passwordResetController)

  await app.register(accountsPlugin, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(invitationsPlugin, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(membersPlugin, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(organizationsPlugin, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(sessionsController, {
    authService: container.authService,
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(usersPlugin, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  // Core controllers
  await app.register(configController, {
    configService: container.configService
  })

  // Billing controllers
  // await app.register(callbacksController, {
  //   configService: container.configService,
  //   stripeService: container.stripeService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(paymentMethodsController, {
  //   stripeService: container.stripeService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(plansController, {
  //   logger: container.loggerService.logger,
  //   stripeService: container.stripeService
  // })

  // await app.register(stripeController, {
  //   stripeService: container.stripeService
  // })

  // await app.register(subscriptionsController, {
  //   stripeService: container.stripeService
  // })

  // Orchestration controllers
  await app.register(artifactsController, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(labelsController, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(pipelinesController, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(runsController, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  await app.register(toolsController, {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })
}
