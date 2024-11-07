id = 1

function request()
    method = "GET"
    path = "/v0/employees?id=" .. tostring(id)
    id = id + 1
    return wrk.format(method, path, body)
end