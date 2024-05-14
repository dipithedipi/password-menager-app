import { Argon2, Argon2Mode } from '@sphereon/isomorphic-argon2';
import * as CryptoJS from 'crypto-js';


// base64
function base64Encode(input: string): string {
    return btoa(input);
}

function base64Decode(input: string): string {
    return atob(input);
}

function base64DecodeBytes(input: string): Uint8Array {
    const str = atob(input);
    const bytes = new Uint8Array(str.length);
    for (let i = 0; i < str.length; i++) {
        bytes[i] = str.charCodeAt(i);
    }
    return bytes;
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

// AES
function encryptAES(text: string, passphrase: string): string {
  return CryptoJS.AES.encrypt(text, passphrase).toString();
}

function decryptAES(ciphertext: string , passphrase: string): string {
  const bytes = CryptoJS.AES.decrypt(ciphertext, passphrase);
  const originalText = bytes.toString(CryptoJS.enc.Utf8);
  return originalText;
}

export { hashPasswordArgon2id, encryptAES, decryptAES, base64Encode, base64Decode, base64DecodeBytes};