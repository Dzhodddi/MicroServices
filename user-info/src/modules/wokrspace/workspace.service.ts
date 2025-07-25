import * as repository from './workspace.repository';

export async function createWorkspace(id: string, data: {name: string, description?: string}) {
    return repository.create({
        name: data.name,
        description: data.description,
        creatorId: id
    })
}
