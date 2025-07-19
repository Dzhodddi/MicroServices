import fastify from "fastify";
import {logger} from "./logger";

export async function buildServer() {
    const app = fastify({
    })
    logger.info("Running server")
    return app
}