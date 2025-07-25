import fastify from "fastify";
import {logger} from "./utils/logger";
import authPlugin from "./plugins/jwt"
import workspaceModule from "./modules/wokrspace";
import {fastifySwaggerUi} from '@fastify/swagger-ui';
import fastifySwagger from "@fastify/swagger";
export async function buildServer() {
    const app = fastify({
    })
    await app.register(fastifySwagger, {
        openapi: {
            info: {
                title: 'User info service Swagger',
                description: 'My info service docs',
                version: '1.0.0'
            },
            components: {
                securitySchemes: {
                    bearerAuth: {
                        type: 'http',
                        scheme: 'bearer',
                        bearerFormat: 'JWT',
                    },
                },
            },
            security: [{ bearerAuth: [] }],
        },

    })
    await app.register(fastifySwaggerUi, {
        routePrefix: '/docs',
        uiConfig: {
            docExpansion: 'full',
        },
    });


    await app.register(authPlugin)
    await app.register(workspaceModule, {
        prefix: "/api"
    })


    logger.info("Running server")
    return app
}