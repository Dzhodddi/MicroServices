import { prisma } from '../../utils/prisma';
import {logger} from "../../utils/logger";

export function create(data: {name: string, description?: string, creatorId: string}) {
    return prisma.workspace.create({data})
}