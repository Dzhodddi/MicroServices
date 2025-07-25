import {buildServer} from "./server";
import {logger} from "./utils/logger";
import {getEnvInt} from "./env";

const port = getEnvInt("PORT", 8081)
async function gracefulShutdown({app}: {app: Awaited<ReturnType<typeof buildServer>>}) {
    await app.close()
}

async function main() {
    const app = await buildServer()
    await app.ready()
    app.swagger()
    await app.listen({
        port: port,
    })
    const signals = [`SIGINT`, `SIGTERM`]
    for (const signal of signals) {
        process.on(signal, () => gracefulShutdown({app}))
    }

    logger.info(`server is running on port ${port}`)
}

main().catch((err) => console.error(err))