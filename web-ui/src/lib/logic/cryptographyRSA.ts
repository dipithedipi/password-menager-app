function ab2str(buf: any) {
    return String.fromCharCode.apply(null, Array.from(new Uint8Array(buf)));
}

/*
    Export the given key and write it into the "exported-key" space.
*/
async function exportCryptoKey(key: any) {
    const exported = await window.crypto.subtle.exportKey("spki", key);
    const exportedAsString = ab2str(exported);
    const exportedAsBase64 = window.btoa(exportedAsString);
    
    //return `-----BEGIN PRIVATE KEY-----\n${exportedAsBase64}\n-----END PRIVATE KEY-----`;
    return `${window.btoa(exportedAsBase64)}`;
}

async function generateKeyPairRSA(): Promise<CryptoKeyPair> {
    let keyPair = await window.crypto.subtle.generateKey({
            name: "RSA-OAEP",
            modulusLength: 2048,
            publicExponent: new Uint8Array([1, 0, 1]),
            hash: "SHA-256",
        }, true, ["encrypt", "decrypt"])

    return keyPair
}

async function decryptRSA(privateKey: CryptoKey, data: ArrayBuffer): Promise<string> {
    console.log("decryptRSA");
    console.log(privateKey);
    console.log(data);
    try {
        let decrypted = await window.crypto.subtle.decrypt(
            {
                name: "RSA-OAEP"
            },
            privateKey,
            data
        );
        return new TextDecoder().decode(decrypted);
    } catch (error) {
        console.error("Errore durante la decrittografia");
        console.error(error);
        throw error;
    }

}

export { exportCryptoKey, generateKeyPairRSA, decryptRSA }