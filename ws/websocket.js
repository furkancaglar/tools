const net = require("net")
const request = require("request")
const arglist = process.argv
if (arglist.length < 3) {
    throw new Error("no port name is specified")
}
else if (!/^\d{2,5}$/.test(arglist[2])) {
    console.log(arglist[2])
    throw new Error("invalid port")
}
let port = parseInt(arglist[2])
if (port > 65000) {
    throw new Error("invalid port")
}
const sock = new net.Socket()
const io = require("socket.io")(port)
const cmds = [
    "newgame",
    "ball",
    "winning",
    "endgame"
]
let token = ""
let balls = []
let stop = false

function retry(e) {
    inc++

    setTimeout(function () {
        if (0 == inc) return
        console.log(e)
        sock.connect({host: process.env.SOCKET_HOST || "localhost", port: process.env.SOCKET_PORT || 1111})
    }, 1000 * inc)

}

sock.connect({host: process.env.SOCKET_HOST || "localhost", port: process.env.SOCKET_PORT || 1111})
sock.on("data", function (d) {

    let data
    try {
        data = JSON.parse(d.toString())

        if (0 == data.rooms.length) {
            io.emit(data.event, data.data)
        }
        data.rooms.forEach(function (room) {
            io.of("/" + room).emit(data.event, data.data)
        });
    } catch (e) {
        let dt = d.toString().split("}{")
        if (dt.length <= 1) return
        console.log("field split  ", dt)
        dt.forEach(field => {
            if (field[0] != "{") field = "{" + field
            if (field[field.length - 1] != "}") field += "}"
            try {
                data = JSON.parse(field)
                io.emit(data.type, field)
            } catch (e) {
                console.error(e)
            }
        })
    }
})

sock.on("connect", function () {
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

function sendNewGameReq(token) {
    request.get("http://localhost:8080/api/tombala/game/new?token=" + token)
}

function sendBall(token) {
    var ball = getRandomBall()
    request.get(`http://localhost:8080/api/tombala/game/newball?ball=${ball}&token=${token}`)
}

function getRandomBall() {
    if (balls.length > 89) {
        balls = []
    }
    var ball
    do {
        ball = Math.floor(Math.random() * 90) + 1
    }
    while (balls.indexOf(ball) != -1)
    balls.push(ball)
    return ball
}

function getToken() {
    request.post("http://localhost:8080/api/tombala/user/login", {form: {username: "user1", password: "user1"}})
        .on("data", function (d) {
            const data = JSON.parse(d.toString())
            if (!data.token) throw new Error("no token")
            token = data.token
            sendNewGameReq(token)
        })
}
