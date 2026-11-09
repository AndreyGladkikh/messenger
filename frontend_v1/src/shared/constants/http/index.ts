const method = {
    get: 'GET',
    post: 'POST',
    put: 'PUT',
    patch: 'PATCH',
    delete: 'DELETE',
}

const status = {
    badRequest: 400,
    unauthorized: 401,
    notFound: 404,
    conflict: 409,
    noContent: 204,
}

const header = {
    contentType: 'Content-Type',
    authorization: 'Authorization',
}

const contentType = {
    json: 'application/json',
    problemPlusJson: 'application/problem+json',
}

const cookie = {
    accessToken: 'access_token',
}

export default Object.freeze({
    method,
    status,
    header,
    contentType,
    cookie,
})
