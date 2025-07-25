import {FastifyPluginAsync} from "fastify";
import {workspaceRoutes} from "./workspace.controller";

const workspaceModule: FastifyPluginAsync = async (app) => {
    await app.register(workspaceRoutes, { prefix: '/workspace' });
};

export default workspaceModule;