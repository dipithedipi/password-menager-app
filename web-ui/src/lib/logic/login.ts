import {hashPasswordArgon2id, encryptAES, decryptAES } from './cryptography';
import { setCookie, getCookie } from './utils';

async function getSalt(email: string): (Promise<string|boolean>) {
    try {
        const response = await fetch('http://127.0.0.1:8000/user/login/salt', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error status: ${response.status}`);
        }

        const data = await response.json();
        return data.salt;
    } catch (error) {
        console.error('Error fetching salt:', error);
        return false;
    }
}

async function login(email: string, password: string, salt: string, otpCode: string): Promise<{ success: boolean; message: string}> {    
    // hash password
    const hashedPassword = await hashPasswordArgon2id(password, salt);
    console.log('hashedPassword:', hashedPassword);
    try {
        const response = await fetch('http://127.0.0.1:8000/user/login', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
                'Access-Control-Allow-Credentials': 'true',
            },
            body: JSON.stringify({ email, password: hashedPassword, otp:otpCode }),
        });

        const message = await response.text();
        const parsed = JSON.parse(message);
        if (!response.ok) {
            return { success: false, message: parsed.message };
        }

        return { success: true, message: message };
    } catch (error) {
        console.error('Error logging in:', error);
        return { success: false, message: 'Error logging in' };
    }
}

// export
export { getSalt, login };