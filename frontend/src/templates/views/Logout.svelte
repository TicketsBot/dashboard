<script>
  import {clearLocalStorage, setDefaultHeaders} from '../includes/Auth.svelte'
  import axios from "axios";
  import {API_URL} from "../js/constants";
  import {errorPage} from "../js/util";

  setDefaultHeaders();

  async function process() {
      const res = await axios.post(`${API_URL}/logout`);
      if (res.status !== 204) {
          errorPage(res.status, res.data.error)
          return
      }

      clearLocalStorage()
      window.location.href = 'https://ticketsbot.net'
  }

  process()
</script>