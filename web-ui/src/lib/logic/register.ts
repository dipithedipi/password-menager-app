import {hashPasswordArgon2id, base64DecodeBytes } from './cryptography';
import {exportCryptoKey, generateKeyPairRSA, decryptRSA} from './cryptographyRSA'

async function getSalt(email: string): (Promise<string|boolean>) {
    try {
        const response = await fetch('http://127.0.0.1:8000/user/register/salt', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
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

async function checkUsername(username: string): (Promise<boolean>) {
    try {
        const response = await fetch('http://127.0.0.1:8000/user/register/username', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error status: ${response.status}`);
        }

        const data = await response.json();
        return data.available;
    } catch (error) {
        console.error('Error checking username:', error);
        return false;
    }
}

async function register(email: string, username: string, password: string): Promise<{ success: boolean; message: string; otp:string }> {    
    // get salt
    const salt = await getSalt(email);
    if (salt === false) {
        return { success: false, message: 'Error fetching registation salt', otp:"" };
    }

    // generate RSA key pair
    const keyPair = await generateKeyPairRSA();
    const publicKey = await exportCryptoKey(keyPair.publicKey);
    console.log('publicKey:', publicKey);

    // hash password
    const hashedPassword = await hashPasswordArgon2id(password, salt.toString());
    console.log('hashedPassword:', hashedPassword);
    try {
        const response = await fetch('http://127.0.0.1:8000/user/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                username,
                email, 
                password: hashedPassword, 
                salt,
                publicKey: publicKey
            }),
        });

        const message = await response.text();
        if (!response.ok) {
            return { success: false, message: JSON.parse(message).message, otp:"" };
        }
        
        // decode the message to get the otpSecretUri
        let otpSecretUri = JSON.parse(message).otpSecretUri
        console.log("otpSecretUri:", otpSecretUri);
        
        // form base64 to byte array
        let otpSecretUriEnc = base64DecodeBytes(otpSecretUri);
        console.log("otpSecretUriEnc:", otpSecretUriEnc);
        const decryptOtpUri = await decryptRSA(keyPair.privateKey, otpSecretUriEnc);
        console.log("decryptOtpUri:", decryptOtpUri);

        return { success: true, message: message, otp: decryptOtpUri};
    } catch (error) {
        console.error('Error registering in:', error);
        return { success: false, message: 'Error registering in', otp:"" };
    }
}

export { register, checkUsername };