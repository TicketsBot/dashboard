<script>
    import axios from "axios";
    import {redirectLogin, setToken} from '../includes/Auth.svelte'
    import {API_URL} from "../js/constants";
    import {errorPage} from '../js/util'

    async function process() {
        const code = new URLSearchParams(window.location.search).get('code')
        if (code === null) {
            redirectLogin()
        }

        axios.defaults.validateStatus = false
        axios.defaults.headers.common['x-tickets'] = 'true'
        const res = await axios.post(`${API_URL}/callback?code=${code}`)
        if (res.status !== 200) {
            errorPage(res.status, res.data.error)
            return
        }

        setToken(res.data.token)
        window.location.href = '/'
    }

    process()
</script>

<style>
  body {
      background-color: #121212;
  }
</style>