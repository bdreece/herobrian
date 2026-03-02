import axios from 'axios';
import { useToken } from '../composables/auth';

axios.defaults.baseURL = import.meta.env.BASE_URL;
axios.interceptors.request.use(config => {
    const token = useToken();
    if (token.value) {
        config.headers.setAuthorization(`Bearer ${token.value}`, true);
    }

    return config;
});
