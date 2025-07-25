import fp from 'fastify-plugin';
import {FastifyPluginAsync} from 'fastify';
import jwt from 'jsonwebtoken';
import {prisma} from "../utils/prisma";

interface JwtPayload {
    email: string;
    role: string
}

declare module 'fastify' {
    interface FastifyRequest {
        user: JwtPayload;
    }

    interface FastifyInstance {
        auth: {
            verifyJWT: (request: any, reply: any) => Promise<void>;
        };
    }
}

const authPlugin: FastifyPluginAsync = async (fastify) => {
    fastify.decorate('auth', {
        verifyJWT: async (req, res) => {
            try {
                const authHeader = req.headers.authorization;

                if (!authHeader || !authHeader.startsWith('Bearer ')) {
                    return res.code(401).send({ message: 'Missing or invalid Authorization header' });
                }

                const token = authHeader.split(' ')[1];
                const decode = jwt.verify(token, process.env.JWT_SECRET!) as JwtPayload;
                const user = await prisma.user.findUnique({
                    where: {
                        email: decode.email
                    }
                })
                if (!user) return res.status(404).send({ error: "Not found" })
                req.user = user
            } catch (err) {
                return res.code(401).send({message:"JWT verification failed"});
            }
        },
    });
};

export default fp(authPlugin);