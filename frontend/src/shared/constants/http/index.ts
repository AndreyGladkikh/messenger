const method = {
    get: 'GET',
    post: 'POST',
    put: 'PUT',
    patch: 'PATCH',
    delete: 'DELETE',
}

const status = {
    unauthorized: 401,
}

const header = {
    contentType: 'Content-Type',
}

const contentType = {
    json: 'application/json',
}

export default Object.freeze({
    method,
    status,
    header,
    contentType,
})