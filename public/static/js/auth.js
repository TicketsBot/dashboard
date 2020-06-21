async function getToken() {
    let token = window.localStorage.getItem('token');
    if (token == null) {
        let res = await axios.post('/token', {
            withCredentials: true
        });

        if (res.status !== 200 || !res.data.success) {
            console.log("An error occurred whilst retrieving an authentication token. Please contact the developer");
            console.log(res);
            return;
        }

        token = res.data.token;
        localStorage.setItem('token', token);
    }

    return token;
}

function clearLocalStorage() {
    window.localStorage.clear();
}

async function setDefaultHeader() {
    const token = await getToken();
    axios.defaults.headers.common['Authorization'] = token;
    axios.defaults.validateStatus = false;
}

setDefaultHeader();