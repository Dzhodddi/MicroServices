import {FastifyInstance, FastifyRequest} from "fastify";


export async function usersRouts(app: FastifyInstance) {
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

    })
}