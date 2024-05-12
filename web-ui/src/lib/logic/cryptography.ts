import { Argon2, Argon2Mode } from '@sphereon/isomorphic-argon2';
import * as CryptoJS from 'crypto-js';
import rsa from 'js-crypto-rsa';

// base64
function base64Encode(input: string): string {
    return btoa(input);
}

// argon2id
async function hashPasswordArgon2id(password: string, salt: string): Promise<string> {
    console.log(salt)

    const hash = await Argon2.hash(password, salt, {
        hashLength: 64,
        memory: 64 * 1024,
        parallelism: 2,
        mode: Argon2Mode.Argon2id,
        iterations: 3,
    })

    return hash.hex;
}

// RSA
async function generateRSAKeyPair(): Promise<{publicKeyRSA: JsonWebKey, privateKeyRSA: JsonWebKey}>{
    const keys = await rsa.generateKey(2048)
    return {publicKeyRSA:keys.publicKey, privateKeyRSA:keys.privateKey};
}

async function encryptRSA(text: string, publicKey: JsonWebKey): Promise<Uint8Array> {
    const textBytes = new TextEncoder().encode(text);
    const encrypted = await rsa.encrypt(
        textBytes,
        publicKey,
        'SHA-256', 
    )
    return encrypted;
}

async function decryptRSA(encrypted: Uint8Array, privateKey: JsonWebKey): Promise<string> {
    const decrypted = await rsa.decrypt(
        encrypted,
        privateKey,
        'SHA-256',
    )
    return new TextDecoder().decode(decrypted);
}

// AES
function encryptAES(text: string, passphrase: string): string {
  return CryptoJS.AES.encrypt(text, passphrase).toString();
}

function decryptAES(ciphertext: string , passphrase: string): string {
  const bytes = CryptoJS.AES.decrypt(ciphertext, passphrase);
  const originalText = bytes.toString(CryptoJS.enc.Utf8);
  return originalText;
}

export { hashPasswordArgon2id, encryptAES, decryptAES, encryptRSA, decryptRSA, generateRSAKeyPair, base64Encode};