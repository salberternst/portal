const SessionStorageThingsboardCredentials = "thingsboard_credentials"

const fetchThingsboardCredentials = async () => {
    let data = sessionStorage.getItem(SessionStorageThingsboardCredentials)
    if (data === null) {
        const response = await fetch('/api/things/credentials')
        const credentials = await response.json()

        if (!response.ok) {
            //todo: handle
        }

        sessionStorage.setItem(SessionStorageThingsboardCredentials, JSON.stringify(credentials))

        return credentials
    }

    return JSON.parse(data)
}

const fetchThings = async (pagination: any) => {
    const { page, perPage }: { page: number; perPage: number } = pagination
    const credentials = await fetchThingsboardCredentials()
    const response = await fetch(`http://thingsboard.192-168-178-60.nip.io/api/tenant/devices?page=${page - 1}&pageSize=${perPage}`, {
        headers: {
            Authorization: `Bearer ${credentials.token}`
        }
    })
    return response.json()
}

const fetchThing = async (id: string) => {
    const credentials = await fetchThingsboardCredentials()
    const response = await fetch(`http://thingsboard.192-168-178-60.nip.io/api/device/${id}`, {
        headers: {
            Authorization: `Bearer ${credentials.token}`
        }
    })
    return response.json()
}

export default {
    getList: async (resource: any, params: any) => {
        if (resource === "thing") {
            const result = await fetchThings(params.pagination)
            return {
                data: result.data.map((thing: any) => ({ ...thing, id: thing.id.id })),
                total: result.totalElements
            }
        }
    },
    getOne: async (resource: any, params: any) => {
        if (resource === "thing") {
            const result = await fetchThing(params.id)
            return {
                data: { 
                    ...result,
                    id: params.id
                }
            }
        }
    },
}