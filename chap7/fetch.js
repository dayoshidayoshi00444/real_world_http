fetch("news.json", {
    method: 'GET',
    mode: 'cors',
    credentials: 'include',
    cache: 'default',
    headers: {
        'content-Type': 'application/json'
    }
}).then((response) => {
    return response.json();
}).then((json) => {
    console.log(json);
});
