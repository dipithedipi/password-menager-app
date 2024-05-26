<script lang="ts">
	import { masterPassword } from '$lib/store/passwordStore';
	import { decryptAES, hashSHA256 } from '$lib/logic/cryptography';
	import { waitfetchData } from '$lib/logic/fetch';
    import { onMount } from 'svelte';
    
    let compromised = 0;
    let same = 0;
    let totalPassword = 0;

    let compromisedPasswords:any = [];
    let samePasswords:any = []

    async function healthCheck() {
        console.log('Health check');
        let {data, success} = await waitfetchData('http://127.0.0.1:8000/password/health', 'GET', {});
        
        if (!success) {
            console.log('Failed to fetch data');
            return;
        }

        compromised = data.compromisedPasswordCount;
        same = data.samePasswordCount;
        totalPassword = data.totalPassword;
        compromisedPasswords = data.compromisedPassword;
        samePasswords = data.samePassword;
    }

    function checkFalseResultCompromised() {
        if (compromisedPasswords.length == 0) {
            return
        }

        compromisedPasswords.forEach((compromisedPassword: any) => {
            //console.log("password:", compromisedPassword.password);
            if (hashSHA256(decryptAES(compromisedPassword.password, $masterPassword)) !== compromisedPassword.LeakedHash) {
                //console.log("triggered");
                compromised -= 1;
                compromisedPasswords = compromisedPasswords.filter((item: any) => item !== compromisedPassword)
            }
            // console.log("password:", decryptAES(compromisedPassword.password, $masterPassword), 
            // "hash:", hashSHA256(decryptAES(compromisedPassword.password, $masterPassword)), 
            // "partial hash:", compromisedPassword.Partialhash,
            // "leaked hash:", compromisedPassword.LeakedHash);
        });
    }

    function checkFalseResultSame(): void {
        if (samePasswords.length === 0) {
            return;
        }

        // flat the nested array
        let samePasswordsFlatted = samePasswords.flat()
        console.log(samePasswordsFlatted)

        samePasswordsFlatted.forEach((samePassword: any) => {
            samePassword.fullHash = hashSHA256(decryptAES(samePassword.password, $masterPassword))
        })

        // remove the hashes that are unique
        let notUniqueHash = samePasswordsFlatted.filter((item: any) => samePasswordsFlatted.filter((item2: any) => item2.fullHash === item.fullHash).length > 1)

        // split the array in sub array where the hash is the same
        let samePasswordsGrouped = notUniqueHash.reduce((acc: any, item: any) => {
            let key = item.fullHash;
            if (!acc[key]) {
                acc[key] = [];
            }
            acc[key].push(item);
            return acc;
        }, {});

        samePasswords = Object.values(samePasswordsGrouped)
    }


    onMount(async () => {
        await healthCheck();
        checkFalseResultCompromised();
        checkFalseResultSame();
    });
</script>

<div class="h-screen bg-gray-700">
    <section class="pt-6 dark:bg-gray-700 flex flex-col justify-center items-center">
        <div class="max-w-screen-xl px-4 py-4 mx-auto text-center lg:px-6">
            <dl class="grid max-w-screen-md gap-7 mx-auto text-gray-900 sm:grid-cols-3 dark:text-white">
                <div class="flex flex-col dark:bg-gray-800 items-center justify-center p-2 border-2 border-dashed border-slate-500 rounded-md">
                    <dt class="mb-2 text-3xl md:text-4xl font-extrabold">{totalPassword}</dt>
                    <dd class="font-light text-gray-500 dark:text-gray-400">Total password</dd>
                </div>
                <div class="flex flex-col items-center dark:bg-gray-800 justify-center p-2 text-yellow-300 border-2 border-dashed border-slate-500 rounded-md">
                    <dt class="mb-2 text-3xl md:text-4xl font-extrabold">{same}</dt>
                    <dd class="font-light text-gray-500 dark:text-gray-400">Reused password</dd>
                </div>
                <div class="flex flex-col items-center dark:bg-gray-800 justify-center p-2 text-red-600 border-2 border-dashed border-slate-500 rounded-md">
                    <dt class="mb-2 text-3xl md:text-4xl font-extrabold">{compromised}</dt>
                    <dd class="font-light text-gray-500 dark:text-gray-400">Compromised</dd>
                </div>

                <div class="pl-2 flex items-center justify-center col-span-3">
                    <button type="button" class="px-3 py-2 mb-3 w-full mr-2 sm:mr:1 text-sm font-medium text-center text-white bg-blue-800 rounded-lg hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-400 dark:bg-blue-700 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
                        Check all passwords now
                    </button>
                </div>
            </dl>
        </div>
        <div>
            <h2 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">Compromised:</h2>
            {#each compromisedPasswords as compromisedPassword}
                <ul class="max-w-md space-y-1 pb-2 text-gray-500 list-disc list-inside dark:text-gray-400">
                    <li class="flex dark:bg-gray-800 items-center border-2 rounded-lg border-red-500 p-4">
                        <svg class="scale-125 w-3.5 h-3.5 me-2 text-red-500 dark:text-red-400 flex-shrink-0" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24">
                            <path d="M12.8659 3.00017L22.3922 19.5002C22.6684 19.9785 22.5045 20.5901 22.0262 20.8662C21.8742 20.954 21.7017 21.0002 21.5262 21.0002H2.47363C1.92135 21.0002 1.47363 20.5525 1.47363 20.0002C1.47363 19.8246 1.51984 19.6522 1.60761 19.5002L11.1339 3.00017C11.41 2.52187 12.0216 2.358 12.4999 2.63414C12.6519 2.72191 12.7782 2.84815 12.8659 3.00017ZM4.20568 19.0002H19.7941L11.9999 5.50017L4.20568 19.0002ZM10.9999 16.0002H12.9999V18.0002H10.9999V16.0002ZM10.9999 9.00017H12.9999V14.0002H10.9999V9.00017Z"></path>
                        </svg>
                        {decryptAES(compromisedPassword.username, $masterPassword)} for {compromisedPassword.website}
                    </li>
                </ul>
            {/each}
        </div>
        <div class="pt-6">
            {#each samePasswords as samePassword}
            <h2 class="mb-2 mt-2 text-lg font-semibold text-gray-900 dark:text-white">Same passwords:</h2>
                {#each samePassword as singlePassword}
                    <ul class="max-w-md space-y-1 pb-2 text-gray-500 list-disc list-inside dark:text-gray-400">
                        <li class="flex dark:bg-gray-800 items-center border-2 rounded-lg border-yellow-300 p-4">
                            <svg class="scale-125 w-3.5 h-3.5 me-2 text-yellow-500 dark:text-yellow-400 flex-shrink-0" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 24 24">
                                <path d="M12.8659 3.00017L22.3922 19.5002C22.6684 19.9785 22.5045 20.5901 22.0262 20.8662C21.8742 20.954 21.7017 21.0002 21.5262 21.0002H2.47363C1.92135 21.0002 1.47363 20.5525 1.47363 20.0002C1.47363 19.8246 1.51984 19.6522 1.60761 19.5002L11.1339 3.00017C11.41 2.52187 12.0216 2.358 12.4999 2.63414C12.6519 2.72191 12.7782 2.84815 12.8659 3.00017ZM4.20568 19.0002H19.7941L11.9999 5.50017L4.20568 19.0002ZM10.9999 16.0002H12.9999V18.0002H10.9999V16.0002ZM10.9999 9.00017H12.9999V14.0002H10.9999V9.00017Z"></path>
                            </svg>
                            {decryptAES(singlePassword.username, $masterPassword)} for {singlePassword.website}
                        </li>
                    </ul>
                {/each}
            {/each}
        </div>
    </section>
</div>