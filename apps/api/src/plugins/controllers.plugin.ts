import type { FastifyPluginAsync } from 'fastify'

import fp from 'fastify-plugin'

// import {
//   accountsPlugin,
//   invitationsPlugin,
//   membersPlugin,
//   organizationsPlugin,
//   sessionsPlugin,
//   usersPlugin
// } from '@archesai/auth'
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
  // Core controllers
  await app.register(configController, {
    configService: container.configService
  })

  // Auth controllers
  // await app.register(accountsPlugin, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(invitationsPlugin, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(membersPlugin, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(organizationsPlugin, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(sessionsPlugin, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(usersPlugin, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // Billing controllers
  // await app.register(callbacksController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(paymentMethodsController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(plansController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(stripeController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(subscriptionsController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // Orchestration controllers
  await app.register(fp(artifactsController), {
    databaseService: container.databaseService,
    websocketsService: container.websocketsService
  })

  // await app.register(labelsController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(pipelinesController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(runsController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })

  // await app.register(toolsController, {
  //   databaseService: container.databaseService,
  //   websocketsService: container.websocketsService
  // })
}
