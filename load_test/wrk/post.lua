function randomString(length)
    local chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
    local result = ''
    for i = 1, length do
        local index = math.random(#chars)
        result = result .. chars:sub(index, index)
    end
    return result
end

function request()
    local method = "POST"
    local path = "/v0/entity"
    local name = randomString(25)
    local salary = math.random(1, 100000)
    local body = '{"name": "' .. name .. '", "salary": ' .. salary .. '}'

    return wrk.format(method, path, headers, body)
end
