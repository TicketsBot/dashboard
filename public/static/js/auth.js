const _tokenKey = 'token';

async function getToken() {
    let token = window.localStorage.getItem(_tokenKey);
    if (token == null) {
        let res = await axios.post('/token', undefined, {
            withCredentials: true,
            headers: {
                'x-tickets': 'true'
            }
        });

        if (res.status !== 200 || !res.data.success) {
            console.log("An error occurred whilst retrieving an authentication token. Please contact the developer");
            console.log(res);
            return;
        }

        token = res.data.token;
        localStorage.setItem(_tokenKey, token);
    }

    return token;
}

function clearLocalStorage() {
    window.localStorage.clear();
}

async function setDefaultHeader() {
    axios.defaults.headers.common['Authorization'] = await getToken();
    axios.defaults.headers.common['x-tickets'] = 'true'; // arbitrary header name and value
    axios.defaults.validateStatus = false;
}

async function _refreshToken() {
    window.localStorage.removeItem(_tokenKey);
    await getToken();
}

function addRefreshInterceptor() {
    axios.interceptors.response.use(async (res) => { // we set validateStatus to false
        if (res.status === 401) {
            await _refreshToken();
        }
        return res;
    }, async (err) => {
        if (err.response.status === 401) {
            await _refreshToken();
        }
        return err.response;
    });
}

setDefaultHeader();
//addRefreshInterceptor();