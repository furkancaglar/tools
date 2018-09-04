const net = require("net")
const flags = require("flags")

flags.defineInteger('port', 1112, 'port to serve socket.io');
flags.defineInteger('timeout', 1000, 'timeout for pong');
flags.parse()

if (flags.get("port") > 65000)
{
    throw new Error("invalid port")
}
const sock = new net.Socket()
const io = require("socket.io")(flags.get("port"))

function retry(e)
{
    inc++

    setTimeout(function ()
    {
        if (0 == inc) return
        console.log("error : ", e)
        sock.connect({host: process.env.SOCKET_HOST || "localhost", port: process.env.SOCKET_PORT || 1111})
    }, 1000 * inc)

}

function emit__data(data)
{
    if (!data.rooms || !data.rooms.length)
    {
        io.emit(data.event, data.data)
    }
    data.rooms.forEach(function (room)
    {
        io.of("/" + room).emit(data.event, data.data)
    });
}

sock.connect({host: process.env.SOCKET_HOST || "localhost", port: process.env.SOCKET_PORT || 1111})
sock.on("data", function (d)
{

    let data
    try
    {
        data = JSON.parse(d.toString())
        emit__data(data)
    } catch (e)
    {
        let dt = d.toString().split("}{")
        if (dt.length <= 1) return
        dt.forEach(field => {
            if (field[0] != "{") field = "{" + field
            if (field[field.length - 1] != "}") field += "}"
            try
            {
                data = JSON.parse(field)
                emit__data(data)
            } catch (e)
            {
                console.error(e)
            }
        })
    }
})

sock.on("connect", function ()
{
    inc = 0
});
sock.on("end", retry)

var inc = 0
sock.on("error", retry)

setInterval(_ => {

        sock.write("pong")
        io.clients((_, clients) => {
            process.stdout.write("\rconnected clients num:" + clients.length)
        })
    }
    , flags.get("timeout"))
