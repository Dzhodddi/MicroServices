import {FastifyReply, FastifyRequest} from "fastify";


export function roleCheckHandler(role: 'teacher'| "student" | "visitor") {
    return async function (req: FastifyRequest, res: FastifyReply) {
        if (req.user?.role !== role) {
            return res.status(403).send({ error: 'Forbidden: Insufficient role' });
        }
    }
}