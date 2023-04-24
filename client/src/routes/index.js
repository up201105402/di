import Dashboard from '../pages/Dashboard.vue';
import SignIn from '../pages/LogIn.vue';
import SignUp from '../pages/SignUp.vue';

export const routes = [
    { path: '/', name: 'Dashboard', icon: 'glyphicon', component: Dashboard },
    { path: '/login', name: 'Log In', icon: 'glyphicon', component: SignIn },
    { path: '/signup', name: 'Sign Up', icon: 'glyphicon', component: SignUp },
];