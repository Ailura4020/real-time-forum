// src/router.js
const routes = {
    '/': 'Home',
    '/about': 'About',
};

export function router() {
    const path = window.location.pathname;
    const page = routes[path] || 'Home';
    return page;
}
