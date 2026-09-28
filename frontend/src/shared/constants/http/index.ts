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