// src/pages/Home.js
// import styles from "../styles/Main.module.css";
import styles from '../styles/Home.module.css';

export default function Home() {
    return `
        <div class="${styles.home}">
            <h1>Home Page</h1>
            <p>Welcome to the home page!</p>
            <a href="/about">Go to About</a>
        </div>
    `;
}
