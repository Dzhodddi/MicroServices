import {Pool} from "pg"
import {getEnvString} from "../env";

const pool = new Pool({
    connectionString: getEnvString("DATABASE_URL", ""),
    ssl: false,

})

