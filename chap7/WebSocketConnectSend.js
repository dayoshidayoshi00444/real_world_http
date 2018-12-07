var socket = new WebSocket('ws://game.example.com:12010/updates');
socket.onopen = () => {
    setInterval(() => {
        if (socket.bufferedAmonut === 0) {
            socket.send(getUpdateData());
        }
    }, 50);
};
