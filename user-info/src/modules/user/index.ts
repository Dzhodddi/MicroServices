import {FastifyPluginAsync} from "fastify";

const usersModule: FastifyPluginAsync = async (app) => {
    await app.register(undefined, { prefix: '/users' });
};

export default usersModule;