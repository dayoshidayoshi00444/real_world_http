const evtSource = new EventSource("ssedemo.php");

// メッセージのイベントハンドラ
evtSource.onmessage = (e) => {
    const newElement = document.createElement("li");
    newElemnt.innerHTML = "message: " + e.data;{
        eventList.appendChild(newElement);
    }
};

evtSource.addEventListener("ping", (e) => {
    const newElement = document.createElement("li");

    const obj = JSON.parse(e.data);
    newElemnt.innerHTML = "ping at " + obj.time;
    eventList.appendChild(newElement);
}, false);


