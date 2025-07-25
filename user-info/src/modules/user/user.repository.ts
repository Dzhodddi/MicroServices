import {prisma} from "../../utils/prisma";


export function create(data: {name: string, description?: string, creatorId: string}) {
    return prisma.workspace.create({data})
}