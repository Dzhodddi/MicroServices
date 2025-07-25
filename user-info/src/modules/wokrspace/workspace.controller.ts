import {FastifyInstance, FastifyRequest} from "fastify";
import {CreateWorkspaceSchema} from "./workspace.schema";
import * as workspaceService from './workspace.service';
import {roleCheckHandler} from "../../plugins/rolehandler";

export async function workspaceRoutes(app: FastifyInstance) {
    //app.delete("/delete", {})
    app.post("/create", {
        schema: {
            tags: ['Workspace'],
            security: [{ bearerAuth: []}],
            summary: 'Create a new workspace',
            body: {
                type: 'object',
                required: ['name'],
                properties: {
                    name: { type: 'string' },
                    description: { type: 'string' },

                },
            },
            response: {
                201: {
                    type: 'object',
                    properties: {
                        id: { type: 'string' },
                        name: { type: 'string' }
                    }
                },
                500: {
                    type: 'object',
                    properties: {
                        error: {type: "string"},
                    }
                },
                401: {
                    type: 'object',
                    properties: {
                        message: {type: "string"},
                    }
                },
                403: {
                    type: 'object',
                    properties: {
                        message: {type: "string"}
                    }
                }
            }
        },
        //preHandler: [app.auth.verifyJWT, roleCheckHandler("teacher")],
        handler: async (req: FastifyRequest<{Body: CreateWorkspaceSchema}>, res) => {
            const workspace = await workspaceService.createWorkspace("test", req.body);
            res.code(201).send(workspace)
        }
    })
}