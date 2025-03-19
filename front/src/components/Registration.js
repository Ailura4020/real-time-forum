// src/pages/RegisterModal.js
import styles from '../styles/App.module.css';

export default function RegisterModal({ onRegister, onClose }) {
    return `
        <div class="${styles.modal}">
            <div class="${styles.modalContent}">
                <span class="${styles.close}" onclick="${onClose}">&times;</span>
                <h2>Register</h2>
                <form id="registerForm">
                    <input type="text" name="nickname" placeholder="Nickname" required />
                    <input type="number" name="age" placeholder="Age" required />
                    <select name="gender" required>
                        <option value="">Select Gender</option>
                        <option value="male">Male</option>
                        <option value="female">Female</option>
                    </select>
                    <input type="text" name="first_name" placeholder="First Name" required />
                    <input type="text" name="last_name" placeholder="Last Name" required />
                    <input type="email" name="email" placeholder="Email" required />
                    <input type="password" name="password" placeholder="Password" required />
                    <button type="submit">Register</button>
                </form>
                <div id="registerMessage"></div>
            </div>
        </div>
    `;
}
