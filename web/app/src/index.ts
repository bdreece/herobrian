import 'htmx.org';
import 'htmx.org/dist/ext/head-support';
import 'htmx.org/dist/ext/sse';

import './index.css';

const cookies = document.cookie
    .split(';')
    .map(v => v.split('='))
    .reduce((acc, v) => {
        acc[decodeURIComponent(v[0].trim())] = decodeURIComponent(v[1].trim());
        return acc;
    }, {} as Record<string, string>)

window.addEventListener('htmx:configRequest', e => {
    (e as CustomEvent<{ headers: Record<string, string> }>).detail.headers['X-CSRF-Token'] = cookies['_csrf'];
})


export { }
