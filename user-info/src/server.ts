import fastify from "fastify";
import {logger} from "./logger";
import from ""
export async function buildServer() {
    const app = fastify({
    })
    app.register(authPlugin)
    logger.info("Running server")
    return app
}