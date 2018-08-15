const net = require("net")
const request = require("request")
const arglist = process.argv
if (arglist.length < 3)
{
    throw new Error("no port name is specified")
}
else if (!/^\d{2,5}$/.test(arglist[2]))
{
    console.log(arglist[2])
    throw new Error("invalid port")
}
let port = parseInt(arglist[2])
if (port > 65000)
{
    throw new Error("invalid port")
}
const sock = new net.Socket()
const io = require("socket.io")(port)
function retry(e)
{
    inc++

    setTimeout(function ()
{
        if (0 == inc) return
        console.log(e)
        sock.connect({host: process.env.SOCKET_HOST || "localhost", port: process.env.SOCKET_PORT || 1111})
    }, 1000 * inc)

}

sock.connect({host: process.env.SOCKET_HOST || "localhost", port: process.env.SOCKET_PORT || 1111})
sock.on("data", function (d)
{
    let data
    try {
        data = JSON.parse(d.toString())
        if (!data.rooms||!data.rooms.length)
{
            io.emit(data.event, data.data)
        }
        data.rooms.forEach(function (room)
{
            io.of("/" + room).emit(data.event, data.data)
        });
    } catch (e)
{
        let dt = d.toString().split("}{")
        if (dt.length <= 1) return
        console.log("field split  ", dt)
        dt.forEach(field => {
            if (field[0] != "{") field = "{" + field
            if (field[field.length - 1] != "}") field += "}"
            try {
                data = JSON.parse(field)
                io.emit(data.type, field)
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
    , 500)
