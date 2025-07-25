import * as repository from "../wokrspace/workspace.repository";

export async function createUser(id: string, data: {name: string, description?: string}) {
    return repository.create({
        name: data.name,
        description: data.description,
        creatorId: id
    })
}
